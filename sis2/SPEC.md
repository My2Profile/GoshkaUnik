# SIS - build plan (staged, like the Dara task)

The same text-analyzer task service, grown one stage at a time. Every stage
ends with something that runs and is tested. Stages 1-4 are exactly what the
course grades (CP1 hardening, CP2, CP3). Stage 5 is optional extra learning.

Working rule: YOU type every line. Claude explains, reviews, and asks defense
questions. Turn off inline AI suggestions while you do the stages.

Spring -> Go map (you already know the Spring side):

| Spring (Daras)         | Go (SIS)                                        |
|------------------------|-------------------------------------------------|
| Entity                 | `Task`, `Result` structs                        |
| DTO                    | request/response structs with json tags         |
| Repository             | `Store` interface + in-memory implementation    |
| Service                | `Service`                                       |
| Controller             | `http.Handler` functions                        |
| @ExceptionHandler      | one function that maps errors to HTTP statuses  |
| application.properties | `config.go` (env vars)                          |
| docker-compose         | docker-compose.yml (stage 5)                    |

---

## Stage 1 - Harden CP1 (Week 6 + Week 7 material)

Goal: a task can never crash the program, errors carry detail, memory use is
measured.

1. **Custom error type.** `*ValidationError{Field string, Reason string}`
   returned by `CreateTask` for bad input, still matching `ErrEmptyInput`
   (implement `Unwrap`). Callers use `errors.As` to read the field.
2. **Panic safety.** `ProcessTask` must survive a `Processor` that panics.
   Use `defer` + `recover` and a named result `err`. The task goes back to
   `Pending` and the caller gets `process task %d: panic: <value>`.
3. **Function values.** Add `type ProcessorFunc func(string) (*Result, error)`
   with a `Process` method, so any plain function is a `Processor`
   (this is how `http.HandlerFunc` works).
4. **Tests you do not have yet:**
   - analyzer: `"Go go GO"` -> 1 unique word, empty-after-trim, multi-line
   - fake `Processor` that returns an error -> task is `Pending` again
   - fake `Processor` that panics -> no crash, task is `Pending` again
   - `errors.Is` for every sentinel, `errors.As` for `*ValidationError`
5. **Memory.** Benchmark `TextAnalyzer.Process` with `-benchmem`. Use
   `testing.AllocsPerRun`. Run `go build -gcflags='-m -l' ./...` and
   explain every `escapes to heap` line in `internal/task`.
   Then reduce allocations (hint: `make(map[string]struct{}, n)`), re-measure,
   write the before/after numbers in the README.

Definition of done: `gofmt`, `go vet`, `go test ./...` all green.

## Stage 2 - Storage behind an interface

1. `Store` interface: `Save(*Task)`, `Get(id) (*Task, error)`, `List() []*Task`.
2. `MemoryStore` implements it (map inside). `Service` depends on `Store`,
   not on a map.
3. Decide and document: does `Get` return the stored pointer or a copy?
   (Copy is safer once goroutines appear. Know why.)
4. Tests run against `MemoryStore`; add a compile-time check
   `var _ Store = (*MemoryStore)(nil)`.

## Stage 3 - CP2: concurrent processing

1. Worker pool: N goroutines reading from a **bounded** channel (the queue).
2. `Submit(id)` never blocks forever: queue full -> return `ErrQueueFull`.
3. `sync.RWMutex` (or channel ownership) around shared state. State
   transitions must be safe when two workers touch different tasks.
4. Lifecycle: `Start()`, `Shutdown()` that stops accepting work, lets
   workers drain, waits with `WaitGroup`. Nobody leaks.
5. Prove it: `go test -race ./...` clean. Write a test that submits 1000
   tasks from many goroutines.

Defense questions to be able to answer: who owns each channel, who closes it,
when does each goroutine end, what would the race detector say if the mutex
were removed.

## Stage 4 - CP3: context, HTTP API, graceful shutdown

1. `ProcessTask(ctx, id)`: honors cancellation and a per-task timeout
   (`context.WithTimeout`).
2. `net/http` API (standard library only):
   - `POST /tasks` `{"input": "..."}` -> 201 + task JSON
   - `GET /tasks/{id}` -> 200 / 404
   - `GET /tasks` -> 200 list
   - `GET /healthz`
3. One error-to-status function with `errors.Is` / `errors.As`:
   `ErrEmptyInput`/`*ValidationError` -> 400, `ErrNotFound` -> 404,
   `ErrQueueFull` -> 503, anything else -> 500. JSON error body.
4. Graceful shutdown: `signal.NotifyContext`, `server.Shutdown(ctx)`, then
   pool `Shutdown()`.
5. `httptest` tests for each endpoint and each error status.

## Stage 5 - Optional (the Dara-style extras)

- Config from env vars (`PORT`, `WORKERS`, `QUEUE_SIZE`).
- `Dockerfile` (multi-stage, like the teacher's go-k8s-service).
- Postgres `Store` with `database/sql`, `docker-compose.yml` with a
  postgres service, a `Makefile`.
- `GOMEMLIMIT` in the container, `pprof`, `GODEBUG=gctrace=1` on a load test.

---

## How each stage goes (same loop as Daras)

1. Claude explains the idea in plain words and what the finished piece
   looks like.
2. You write it in the real project.
3. You run `gofmt`, `go vet`, `go test`; paste the output or the code.
4. Claude reviews, points at the line, lets you fix it.
5. Claude asks 3 defense questions. Next stage starts only when you can
   answer them.
