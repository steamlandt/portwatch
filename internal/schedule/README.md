# schedule

Provides a simple interval-based scheduler used by the portwatch daemon to trigger periodic port scans.

## Usage

```go
s := schedule.New(30 * time.Second)
go s.Start(func() {
    // run a scan cycle
})

// later...
s.Stop()
```

## API

| Function | Description |
|----------|-------------|
| `New(interval)` | Create a new Scheduler |
| `Start(fn)` | Begin ticking; blocks until stopped |
| `Stop()` | Halt the scheduler |
| `Interval()` | Return the configured interval |

The scheduler is safe to stop exactly once.
