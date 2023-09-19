package manager

import (
	"encoding/json"
	"facility-vehicle-projection/internal/models"
	"fmt"
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
	if eventType == models.TransportStartedEvent {
		specializedEvent := models.TransportStarted{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newTransportStartedManager(specializedEvent), nil
	}
	if eventType == models.TransportEndedEvent {
		specializedEvent := models.TransportEnded{}
		err := json.Unmarshal([]byte(event.Data()), &specializedEvent)
		if err != nil {
			return nil, err
		}
		return newTransportEndedManager(specializedEvent), nil
	}
	return nil, fmt.Errorf("unexpected event type %s", eventType)
}
