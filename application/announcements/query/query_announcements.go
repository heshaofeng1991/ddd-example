package query

import (
	"context"

	domainAnnouncements "github.com/heshaofeng1991/ddd-johnny/domain/announcements"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Conditions struct {
	Status int
}

type GetAnnouncementsHandler struct {
	repo domainAnnouncements.Repository
}

func NewGetAnnouncementsHandler(repo domainAnnouncements.Repository) GetAnnouncementsHandler {
	var getAnnouncementsHandler GetAnnouncementsHandler

	if repo == nil {
		return getAnnouncementsHandler
	}

	return GetAnnouncementsHandler{repo}
}

func (h GetAnnouncementsHandler) Handle(ctx context.Context) (
	[]*domainAnnouncements.Announcement, error,
) {
	rsp, err := h.repo.GetAnnouncements(ctx)
	if err != nil {
		logrus.Errorf("GetAnnouncements Handle error: %v", err)

		return rsp, errors.Wrap(err, "query announcements failed")
	}

	return rsp, nil
}
