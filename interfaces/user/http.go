package interfaces

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go/request"
	"github.com/go-chi/render"
	jwtAuth "github.com/heshaofeng1991/common/util/auth"
	"github.com/heshaofeng1991/common/util/env"
	"github.com/heshaofeng1991/common/util/httpresponse"
	userApplication "github.com/heshaofeng1991/ddd-johnny/application/user"
	"github.com/heshaofeng1991/ddd-johnny/application/user/command"
	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/heshaofeng1991/entgo/ent/viewer"
	"github.com/sirupsen/logrus"
)

type HTTPServer struct {
	application userApplication.Application
}

func (h HTTPServer) LinkStore(w http.ResponseWriter, r *http.Request) {
	reqBody := LinkStoreReq{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		httpresponse.ErrRsp(w, r, http.StatusBadRequest, "invalid request params")

		return
	}

	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: -1}})

	// step 1: login
	err := h.application.Commands.LinkStore.Handle(
		ctx, command.LinkStore{
			UserID:    reqBody.UserId,
			StoreCode: reqBody.StoreCode,
		},
	)
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	render.Respond(
		w, r, LinkStoreResp{
			Code:    0,
			Data:    nil,
			Message: "link store success",
		},
	)
}

func NewHTTPServer(application userApplication.Application) HTTPServer {
	return HTTPServer{
		application: application,
	}
}

func (h HTTPServer) CheckEmail(w http.ResponseWriter, r *http.Request) {
	resp := CheckEmailResp{
		Code:    0,
		Message: "",
	}

	body := CheckEmailReq{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.ErrRsp(w, r, http.StatusBadRequest, "invalid request params")

		return
	}

	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: -1}})

	exist, err := h.application.Commands.CheckEmail.Handle(ctx, string(body.Email))
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error"+err.Error())

		return
	}

	if exist {
		resp.Code = 1
		resp.Message = "email already exist"
	}

	render.Respond(w, r, resp)
}

func (h HTTPServer) Login(w http.ResponseWriter, r *http.Request) {
	resp := TokenResp{
		Code:    0,
		Message: "",
	}

	body := LoginReq{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.ErrRsp(w, r, http.StatusBadRequest, "invalid request params")

		return
	}

	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: -1}})

	// step 1: login
	token, err := h.application.Commands.Login.Handle(
		ctx, command.LoginCommand{
			Email:     string(body.Email),
			Password:  body.Password,
			StoreCode: body.StoreCode,
		},
	)
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	resp.Data = struct {
		Token *string `json:"token,omitempty"`
	}(struct{ Token *string }{Token: &token})

	render.Respond(w, r, resp)
}

func (h HTTPServer) Signup(w http.ResponseWriter, r *http.Request) {
	resp := TokenResp{
		Code:    0,
		Message: "",
	}

	body := SignupReq{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.ErrRsp(w, r, http.StatusBadRequest, "invalid request params")

		return
	}

	req := command.SignupReq{
		BusinessPlatform: body.BusinessPlatform,
		Concerns:         body.Concerns,
		Email:            string(body.Email),
		Password:         body.Password,
		Phone:            body.Phone,
		Referrer:         body.Referrer,
		Source:           body.Source,
		SourceTag:        body.SourceTag,
		StoreCode:        body.StoreCode,
		Username:         body.Username,
		Website:          body.Website,
	}

	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: -1}})

	// step 1: signup
	token, err := h.application.Commands.Signup.Handle(ctx, req)
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, err.Error())

		return
	}

	//
	// // step 3: user referrer
	// if body.Referrer != nil {
	// 	if err := h.application.Commands.Referral.Handle(
	// 		ctx, command.ReferralCommand{
	// 			UserID:          userID,
	// 			InvitedByUserID: int64(*body.Referrer),
	// 		},
	// 	); err != nil {
	// 		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, err.Error())
	// 	}
	// }

	resp.Data = struct {
		Token *string `json:"token,omitempty"`
	}(struct{ Token *string }{Token: &token})

	render.Respond(w, r, resp)
}

func (h HTTPServer) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := ForgotPasswordResp{
		Code:    0,
		Message: "",
	}

	body := ForgotPasswordReq{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.ErrRsp(w, r, http.StatusBadRequest, "invalid request params")

		return
	}

	err := h.application.Commands.ForgotPassword.Handle(context.Background(), string(body.Email))
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error")

		return
	}

	render.Respond(w, r, resp)
}

func (h HTTPServer) GetGuideInfo(w http.ResponseWriter, r *http.Request) {
	resp := GuideInfoResp{
		Code:    0,
		Message: "",
	}

	var claims jwtAuth.OMSClaims

	keyFunc := func(token *jwt.Token) (i interface{}, e error) { return []byte(env.JwtSecret), nil }

	token, err := request.ParseFromRequest(
		r,
		request.AuthorizationHeaderExtractor,
		keyFunc,
		request.WithClaims(&claims),
	)
	if err != nil {
		logrus.Errorf("ParseFromRequest error: %v", err)

		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error"+err.Error())
	}

	userID, tenantID, err := jwtAuth.OMSAuthenticate(token.Raw)
	if err != nil {
		logrus.Errorf("Authenticate error: %v", err)
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error"+err.Error())
	}

	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: tenantID}})

	guideInfo, err := h.application.Queries.GetGuideInfo.Handle(ctx, userID)
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error")

		return
	}

	steps, err := json.Marshal(guideInfo.Steps())
	if err != nil {
		return
	}

	if err = json.Unmarshal(steps, &resp.Data.Steps); err != nil {
		return
	}

	questions, err := json.Marshal(guideInfo.Questions())
	if err != nil {
		return
	}

	if err = json.Unmarshal(questions, &resp.Data.Questions); err != nil {
		return
	}

	resp.Data.Finished = guideInfo.Finished()

	render.Respond(w, r, resp)
}

func (h HTTPServer) UpdateGuideInfo(w http.ResponseWriter, r *http.Request) {
	var claims jwtAuth.OMSClaims
	keyFunc := func(token *jwt.Token) (i interface{}, e error) { return []byte(env.JwtSecret), nil }
	token, err := request.ParseFromRequest(
		r,
		request.AuthorizationHeaderExtractor,
		keyFunc,
		request.WithClaims(&claims),
	)
	if err != nil {
		logrus.Errorf("ParseFromRequest error: %v", err)
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error"+err.Error())

		return
	}

	userID, tenantID, err := jwtAuth.OMSAuthenticate(token.Raw)
	if err != nil {
		logrus.Errorf("Authenticate error: %v", err)
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error"+err.Error())

		return
	}

	resp := GuideInfoResp{
		Code:    0,
		Message: "",
	}

	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: tenantID}})

	req := UpdateGuideInfoReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrRsp(w, r, http.StatusBadRequest, "invalid request params")
		return
	}

	domainQuestions := make([]*domainUser.Question, 0)
	questions, err := json.Marshal(req.Questions)
	if err != nil {
		return
	}

	if err = json.Unmarshal(questions, &domainQuestions); err != nil {
		return
	}

	guideInfo, err := h.application.Commands.UpdateGuideInfo.Handle(
		ctx, command.UpdateGuideInfo{
			UserID:          userID,
			Questions:       domainQuestions,
			SkipIntegration: req.SkipIntegration,
		},
	)
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error")

		return
	}

	steps, err := json.Marshal(guideInfo.Steps())
	if err != nil {
		return
	}

	if err = json.Unmarshal(steps, &resp.Data.Steps); err != nil {
		return
	}

	questions, err = json.Marshal(guideInfo.Questions())
	if err != nil {
		return
	}

	if err = json.Unmarshal(questions, &resp.Data.Questions); err != nil {
		return
	}

	resp.Data.Finished = guideInfo.Finished()

	render.Respond(w, r, resp)
}

func (h HTTPServer) ModifyPassword(w http.ResponseWriter, r *http.Request) {
	resp := ModifyPasswordResp{
		Code:    0,
		Message: "",
	}

	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: -1}})

	body := ModifyPasswordReq{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.ErrRsp(w, r, http.StatusBadRequest, "invalid request params")

		return
	}

	err := h.application.Commands.ModifyPassword.Handle(ctx, body)
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error")

		return
	}

	render.Respond(w, r, resp)
}

func (h HTTPServer) UserProfile(w http.ResponseWriter, r *http.Request) {
	resp := UserInfoResp{
		Code:    0,
		Data:    UserInfo{},
		Message: "",
	}

	var claims jwtAuth.OMSClaims

	keyFunc := func(token *jwt.Token) (i interface{}, e error) { return []byte(env.JwtSecret), nil }

	token, err := request.ParseFromRequest(
		r,
		request.AuthorizationHeaderExtractor,
		keyFunc,
		request.WithClaims(&claims),
	)
	if err != nil {
		logrus.Errorf("ParseFromRequest error: %v", err)

		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error"+err.Error())
	}

	userID, tenantID, err := jwtAuth.OMSAuthenticate(token.Raw)
	if err != nil {
		logrus.Errorf("Authenticate error: %v", err)

		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error"+err.Error())
	}

	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: tenantID}})

	user, err := h.application.Queries.GetUserProfile.Handle(ctx, userID)
	if err != nil {
		httpresponse.ErrRsp(w, r, http.StatusInternalServerError, "internal server error")

		return
	}

	resp.Data = UserInfo{
		UserId:        user.UserID(),
		Email:         user.Email(),
		Phone:         user.Phone(),
		Username:      user.Username(),
		GuideFinished: user.GuideInfo().Finished(),
	}

	render.Respond(w, r, resp)
}
