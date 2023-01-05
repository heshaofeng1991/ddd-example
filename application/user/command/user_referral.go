package command

import (
	"context"

	domainUserReferral "github.com/heshaofeng1991/ddd-johnny/domain/userReferral"
	"github.com/pkg/errors"
)

type ReferralHandler struct {
	userReferralRepo domainUserReferral.Repository
}

type ReferralCommand struct {
	UserID          int64
	InvitedByUserID int64
}

func NewReferralHandler(userReferralRepo domainUserReferral.Repository) ReferralHandler {
	var handler ReferralHandler

	if userReferralRepo == nil {
		return handler
	}

	return ReferralHandler{
		userReferralRepo: userReferralRepo,
	}
}

func (h ReferralHandler) Handle(ctx context.Context, cmd ReferralCommand) error {
	if err := h.userReferralRepo.Create(ctx, cmd.UserID, cmd.InvitedByUserID); err != nil {
		return errors.Wrap(err, "referralHandler.Handle")
	}

	return nil
}
