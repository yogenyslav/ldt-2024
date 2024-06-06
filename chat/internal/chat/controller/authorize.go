package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg/secure"
)

func (ctrl *Controller) Authorize(ctx context.Context, token string) (string, error) {
	authToken, err := secure.Decrypt(token, ctrl.cipherKey)
	if err != nil {
		return "", err
	}
	userInfo, err := ctrl.kc.GetUserInfo(ctx, authToken, ctrl.realm)
	if err != nil || userInfo.PreferredUsername == nil {
		return "", shared.ErrInvalidJWT
	}

	return *userInfo.PreferredUsername, nil
}
