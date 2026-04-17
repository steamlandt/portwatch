package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Digest computes a stable hash over a set of ports for change detection.
type Digest struct{}

// New returns a new Digest.
func New() *Digest {
	return &Digest{}
}

// Port represents a minimal port entry used for hashing.
type Port struct {
	Proto  string
	Number int
	State  string
}

// Compute returns a hex-encoded SHA-256 hash of the sorted port list.
// The input slice is marshalled to JSON to produce a deterministic byte sequence.
func (d *Digest) Compute(ports []Port) (string, error) {
	if len(ports) == 0 {
		return emptyHash(), nil
	}

	b, err := json.Marshal(ports)
	if err != nil {
		return "", fmt.Errorf("digest: marshal failed: %w", err)
	}

	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// Equal returns true when two digests match.
func (d *Digest) Equal(a, b string) bool {
	return a != "" && a == b
}

func emptyHash() string {
	sum := sha256.Sum256([]byte{})
	return hex.EncodeToString(sum[:])
}
