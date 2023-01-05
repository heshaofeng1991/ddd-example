package query

import (
	"context"

	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	"github.com/pkg/errors"
)

type GetUserProfileHandler struct {
	repo domainUser.Repository
}

func NewGetUserProfileHandler(repo domainUser.Repository) GetUserProfileHandler {
	var handler GetUserProfileHandler
	if repo == nil {
		return handler
	}

	return GetUserProfileHandler{repo: repo}
}

func (q GetUserProfileHandler) Handle(ctx context.Context, userID int64) (*domainUser.User, error) {
	result, err := q.repo.GetUserInfo(ctx, userID)

	return result, errors.Wrap(err, "")
}
