package command

import (
	"context"

	"github.com/heshaofeng1991/common/util/auth"
	"github.com/heshaofeng1991/common/util/env"
	"github.com/pkg/errors"
)

type GenerateTokenHandler struct{}

type GenerateTokenCommand struct {
	ID       int64
	TenantID int64
}

func NewGenerateTokenHandler() GenerateTokenHandler {
	return GenerateTokenHandler{}
}

func (h GenerateTokenHandler) Handle(ctx context.Context, cmd GenerateTokenCommand) (token string, err error) {
	secret := env.JwtSecret

	token, err = auth.GenerateOMSToken(cmd.ID, cmd.TenantID, secret, jwt.SigningMethodHS256)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	return token, nil
}
