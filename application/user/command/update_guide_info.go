package command

import (
	"context"

	"github.com/heshaofeng1991/ddd-johnny/domain"
	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	"github.com/pkg/errors"
)

type UpdateGuideInfo struct {
	UserID          int64
	Questions       []*domainUser.Question
	SkipIntegration *bool
}

type UpdateGuideInfoHandler struct {
	repo      domainUser.Repository
	publisher domain.EventPublisher
}

func NewUpdateGuideInfoHandler(repo domainUser.Repository) UpdateGuideInfoHandler {
	return UpdateGuideInfoHandler{repo: repo}
}

func (h UpdateGuideInfoHandler) Handle(ctx context.Context, cmd UpdateGuideInfo) (*domainUser.GuideInfo, error) {
	user, err := h.repo.GetUserInfo(ctx, cmd.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if len(cmd.Questions) == 6 { //nolint:gomnd
		user.GuideInfo().SetStatus(user.GuideInfo().Status() | 2) //nolint:gomnd
	}

	for _, question := range cmd.Questions {
		if question.Answer == "" {
			user.GuideInfo().SetStatus(user.GuideInfo().Status() ^ 2) //nolint:gomnd
		}
	}

	user.GuideInfo().SetQuestions(cmd.Questions)

	// 完成了第二步，并且跳过集成店铺，则设置状态为已完成指引
	if (user.GuideInfo().Status()&2 > 0) && cmd.SkipIntegration != nil && *cmd.SkipIntegration {
		user.GuideInfo().SetFinished(true)
	}

	if user.GuideInfo().Status() == 7 {
		user.GuideInfo().SetFinished(true)
	}

	user, err = h.repo.Save(ctx, user, nil)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	h.publisher.Notify(domainUser.GuideInfoUpdated{UserID: user.UserID()})

	for _, item := range user.GuideInfo().Steps() {
		if user.GuideInfo().Status()&(1<<(item.Step-1)) > 0 {
			item.Status = "done"
		} else {
			item.Status = "pending"
		}
	}

	return user.GuideInfo(), nil
}
