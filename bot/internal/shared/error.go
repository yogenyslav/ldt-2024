package shared

import (
	"errors"
)

// 500
var (
	// ErrCipherTooShort is returned when the cipher text is too short.
	ErrCipherTooShort = errors.New("cipher text is too short")
)