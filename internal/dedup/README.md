# dedup

The `dedup` package provides simple deduplication for port change events.

## Purpose

When the monitor polls ports repeatedly, the same open/closed state may be
detected across multiple cycles. Without deduplication this would fire
redundant alerts for ports whose state has not actually changed.

`Deduplicator` tracks the last reported state per `proto:port` key and
suppresses events where the state is unchanged since the last report.

## Usage

```go
d := dedup.New()

if !d.IsDuplicate("tcp", 8080, "open") {
    // fire alert
}
```

## Thread Safety

All methods are safe for concurrent use.
