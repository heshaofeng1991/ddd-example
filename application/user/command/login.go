package command

import (
	"context"
	"time"

	"github.com/heshaofeng1991/common/util"
	"github.com/heshaofeng1991/common/util/auth"
	"github.com/heshaofeng1991/common/util/env"
	"github.com/heshaofeng1991/ddd-johnny/domain"
	domainListing "github.com/heshaofeng1991/ddd-johnny/domain/listing"
	domainOrder "github.com/heshaofeng1991/ddd-johnny/domain/order"
	domainStore "github.com/heshaofeng1991/ddd-johnny/domain/store"
	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	userRepo    domainUser.Repository
	storeRepo   domainStore.Repository
	orderRepo   domainOrder.Repository
	listingRepo domainListing.Repository
	publisher   domain.EventPublisher
}

type LoginCommand struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	StoreCode *string `json:"storeCode"`
}

func NewLoginHandler(
	userRepo domainUser.Repository,
	storeRepo domainStore.Repository,
	orderRepo domainOrder.Repository,
	listingRepo domainListing.Repository,
) LoginHandler {
	if userRepo == nil || storeRepo == nil || orderRepo == nil || listingRepo == nil {
		panic("userRepo, storeRepo, orderRepo, listingRepo must not be nil")
	}

	return LoginHandler{
		userRepo:    userRepo,
		storeRepo:   storeRepo,
		orderRepo:   orderRepo,
		listingRepo: listingRepo,
		publisher:   domain.EventPublisher{},
	}
}

func (h LoginHandler) Handle(ctx context.Context, cmd LoginCommand) (string, error) {
	user, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(cmd.Password)); err != nil {
		return "", domainUser.ErrInvalidEmailOrPassword
	}

	user.SetLastLoggedTime(time.Now())

	store := (*domainStore.Store)(nil)
	if cmd.StoreCode != nil && *cmd.StoreCode != "" {
		store, _ = h.storeRepo.GetByStoreCode(ctx, *cmd.StoreCode)
	}

	updateUserFn := (func(ctx context.Context, user *domainUser.User) (*domainUser.User, error))(nil)

	if store != nil && store.StoreTenant() == util.SystemTenantID {
		updateUserFn = func(ctx context.Context, user *domainUser.User) (*domainUser.User, error) {
			user.SetStoreCode(*cmd.StoreCode)
			h.publisher.Notify(domainUser.NewLinkStoreRequestedEvent(user.UserID()))

			return user, nil
		}
	}

	if _, err := h.userRepo.Save(ctx, user, updateUserFn); err != nil {
		return "", errors.Wrap(err, "")
	}

	result, err := auth.GenerateOMSToken(user.UserID(), user.Tenant().TenantID(), env.JwtSecret, jwt.SigningMethodHS256)

	return result, errors.Wrap(err, "")
}
