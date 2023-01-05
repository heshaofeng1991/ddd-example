package command

import (
	"context"

	"github.com/heshaofeng1991/common/util"
	domainListing "github.com/heshaofeng1991/ddd-johnny/domain/listing"
	domainOrder "github.com/heshaofeng1991/ddd-johnny/domain/order"
	domainStore "github.com/heshaofeng1991/ddd-johnny/domain/store"
	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	"github.com/pkg/errors"
)

type LinkStoreHandler struct {
	userRepo    domainUser.Repository
	storeRepo   domainStore.Repository
	orderRepo   domainOrder.Repository
	listingRepo domainListing.Repository
}

type LinkStore struct {
	UserID    int64
	StoreCode string
}

func NewLinkStoreHandler(
	userRepo domainUser.Repository,
	storeRepo domainStore.Repository,
	orderRepo domainOrder.Repository,
	listingRepo domainListing.Repository,
) LinkStoreHandler {
	if userRepo == nil || storeRepo == nil || orderRepo == nil || listingRepo == nil {
		panic("userRepo, storeRepo, orderRepo, listingRepo must not be nil")
	}

	return LinkStoreHandler{
		userRepo:    userRepo,
		storeRepo:   storeRepo,
		orderRepo:   orderRepo,
		listingRepo: listingRepo,
	}
}

func (h LinkStoreHandler) Handle(ctx context.Context, cmd LinkStore) error {
	user, err := h.userRepo.GetUserInfo(ctx, cmd.UserID)
	if err != nil {
		return errors.Wrap(err, "")
	}

	store, err := h.storeRepo.GetByStoreCode(ctx, cmd.StoreCode)
	if err != nil {
		return err
	}

	user.GuideInfo().SetStatus(user.GuideInfo().Status() | 4)
	if user.GuideInfo().Status() == 7 {
		user.GuideInfo().SetFinished(true)
	}

	updateUserFn := (func(ctx context.Context, user *domainUser.User) (*domainUser.User, error))(nil)
	if store.StoreTenant() == util.SystemTenantID { //nolint:nestif
		updateUserFn = func(ctx context.Context, user *domainUser.User) (*domainUser.User, error) {
			if err := h.orderRepo.UpdateTenant(ctx, store, user.Tenant()); err != nil {
				return nil, errors.Wrap(err, "failed to update order tenant")
			}

			if err := h.listingRepo.UpdateTenant(ctx, store, user.Tenant()); err != nil {
				return nil, errors.Wrap(err, "failed to update listing tenant")
			}

			if err := store.CanUpdateTenant(); err == nil {
				if err := h.storeRepo.UpdateTenantAndClearStoreCode(ctx, store, user.Tenant()); err != nil {
					return nil, errors.Wrap(err, "failed to update store tenant")
				}
			}

			return user, nil
		}
	}

	if _, err := h.userRepo.Save(ctx, user, updateUserFn); err != nil {
		return errors.Wrap(err, "")
	}

	// TODO: push "User-Link-Store" event.

	return nil
}
