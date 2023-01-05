package domain

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	AWS "github.com/heshaofeng1991/common/util/aws"
	"github.com/heshaofeng1991/common/util/env"
	"github.com/sirupsen/logrus"
)

type Event interface {
	Name() string
	ID() int64
}

type SQSEvent struct {
	Name string `json:"name"`
	ID   int64  `json:"id,string"`
}

// GeneralError actual event
type GeneralError string

func (e GeneralError) ID() int64 {
	return 0
}

func NewGeneralError(err error) Event {
	return GeneralError(err.Error())
}

func (e GeneralError) Name() string {
	return "event.general.error"
}

type EventHandler interface {
	Notify(event Event)
}

type EventPublisher struct {
	handlers map[string][]EventHandler
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{}
}

func (e *EventPublisher) Subscribe(handler EventHandler, events ...Event) {
	for _, event := range events {
		e.handlers[event.Name()] = append(e.handlers[event.Name()], handler)
	}
}

func (e *EventPublisher) Notify(event Event) {
	e.notify(event)

	sess, err := AWS.NewS3Session(env.AwsRegion)
	if err != nil {
		logrus.Errorf("NotificationJDTrackDetail NewS3Session err %v", err)
	}

	svc := sqs.New(sess)
	messageBody, err := json.Marshal(struct {
		Type   string `json:"name"`
		UserID int64  `json:"user_id"`
	}{
		Type:   event.Name(),
		UserID: event.ID(),
	})
	if err != nil {
		return
	}

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBody)),
		QueueUrl:    aws.String(os.Getenv("QUEUE_URL")),
	})
	if err != nil {
		logrus.Errorf("SendMessage err %v", err)
		return
	}

	logrus.Infof("send message body: %s", string(messageBody))
}

func (e *EventPublisher) notify(event Event) {
	for _, handler := range e.handlers[event.Name()] {
		handler.Notify(event)
	}
}
