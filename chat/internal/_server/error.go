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
	// 401
	shared.ErrLoginFailed: {
		Msg:    "login failed, invalid email or password",
		Status: http.StatusUnauthorized,
	},
}
