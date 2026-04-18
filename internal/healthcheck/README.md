# healthcheck

Aggregates named health probes for the portwatch daemon and exposes a JSON report.

## Usage

```go
hc := healthcheck.New()
hc.Set("scanner", healthcheck.StatusOK, "scanning ports 1-1024")
hc.Set("state",   healthcheck.StatusDegraded, "slow disk")

fmt.Println(hc.Overall()) // degraded
```

## Reporter

`Reporter` serialises the current state to any `io.Writer` as indented JSON:

```go
r := healthcheck.NewReporter(hc, os.Stdout)
r.Write()
```

Example output:

```json
{
  "overall": "degraded",
  "checks": [
    {"name": "scanner", "status": "ok",       "message": "scanning ports 1-1024", "at": "..."},
    {"name": "state",   "status": "degraded", "message": "slow disk",            "at": "..."}
  ],
  "at": "..."
}
```

## Status levels

| Constant         | String       |
|------------------|--------------|
| `StatusOK`       | `ok`         |
| `StatusDegraded` | `degraded`   |
| `StatusDown`     | `down`       |
