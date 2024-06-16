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
	shared.ErrInvalidUUID: {
		Status: http.StatusBadRequest,
	},
	shared.ErrEmptyQueryHint: {
		Status: http.StatusBadRequest,
	},
	shared.ErrWsProtocolRequired: {
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
	// 404
	shared.ErrNoSessionWithID: {
		Status: http.StatusNotFound,
	},
	shared.ErrNoQueryWithID: {
		Status: http.StatusNotFound,
	},
	shared.ErrNoResponseWithID: {
		Status: http.StatusNotFound,
	},
	shared.ErrNoFavoriteWithID: {
		Status: http.StatusNotFound,
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
	shared.ErrUpdateSession: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrDeleteSession: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrBeginTx: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrCommitTx: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrCreateQuery: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrCreateResponse: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrUpdateResponse: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrUpdateQuery: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrGetQuery: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrUpdateFavorite: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrDeleteFavorite: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrCreateFavorite: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrGetFavorite: {
		Status: http.StatusInternalServerError,
	},
}
