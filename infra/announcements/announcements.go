package announcements

import (
	"context"

	"entgo.io/ent/dialect/sql"
	domainAnnouncements "github.com/heshaofeng1991/ddd-johnny/domain/announcements"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/heshaofeng1991/entgo/ent/gen/announcements"
	"github.com/heshaofeng1991/entgo/ent/gen/predicate"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EntRepository struct {
	entClient *ent.Client
}

func NewEntRepository(entClient *ent.Client) EntRepository {
	return EntRepository{
		entClient: entClient,
	}
}

func (r EntRepository) GetAnnouncements(ctx context.Context) ([]*domainAnnouncements.Announcement, error) {
	var conditions []predicate.Announcements

	conditions = append(conditions, announcements.StatusEQ(1))
	conditions = append(conditions, announcements.DeletedAtIsNil())

	announcement := r.entClient.Announcements.Query().Where(
		announcements.And(
			conditions...,
		),
	)

	result, err := announcement.Order(func(s *sql.Selector) {
		s.OrderBy(sql.Desc(announcements.FieldIndex))
	}).All(ctx)
	if err != nil {
		logrus.Errorf(" GetAnnouncements error: %v", err)

		return nil, errors.Wrap(err, "")
	}

	announcementArr := make([]*domainAnnouncements.Announcement, 0)

	for _, announcement := range result {
		item := &domainAnnouncements.Announcement{}
		item.ID = announcement.ID
		item.Title = announcement.Title
		item.Content = announcement.Content
		item.CreatedAt = announcement.CreatedAt
		item.UpdatedAt = announcement.UpdatedAt
		item.Status = announcement.Status

		if announcement.Edges.Users != nil {
			item.CreateBy = announcement.Edges.Users.Name
		}

		item.Index = announcement.Index
		item.Expiration = announcement.Expiration
		announcementArr = append(announcementArr, item)
	}

	return announcementArr, nil
}
