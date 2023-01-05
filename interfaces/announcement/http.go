package interfaces

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	internal "github.com/heshaofeng1991/common"
	"github.com/heshaofeng1991/common/util/httpresponse"
	applicationAnnouncements "github.com/heshaofeng1991/ddd-johnny/application/announcements"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/heshaofeng1991/entgo/ent/viewer"
)

type HTTPServer struct {
	application applicationAnnouncements.Application
}

func NewHTTPServer(application applicationAnnouncements.Application) HTTPServer {
	return HTTPServer{
		application: application,
	}
}

func (h HTTPServer) Get(w http.ResponseWriter, r *http.Request) {
	rsp := AnnouncementsResp{}

	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: -1}})

	announcements, err := h.application.Queries.GetAnnouncements.Handle(ctx)
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "get announcements failed"+err.Error())

		return
	}

	announcementArr := make([]Announcement, 0)

	for _, announcement := range announcements {
		item := Announcement{}
		item.AnnouncementId = int(announcement.ID)
		item.Title = announcement.Title
		item.Content = announcement.Content
		item.CreatedAt = announcement.CreatedAt
		item.UpdatedAt = announcement.UpdatedAt
		item.Status = announcement.Status
		item.Expiration = announcement.Expiration.Format("2006-01-02")
		announcementArr = append(announcementArr, item)
	}

	rsp.Data = announcementArr
	rsp.Message = internal.Success
	rsp.Code = 0

	render.Respond(w, r, rsp)

	return
}
