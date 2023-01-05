package user

import (
	applicationUser "github.com/heshaofeng1991/ddd-johnny/application/user"
	"github.com/heshaofeng1991/ddd-johnny/application/user/command"
	"github.com/heshaofeng1991/ddd-johnny/application/user/query"
	"github.com/heshaofeng1991/ddd-johnny/infra/adapters"
	infraUserReferral "github.com/heshaofeng1991/ddd-johnny/infra/userReferral"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
)

func NewApplication(entClient *ent.Client) applicationUser.Application {
	return newApplication(entClient)
}

func newApplication(entClient *ent.Client) applicationUser.Application {
	userRepository := adapters.NewUserRepository(entClient)
	storeRepository := adapters.NewStoreRepository(entClient)
	tenantRepository := adapters.NewTenantRepository(entClient)

	listingRepository := adapters.NewListingRepository(entClient)
	orderRepository := adapters.NewOrderRepository(entClient)

	userReferralRepository := infraUserReferral.NewEntRepository(entClient)

	return applicationUser.Application{
		Commands: applicationUser.Commands{
			CheckEmail: command.NewCheckEmailHandler(userRepository),
			Login:      command.NewLoginHandler(userRepository, storeRepository, orderRepository, listingRepository),
			Signup: command.NewSignupHandler(
				userRepository,
				storeRepository,
				tenantRepository,
				orderRepository,
				listingRepository,
			),
			ForgotPassword:  command.NewForgotPasswordHandler(userRepository),
			ModifyPassword:  command.NewModifyPasswordHandler(userRepository),
			UpdateGuideInfo: command.NewUpdateGuideInfoHandler(userRepository),
			GenerateToken:   command.NewGenerateTokenHandler(),
			LinkStore: command.NewLinkStoreHandler(
				userRepository,
				storeRepository,
				orderRepository,
				listingRepository,
			),
			Referral:          command.NewReferralHandler(userReferralRepository),
			SyncUserToHubspot: command.NewSyncUserToHubspotHandler(userRepository),
		},
		Queries: applicationUser.Queries{
			GetGuideInfo:   query.NewGetGuideInfoHandler(userRepository),
			GetUserProfile: query.NewGetUserProfileHandler(userRepository),
		},
	}
}
