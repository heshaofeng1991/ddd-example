package mq

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	AWS "github.com/heshaofeng1991/common/util/aws"
	"github.com/heshaofeng1991/common/util/env"
	svrUser "github.com/heshaofeng1991/ddd-johnny/service/user"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func hubspotQueueURL() string {
	if value, ok := os.LookupEnv("QUEUE_URL"); ok {
		return value
	}

	return "https://sqs.ap-east-1.amazonaws.com/877499521494/dev-hubspot"
}

type SQSCRMMessage struct {
	Type   string `json:"name"`
	UserID int64  `json:"user_id"`
}

func SQSOrderShippedToString(userID int64) string {
	body := SQSCRMMessage{
		UserID: userID,
	}

	b, _ := json.Marshal(body)

	return string(b)
}

func SendCrmTestMsg() {
	sess, err := AWS.NewS3Session(env.AwsRegion)
	if err != nil {
		logrus.Errorf("NotificationJDTrackDetail NewS3Session err %v", err)
	}

	svc := sqs.New(sess)
	c := []int64{1246}

	for _, item := range c {
		_, err := svc.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String(SQSOrderShippedToString(item)),
			QueueUrl:    aws.String(hubspotQueueURL()),
		})
		if err != nil {
			return
		}
	}

	return
}

func ConsumeSQSMessage(entClient *ent.Client) {
	var err error

	sess, err := AWS.NewS3Session(env.AwsRegion)
	if err != nil {
		logrus.Errorf("ConsumeSQSMessage NewS3Session err %v", err)
		panic(any("ConsumeSQSMessage NewS3Session err"))
	}

	app := svrUser.NewApplication(entClient)
	svc := sqs.New(sess)

	var maxNumberOfMessages int64 = 10

	for {
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			MaxNumberOfMessages: aws.Int64(maxNumberOfMessages),
			QueueUrl:            aws.String(hubspotQueueURL()),
		})
		if err != nil {
			logrus.Errorf("ReceiveMessage err %v", err)

			continue
		}

		msg := SQSCRMMessage{}

		for _, item := range result.Messages {
			if item.Body == nil {
				continue
			}

			if err = json.Unmarshal([]byte(*item.Body), &msg); err != nil {
				continue
			}

			logrus.Infof("receive message body: %s, %d", msg.Type, msg.UserID)

			switch msg.Type {
			case "event.user.signed_up":
				err = app.Commands.SyncUserToHubspot.Handle(msg.UserID)
				// TODO: send welcome email
			case "event.user.guide_info_updated":
				err = app.Commands.SyncUserToHubspot.Handle(msg.UserID)
			default:
				err = errors.Errorf("unknown message type %s", msg.Type)
			}

			if err == nil {
				_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
					ReceiptHandle: item.ReceiptHandle,
					QueueUrl:      aws.String(hubspotQueueURL()),
				})
				if err != nil {
					logrus.Errorf("DeleteMessage err %v", err)
				}
			} else {
				logrus.Errorf("handle message err %v", err)
			}
		}
	}
}
