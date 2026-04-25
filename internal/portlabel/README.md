# portlabel

The `portlabel` package maps port numbers to human-readable service names.

## Usage

```go
labeler := portlabel.New(nil)
fmt.Println(labeler.Label(80))          // "http"
fmt.Println(labeler.LabelOrPort(9999))  // "9999"
```

## Custom Overrides

Pass a `map[uint16]string` to `New` to override or extend the default labels:

```go
custom := map[uint16]string{8080: "my-app"}
labeler := portlabel.New(custom)
fmt.Println(labeler.Label(8080)) // "my-app"
```

## Defaults

The default table includes common well-known ports (HTTP, HTTPS, SSH, DNS, etc.).
Custom entries always take precedence over defaults.
