# portevict

The `portevict` package tracks **port eviction events** — ports that were observed as open and subsequently closed. It maintains a rolling log of these events with timing metadata.

## Concepts

- **Eviction**: A port that transitions from open → closed after having been open for a measurable period.
- **Entry**: A single eviction record containing the port, when it was first seen open, when it closed, and the computed duration.

## Usage

```go
tracker, err := portevict.New("/var/lib/portwatch/evictions.json")
if err != nil {
    log.Fatal(err)
}

// When a port closes:
err = tracker.Record(port, openedAt, time.Now())

// Retrieve all eviction events:
events := tracker.All()
for _, e := range events {
    fmt.Printf("port %d was open for %s\n", e.Port.Number, e.Duration)
}

// Clear the log:
tracker.Reset()
```

## Persistence

Pass a file path to `New` to enable JSON persistence. Pass an empty string for in-memory-only operation.

## Thread Safety

All methods are safe for concurrent use.
