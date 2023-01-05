/*
	@Author  johnny
	@Author  heshaofeng1991@gmail.com
	@Version v1.0.0
	@File    main
	@Date    2022/4/13 15:41
	@Desc
*/

package main

import (
	"net/http"
	"strings"

	"github.com/heshaofeng1991/common/dao"
	"github.com/heshaofeng1991/common/util/env"
	"github.com/heshaofeng1991/common/util/log"
	"github.com/heshaofeng1991/common/util/sentry"
	interfaceAnnouncement "github.com/heshaofeng1991/ddd-johnny/interfaces/announcement"
	interfaceHealth "github.com/heshaofeng1991/ddd-johnny/interfaces/health-check"
	interfaceUser "github.com/heshaofeng1991/ddd-johnny/interfaces/user"
	"github.com/heshaofeng1991/ddd-johnny/mq"
	svrAnnouncement "github.com/heshaofeng1991/ddd-johnny/service/announcement"
	svrUser "github.com/heshaofeng1991/ddd-johnny/service/user"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	_ "github.com/heshaofeng1991/entgo/ent/gen/runtime"
	"github.com/sirupsen/logrus"
)

func main() {
	if env.SentryDsn != "" {
		sentry.Init()

		log.InitLog()
	}

	// Initialize the db.
	entClient, err := dao.Open()
	if err != nil {
		logrus.Infof("init db error: %v", err)

		return
	}

	go mq.ConsumeSQSMessage(entClient)
	//
	// 	go func() {
	// 		err := consumers.InitSQS()
	// 		if err != nil {
	// 			logrus.Errorf("init sqs error: %v", err)
	// 		}
	// 	}()

	switch strings.ToLower(env.ServerType) {
	case "http":
		RunHTTPServer(
			HealthHandler(),
			UserHandler(entClient),
			AnnouncementHandler(entClient),
		)
	case "grpc":
	default:
	}
}

func HealthHandler() http.Handler {
	return interfaceHealth.Handler(interfaceHealth.NewHTTPServer())
}

func UserHandler(entClient *ent.Client) http.Handler {
	app := svrUser.NewApplication(entClient)

	return interfaceUser.Handler(interfaceUser.NewHTTPServer(app))
}

// AnnouncementHandler 发布公告.
func AnnouncementHandler(entClient *ent.Client) http.Handler {
	app := svrAnnouncement.NewApplication(entClient)

	return interfaceAnnouncement.Handler(interfaceAnnouncement.NewHTTPServer(app))
}
