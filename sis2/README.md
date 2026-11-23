# SIS - Checkpoint 2: Concurrent Text Analyzer Tasks

An in-memory Go service that analyzes text. A **task** is one piece of text.
CP2 adds a **worker pool**: tasks are submitted to a bounded queue and processed
concurrently by several workers, with safe shutdown.
Golang Application Development course (KBTU, Fall 2026).

## Result of a task

| Field         | Meaning                                               |
|---------------|-------------------------------------------------------|
| `Words`       | whitespace-separated words                            |
| `Runes`       | characters (runes, not bytes)                         |
| `Lines`       | lines (count of `\n` plus one)                        |
| `UniqueWords` | distinct words, case-insensitive                      |

## Project structure

```
cmd/sis/main.go          demo: start pool, submit tasks, shutdown, print results
internal/task/
  errors.go              sentinel errors
  task.go                Task, State, Result, transition rules, clone
  analyzer.go            Processor interface and TextAnalyzer
  service.go             Service: create/get/process, worker pool, Submit, Shutdown
  *_test.go              unit and concurrency tests
```

## How to run

Requires Go 1.27 (and gcc for the race detector on Windows).

```
go run ./cmd/sis
go test ./... -race
```

## Task lifecycle

```
Pending -> Queued -> Processing -> Completed
   ^         |            |
   +---------+            +-----> Failed
```

| From         | Allowed next        |
|--------------|---------------------|
| `Pending`    | `Queued`, `Processing` |
| `Queued`     | `Processing`, `Pending` |
| `Processing` | `Completed`, `Failed`  |
| `Completed`, `Failed` | none (final)  |

`Pending -> Processing` keeps the synchronous `ProcessTask` from CP1 working.
`Queued -> Pending` is the roll-back when the queue is full.

## Concurrency design

- **Queue:** a buffered `chan int` of task IDs. Its capacity bounds the
  backlog (backpressure). `Submit` never blocks: it uses `select` with
  `default` and returns `ErrQueueFull` when the queue is full, rolling the
  task back to `Pending`.
- **Workers:** `Start(workers, queueSize)` launches goroutines that
  `range` over the queue. Each worker calls `ProcessTask`.
- **Synchronization:** one `sync.RWMutex` protects the task map and every
  task. Writers take `Lock`, readers take `RLock`.
- **No lock during processing:** `ProcessTask` has three phases - claim the
  task under the lock (`Pending/Queued -> Processing`), run the `Processor`
  with **no lock held** (so workers really overlap), then store the result
  under the lock. Because the state change happens under the lock, two
  workers can never process the same task.
- **Snapshots:** `CreateTask`, `GetTaskByID` and `ProcessTask` return deep
  copies (`clone`, including the `Result`), so callers never share memory
  with the workers.
- **Panic safety:** the `Processor` call is wrapped in `recover`. A panic
  becomes an error, the task is marked `Failed` with the message, and the
  worker keeps running.
- **Shutdown:** sets `stopped` and closes the queue under the same mutex that
  `Submit` uses (so nobody can send on a closed channel), then **releases the
  lock** and waits on a `sync.WaitGroup` for all workers. Tasks already
  queued are finished. It is idempotent. `Submit` after `Shutdown` returns
  `ErrStopped`; `Submit` before `Start` returns `ErrNotStarted`.

## Errors

Wrapped with `%w`; use `errors.Is`:

`ErrEmptyInput`, `ErrNotFound`, `ErrInvalidTransition`, `ErrQueueFull`,
`ErrStopped`, `ErrNotStarted`.

## Tests

`go test ./... -race` covers: state transitions, analyzer cases (incl.
Cyrillic), snapshot semantics, workers overlapping in time (gate processor),
queue full + rollback, submit before start / after shutdown, idempotent
shutdown, shutdown draining the queue, panic containment, and 50 goroutines
using the service at once.
