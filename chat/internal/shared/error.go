package shared

import (
	"github.com/pkg/errors"
)

// 400
var (
	// ErrLoginFailed is an error when login failed.
	ErrLoginFailed = errors.New("login failed")
	// ErrParseBody is an error when failed to parse body.
	ErrParseBody = errors.New("failed to parse body")
	// ErrInvalidUUID is an error when an invalid uuid was passed.
	ErrInvalidUUID = errors.New("invalid uuid")
	// ErrWsProtocolRequired is an error when server requirement for ws is not met.
	ErrWsProtocolRequired = errors.New("upgrade to websocket is required")
	// ErrEmptyQueryHint is an error when get hint for the query with empty prompt value.
	ErrEmptyQueryHint = errors.New("query hint can't be empty")
)

// 401
var (
	// ErrMissingJWT is an error when JWT is missing.
	ErrMissingJWT = errors.New("missing JWT")
	// ErrInvalidJWT is an error when JWT is invalid.
	ErrInvalidJWT = errors.New("invalid JWT")
)

// 404
var (
	// ErrNoSessionWithID is an error when session with requested id wasn't found.
	ErrNoSessionWithID = errors.New("no session with such id found")
	// ErrNoResponseWithID is an error when response with requested id wasn't found.
	ErrNoResponseWithID = errors.New("no response with such id found")
	// ErrNoQueryWithID is an error when query with requested id wasn't found.
	ErrNoQueryWithID = errors.New("no query with such id found")
)

// 500
var (
	// ErrCtxConvertType is an error when trying to convert a value from context to a wrong type.
	ErrCtxConvertType = errors.New("wrong type converting in context")
	// ErrCipherTooShort is returned when the cipher text is too short.
	ErrCipherTooShort = errors.New("cipher text is too short")
	// ErrEncryption is an error when failed to encrypt data.
	ErrEncryption = errors.New("encryption failed")
	// ErrBeginTx is an error when failed to being an sql transaction.
	ErrBeginTx = errors.New("failed to being an sql transaction")
	// ErrCommitTx is an error when failed to commit transaction.
	ErrCommitTx = errors.New("failed to commit transaction")

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

	// ErrCreateQuery is an error when failed to create query.
	ErrCreateQuery = errors.New("failed to create query")
	// ErrCreateResponse is an error when failed to create response.
	ErrCreateResponse = errors.New("failed to create response")
	// ErrUpdateResponse is an error when failed to update response.
	ErrUpdateResponse = errors.New("failed to update response")
	// ErrUpdateQuery is an error when failed to update query.
	ErrUpdateQuery = errors.New("failed to update query")
	// ErrGetQuery is an error when failed to get query.
	ErrGetQuery = errors.New("failed to get query")
)
