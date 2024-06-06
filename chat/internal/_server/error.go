package server

import (
	"net/http"

	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	srvresp "github.com/yogenyslav/pkg/response"
)

var errStatus = map[error]srvresp.ErrorResponse{
	// 400
	shared.ErrParseBody: {
		Status: http.StatusBadRequest,
	},
	shared.ErrSessionDuplicateID: {
		Msg:    "request needs to be revoked",
		Status: http.StatusBadRequest,
	},
	// 401
	shared.ErrLoginFailed: {
		Msg:    "invalid username or password",
		Status: http.StatusUnauthorized,
	},
	shared.ErrInvalidJWT: {
		Status: http.StatusUnauthorized,
	},
	shared.ErrMissingJWT: {
		Status: http.StatusUnauthorized,
	},
	// 500
	shared.ErrCreateSession: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrCtxConvertType: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrCipherTooShort: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrEncryption: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrGetSession: {
		Status: http.StatusInternalServerError,
	},
}
