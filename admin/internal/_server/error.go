package _server

import (
	"net/http"

	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
	srvresp "github.com/yogenyslav/pkg/response"
)

var errStatus = map[error]srvresp.ErrorResponse{
	// 400
	shared.ErrParseBody: {
		Status: http.StatusBadRequest,
	},
	shared.ErrParseFormValue: {
		Status: http.StatusBadRequest,
	},
	shared.ErrDuplicateTitle: {
		Status: http.StatusBadRequest,
	},
	// 401
	shared.ErrInvalidJWT: {
		Status: http.StatusUnauthorized,
	},
	shared.ErrMissingJWT: {
		Status: http.StatusUnauthorized,
	},
	shared.ErrLoginFailed: {
		Status: http.StatusUnauthorized,
	},
	// 403
	shared.ErrForbidden: {
		Status: http.StatusForbidden,
	},
	// 404
	shared.ErrNoOrganization: {
		Status: http.StatusNotFound,
	},
	// 500
	shared.ErrCtxConvertType: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrCipherTooShort: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrEncryption: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrCreateOrganization: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrUpdateOrganization: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrGetOrganization: {
		Status: http.StatusInternalServerError,
	},
}
