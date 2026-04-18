package portlabel

// Well-known port labels for common services.
var wellKnown = map[uint16]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Labeler maps port numbers to human-readable service names.
type Labeler struct {
	custom map[uint16]string
}

// New returns a Labeler seeded with well-known labels.
// Extra custom mappings can be supplied and will override defaults.
func New(custom map[uint16]string) *Labeler {
	m := make(map[uint16]string, len(wellKnown)+len(custom))
	for k, v := range wellKnown {
		m[k] = v
	}
	for k, v := range custom {
		m[k] = v
	}
	return &Labeler{custom: m}
}

// Label returns the service name for port, or an empty string if unknown.
func (l *Labeler) Label(port uint16) string {
	return l.custom[port]
}

// LabelOrPort returns the service name for port, or the numeric string if unknown.
func (l *Labeler) LabelOrPort(port uint16) string {
	if name, ok := l.custom[port]; ok {
		return name
	}
	return fmt.Sprintf("%d", port)
}

import "fmt"
