package manager

import (
	"encoding/json"
	"fmt"
	"parcel-sla-projection/internal/models"
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
