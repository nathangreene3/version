package version

import (
	"errors"
)

var (
	ErrInvalidField   = errors.New("invalid field")
	ErrInvalidVersion = errors.New("invalid version")
)
