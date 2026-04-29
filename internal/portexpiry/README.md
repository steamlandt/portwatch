# portexpiry

The `portexpiry` package tracks how long each open port has been continuously
observed and surfaces ports that have exceeded a configured maximum age.

## Purpose

Some ports are expected to be short-lived (e.g. ephemeral services, debug
endpoints). `portexpiry` lets portwatch flag ports that have been open longer
than a policy threshold, enabling operators to investigate or rotate them.

## Usage

```go
import "github.com/user/portwatch/internal/portexpiry"

// Warn when a port has been open for more than 24 hours.
e := portexpiry.New(24 * time.Hour)

// Call Observe each scan cycle for every open port.
for _, p := range openPorts {
    e.Observe(p)
}

// Call Remove when a port closes so its timer resets.
for _, p := range closedPorts {
    e.Remove(p)
}

// Retrieve ports that have been open too long.
expired := e.Expired(openPorts)
for _, entry := range expired {
    fmt.Printf("port %d/%s open for %v\n", entry.Port.Number, entry.Port.Proto, entry.Duration)
}
```

## API

| Function | Description |
|---|---|
| `New(maxAge)` | Create a new Expiry tracker with the given age threshold |
| `Observe(port)` | Record a port as open; no-op if already tracked |
| `Remove(port)` | Stop tracking a port (e.g. after it closes) |
| `Expired(ports)` | Return entries for ports open longer than maxAge |
| `Age(port)` | Return how long a port has been tracked |
