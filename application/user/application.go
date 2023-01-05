package user

import (
	"github.com/heshaofeng1991/ddd-johnny/application/user/command"
	"github.com/heshaofeng1991/ddd-johnny/application/user/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CheckEmail        command.CheckEmailHandler
	Login             command.LoginHandler
	Signup            command.SignupHandler
	ForgotPassword    command.ForgotPasswordHandler
	UpdateGuideInfo   command.UpdateGuideInfoHandler
	ModifyPassword    command.ModifyPasswordHandler
	GenerateToken     command.GenerateTokenHandler
	Referral          command.ReferralHandler
	LinkStore         command.LinkStoreHandler
	SyncUserToHubspot command.SyncUserToHubspotHandler
}

type Queries struct {
	GetGuideInfo   query.GetGuideInfoHandler
	GetUserProfile query.GetUserProfileHandler
}
