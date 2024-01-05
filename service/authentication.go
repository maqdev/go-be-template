package service

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/maqdev/go-be-template/config"
	api "github.com/maqdev/go-be-template/gen/api/authors"
)

func AuthHandler(cfg *config.AppConfig) api.SecurityHandler {
	return &authHandler{cfg: cfg}
}

type authHandler struct {
	cfg *config.AppConfig
}

func (a authHandler) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (
	context.Context, error) {
	if t.Token == "123" {
		return ctx, nil
	}
	return ctx, ErrInvalidToken
}

var ErrInvalidToken = errors.New("Invalid token")
