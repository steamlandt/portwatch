# throttle

The `throttle` package limits how frequently an alert can fire for a given port key, preventing alert storms when a port flaps rapidly.

## Usage

```go
th := throttle.New(5 * time.Minute)

if th.Allow("tcp:8080") {
    // send alert
}
```

## Behaviour

- The first call for a key is always allowed.
- Subsequent calls within the configured interval are blocked.
- After the interval elapses the key is allowed again.
- `Reset(key)` clears state for a single key.
- `ResetAll()` clears all tracked keys.

## Thread Safety

All methods are safe for concurrent use.
