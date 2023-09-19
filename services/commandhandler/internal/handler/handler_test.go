package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"command-handler/internal"
	es "command-handler/internal/eventstore"
	smMock "command-handler/internal/statemachine/mock"
	smConfigMock "command-handler/internal/statemachine/statemachineconfig/mock"
	"command-handler/internal/utils"
	tmMock "command-handler/internal/utils/mock"

	"github.com/aws/aws-lambda-go/events"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/mock/gomock"
)

var prodID string

func init() {
	prodID = "productId"
}

type eventMatcher struct {
	expectedEvent es.Event
}

func EventIgnoreIDMatcher(expectedEvent es.Event) gomock.Matcher {
	return &eventMatcher{expectedEvent}
}
func (m *eventMatcher) Matches(x interface{}) bool {
	actualEvent, ok := x.(es.Event)
	if !ok {
		fmt.Println("ERROR WHILE CASTING!!")
		return false
	}

	// Ignore the eventId field when comparing the events
	return actualEvent.AggregateID == m.expectedEvent.AggregateID &&
		actualEvent.Data == m.expectedEvent.Data &&
		actualEvent.Source == m.expectedEvent.Source &&
		actualEvent.Type == m.expectedEvent.Type &&
		actualEvent.Timestamp == m.expectedEvent.Timestamp &&
		actualEvent.TimestampSent == m.expectedEvent.TimestampSent

}

func (m *eventMatcher) String() string {
	return fmt.Sprintf("%+v", m.expectedEvent)
}
func buildCloudEvent(id string, eventType string, source string, subject string, time time.Time, data map[string]interface{}) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetID(id)
	event.SetSource(source)
	event.SetType(eventType)
	event.SetSubject(subject)
	event.SetTime(time)
	event.SetData(cloudevents.ApplicationJSON, data)
	return event
}

func Test_handle_events_invalid_body(t *testing.T) {

	ctrl := gomock.NewController(t)
	internal.EventStore = es.NewMockIEventStore(ctrl)
	internal.StateMachineConfig = smConfigMock.NewMockIStateMachineConfig(ctrl)
	internal.StateMachine = smMock.NewMockIStateMachine(ctrl)
	internal.Timestamp = tmMock.NewMockITimestampManager(ctrl)
	internal.ConfigLoaded = true

	generatedTimestamp := time.Now().UnixMilli()

	internal.Timestamp.(*tmMock.MockITimestampManager).EXPECT().GenerateTimestamp().Return(generatedTimestamp)
	sqsEventsWrongCloudEvent := events.SQSEvent{Records: []events.SQSMessage{{Body: "a body"}}}

	type args struct {
		ctx      context.Context
		sqsEvent events.SQSEvent
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		startErrString string
	}{
		{name: "Invalid body.", args: args{ctx: context.TODO(), sqsEvent: sqsEventsWrongCloudEvent}, wantErr: true, startErrString: "invalid character"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HandleEvents(tt.args.ctx, tt.args.sqsEvent)
			if len(result.BatchItemFailures) == 0 && tt.wantErr {
				t.Errorf("handleEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && !strings.HasPrefix(err.Error(), tt.startErrString) {
				t.Errorf("handleEvents() errorString = %v, wantErrString = %v", err.Error(), tt.startErrString)
			}
		})
	}
}

func Test_handle_events_cannot_append_command(t *testing.T) {

	ctrl := gomock.NewController(t)
	internal.EventStore = es.NewMockIEventStore(ctrl)
	internal.StateMachineConfig = smConfigMock.NewMockIStateMachineConfig(ctrl)
	internal.StateMachine = smMock.NewMockIStateMachine(ctrl)
	internal.Timestamp = tmMock.NewMockITimestampManager(ctrl)
	internal.ConfigLoaded = true

	generatedTimestamp := time.Now().UnixMilli()

	internal.Timestamp.(*tmMock.MockITimestampManager).EXPECT().GenerateTimestamp().Return(generatedTimestamp)
	commandCEAppendFailed := buildCloudEvent("id", "Logistic.PCL.Product.Accept.Action", "Logistic.PCL.UP.OMP", "subject", time.Now(), map[string]interface{}{"product": map[string]interface{}{"id": prodID}})
	body, err := json.Marshal(commandCEAppendFailed)
	if err != nil {
		t.Errorf("Cannot start test. Wrong test data: %v", err)
	}
	sqsEventsAppendFailed := events.SQSEvent{Records: []events.SQSMessage{{Body: string(body)}}}
	aggregateId := utils.ComputeChecksum(commandCEAppendFailed.Data(), true)
	internal.EventStore.(*es.MockIEventStore).EXPECT().Get(context.TODO(), aggregateId).Return([]es.Event{}, errors.New(""))
	eventAppendFailed := es.Event{AggregateID: aggregateId,
		Timestamp:     generatedTimestamp,
		Type:          commandCEAppendFailed.Type(),
		Data:          string(commandCEAppendFailed.Data()),
		Source:        commandCEAppendFailed.Source(),
		EventID:       commandCEAppendFailed.ID(),
		TimestampSent: commandCEAppendFailed.Time().UnixMilli()}
	appendError := "AppendError"
	internal.EventStore.(*es.MockIEventStore).EXPECT().Append(context.TODO(), EventIgnoreIDMatcher(eventAppendFailed)).Return(errors.New(appendError))

	type args struct {
		ctx      context.Context
		sqsEvent events.SQSEvent
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		startErrString string
	}{
		{name: "Cannot append command.", args: args{ctx: context.TODO(), sqsEvent: sqsEventsAppendFailed}, wantErr: true, startErrString: appendError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HandleEvents(tt.args.ctx, tt.args.sqsEvent)
			if len(result.BatchItemFailures) == 0 && tt.wantErr {
				t.Errorf("handleEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && !strings.HasPrefix(err.Error(), tt.startErrString) {
				t.Errorf("handleEvents() errorString = %v, wantErrString = %v", err.Error(), tt.startErrString)
			}
		})
	}
}

func Test_handle_events_idempotency_duplicated_command(t *testing.T) {

	ctrl := gomock.NewController(t)
	internal.EventStore = es.NewMockIEventStore(ctrl)
	internal.StateMachineConfig = smConfigMock.NewMockIStateMachineConfig(ctrl)
	internal.StateMachine = smMock.NewMockIStateMachine(ctrl)
	internal.Timestamp = tmMock.NewMockITimestampManager(ctrl)
	internal.ConfigLoaded = true

	generatedTimestamp := time.Now().UnixMilli()

	//9
	internal.Timestamp.(*tmMock.MockITimestampManager).EXPECT().GenerateTimestamp().Return(generatedTimestamp)
	commandCEDuplicatedCmd := buildCloudEvent("id", "Logistic.PCL.Product.Accept.Action", "Logistic.PCL.UP.OMP", "subject", time.Now(), map[string]interface{}{"product": map[string]interface{}{"id": prodID}})
	body, err := json.Marshal(commandCEDuplicatedCmd)
	if err != nil {
		t.Errorf("Cannot start test. Wrong test data: %v", err)
	}
	sqsEventsDuplicatedCmd := events.SQSEvent{Records: []events.SQSMessage{{Body: string(body)}}}
	aggregateId := utils.ComputeChecksum(commandCEDuplicatedCmd.Data(), true)
	retrievedEventsDuplicatedCmd := []es.Event{
		{AggregateID: aggregateId, Type: commandCEDuplicatedCmd.Type(), Timestamp: generatedTimestamp, Data: string(commandCEDuplicatedCmd.Data()),
			Source: commandCEDuplicatedCmd.Source(), EventID: "id1", TimestampSent: generatedTimestamp},
		{AggregateID: aggregateId, Type: "Logistic.PCL.Product.Accept.Anything2", Timestamp: generatedTimestamp, Data: "Something",
			Source: "Logistic.PCL.UP.OMP", EventID: "id2", TimestampSent: generatedTimestamp}}
	internal.EventStore.(*es.MockIEventStore).EXPECT().Get(context.TODO(), aggregateId).Return(retrievedEventsDuplicatedCmd, nil)

	type args struct {
		ctx      context.Context
		sqsEvent events.SQSEvent
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		startErrString string
	}{
		{name: "9. Idempotency. Duplicated command", args: args{ctx: context.TODO(), sqsEvent: sqsEventsDuplicatedCmd}, wantErr: false, startErrString: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HandleEvents(tt.args.ctx, tt.args.sqsEvent)
			if len(result.BatchItemFailures) == 0 && tt.wantErr {
				t.Errorf("handleEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err != nil) && !strings.HasPrefix(err.Error(), tt.startErrString) {
				t.Errorf("handleEvents() errorString = %v, wantErrString = %v", err.Error(), tt.startErrString)
			}
		})
	}
}

func Test_handle_events_batch(t *testing.T) {

	ctrl := gomock.NewController(t)
	internal.EventStore = es.NewMockIEventStore(ctrl)
	internal.StateMachineConfig = smConfigMock.NewMockIStateMachineConfig(ctrl)
	internal.StateMachine = smMock.NewMockIStateMachine(ctrl)
	internal.Timestamp = tmMock.NewMockITimestampManager(ctrl)
	internal.ConfigLoaded = true

	generatedTimestamp := time.Now().UnixMilli()

	//Invalid body
	internal.Timestamp.(*tmMock.MockITimestampManager).EXPECT().GenerateTimestamp().Return(generatedTimestamp)
	sqsEventsWrongCloudEvent := events.SQSEvent{Records: []events.SQSMessage{{Body: "a body"}}}
	bodyErr3, err := json.Marshal(sqsEventsWrongCloudEvent)
	if err != nil {
		t.Errorf("Cannot start test. Wrong test data: %v", err)
	}

	//Transition ok
	internal.Timestamp.(*tmMock.MockITimestampManager).EXPECT().GenerateTimestamp().Return(generatedTimestamp)
	commandCEOK := buildCloudEvent("id", "Logistic.PCL.Product.Accept.Action", "Logistic.PCL.UP.OMP", "subject", time.Now(), map[string]interface{}{"product": map[string]interface{}{"id": "prodId"}})
	bodyOk, err := json.Marshal(commandCEOK)
	if err != nil {
		t.Errorf("Cannot start test. Wrong test data: %v", err)
	}
	aggregateId := utils.ComputeChecksum(commandCEOK.Data(), true)
	retrievedEventsOK := []es.Event{
		{AggregateID: aggregateId, Type: "Logistic.PCL.Product.Accept.Anything1", Timestamp: generatedTimestamp, Data: "something", Source: "Logistic.PCL.UP.OMP", EventID: "id1", TimestampSent: time.Now().UnixMilli()},
		{AggregateID: aggregateId, Type: "Logistic.PCL.Product.Accept.Anything2", Timestamp: generatedTimestamp, Data: "Something", Source: "Logistic.PCL.UP.OMP", EventID: "id2", TimestampSent: time.Now().UnixMilli()}}
	internal.EventStore.(*es.MockIEventStore).EXPECT().Get(context.TODO(), aggregateId).Return(retrievedEventsOK, nil)
	eventOK := es.Event{AggregateID: aggregateId,
		Timestamp:     generatedTimestamp,
		Type:          commandCEOK.Type(),
		Data:          string(commandCEOK.Data()),
		Source:        commandCEOK.Source(),
		EventID:       commandCEOK.ID(),
		TimestampSent: commandCEOK.Time().UnixMilli()}
	internal.EventStore.(*es.MockIEventStore).EXPECT().Append(context.TODO(), EventIgnoreIDMatcher(eventOK)).Return(nil)

	events := events.SQSEvent{Records: []events.SQSMessage{{MessageId: "AAA", Body: string(bodyErr3)}, {MessageId: "BBB", Body: string(bodyOk)}}}
	res, _ := HandleEvents(context.TODO(), events)
	if len(res.BatchItemFailures) != 1 {
		t.Errorf("handleEvents(),  expected 1 items failed in batch, but found %v", len(res.BatchItemFailures))
	}
	for i := range res.BatchItemFailures {
		resID := res.BatchItemFailures[i].ItemIdentifier

		if resID == events.Records[i].MessageId {
			continue
		}
		t.Errorf(" handleEvents(), expected failed message ID %v, but received %v", resID, events.Records[i].MessageId)

	}
}
func Test_handle_events_batch_config_off(t *testing.T) {

	ctrl := gomock.NewController(t)
	internal.EventStore = es.NewMockIEventStore(ctrl)
	internal.StateMachineConfig = smConfigMock.NewMockIStateMachineConfig(ctrl)
	internal.StateMachine = smMock.NewMockIStateMachine(ctrl)
	internal.Timestamp = tmMock.NewMockITimestampManager(ctrl)
	internal.ConfigLoaded = false

	//Missing product id
	sqsEventsWrongCloudEvent := events.SQSEvent{Records: []events.SQSMessage{{Body: "a body"}}}
	bodyErr3, err := json.Marshal(sqsEventsWrongCloudEvent)
	if err != nil {
		t.Errorf("Cannot start test. Wrong test data: %v", err)
	}

	events := events.SQSEvent{Records: []events.SQSMessage{{MessageId: "XXX", Body: string(bodyErr3)}}}
	res, _ := HandleEvents(context.TODO(), events)
	if len(res.BatchItemFailures) != 1 {
		t.Errorf("handleEvents(),  expected 4 items failed in batch, but found %v", len(res.BatchItemFailures))
	}

	for i := range res.BatchItemFailures {
		resID := res.BatchItemFailures[i].ItemIdentifier

		if resID != "" {
			t.Errorf("handleEvents(), expected empty failed batch, but found message ID %v", resID)
			continue
		}

	}
}
