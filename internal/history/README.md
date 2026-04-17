# history

The `history` package provides an append-only log of port change events detected by portwatch.

## Usage

```go
h, err := history.New("/var/lib/portwatch/history.json")
if err != nil {
    log.Fatal(err)
}

// Record an event
h.Record(8080, "tcp", "opened")

// Read all events
for _, e := range h.All() {
    fmt.Printf("%s  port=%d proto=%s event=%s\n",
        e.Timestamp.Format(time.RFC3339), e.Port, e.Proto, e.Event)
}
```

## Persistence

Pass a file path to `New` to enable JSON persistence. Entries are appended and
flushed to disk after every `Record` call. Pass an empty string to keep the log
in memory only.

## Entry fields

| Field       | Type     | Description                        |
|-------------|----------|------------------------------------|
| `timestamp` | RFC3339  | UTC time the event was recorded    |
| `port`      | int      | Port number                        |
| `proto`     | string   | Protocol (`tcp` or `udp`)          |
| `event`     | string   | `"opened"` or `"closed"`           |
