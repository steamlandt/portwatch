# porttrend

Tracks how frequently individual ports appear across successive scans, distinguishing stable long-running services from transient or ephemeral listeners.

## Concepts

| Term | Meaning |
|------|---------|
| **Observation** | A single sighting of a port during one scan cycle |
| **Seen count** | Total number of scans in which the port was observed |
| **Transient** | A port whose seen count is below a configurable threshold |

## Usage

```go
tr := porttrend.New()

// Call Observe once per port per scan cycle.
tr.Observe("80",   time.Now())
tr.Observe("443",  time.Now())
tr.Observe("8080", time.Now())

// Retrieve the trend entry for a specific port.
if e, ok := tr.Get("80"); ok {
    fmt.Printf("port 80 seen %d times\n", e.Seen)
}

// List ports seen fewer than 3 times (likely transient).
for _, e := range tr.Transient(3) {
    fmt.Printf("transient: %s (seen %d)\n", e.Port, e.Seen)
}
```

## Reporting

```go
r := porttrend.NewReporter(tr, os.Stdout)
r.Report()              // full table sorted by frequency
r.ReportTransient(3)    // only low-frequency ports
```
