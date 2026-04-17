# suppress

The `suppress` package provides time-based alert suppression for portwatch.

## Purpose

When a port remains in an unexpected state (e.g. stays open), repeated scans
would otherwise generate continuous alerts. The suppressor ensures that an
alert for a given key is only forwarded once per suppression window.

## Usage

```go
s := suppress.New(10 * time.Minute)

key := fmt.Sprintf("port:%d:%s", port.Number, port.State)
if !s.IsSuppressed(key) {
    alert.Send(port)
}
```

## Behaviour

- The first call to `IsSuppressed(key)` always returns `false` and records the time.
- Subsequent calls within the window return `true` (suppressed).
- After the window expires the next call returns `false` and resets the timer.
- `Reset(key)` removes a single key; `ResetAll()` clears all state.

## Thread Safety

All methods are safe for concurrent use.
