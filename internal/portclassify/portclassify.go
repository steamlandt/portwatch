// Package portclassify assigns a traffic classification to a port based on
// its number and protocol, distinguishing between system, registered, and
// dynamic/ephemeral ranges as defined by IANA.
package portclassify

import "fmt"

// Class represents the IANA-defined port range classification.
type Class string

const (
	ClassSystem     Class = "system"     // 0–1023
	ClassRegistered Class = "registered" // 1024–49151
	ClassDynamic    Class = "dynamic"    // 49152–65535
	ClassUnknown    Class = "unknown"
)

// Port is the minimal interface expected from a port value.
type Port struct {
	Number   int
	Protocol string
}

// Result holds the classification outcome for a single port.
type Result struct {
	Port  Port
	Class Class
	Label string
}

// Classifier assigns port range classes to ports.
type Classifier struct{}

// New returns a new Classifier.
func New() *Classifier {
	return &Classifier{}
}

// Classify returns the Class for the given port number.
func (c *Classifier) Classify(p Port) Result {
	cls := classifyNumber(p.Number)
	return Result{
		Port:  p,
		Class: cls,
		Label: fmt.Sprintf("%s/%s", cls, p.Protocol),
	}
}

// ClassifyAll classifies a slice of ports and returns a map keyed by port number.
func (c *Classifier) ClassifyAll(ports []Port) map[int]Result {
	out := make(map[int]Result, len(ports))
	for _, p := range ports {
		out[p.Number] = c.Classify(p)
	}
	return out
}

func classifyNumber(n int) Class {
	switch {
	case n >= 0 && n <= 1023:
		return ClassSystem
	case n >= 1024 && n <= 49151:
		return ClassRegistered
	case n >= 49152 && n <= 65535:
		return ClassDynamic
	default:
		return ClassUnknown
	}
}
