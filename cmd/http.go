/*
	@Author  johnny
	@Author  heshaofeng1991@gmail.com
	@Version v1.0.0
	@File    http
	@Date    2022/4/13 15:43
	@Desc
*/

package main

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go/request"
	"github.com/go-chi/docgen"
	"github.com/heshaofeng1991/common/middleware"
	jwtAuth "github.com/heshaofeng1991/common/util/auth"
	"github.com/heshaofeng1991/common/util/env"
	httperr "github.com/heshaofeng1991/common/util/httpresponse"
	"github.com/sirupsen/logrus"
)

func RunHTTPServer(healthHandler, userHandler, announcementHandler http.Handler) {
	RunHTTPServerOnAddr(":80", healthHandler, userHandler, announcementHandler)
}

func RunHTTPServerOnAddr(addr string, healthHandler, userHandler, announcementHandler http.Handler) {
	apiRouter := chi.NewRouter()

	middleware.SetMiddlewares(apiRouter)
	addAuthMiddleware(apiRouter)

	apiRouter.Mount("/oms/v2", userHandler)
	apiRouter.Mount("/oms/v2/health-check", healthHandler)
	apiRouter.Mount("/oms/v2/announcements", announcementHandler)

	logrus.Info("Starting HTTP server")

	docgen.PrintRoutes(apiRouter)

	err := http.ListenAndServe(addr, apiRouter)

	logrus.Infof("ListenAndServe err %v", err)
}

func addAuthMiddleware(router *chi.Mux) {
	router.Use(auth)
}

func auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rsp http.ResponseWriter, req *http.Request) {
			publicAPI := []string{"/oms/v2/health-check", "/oms/v2/login", "/oms/v2/signup"}
			for _, path := range publicAPI {
				if strings.HasPrefix(req.URL.Path, path) {
					handler.ServeHTTP(rsp, req)

					return
				}
			}

			var claims jwtAuth.OMSClaims

			jwtSecret := env.JwtSecret

			if jwtSecret == "" {
				jwtSecret = "wms"
			}

			token, err := request.ParseFromRequest(
				req,
				request.AuthorizationHeaderExtractor,
				func(token *jwt.Token) (i interface{}, e error) {
					return []byte(jwtSecret), nil
				},
				request.WithClaims(&claims),
			)
			if err != nil {
				httperr.BadRequest("parse jwt token failed", err, rsp, req)

				return
			}

			logrus.Infof("token %v", token)
			logrus.Infof("error %v", err)

			if !token.Valid {
				httperr.BadRequest("invalid jwt signature", nil, rsp, req)

				return
			}

			handler.ServeHTTP(rsp, req)
		},
	)
}
