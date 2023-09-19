package publisher

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Define a mock struct to use in unit tests
type mockSNSClient struct {
	snsiface.SNSAPI
}

func (m *mockSNSClient) PublishWithContext(aws.Context, *sns.PublishInput, ...request.Option) (*sns.PublishOutput, error) {
	mID := "test"
	return &sns.PublishOutput{MessageId: &mID}, nil
}

// func (m *mockSNSClient) PublishWithContext(ctx context.Context, msg string, topicARN string) (*sns.PublishOutput, error) {

// 	if msg == "" {
// 		fmt.Println("Cannot send empty message")
// 		return nil, errors.New("Empty MSG ")
// 	}

// 	if topicARN == "" {
// 		fmt.Println("Topic ARN is empty")
// 		return nil, errors.New("Topic Arn not specified")
// 	}

// 	_ = &sns.PublishInput{
// 		Message:  &msg,
// 		TopicArn: &topicARN,
// 	}

// 	resp := sns.PublishOutput{
// 		MessageId: aws.String("test-message-ID"),
// 	}
// 	return &resp, nil
// }

// func TestPublishMessage(t *testing.T) {
// 	thisTime := time.Now().String()
// 	t.Log("Starting unit test at " + thisTime)

// 	msg := "test-message"
// 	topicARN := "test-topic-ARN"

// 	mockSvc := &mockSNSClient{}

// 	result, err := mockSvc.PublishWithContext(context.TODO(), msg, topicARN)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log("Message ID: " + *result.MessageId)
// }

func TestPublishMessageFromCloudEvent(t *testing.T) {
	thisTime := time.Now()

	event := cloudevents.NewEvent()

	event.SetID("id")
	event.SetData(cloudevents.ApplicationJSON, "data")
	event.SetSource("source")
	event.SetType("type")
	event.SetTime(thisTime)

	topicARN := "test-topic-ARN"

	mockSvc := mockSNSClient{}

	p := NewSNSPublisher(&mockSvc)
	err := p.PublishEvent(context.TODO(), &event, topicARN)

	if err != nil {
		t.Fatal(err)
	}
}

func TestPublishMessageFromCloudEventNil(t *testing.T) {
	topicARN := "test-topic-ARN"

	mockSvc := mockSNSClient{}
	p := NewSNSPublisher(&mockSvc)
	err := p.PublishEvent(context.TODO(), nil, topicARN)

	if err != nil {
		t.Fatal(err)
	}
}

func TestPublishMessageMissingTopicArn(t *testing.T) {
	thisTime := time.Now()

	event := cloudevents.NewEvent()

	event.SetID("id")
	event.SetData(cloudevents.ApplicationJSON, "data")
	event.SetSource("source")
	event.SetType("type")
	event.SetTime(thisTime)

	topicARN := ""

	mockSvc := mockSNSClient{}

	p := NewSNSPublisher(&mockSvc)
	err := p.PublishEvent(context.TODO(), &event, topicARN)

	if err == nil {
		t.Fatal(err)
	}
}
