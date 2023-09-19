package handler

import (
	"context"
	"encoding/json"

	overlog "github.com/Trendyol/overlog"
	"github.com/aws/aws-lambda-go/events"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"

	"facility-vehicle-projection/internal/manager"
	"facility-vehicle-projection/internal/utils"
)

type MessageReceived struct {
	Message string `json:"Message"`
}

func HandleEvents(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	var failures []events.SQSBatchItemFailure

	for _, msg := range sqsEvent.Records {
		xRayTraceID := utils.GetXRayTraceID(msg)
		overlog.MDC().Set("trace id", xRayTraceID)

		log.Debug().Msgf("Got SQS message %q with body %q\n", msg.MessageId, msg.Body)

		// Remove SNS envelope
		message := &MessageReceived{}
		err := json.Unmarshal([]byte(msg.Body), message)
		if err != nil {
			log.Err(err).Msgf("cannot unmarshal message to convert sns format, msg id: %s", msg.MessageId)
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: msg.MessageId})
			continue
		}
		log.Debug().Msgf("Message converted: %s\n", message)

		// Decode to cloud event
		event := cloudevents.Event{}
		err1 := json.Unmarshal([]byte(message.Message), &event)

		if err1 != nil {
			log.Err(err1).Msgf("cannot unmarshal message to convert in cloud event format, msg id: %s", msg.MessageId)
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: msg.MessageId})
			continue
		}
		log.Debug().Msgf("CloudEvent: %s\n", event)

		// Process event
		manager, err := manager.GetEventManager(event)
		if err != nil {
			log.Err(err).Msgf("error during creation of event manager")
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: msg.MessageId})
		}

		err = manager.ManageEvent(ctx)
		if err != nil {
			log.Err(err).Msgf("error while processing event")
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: msg.MessageId})
		}
		log.Debug().Msgf("event %s processed", event)
	}
	return events.SQSEventResponse{BatchItemFailures: failures}, nil
}
