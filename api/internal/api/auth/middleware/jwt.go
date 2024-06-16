package middleware

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/yogenyslav/ldt-2024/api/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserOrganizationRepo interface {
	FindOrganization(ctx context.Context, username string) (string, error)
}

// JWT мидлварь для авторизации по jwt.
func JWT(kc *gocloak.GoCloak, realm string, repo UserOrganizationRepo) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token, err := auth.AuthFromMD(ctx, "Bearer")
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		userInfo, err := kc.GetUserInfo(ctx, token, realm)
		if err != nil || userInfo.PreferredUsername == nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		org, err := repo.FindOrganization(ctx, *userInfo.PreferredUsername)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		ctx = context.WithValue(ctx, shared.OrganizationKey, org)
		return context.WithValue(ctx, shared.UsernameKey, *userInfo.PreferredUsername), nil
	}
}
