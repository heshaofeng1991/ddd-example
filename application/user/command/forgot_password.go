package command

import (
	"context"

	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	"github.com/pkg/errors"
)

type ForgotPasswordHandler struct {
	repo domainUser.Repository
}

func NewForgotPasswordHandler(repo domainUser.Repository) ForgotPasswordHandler {
	var handler ForgotPasswordHandler

	if repo == nil {
		return handler
	}

	return ForgotPasswordHandler{repo: repo}
}

func (h ForgotPasswordHandler) Handle(ctx context.Context, email string) error {
	if email == "" {
		return nil
	}

	findUser, err := h.repo.GetByEmail(ctx, email)
	if err != nil {
		return errors.Wrap(err, "error getting user by email")
	}

	if findUser == nil {
		return nil
	}

	return nil
}
