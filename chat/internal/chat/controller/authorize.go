package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	"github.com/yogenyslav/ldt-2024/chat/pkg/secure"
)

// Authorize validates access token.
func (ctrl *Controller) Authorize(ctx context.Context, token string) (context.Context, string, error) {
	authToken, err := secure.Decrypt(token, ctrl.cipherKey)
	if err != nil {
		return nil, "", err
	}
	userInfo, err := ctrl.kc.GetUserInfo(ctx, authToken, ctrl.realm)
	if err != nil || userInfo.PreferredUsername == nil {
		return nil, "", shared.ErrInvalidJWT
	}
	return pkg.PushToken(ctx, authToken), *userInfo.PreferredUsername, nil
}
