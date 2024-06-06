package shared

import (
	"github.com/pkg/errors"
)

var (
	// ErrLoginFailed is an error when login failed.
	ErrLoginFailed = errors.New("login failed")

	// ErrNoSessionWithID is an error when session with requested id wasn't found.
	ErrNoSessionWithID = errors.New("no session with such id found")

	// ErrParseBody is an error when failed to parse body.
	ErrParseBody = errors.New("failed to parse body")
	// ErrCtxConvertType is an error when trying to convert a value from context to a wrong type.
	ErrCtxConvertType = errors.New("wrong type converting in context")
	// ErrCipherTooShort is returned when the cipher text is too short.
	ErrCipherTooShort = errors.New("cipher text is too short")
	// ErrEncryption is an error when failed to encrypt data.
	ErrEncryption = errors.New("encryption failed")
	// ErrInvalidUUID is an error when an invalid uuid was passed.
	ErrInvalidUUID = errors.New("invalid uuid")

	// ErrMissingJWT is an error when JWT is missing.
	ErrMissingJWT = errors.New("missing JWT")
	// ErrInvalidJWT is an error when JWT is invalid.
	ErrInvalidJWT = errors.New("invalid JWT")

	// ErrCreateSession is an error when unable to create session.
	ErrCreateSession = errors.New("failed to create session")
	// ErrSessionDuplicateID is an error when generated duplicating session uuid.
	ErrSessionDuplicateID = errors.New("got duplicated session id")
	// ErrGetSession is an error when failed to get session.
	ErrGetSession = errors.New("failed to get session")
	// ErrUpdateSession is an error when failed to update session.
	ErrUpdateSession = errors.New("failed to update session")
	// ErrDeleteSession is an error when failed to delete session.
	ErrDeleteSession = errors.New("failed to delete session")
)
