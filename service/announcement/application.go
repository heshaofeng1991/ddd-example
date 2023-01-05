package announcement

import (
	applicationAnnouncements "github.com/heshaofeng1991/ddd-johnny/application/announcements"
	"github.com/heshaofeng1991/ddd-johnny/application/announcements/query"
	infraAnnouncements "github.com/heshaofeng1991/ddd-johnny/infra/announcements"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
)

func NewApplication(entClient *ent.Client) applicationAnnouncements.Application {
	return newApplication(entClient)
}

func newApplication(entClient *ent.Client) applicationAnnouncements.Application {
	repository := infraAnnouncements.NewEntRepository(entClient)

	return applicationAnnouncements.Application{
		Queries: applicationAnnouncements.Queries{
			GetAnnouncements: query.NewGetAnnouncementsHandler(repository),
		},
	}
}
