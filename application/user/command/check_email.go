package command

import (
	"context"

	"github.com/go-errors/errors"
	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
)

type CheckEmailHandler struct {
	repo domainUser.Repository
}

func NewCheckEmailHandler(repo domainUser.Repository) CheckEmailHandler {
	var checkEmailHandler CheckEmailHandler

	if repo == nil {
		return checkEmailHandler
	}

	return CheckEmailHandler{repo}
}

func (h CheckEmailHandler) Handle(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("email is empty")
	}

	findUser, err := h.repo.GetByEmail(ctx, email)
	if err != nil || findUser == nil {
		return false, nil
	}

	return true, nil
}
