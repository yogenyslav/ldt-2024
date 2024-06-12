package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	"github.com/yogenyslav/ldt-2024/chat/pkg/secure"
)

// Authorize выполняет авторизацию пользователя.
func (ctrl *Controller) Authorize(ctx context.Context, token string) (context.Context, string, error) {
	authToken, err := secure.Decrypt(token, ctrl.cipherKey)
	if err != nil {
		return nil, "", err
	}
	userInfo, err := ctrl.kc.GetUserInfo(ctx, authToken, ctrl.realm)
	if err != nil || userInfo.PreferredUsername == nil {
		log.Error().Err(err).Msg("failed to get user info")
		return nil, "", shared.ErrInvalidJWT
	}
	return pkg.PushToken(ctx, authToken), *userInfo.PreferredUsername, nil
}
