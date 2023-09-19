package converter

import (
	sm "command-handler/internal/statemachine"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func CommandToEvent(command *cloudevents.Event, stateMachine sm.IStateMachine) (*cloudevents.Event, error) {
	event := command.Clone()

	// Replace type
	// commandTypeSplitted := strings.Split(command.Type(), ".")
	// index := len(commandTypeSplitted) - 1
	// if index < 0 || len(commandTypeSplitted) != 5 {
	// 	return nil, fmt.Errorf("command cloud event has wrong type field")
	// }
	// eventType, err := stateMachine.CommandToEvent(commandTypeSplitted[index])
	// if err != nil {
	// 	return nil, err
	// }
	// commandTypeSplitted[index] = eventType
	// event.SetType(strings.Join(commandTypeSplitted, "."))

	// Replace subject
	event.SetSubject("Event")

	return &event, nil
}
