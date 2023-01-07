package announcements

import (
	"context"
	"time"
)

type Announcement struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Status     int       `json:"status"`
	CreateBy   string    `json:"create_by"`
	Index      int       `json:"index"`
	Expiration time.Time `json:"expiration"`
}

type Repository interface {
	GetAnnouncements(ctx context.Context) ([]*Announcement, error)
}
