# SIS - Checkpoint 1: Text Analyzer Tasks

A small, synchronous, in-memory Go application. A **task** is one unit of
work: a piece of text that gets analyzed. The program creates tasks, validates
their input, processes them, stores the results, and tracks each task's state.
Written for the Golang Application Development course (KBTU, Fall 2026).

## What it does

Input: a text string, for example `"Hello Go world"`.
Processing: text analysis. Result:

| Field         | Meaning                                              |
|---------------|------------------------------------------------------|
| `Words`       | number of whitespace-separated words                 |
| `Runes`       | number of characters (runes, not bytes)              |
| `Lines`       | number of lines (count of `\n` plus one)             |
| `UniqueWords` | number of distinct words, compared case-insensitively |

## Project structure

```
cmd/sis/main.go          small demo program
internal/task/
  errors.go              sentinel errors
  task.go                Task, State, Result, state transition rules
  analyzer.go            Processor interface and TextAnalyzer
  service.go             in-memory Service: create, get, process tasks
  *_test.go              unit tests
```

## How to run

Requires Go 1.27.

```
go run ./cmd/sis
go test ./...
```

## Task lifecycle

```
Pending  ->  Processing  ->  Completed
```

| From         | Allowed next state |
|--------------|--------------------|
| `Pending`    | `Processing`       |
| `Processing` | `Completed`        |
| `Completed`  | none (final)       |

- A new task starts as `Pending` with no result.
- `ProcessTask` moves it to `Processing`, runs the analyzer, stores the
  result, then moves it to `Completed`.
- Any other transition is rejected with `ErrInvalidTransition`. For example,
  processing a `Completed` task again returns an error.
- If the analyzer returns an error, the task goes back to `Pending` and the
  error is returned.

## Errors

Expected failures are returned as errors (no `panic`). They are wrapped with
`%w`, so callers can use `errors.Is`:

- `ErrEmptyInput` - the input is empty or only whitespace
- `ErrNotFound` - no task with the requested ID
- `ErrInvalidTransition` - the requested state change is not allowed

## Design notes

- `Processor` is an interface (`Process(input string) (*Result, error)`);
  `TextAnalyzer` implements it, so other task domains could be plugged in.
- Tasks are stored as `map[int]*Task`. Pointers are used so that state
  changes made by the service persist in the map.
- `Task.transition` has a pointer receiver because it modifies the task.
- IDs come from a simple counter. The service is not safe for concurrent
  use; concurrency is part of a later checkpoint.
