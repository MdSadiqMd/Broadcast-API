package utils

import (
	"context"

	"github.com/MdSadiqMd/Broadcast-API/internal/models"
)

type contextKey string

const userContextKey contextKey = "user"

func SetUserInContext(ctx context.Context, user *models.JWTClaims) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) *models.JWTClaims {
	user, ok := ctx.Value(userContextKey).(*models.JWTClaims)
	if !ok {
		return nil
	}
	return user
}
