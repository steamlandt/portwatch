package scanner

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// Port represents an open port detected on the system.
type Port struct {
	Protocol string
	Address  string
	Port     int
}

func (p Port) String() string {
	return fmt.Sprintf("%s:%d/%s", p.Address, p.Port, p.Protocol)
}

// Scanner scans a range of ports on a given host.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a new Scanner with the given host and timeout.
func New(host string, timeout time.Duration) *Scanner {
	return &Scanner{Host: host, Timeout: timeout}
}

// Scan checks ports in [start, end] and returns those that are open.
func (s *Scanner) Scan(start, end int) ([]Port, error) {
	if start < 1 || end > 65535 || start > end {
		return nil, fmt.Errorf("invalid port range: %d-%d", start, end)
	}

	var open []Port
	for port := start; port <= end; port++ {
		addr := net.JoinHostPort(s.Host, strconv.Itoa(port))
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err == nil {
			conn.Close()
			open = append(open, Port{
				Protocol: "tcp",
				Address:  s.Host,
				Port:     port,
			})
		}
	}
	return open, nil
}
