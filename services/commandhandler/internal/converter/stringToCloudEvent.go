package converter

import (
	"encoding/json"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func StringToCloudEvent(body string) (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	err := json.Unmarshal([]byte(body), &event)

	if nil != err {
		return nil, err
	}
	return &event, nil
}
