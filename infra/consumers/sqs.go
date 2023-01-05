package consumers

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/heshaofeng1991/common/util/env"
	"github.com/heshaofeng1991/ddd-johnny/domain"
	"github.com/heshaofeng1991/ddd-johnny/domain/user"
)

type SQSService struct {
	svs         *sqs.SQS
	publisher   *domain.EventPublisher
	stopChannel chan bool
}

func InitSQS() error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-east-1"),
		Credentials: credentials.NewStaticCredentials(env.AwsAccessKeyID, env.AwsSecretAccessKey, ""),
	}))

	stopChannel := make(chan bool)

	publisher := domain.NewEventPublisher()
	sqsSvc := &SQSService{
		svs:         sqs.New(sess),
		publisher:   publisher,
		stopChannel: stopChannel,
	}

	sqsSvc.Run()
	return nil
}

// Run starts SQS message listening
func (s *SQSService) Run() {
	eventChan := make(chan domain.Event)

MessageLoop:
	for {
		s.listen(eventChan)

		select {
		case event := <-eventChan:
			s.publisher.Notify(event)
		case <-s.stopChannel:
			break MessageLoop
		}
	}

	close(eventChan)
	close(s.stopChannel)
}

// Stop stops SQS message listening
func (s *SQSService) Stop() {
	s.stopChannel <- true
}

func (s *SQSService) listen(eventChan chan domain.Event) {
	go func() {
		receiveMessageOutput, err := s.svs.ReceiveMessage(
			&sqs.ReceiveMessageInput{
				MaxNumberOfMessages: aws.Int64(10),
				QueueUrl:            aws.String("https://sqs.ap-east-1.amazonaws.com/877499521494/local-oms"),
			},
		)
		if err != nil {
			return
		}

		for _, message := range receiveMessageOutput.Messages {
			if message.Body == nil {
				continue
			}

			//       var event domain.Event
			ev := domain.SQSEvent{}
			if err = json.Unmarshal([]byte(*message.Body), &ev); err != nil {
				continue
			}

			event := user.NewLinkStoreRequestedEvent(ev.ID)
			eventChan <- event
		}
	}()
}
