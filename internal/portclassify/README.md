# portclassify

Assigns IANA port range classifications to observed ports.

## Classes

| Class        | Range         | Description                        |
|--------------|---------------|------------------------------------|
| `system`     | 0 – 1023      | Well-known / privileged ports      |
| `registered` | 1024 – 49151  | Registered application ports       |
| `dynamic`    | 49152 – 65535 | Ephemeral / dynamic ports          |

## Usage

```go
c := portclassify.New()
result := c.Classify(portclassify.Port{Number: 443, Protocol: "tcp"})
fmt.Println(result.Class) // system
```

### Bulk classification

```go
results := c.ClassifyAll(ports)
```

### Reporting

```go
r := portclassify.NewReporter(c, os.Stdout)
r.Report(ports)

summary := r.Summary(ports)
fmt.Println(summary[portclassify.ClassSystem])
```

## Integration

`portclassify` is consumed by `portmeta` to enrich port metadata with range
information, and by the monitor pipeline to surface unexpected system-range
ports that appear after a baseline has been established.
