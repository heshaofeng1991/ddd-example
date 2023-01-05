package announcements

import (
	"context"
	"time"
)

type Announcement struct {
	ID         int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	Title      string
	Content    string
	Status     int
	CreateBy   string
	Index      int
	Expiration time.Time
}

type Repository interface {
	GetAnnouncements(ctx context.Context) ([]*Announcement, error)
}
