package adapters

import (
	"context"
	"encoding/json"

	"github.com/heshaofeng1991/common/dao"
	"github.com/heshaofeng1991/common/util"
	domainTenant "github.com/heshaofeng1991/ddd-johnny/domain/tenant"
	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/heshaofeng1991/entgo/ent/gen/tenant"
	"github.com/heshaofeng1991/entgo/ent/gen/user"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	entClient *ent.Client
}

func NewUserRepository(entClient *ent.Client) *UserRepository {
	return &UserRepository{entClient: entClient}
}

func (u UserRepository) GetByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	userModel, err := u.entClient.User.Query().
		Where(user.EmailEQ(email)).
		Where(user.TypeEQ(util.OMSUserType)).
		Where(user.DeletedAtIsNil()).
		WithTenant().
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, domainUser.ErrUserNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	tenant, err := domainTenant.UnmarshalTenantFromDatabase(userModel.Edges.Tenant.ID, userModel.Edges.Tenant.Code)
	if err != nil {
		return nil, err
	}

	questions := make([]*domainUser.Question, 0)
	if userModel.Questions != "" {
		if err := json.Unmarshal([]byte(userModel.Questions), &questions); err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	result, err := domainUser.UnmarshalUserFromDatabase(
		userModel.ID,
		userModel.Name,
		userModel.Email,
		userModel.Password,
		userModel.Phone,
		tenant,
		domainUser.WithLastLoggedTime(userModel.LastLoggedTime),
		domainUser.WithSource(&userModel.Source),
		domainUser.WithSourceTag(userModel.SourceTag),
		domainUser.WithPlatform(&userModel.Platform),
		domainUser.WithWebsite(&userModel.Website),
		domainUser.WithConcerns(&userModel.Concerns),
		domainUser.WithGuideInfo(
			domainUser.NewGuideInfo(
				userModel.GuideStatus,
				userModel.GuideFinished,
				userModel.HsObjectID,
				questions,
			),
		),
	)

	return result, errors.Wrap(err, "")
}

func (u UserRepository) Create(
	ctx context.Context,
	user *domainUser.User,
	createFn func(ctx context.Context, u *domainUser.User) (*domainUser.User, error),
) (*domainUser.User, error) {
	tenantModel, err := u.entClient.Tenant.Query().
		Where(tenant.ID(user.Tenant().TenantID())).
		Where(tenant.DeletedAtIsNil()).
		First(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	newUser := (*ent.User)(nil)

	if err := dao.WithTx(
		ctx,
		u.entClient,
		func(tx *ent.Tx) error {
			// 事务如何实现.
			if createFn != nil {
				if _, err = createFn(ctx, user); err != nil {
					return errors.Wrap(err, "")
				}
			}

			newUser, err = u.entClient.User.
				Create().
				SetName(user.Username()).
				SetEmail(user.Email()).
				SetType(util.OMSUserType).
				SetPassword(user.Password()).
				SetPhone(user.Phone()).
				SetStatus(user.Status()).
				SetSource(user.Source()).
				SetStoreCode(user.StoreCode()).
				SetSourceTag(user.SourceTag()).
				SetPlatform(user.Platform()).
				SetWebsite(user.Website()).
				SetConcerns(user.Concerns()).
				SetGuideFinished(user.GuideInfo().Finished()).
				SetGuideStatus(user.GuideInfo().Status()).
				SetTenant(tenantModel).
				Save(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to create user")
			}

			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "")
	}

	if newUser == nil {
		return nil, errors.New("failed to create user")
	}

	result, err := domainUser.UnmarshalUserFromDatabase(
		newUser.ID,
		newUser.Name,
		newUser.Email,
		newUser.Password,
		newUser.Phone,
		user.Tenant(),
		domainUser.WithLastLoggedTime(newUser.LastLoggedTime),
		domainUser.WithSource(&newUser.Source),
		domainUser.WithSourceTag(newUser.SourceTag),
		domainUser.WithWebsite(&newUser.Website),
		domainUser.WithPlatform(&newUser.Platform),
	)

	return result, errors.Wrap(err, "")
}

func (u UserRepository) Save(
	ctx context.Context,
	user *domainUser.User,
	updateFn func(ctx context.Context, u *domainUser.User) (*domainUser.User, error),
) (*domainUser.User, error) {
	questions, err := json.Marshal(user.GuideInfo().Questions())
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if err := dao.WithTx(
		ctx,
		u.entClient,
		func(transaction *ent.Tx) error {
			if updateFn != nil {
				if _, err = updateFn(ctx, user); err != nil {
					return err
				}
			}

			updateHandler := transaction.User.UpdateOneID(user.UserID())
			if user.Username() != "" {
				updateHandler.SetName(user.Username())
			}
			if user.Email() != "" {
				updateHandler.SetEmail(user.Email())
			}
			if user.Password() != "" {
				updateHandler.SetPassword(user.Password())
			}

			updateHandler.SetPhone(user.Phone())
			updateHandler.SetStatus(user.Status())
			updateHandler.SetLastLoggedTime(user.LastLoggedTime())
			if user.GuideInfo().Questions() != nil {
				updateHandler.SetQuestions(string(questions))
			}
			if user.GuideInfo().HsObjectID != "" {
				updateHandler.SetHsObjectID(user.GuideInfo().HsObjectID)
			}
			updateHandler.SetGuideStatus(user.GuideInfo().Status())
			updateHandler.SetGuideFinished(user.GuideInfo().Finished())
			updateHandler.SetStoreCode(user.StoreCode())
			_, err := updateHandler.Save(ctx)

			return errors.Wrap(err, "save user failed")
		},
	); err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	return u.GetUserInfo(ctx, user.UserID())
}

func (u UserRepository) GetUserInfo(ctx context.Context, userID int64) (*domainUser.User, error) {
	userModel, err := u.entClient.User.Query().
		Where(user.ID(userID)).
		Where(user.TypeEQ(util.OMSUserType)).
		Where(user.DeletedAtIsNil()).
		WithTenant().
		First(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	tenant, err := domainTenant.UnmarshalTenantFromDatabase(userModel.Edges.Tenant.ID, userModel.Edges.Tenant.Code)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	var questions []*domainUser.Question
	if userModel.Questions != "" {
		if err := json.Unmarshal([]byte(userModel.Questions), &questions); err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	info := domainUser.NewGuideInfo(userModel.GuideStatus, userModel.GuideFinished, userModel.HsObjectID, questions)

	result, err := domainUser.UnmarshalUserFromDatabase(
		userModel.ID,
		userModel.Name,
		userModel.Email,
		userModel.Password,
		userModel.Phone,
		tenant,
		domainUser.WithGuideInfo(info),
		domainUser.WithLastLoggedTime(userModel.LastLoggedTime),
		domainUser.WithSource(&userModel.Source),
		domainUser.WithSourceTag(userModel.SourceTag),
		domainUser.WithWebsite(&userModel.Website),
		domainUser.WithPlatform(&userModel.Platform),
		domainUser.WithConcerns(&userModel.Concerns),
	)

	return result, errors.Wrap(err, "")
}

func (u UserRepository) AddHsObjectID(ctx context.Context, userID int64, hsObjectID string) error {
	if _, err := u.entClient.User.Update().SetHsObjectID(hsObjectID).
		Where(user.IDEQ(userID)).Save(ctx); err != nil {
		logrus.Errorf("AddHsObjectID err %v", err)

		return errors.Wrap(err, "")
	}

	return nil
}
