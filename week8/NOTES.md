# Week 8 — Tasks 1, 3, 4, 5: results and explanations (my machine: 2 cores, Go 1.24; yours will differ)

Layout: `task1/ task2/ task3/ task4/{sched,preempt,threads}/ task5/`, one module `week8` (run from the module root).

## Task 1 — concurrency vs parallelism (`go run ./task1`, `GOMAXPROCS=1 go run ./task1`)
NumCPU = 2.

| Work | Default (GOMAXPROCS=2) | GOMAXPROCS=1 |
|---|---|---|
| CPU, 1 goroutine | 9.80 s | 9.77 s |
| CPU, 2 goroutines (= cores) | 7.22 s | 9.83 s |
| CPU, 4 goroutines (2x cores) | 5.55 s | 9.87 s |
| CPU, 8 goroutines | 5.10 s | 9.89 s |
| I/O, 1 / 10 / 1000 goroutines | 100 / 101 / 102 ms | 100 / 100 / 103 ms |

- GOMAXPROCS=1: only one goroutine runs Go code at a time, so CPU work can't speed up. 1000 sleeping goroutines don't use the CPU, they are parked, so all wait the same 100 ms together.
- Speed-up with 2 goroutines is not 2x because the ranges are unequal: checking bigger numbers is harder, so the goroutine with the upper half has more work and finishes last. With more (smaller) chunks the load balances better (4 and 8 are faster than 2), but past the number of cores the gain stops: only 2 can run at once.
- **Participation**: concurrency without parallelism — a web server on 1 core handles many requests that mostly wait for DB/network. Concurrency with parallelism — the same server on many cores, or a request that resizes images, runs several handlers truly at the same time on different cores.

## Task 3 — cost of a goroutine (`go run ./task3`, `go run ./task3 1000000`, `go test ./task3 -bench . -benchmem`)
- 100 000 goroutines: NumGoroutine = 100001, StackSys grew 205 MB = **2.00 KB per goroutine**, started in **~228 ms**.
- 1 000 000 goroutines: ~2 GB stack, started in **~2.2 s**. After `close(stop)` only 1 goroutine remains.
- Benchmark: `BenchmarkGoroutineStartStop` 397.6 ns/op (16 B, 1 alloc); `BenchmarkFuncCall` 1.84 ns/op (0 B). A goroutine is **~200x more expensive than a plain call** (in the lecture 509 vs 19.7 ns ≈ 26x because the called function there is not trivial).
- Not worth starting one when the work is tiny (a few ns to a microsecond): the start/stop and scheduling cost is bigger than the work (see Task 7), or when the result is needed immediately and nothing runs in parallel.

## Task 4 — the scheduler
**4.1 schedtrace** (`GODEBUG=schedtrace=500 ./sched`, binary built first so `go build` output isn't mixed in):
`SCHED 501ms: gomaxprocs=2 idleprocs=0 threads=4 spinningthreads=0 needspinning=1 idlethreads=1 runqueue=4 [3 0]`
- `501ms` time since start; `gomaxprocs=2` number of P; `idleprocs=0` both P are busy running a G; `threads=4` OS threads (M) created; `spinningthreads=0` no M is looking for work now; `idlethreads=1` one M is sleeping; `runqueue=4` goroutines in the **global** queue; `[3 0]` goroutines in the **local** queue of P0 and P1.
- 2 running + 4 global + 3 local = 8 busy goroutines on 2 cores: only 2 run, the rest wait.
- With `GOMAXPROCS=1`: `gomaxprocs=1`, `threads=3`, one local queue `[5]`, the global queue holds the rest. Only one G runs; others wait. Total time grows (8 x 2 s = 16 s instead of ~8 s).

**4.2 preemption**: normally prints `main is still alive` (the runtime preempts the busy loop after ~10 ms and main wakes up). With `GODEBUG=asyncpreemptoff=1` it **hangs forever** (my run: killed by `timeout`, exit 124): with one P and no function calls in the loop there is no point where the goroutine can yield. Before Go 1.14 preemption was only cooperative (at function calls) and this always hung; since 1.14 the runtime sends a signal to the thread.

**4.3 threads** (`go run ./task4/threads`): start 4 → after 100 x `time.Sleep` still 4 → after 100 x `syscall.Nanosleep` **102**.
`time.Sleep`: the G is **parked** by the runtime timer, the M immediately takes another G, no thread is blocked. A blocking syscall blocks the **whole M** in the kernel; the runtime hands the **P** to another M so other Gs can run — but if 100 Gs are blocked, 100 M are stuck, so it creates ~100 new threads.

**4.4 "worker 5 first"**: each P has a `runnext` slot for the most recently created goroutine. In the loop, every `go` puts the new G in `runnext` and moves the previous one to the local queue. After the loop, `runnext` = worker 5, queue = 1,2,3,4. When main blocks in `wg.Wait()` the P runs `runnext` first (5), then the queue in order (1 2 3 4). With `GOMAXPROCS=1` it is always `5 1 2 3 4`. With several P an idle P can steal 1–4, so the order can vary. Don't rely on it: the spec promises nothing about order.

## Task 5 — leaks (`go run ./task5 leak|profile|fixed`, `go test ./task5`)
1. **Count**: start 1 → after 100 `fetchLeaky` **101** → after 10 leaky tickers **111**.
2. **Find**: `debug=1` shows `100 @ ... main.fetchLeaky.func1 .../leak.go:19` (the line `ch <- slowQuery()`) and `10 @ ... time.Sleep ... startTickerLeaky.func1 leak.go:47`. `debug=2` shows the state: `goroutine 18 [chan send]:` (blocked sending, nobody receives) for the fetch leak; ticker goroutines are in `[sleep]`.
3. **Fix**: `fetchFixed` uses `make(chan int, 1)`; `startTicker(ctx)` returns on `ctx.Done()`. Output: start 1 → after 100 `fetchFixed` **1** → 10 tickers 11 → after `cancel()` **1**.
4. **Test** `TestFetchDoesNotLeak`: NumGoroutine before, call 100 times with 1 ms timeout, sleep 200 ms, compare. Passes for `fetchFixed` (also with `-race`). With `LEAKY=1 go test ./task5` it **fails**: `goroutine leak: before=2 after=102 (100 leaked)`.
