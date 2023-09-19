package publisher

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Publisher interface {
	PublishEvent(ctx context.Context, msg *cloudevents.Event, topicARN string) error
}
