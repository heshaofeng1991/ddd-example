package command

import (
	"context"

	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
)

type ModifyPasswordHandler struct {
	repo domainUser.Repository
}

func NewModifyPasswordHandler(repo domainUser.Repository) ModifyPasswordHandler {
	var handler ModifyPasswordHandler

	if repo == nil {
		return handler
	}

	return ModifyPasswordHandler{repo: repo}
}

func (h ModifyPasswordHandler) Handle(ctx context.Context, req interface{}) error {
	return nil
}
