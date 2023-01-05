package query

import (
	"context"

	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	"github.com/pkg/errors"
)

type GetGuideInfoHandler struct {
	userRepo domainUser.Repository
}

func NewGetGuideInfoHandler(userRepo domainUser.Repository) GetGuideInfoHandler {
	return GetGuideInfoHandler{
		userRepo: userRepo,
	}
}

func (q GetGuideInfoHandler) Handle(ctx context.Context, userID int64) (*domainUser.GuideInfo, error) {
	user, err := q.userRepo.GetUserInfo(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	for _, item := range user.GuideInfo().Steps() {
		if user.GuideInfo().Status()&(1<<(item.Step-1)) > 0 {
			item.Status = "done"
		} else {
			item.Status = "pending"
		}
	}

	return user.GuideInfo(), nil
}
