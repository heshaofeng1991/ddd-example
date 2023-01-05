package command

import (
	"context"

	"github.com/heshaofeng1991/common/util"
	"github.com/heshaofeng1991/common/util/auth"
	"github.com/heshaofeng1991/common/util/env"
	"github.com/heshaofeng1991/ddd-johnny/domain"
	domainListing "github.com/heshaofeng1991/ddd-johnny/domain/listing"
	domainOrder "github.com/heshaofeng1991/ddd-johnny/domain/order"
	domainStore "github.com/heshaofeng1991/ddd-johnny/domain/store"
	domainTenant "github.com/heshaofeng1991/ddd-johnny/domain/tenant"
	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type SignupHandler struct {
	userRepo    domainUser.Repository
	storeRepo   domainStore.Repository
	tenantRepo  domainTenant.Repository
	orderRepo   domainOrder.Repository
	listingRepo domainListing.Repository
	publisher   domain.EventPublisher
}

type SignupReq struct {
	// 用户店铺所在平台
	BusinessPlatform *string

	// 用户关心的问题
	Concerns *string

	// 用户注册邮箱
	Email string

	// 用户注册密码
	Password string

	// 用户手机号
	Phone string

	// 推荐人 ID
	Referrer *int

	// 用户注册来源
	Source *string

	// 用户注册来源标签
	SourceTag *string

	// 关联店铺编码
	StoreCode *string

	// 用户名
	Username string

	// 用户店铺链接
	Website *string
}

func NewSignupHandler(
	userRepo domainUser.Repository,
	storeRepo domainStore.Repository,
	tenantRepo domainTenant.Repository,
	orderRepo domainOrder.Repository,
	listingRepo domainListing.Repository,
) SignupHandler {
	if userRepo == nil || storeRepo == nil || tenantRepo == nil || orderRepo == nil || listingRepo == nil {
		panic("userRepo, storeRepo, tenantRepo, orderRepo, listingRepo cannot be nil")
	}

	return SignupHandler{
		userRepo:    userRepo,
		storeRepo:   storeRepo,
		tenantRepo:  tenantRepo,
		orderRepo:   orderRepo,
		listingRepo: listingRepo,
	}
}

func (h SignupHandler) Handle(ctx context.Context, cmd SignupReq) (token string, err error) {
	found, _ := h.userRepo.GetByEmail(ctx, cmd.Email)
	if found != nil {
		return "", errors.New("email already exists")
	}

	store := (*domainStore.Store)(nil)
	tenantID := int64(0)

	if cmd.StoreCode != nil && *cmd.StoreCode != "" {
		store, _ = h.storeRepo.GetByStoreCode(ctx, *cmd.StoreCode)
		if store != nil {
			tenantID = store.StoreTenant()
		}
	}

	if tenantID <= util.SystemTenantID {
		tenant, err := h.tenantRepo.Create(ctx)
		if err != nil {
			return "", errors.Wrap(err, "")
		}

		tenantID = tenant.TenantID()
	}

	tenant, err := domainTenant.NewTenant(tenantID, "")
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "bcrypt password error")
	}

	user, err := domainUser.NewUser(
		cmd.Username,
		cmd.Email,
		string(bcryptPassword),
		cmd.Phone,
		tenant,
		domainUser.WithSource(cmd.Source),
		domainUser.WithWebsite(cmd.Website),
		domainUser.WithConcerns(cmd.Concerns),
		domainUser.WithPlatform(cmd.BusinessPlatform),
	)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	user.GuideInfo().SetStatus(1)

	createUserFn := (func(ctx context.Context, user *domainUser.User) (*domainUser.User, error))(nil)

	if store != nil && store.StoreTenant() == util.SystemTenantID {
		createUserFn = func(ctx context.Context, user *domainUser.User) (*domainUser.User, error) {
			user.SetStoreCode(*cmd.StoreCode)
			h.publisher.Notify(domainUser.NewLinkStoreRequestedEvent(user.UserID()))

			return user, nil
		}
	}

	user, err = h.userRepo.Create(ctx, user, createUserFn)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	// 创建用户后, 发送用户注册事件到队列.
	h.publisher.Notify(domainUser.SignedUpEvent{
		UserID: user.UserID(),
	})
	result, err := auth.GenerateOMSToken(user.UserID(), user.Tenant().TenantID(), env.JwtSecret, jwt.SigningMethodHS256)

	return result, errors.Wrap(err, "")
}
