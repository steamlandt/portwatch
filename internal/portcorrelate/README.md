# portcorrelate

Tracks which ports tend to change state together within the same scan cycle.

## Purpose

When multiple ports open or close simultaneously it often indicates a planned
event such as a service restart rather than an unexpected intrusion. By
recording co-occurrence counts, `portcorrelate` lets the monitor distinguish
correlated, expected changes from isolated anomalies.

## Usage

```go
corr := portcorrelate.New(time.Second)

// Call Observe for every pair of ports that changed in the same cycle.
corr.Observe(80, 443)
corr.Observe(80, 443)

// Retrieve pairs seen together at least twice.
strong := corr.Strong(2)
for _, p := range strong {
    fmt.Printf("ports %d and %d co-occurred %d time(s)\n", p.A, p.B, p.Count)
}
```

## API

| Method | Description |
|---|---|
| `New(window)` | Create a new Correlator |
| `Observe(a, b)` | Record that ports a and b changed together |
| `All()` | Return all tracked pairs |
| `Strong(n)` | Return pairs with count ≥ n |
| `Reset()` | Clear all recorded data |
