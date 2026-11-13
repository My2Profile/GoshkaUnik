# Week 8 — Goroutines & the Go Scheduler: 5-minute defence

## 1. The lecture in 60 seconds
- **Concurrency = structure** (independent tasks), **parallelism = execution** (tasks really run at the same time on different cores).
- **I/O-bound** work (waiting) speeds up with many goroutines even on 1 core. **CPU-bound** work speeds up only with more cores.
- `go f()` starts a goroutine. **main does not wait**: when `main` returns, everything stops, no `defer` runs. Use `sync.WaitGroup`.
- Goroutine = 2 KB growing stack, created and switched by the Go runtime in user space (cheap, ~500 ns). OS thread = 1–8 MB, switched by the kernel.
- **Scheduler G-M-P**: G = goroutine, M = OS thread, P = permission to run Go code with a local run queue. Exactly `GOMAXPROCS` P's (default = number of cores).
- Idle P: checks global queue, then **steals half** of another P's queue (work stealing).
- G blocks on channel / mutex / `time.Sleep` → G is parked, M takes another G. G in a blocking **syscall** → M is stuck in the kernel, P is handed to another M (so new threads get created).
- **Preemption** (Go 1.14+): a goroutine running ~10 ms is interrupted by a signal, so a busy loop cannot freeze the program.
- **Leak**: a goroutine that waits forever (nobody reads the channel, no way to stop). Stack and referenced memory are never freed. Fix: buffered channel, `context`/done channel. Find: `runtime.NumGoroutine()`, goroutine pprof profile.

## 2. Task 2 — what I did and why
**Part 1 (`go run . say`)** — original prints only `main 0..2`: main finishes and kills the goroutine. Fix: `wg.Add(1)` before `go`, `defer wg.Done()` inside, `wg.Wait()` at the end of main. Now both loops always finish (order of lines may vary — that's fine).

**Part 2 (`go run . broken`, `-race`, `atomic`, `slice`)**
- Mistake 1: `wg.Add(1)` inside the goroutine → `Wait()` may see counter 0 and return before the goroutines even started. Fix: `Add` before `go`.
- Mistake 2: `counter++` = read, add, write; two goroutines read the same value and both write n+1 → lost updates (I got 920, 921, 897 instead of 1000). `-race` prints `WARNING: DATA RACE`.
- Fix A: `atomic.Int64` with `counter.Add(1)` → always 1000, no race.
- Fix B: `results := make([]int, 1000)`, goroutine `i` writes only `results[i]`, sum after `Wait()` → always 1000, no race (different elements = different memory).

**Part 3 (`go run . panic`, `go run . recover`)** — a panic in a goroutine crashes the whole program; `recover` in `main` does nothing because `recover` only catches panics of its **own** goroutine. The `defer func(){ recover() }()` must be **inside the goroutine** that may panic. Then the goroutine ends, main stays alive.

**Participation: why is "worker 5" printed first?** Every P has a special slot **`runnext`** for the most recently created goroutine. The loop creates workers 1..5; each new `go` kicks the previous goroutine from `runnext` into the local queue, so the last one (5) sits in `runnext` and runs first when main blocks in `wg.Wait()`. Then 1, 2, 3, 4 from the queue in FIFO order. With `GOMAXPROCS=1` it is the same every time (5 1 2 3 4) because there is no other P to steal work. With 2+ cores other P's can steal 1–4, so the order can change a bit, but 5 is still usually first. **This order is not guaranteed by the Go spec** — relying on it is a bug.

## 3. Quick answers to the theory questions
- **Concurrency vs parallelism**: concurrency — web server handles 1000 requests waiting for the DB on 1 core; parallelism — the same server on 8 cores really computes 8 requests at once.
- **main returns** → program exits, all goroutines are stopped immediately, no defers run.
- **`wg.Add(1)` before `go`**: otherwise `wg.Wait()` can run first, see 0 and return early.
- **Data race**: two goroutines access the same variable at once, at least one writes, no synchronization. Find: `go run -race`, `go test -race`.
- **Why cheaper than threads**: 2 KB vs MBs of stack, switching in user space (no kernel), created by runtime.
- **G, M, P / GOMAXPROCS**: see above; GOMAXPROCS = number of P = how many goroutines run Go code at the same time.
- **Work stealing / blocking syscall**: idle P steals half of another queue; syscall blocks M, P is handed off to another M.
- **Preemption**: before 1.14 a busy loop without function calls could hang the program; since 1.14 signal-based preemption (~10 ms).
- **Leak**: causes — unbuffered send nobody reads, worker loop without stop; find — `NumGoroutine()`, `pprof.Lookup("goroutine")`, `goleak`.
