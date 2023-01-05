package announcements

import (
	"github.com/heshaofeng1991/ddd-johnny/application/announcements/query"
)

type Application struct {
	Queries Queries
}

type Queries struct {
	GetAnnouncements query.GetAnnouncementsHandler
}
