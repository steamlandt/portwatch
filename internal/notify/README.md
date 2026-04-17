# notify

The `notify` package provides pluggable notification delivery for portwatch.

## Methods

| Method    | Description                                      |
|-----------|--------------------------------------------------|
| `stdout`  | Print alert to standard output (default)         |
| `webhook` | HTTP POST JSON payload to a URL via `curl`       |
| `exec`    | Run an external command with subject and body    |

## Configuration

```json
{
  "notify": {
    "method": "webhook",
    "target": "https://hooks.example.com/portwatch"
  }
}
```

## Usage

```go
n := notify.New(notify.Config{
    Method: notify.MethodWebhook,
    Target: "https://hooks.example.com/portwatch",
})
n.Send("port opened", "tcp/9200")
```
