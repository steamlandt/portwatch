# audit

The `audit` package provides a persistent, append-only JSONL log of every port
change event detected by portwatch.

## Format

Each line in the log file is a self-contained JSON object:

```json
{"timestamp":"2024-01-15T10:23:00Z","event":"opened","port":8080,"protocol":"tcp","label":"http-alt"}
```

## Usage

```go
log := audit.New("/var/log/portwatch/audit.jsonl")

// Record an event
log.Record("opened", "tcp", 8080, "http-alt")

// Read all recorded entries
entries, err := log.All()
```

## Fields

| Field       | Type   | Description                          |
|-------------|--------|--------------------------------------|
| `timestamp` | string | RFC3339 UTC time of the event        |
| `event`     | string | `"opened"` or `"closed"`             |
| `port`      | number | Port number (1–65535)                |
| `protocol`  | string | Transport protocol, e.g. `"tcp"`     |
| `label`     | string | Human-readable service name, if any  |

## Notes

- The file is created automatically on first write.
- Concurrent writes are safe; a mutex serialises access.
- `All` returns an empty slice (no error) when the log file does not yet exist.
