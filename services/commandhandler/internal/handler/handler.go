package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"command-handler/internal"
	"command-handler/internal/converter"
	es "command-handler/internal/eventstore"
	"command-handler/internal/model"
	"command-handler/internal/utils"

	overlog "github.com/Trendyol/overlog"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-xray-sdk-go/header"
	"github.com/aws/aws-xray-sdk-go/xray"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func HandleEvents(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	var failures []events.SQSBatchItemFailure

	if !internal.ConfigLoaded {
		errLoading := fmt.Errorf("configuration not loaded correctly")
		log.Err(errLoading).Msg("state machine loading or validation error")
		failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: ""})
		return events.SQSEventResponse{BatchItemFailures: failures}, nil
	}

	for _, msg := range sqsEvent.Records {

		xRayTraceID := getXRayTraceID(msg)
		overlog.MDC().Set("trace id", xRayTraceID)

		log.Debug().Msgf("Got SQS message %q with body %q\n", msg.MessageId, msg.Body)

		//Generate serverside timestamp
		generatedTimestamp := internal.Timestamp.GenerateTimestamp()

		//Convert sqs body to cloud events
		commandReceived, errConvert := converter.StringToCloudEvent(msg.Body)
		if errConvert != nil {
			log.Err(errConvert).Msgf("cloud event conversion failed ")
			// If processing fails, add the failed message identifier to the batchItemFailures slice
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: msg.MessageId})
			continue
		}

		// Prepare event
		eventCE, err := converter.CommandToEvent(commandReceived, internal.StateMachine)
		if err != nil {
			log.Err(err).Msgf("command to event conversion failed")
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: msg.MessageId})
			continue
		}

		// id, errorRetrieveAggregateID := retrieveAggregateID(commandReceived)

		// if errorRetrieveAggregateID != nil {
		// 	log.Err(err).Msgf("cannot unmarshal message to retrieve aggregate id")
		// 	failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: msg.MessageId})
		// 	continue
		// }

		checksum := utils.ComputeChecksum(eventCE.Data(), true)

		newEvent := es.Event{AggregateID: checksum,
			Timestamp:     generatedTimestamp,
			Type:          eventCE.Type(),
			Data:          string(eventCE.Data()),
			Source:        eventCE.Source(),
			EventID:       uuid.New().String(),
			TimestampSent: eventCE.Time().UnixMilli()}

		// Check if event already exists
		items, err := internal.EventStore.Get(ctx, checksum)
		isDuplicated := false
		if err != nil {
			log.Err(err).Msg("failed check if event already exists.")
		} else {
			if len(items) != 0 {
				for _, item := range items {
					// Check idempotency
					isDuplicated = item.Type == newEvent.Type && item.AggregateID == newEvent.AggregateID
					if isDuplicated {
						log.Debug().Msgf("ignoring event. An event with type %v, body checksum %v has been already received", newEvent.Type, newEvent.AggregateID)
						// eventExists = true
						break
					}
				}
			}
		}
		if !isDuplicated {
			//Write event to event store
			err = internal.EventStore.Append(ctx, newEvent)
			if err != nil {
				log.Err(err).Msg("cannot append event")
				failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: msg.MessageId})
			}
		}

		log.Debug().Msgf("eventID: %s", newEvent.EventID)
		err = sendXRaySegment(ctx, newEvent.EventID)
		if err != nil {
			log.Err(err).Msg("Error sending Xray annotated segment")
		}
	}
	return events.SQSEventResponse{BatchItemFailures: failures}, nil
}

func getXRayTraceID(event events.SQSMessage) string {
	traceID := header.FromString(event.Attributes["AWSTraceHeader"]).TraceID
	return traceID
}

func sendXRaySegment(ctx context.Context, eventID string) error {
	ctx, subsegment := xray.BeginSubsegment(ctx, "MySubsegment")
	xray.AddAnnotation(ctx, "eventId", eventID)
	defer subsegment.Close(nil)
	log.Debug().Msgf("Subsegment created with eventID: %s", eventID)
	return nil
}

func retrieveAggregateID(commandReceived *cloudevents.Event) (string, error) {

	commandType := strings.Split(commandReceived.Type(), ".")
	aggregateType := commandType[len(commandType)-3]
	var id string

	typeLower := strings.ToLower(aggregateType)
	switch typeLower {
	case "cluster":
		// Get cluster id
		var cluster model.ClusterCommand
		errClusterUnmarshal := json.Unmarshal(commandReceived.Data(), &cluster)
		if errClusterUnmarshal != nil {
			return "", errClusterUnmarshal
		}
		if cluster.Cluster.ID == "" {
			errClusterIDEmpty := fmt.Errorf("cannot retrieve cluster id from message")
			log.Err(errClusterIDEmpty).Msgf("cluster id is empty")
			return "", errClusterIDEmpty
		}
		id = cluster.Cluster.ID
	case "product":
		// Get product id
		var product model.ProductCommand
		errProductUnmarshal := json.Unmarshal(commandReceived.Data(), &product)
		if errProductUnmarshal != nil {
			return "", errProductUnmarshal
		}
		if product.Product.ID == "" {
			errProductIDEmpty := fmt.Errorf("cannot retrieve product id from message")
			log.Err(errProductIDEmpty).Msgf("product id is empty")
			return "", errProductIDEmpty
		}
		id = product.Product.ID
	case "transport":
		// Get transport id
		var transport model.TransportCommand
		errTransportUnmarshal := json.Unmarshal(commandReceived.Data(), &transport)
		if errTransportUnmarshal != nil {
			return "", errTransportUnmarshal
		}
		if transport.Transport.ID == "" {
			errTransportIDEmpty := fmt.Errorf("cannot retrieve transport id from message")
			log.Err(errTransportIDEmpty).Msgf("transport id is empty")
			return "", errTransportIDEmpty
		}
		id = transport.Transport.ID
	}

	return id, nil
}
