package user

import (
	"context"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User, createFn func(ctx context.Context, u *User) (*User, error)) (*User, error)
	Save(ctx context.Context, user *User, updateFn func(ctx context.Context, u *User) (*User, error)) (*User, error)
	GetUserInfo(ctx context.Context, userID int64) (*User, error)
	AddHsObjectID(ctx context.Context, userID int64, hsObjectID string) error
	// GetGuideInfos(ctx context.Context, userID int64) (*GuideInfos, error)
	// SaveGuideInfos(ctx context.Context, userID int64, guideInfos *GuideInfos) (*GuideInfos, error)
}
