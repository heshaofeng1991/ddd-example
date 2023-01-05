package userReferral

import "context"

type Repository interface {
	Create(ctx context.Context, userID int64, InvitedByUserID int64) error
}
