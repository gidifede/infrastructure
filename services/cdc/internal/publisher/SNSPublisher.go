package publisher

import (
	"cdc/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
)

type SNSPublisher struct {
	SNSEngine snsiface.SNSAPI
}

func NewSNSPublisher(engine snsiface.SNSAPI) Publisher {
	return SNSPublisher{SNSEngine: engine}
}

func (snsPublisher SNSPublisher) publishEvent(ctx context.Context, msg string, topicARN string, msgType string) (*sns.PublishOutput, error) {

	if msg == "" {
		err := errors.New("empty MSG ")
		log.Err(err).Msg("Cannot send empty message")
		return nil, err
	}

	if topicARN == "" {
		err := errors.New("topic Arn not specified")
		log.Err(err).Msg("Topic ARN is empty")
		return nil, err
	}

	attributes := map[string]*sns.MessageAttributeValue{
		"Type": {
			DataType:    aws.String("String"),
			StringValue: aws.String(msgType),
		},
	}

	input := &sns.PublishInput{
		Message:           &msg,
		TopicArn:          &topicARN,
		MessageAttributes: attributes,
	}

	result, err := snsPublisher.SNSEngine.PublishWithContext(ctx, input)
	return result, err
}

func (snsPublisher SNSPublisher) PublishEvent(ctx context.Context, msg *cloudevents.Event, topicARN string) error {
	if msg == nil {
		return nil
	}
	msgTypeSplitted := strings.Split(msg.Type(), ".")
	msgType := msgTypeSplitted[len(msgTypeSplitted)-1]

	utils.AddClassAndMethodToMDC(snsPublisher)

	bytes, _ := json.Marshal(msg)

	result, err := snsPublisher.publishEvent(ctx, string(bytes), topicARN, msgType)

	if err != nil {
		log.Err(err).Msg("Got an error publishing the message: ")
		return err
	}

	log.Info().Msgf("Message sent with following ID: %s ", *result.MessageId)

	return nil
}
