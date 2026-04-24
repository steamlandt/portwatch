// Package portlabel maps well-known port numbers to human-readable service names.
package portlabel

import "fmt"

// defaultLabels contains a curated set of well-known port-to-service mappings.
var defaultLabels = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Labeler resolves port numbers to service name strings.
type Labeler struct {
	labels map[int]string
}

// New creates a Labeler. Custom overrides are merged on top of the built-in
// defaults; a nil map is safe and uses only defaults.
func New(custom map[int]string) *Labeler {
	merged := make(map[int]string, len(defaultLabels)+len(custom))
	for k, v := range defaultLabels {
		merged[k] = v
	}
	for k, v := range custom {
		merged[k] = v
	}
	return &Labeler{labels: merged}
}

// Label returns the service name for the given port number, or an empty string
// if the port is not recognised.
func (l *Labeler) Label(port int) string {
	return l.labels[port]
}

// LabelOrPort returns the service name if known, otherwise the port number
// formatted as a decimal string.
func (l *Labeler) LabelOrPort(port int) string {
	if name := l.labels[port]; name != "" {
		return name
	}
	return fmt.Sprintf("%d", port)
}
