# watchdog

The `watchdog` package provides a lightweight liveness tracker for the portwatch monitor loop.

## Overview

A `Watchdog` expects periodic `Beat()` calls from the monitor goroutine. If no beat is received within the configured timeout, the status transitions to `stale`, indicating the loop may have hung or crashed.

## Usage

```go
dog := watchdog.New(30 * time.Second)

// Inside the monitor loop:
dog.Beat()

// Check health elsewhere:
if dog.Status() == watchdog.StatusStale {
    log.Println("monitor loop appears stuck")
}
```

## Reporter

`Reporter` periodically logs the watchdog status to stdout (or any `io.Writer`):

```go
rep := watchdog.NewReporter(dog, 10*time.Second)
rep.Start()
defer rep.Stop()
```

## Statuses

| Status    | Meaning                              |
|-----------|--------------------------------------|
| `unknown` | No heartbeat received yet            |
| `healthy` | Last beat within timeout window      |
| `stale`   | Last beat exceeded timeout threshold |
