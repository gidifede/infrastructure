package manager

import (
	"encoding/json"
	"fmt"
	"parcel-processing-projection/internal/models"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
)

func GetEventManager(event cloudevents.Event) (IEventManager, error) {
	// Get event type
	records := strings.Split(event.Type(), ".")
	eventType := records[len(records)-1]
	log.Debug().Msgf("event type: %s", event.Type())
	log.Debug().Msgf("extracted event type received: %s", eventType)

	// Build manager
	if eventType == models.AcceptedEvent {
		specializedEvent := models.Accepted{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newAcceptedManager(specializedEvent), nil
	}
	if eventType == models.ParcelProcessedEvent {
		specializedEvent := models.ParcelProcessed{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newParcelProcessedManager(specializedEvent), nil
	}
	if eventType == models.ParcelProcessingFailedEvent {
		specializedEvent := models.ParcelProcessingFailed{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newParcelProcessingFailedManager(specializedEvent), nil
	}
	if eventType == models.ParcelLoadedEvent {
		specializedEvent := models.ParcelLoaded{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newParcelLoadedManager(specializedEvent), nil
	}
	if eventType == models.ParcelUnloadedEvent {
		specializedEvent := models.ParcelUnloaded{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newParcelUnloadedManager(specializedEvent), nil
	}
	if eventType == models.TransportEndedEvent {
		specializedEvent := models.TransportEnded{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newTransportEndedManager(specializedEvent), nil
	}
	if eventType == models.TransportStartedEvent {
		specializedEvent := models.TransportStarted{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newTransportStartedManager(specializedEvent), nil
	}
	if eventType == models.DeliveryCompletedEvent {
		specializedEvent := models.DeliveryCompleted{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newDeliveryCompletedManager(specializedEvent), nil
	}

	return nil, fmt.Errorf("unexpected event type %s", eventType)
}
