package producers

import (
	"encoding/json"

	"github.com/heshaofeng1991/ddd-johnny/domain"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// EventSQSHandler publishes internal events to external world
type EventSQSHandler struct {
	svc *sqs.SQS
}

// Notify publishes Event over SQS
func (e *EventSQSHandler) Notify(event domain.Event) {
	data := map[string]string{
		"event": event.Name(),
		"id":    event.ID(),
	}

	body, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = e.svc.SendMessage(
		&sqs.SendMessageInput{
			MessageBody: aws.String(string(body)),
			QueueUrl:    &e.svc.Endpoint,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
