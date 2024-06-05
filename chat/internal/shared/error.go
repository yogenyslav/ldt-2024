package shared

import (
	"github.com/pkg/errors"
)

var (
	// ErrLoginFailed is an error when login failed.
	ErrLoginFailed = errors.New("login failed")
	// ErrParseBody is an error when failed to parse body.
	ErrParseBody = errors.New("failed to parse body")
)
