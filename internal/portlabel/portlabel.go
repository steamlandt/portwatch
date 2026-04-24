package portlabel

import "fmt"

// wellKnown maps common port numbers to their service names.
var wellKnown = map[int]string{
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

// Labeler resolves port numbers to human-readable service names.
type Labeler struct {
	custom map[int]string
}

// New creates a Labeler. Custom overrides any built-in mappings.
func New(custom map[int]string) *Labeler {
	c := make(map[int]string, len(custom))
	for k, v := range custom {
		c[k] = v
	}
	return &Labeler{custom: c}
}

// Label returns the service name for the given port, or an empty string if unknown.
func (l *Labeler) Label(port int) string {
	if name, ok := l.custom[port]; ok {
		return name
	}
	if name, ok := wellKnown[port]; ok {
		return name
	}
	return ""
}

// LabelOrPort returns the service name if known, otherwise the port number as a string.
func (l *Labeler) LabelOrPort(port int) string {
	if name := l.Label(port); name != "" {
		return name
	}
	return fmt.Sprintf("%d", port)
}
