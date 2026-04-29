package portcorrelate

import "fmt"

// This file exists solely to provide the fmt import used by buildKey.
// It is kept separate so the main file stays readable.

var _ = fmt.Sprintf // ensure import is used
