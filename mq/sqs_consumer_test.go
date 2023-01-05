package mq_test

import (
	"database/sql"
	"testing"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/heshaofeng1991/common/util/log"
	"github.com/heshaofeng1991/ddd-johnny/mq"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	_ "github.com/heshaofeng1991/entgo/ent/gen/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func TestSend(t *testing.T) {
	t.Parallel()

	mq.SendCrmTestMsg()
}

func TestReceive(t *testing.T) {
	t.Parallel()

	log.InitLog()

	entClient, err := Open()
	if err != nil {
		logrus.Infof("init db error: %v", err)

		return
	}

	mq.ConsumeSQSMessage(entClient)
}

func Open() (entClient *ent.Client, err error) {
	mysqlDSN := "developer:3rN#m9UNrhPhn+kcW(VLVpMG@tcp(development.c51qjr1ick8t.ap-east-1.rds.amazonaws.com)/nssdb_dev"

	dbConn, err := sql.Open("mysql", mysqlDSN+"?charset=utf8mb4&parseTime=True&loc=UTC")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if dbConn.Ping() != nil {
		return nil, errors.Wrap(err, "")
	}

	dbConn.SetMaxIdleConns(10)

	drv := entsql.OpenDB("mysql", dbConn)

	entClient = ent.NewClient(ent.Driver(drv)).Debug()

	return entClient, nil
}
