package server

import (
	"net/http"

	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	srvresp "github.com/yogenyslav/pkg/response"
)

var errStatus = map[error]srvresp.ErrorResponse{
	// 400
	shared.ErrParseBody: {
		Msg:    "failed to parse request body",
		Status: http.StatusBadRequest,
	},
	shared.ErrSessionDuplicateID: {
		Msg:    "duplicated session id, request needs to be revoked",
		Status: http.StatusBadRequest,
	},
	// 401
	shared.ErrLoginFailed: {
		Msg:    "login failed, invalid username or password",
		Status: http.StatusUnauthorized,
	},
	shared.ErrInvalidJWT: {
		Msg:    "invalid JWT",
		Status: http.StatusUnauthorized,
	},
	shared.ErrMissingJWT: {
		Msg:    "missing JWT",
		Status: http.StatusUnauthorized,
	},
	// 500
	shared.ErrCreateSession: {
		Msg:    "got an error when creating session",
		Status: http.StatusInternalServerError,
	},
	shared.ErrCtxConvertType: {
		Msg:    "got wrong type from context",
		Status: http.StatusInternalServerError,
	},
	shared.ErrCipherTooShort: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrEncryption: {
		Msg:    "failed to encrypt data",
		Status: http.StatusInternalServerError,
	},
}
