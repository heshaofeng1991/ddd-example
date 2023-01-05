package user

import "github.com/heshaofeng1991/ddd-johnny/domain"

// Event interface for describing Order relevant Domain Event
type Event interface {
	domain.Event
}

type SignedUpEvent struct {
	UserID int64
}

func (s SignedUpEvent) Name() string {
	return "event.user.signed_up"
}

func (s SignedUpEvent) ID() int64 {
	return s.UserID
}

type LoggedInEvent struct {
	userID int64
}

func (l LoggedInEvent) Name() string {
	return "event.user.logged_in"
}

func (l LoggedInEvent) ID() int64 {
	return l.userID
}

type LinkStoreRequestedEvent struct {
	userID int64
}

func NewLinkStoreRequestedEvent(userID int64) *LinkStoreRequestedEvent {
	return &LinkStoreRequestedEvent{userID: userID}
}

func (l LinkStoreRequestedEvent) Name() string {
	return "event.user.link_store_requested"
}

func (l LinkStoreRequestedEvent) ID() int64 {
	return l.userID
}

type GuideInfoUpdated struct {
	UserID int64
}

func (g GuideInfoUpdated) Name() string {
	return "event.user.guide_info_updated"
}

func (g GuideInfoUpdated) ID() int64 {
	return g.UserID
}
