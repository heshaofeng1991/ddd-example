package userReferral

import (
	"context"

	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/heshaofeng1991/entgo/ent/gen/userreferral"
	"github.com/pkg/errors"
)

type EntRepository struct {
	entClient *ent.Client
}

func NewEntRepository(entClient *ent.Client) *EntRepository {
	return &EntRepository{entClient: entClient}
}

func (r EntRepository) Create(ctx context.Context, userID int64, invitedByUserID int64) error {
	_, err := r.entClient.UserReferral.
		Query().
		Where(userreferral.UserIDEQ(userID)).
		Where(userreferral.InvitedByUserIDEQ(invitedByUserID)).
		First(ctx)

	if ent.IsNotFound(err) {
		_, err = r.entClient.UserReferral.
			Create().
			SetUserID(userID).
			SetInvitedByUserID(invitedByUserID).
			Save(ctx)

		return errors.Wrap(err, "create user referral")
	}

	return nil
}
