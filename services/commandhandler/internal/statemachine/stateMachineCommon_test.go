package statemachine

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestStateMachine_commandToEvent(t *testing.T) {
	fakeSM := StateMachine{
		States: []string{"A",
			"B",
			"C"},
		Actions: []string{"Action1",
			"Action2"},
		TransitionMatrix: [][]string{{"B", "None"},
			{"None", "C"},
			{"None", "None"}},
		Events: []string{"Event1",
			"Event2"},
	}

	type fields struct {
		States           []string
		Actions          []string
		TransitionMatrix [][]string
		Events           []string
	}
	type args struct {
		command string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantErr       bool
		wantErrString string
	}{
		{name: "Mapping OK.", fields: fields{Actions: fakeSM.Actions, Events: fakeSM.Events}, args: args{command: "Action1"},
			want: "Event1", wantErr: false, wantErrString: ""},
		{name: "Mapping KO. Given command doesn't exist", fields: fields{Actions: fakeSM.Actions, Events: fakeSM.Events}, args: args{command: "Anyone"},
			want: "", wantErr: true, wantErrString: "Anyone is not a valid action for the loaded state machine."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StateMachine{
				States:           tt.fields.States,
				Actions:          tt.fields.Actions,
				TransitionMatrix: tt.fields.TransitionMatrix,
				Events:           tt.fields.Events,
			}
			got, err := s.CommandToEvent(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateMachine.commandToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StateMachine.commandToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateMachine_getStates(t *testing.T) {
	states := []string{"S", "S"}
	actions := []string{"A", "A"}
	events := []string{"E", "E"}
	tm := [][]string{{"T", "T"}, {"T", "T"}}

	type fields struct {
		States           []string
		Actions          []string
		TransitionMatrix [][]string
		Events           []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "Get OK.", fields: fields{States: states, Actions: actions, Events: events, TransitionMatrix: tm}, want: states},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StateMachine{
				States:           tt.fields.States,
				Actions:          tt.fields.Actions,
				TransitionMatrix: tt.fields.TransitionMatrix,
				Events:           tt.fields.Events,
			}
			if got := s.GetStates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateMachine.getStates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateMachine_getActions(t *testing.T) {
	states := []string{"S", "S"}
	actions := []string{"A", "A"}
	events := []string{"E", "E"}
	tm := [][]string{{"T", "T"}, {"T", "T"}}

	type fields struct {
		States           []string
		Actions          []string
		TransitionMatrix [][]string
		Events           []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "Get OK.", fields: fields{States: states, Actions: actions, Events: events, TransitionMatrix: tm}, want: actions},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StateMachine{
				States:           tt.fields.States,
				Actions:          tt.fields.Actions,
				TransitionMatrix: tt.fields.TransitionMatrix,
				Events:           tt.fields.Events,
			}
			if got := s.GetActions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateMachine.GetActions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateMachine_getEvents(t *testing.T) {
	states := []string{"S", "S"}
	actions := []string{"A", "A"}
	events := []string{"E", "E"}
	tm := [][]string{{"T", "T"}, {"T", "T"}}

	type fields struct {
		States           []string
		Actions          []string
		TransitionMatrix [][]string
		Events           []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "Get OK.", fields: fields{States: states, Actions: actions, Events: events, TransitionMatrix: tm}, want: events},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StateMachine{
				States:           tt.fields.States,
				Actions:          tt.fields.Actions,
				TransitionMatrix: tt.fields.TransitionMatrix,
				Events:           tt.fields.Events,
			}
			if got := s.GetEvents(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateMachine.getEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateMachine_getTransitionMatrix(t *testing.T) {
	states := []string{"S", "S"}
	actions := []string{"A", "A"}
	events := []string{"E", "E"}
	tm := [][]string{{"T", "T"}, {"T", "T"}}

	type fields struct {
		States           []string
		Actions          []string
		TransitionMatrix [][]string
		Events           []string
	}
	tests := []struct {
		name   string
		fields fields
		want   [][]string
	}{
		{name: "Get OK.", fields: fields{States: states, Actions: actions, Events: events, TransitionMatrix: tm}, want: tm},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StateMachine{
				States:           tt.fields.States,
				Actions:          tt.fields.Actions,
				TransitionMatrix: tt.fields.TransitionMatrix,
				Events:           tt.fields.Events,
			}
			if got := s.GetTransitionMatrix(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateMachine.GetTransitionMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateMachine_validate(t *testing.T) {

	states := []string{"InitialState", "IntermediateState", "FinalState"}

	actions := []string{"MoveToIntermediateState", "MoveToFinalState"}

	events := []string{"IntermediateStateReached", "FinalStateReached"}

	transitionMatrix := [][]string{
		{"IntermediateState", "None"},
		{"None", "FinalState"},
		{"None", "None"},
	}

	invalidNumberOfRowsTransitionMatrix := [][]string{{"IntermediateState", "None"}}

	invalidNumberOfColumnsTransitionMatrix := [][]string{
		{"IntermediateState"},
		{"None"},
		{"None"},
	}

	noFinalStatesTransitionMatrix := [][]string{
		{"IntermediateState", "None"},
		{"None", "FinalState"},
		{"IntermediateState", "None"},
	}

	invalidStatetransitionMatrix := [][]string{
		{"IntermediateState", "None"},
		{"InvalidState", "FinalState"},
		{"None", "None"},
	}

	noInitialStateTransitionMatrix := [][]string{
		{"IntermediateState", "None"},
		{"InitialState", "FinalState"},
		{"None", "None"},
	}

	invalidNumberOfEventsArray := []string{"FinalStateReached"}

	eventStateRelation := []string{"InitialState", "IntermediateState"}
	invalidEventStateRelation := []string{"InitialState", "InvalidState"}
	wrongLengthEventStateRelation := []string{"InitialState"}

	tests := []struct {
		name          string
		s             *StateMachine
		want          bool
		wantErr       bool
		wantErrString string
	}{
		{"Validate OK", &StateMachine{states, actions, transitionMatrix, events, eventStateRelation}, true, false, ""},
		{"Validate KO, invalid number of rows", &StateMachine{states, actions, invalidNumberOfRowsTransitionMatrix, events, eventStateRelation}, false, true, "the number of states in the transition matrix is not equal to the number of states in the states array"},
		{"Validate KO, invalid number of columns", &StateMachine{states, actions, invalidNumberOfColumnsTransitionMatrix, events, eventStateRelation}, false, true, "the number of actions in the transition matrix is not equal to the number of actions in the actions array"},
		{"Validate KO, no final state", &StateMachine{states, actions, noFinalStatesTransitionMatrix, events, eventStateRelation}, false, true, "the transition matrix does not contain any final state"},
		{"Validate KO, invalid state", &StateMachine{states, actions, invalidStatetransitionMatrix, events, eventStateRelation}, false, true, "state InvalidState is not contained in the transition matrix and is not a valid state"},
		{"Validate KO, no initial state", &StateMachine{states, actions, noInitialStateTransitionMatrix, events, eventStateRelation}, false, true, "the transition matrix does not contain any initial state"},
		{"Validate KO, invalid number of events", &StateMachine{states, actions, transitionMatrix, invalidNumberOfEventsArray, eventStateRelation}, false, true, "the number of actions in the actions array is not equal to the number of events in the events array"},
		{"Validate KO, invalid value in eventStateRelation", &StateMachine{states, actions, transitionMatrix, events, invalidEventStateRelation}, false, true, "state InvalidState in the EventStateRelation doesn't exist in the state array"},
		{"Validate KO, invalid length of eventStateRelation", &StateMachine{states, actions, transitionMatrix, events, wrongLengthEventStateRelation}, false, true, "eventStateRelation must have same length of event array"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("StateMachine.validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && err.Error() != tt.wantErrString {
				t.Errorf("StateMachine.validate() errorString = %v, wantErrString = %v", err.Error(), tt.wantErrString)
			}
			if got != tt.want {
				t.Errorf("StateMachine.validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containInMatrix(t *testing.T) {
	type args struct {
		matrix [][]string
		str    string
	}
	matrix := [][]string{
		{"A", "B"},
		{"c", "D"},
	}
	elementContained := "A"

	elementNotContained := "E"

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Test Contain in Matrix OK", args{matrix, elementContained}, true},
		{"Test Contain in Matrix KO", args{matrix, elementNotContained}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containInMatrix(tt.args.matrix, tt.args.str); got != tt.want {
				t.Errorf("containInMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateMachine_GetEventStateRelation(t *testing.T) {
	states := []string{"S", "S"}
	actions := []string{"A", "A"}
	events := []string{"E", "E"}
	tm := [][]string{{"T", "T"}, {"T", "T"}}
	esRel := []string{"AA", "AA"}

	type fields struct {
		States             []string
		Actions            []string
		TransitionMatrix   [][]string
		Events             []string
		EventStateRelation []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "Get OK.", fields: fields{States: states, Actions: actions, Events: events, TransitionMatrix: tm, EventStateRelation: esRel}, want: esRel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StateMachine{
				States:             tt.fields.States,
				Actions:            tt.fields.Actions,
				TransitionMatrix:   tt.fields.TransitionMatrix,
				Events:             tt.fields.Events,
				EventStateRelation: tt.fields.EventStateRelation,
			}
			if got := s.GetEventStateRelation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateMachine.GetEventStateRelation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateMachine_EventsToState(t *testing.T) {
	fakeSM := StateMachine{
		States: []string{"A",
			"B",
			"C"},
		Actions: []string{"Action1",
			"Action2"},
		TransitionMatrix: [][]string{{"B", "None"},
			{"None", "C"},
			{"None", "None"}},
		Events: []string{"Event1",
			"Event2"},
		EventStateRelation: []string{"B", "C"},
	}

	type fields struct {
		States             []string
		Actions            []string
		TransitionMatrix   [][]string
		Events             []string
		EventStateRelation []string
	}
	type args struct {
		events []string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantErr       bool
		wantErrString string
	}{
		{name: "Mapping OK. First", fields: fields{Actions: fakeSM.Actions, Events: fakeSM.Events, EventStateRelation: fakeSM.EventStateRelation}, args: args{events: []string{"Event1", "Event2"}},
			want: "C", wantErr: false, wantErrString: ""},
		{name: "Mapping OK. Second", fields: fields{Actions: fakeSM.Actions, Events: fakeSM.Events, EventStateRelation: fakeSM.EventStateRelation}, args: args{events: []string{"Event1"}},
			want: "B", wantErr: false, wantErrString: ""},
		{name: "Mapping KO. Given event doesn't exist", fields: fields{Actions: fakeSM.Actions, Events: fakeSM.Events, EventStateRelation: fakeSM.EventStateRelation}, args: args{events: []string{"Event3"}},
			want: "", wantErr: true, wantErrString: "Z is not a valid event for the loaded state machine."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StateMachine{
				States:             tt.fields.States,
				Actions:            tt.fields.Actions,
				TransitionMatrix:   tt.fields.TransitionMatrix,
				Events:             tt.fields.Events,
				EventStateRelation: tt.fields.EventStateRelation,
			}
			got, err := s.EventsToState(tt.args.events)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateMachine.EventsToState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StateMachine.EventsToState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func loadStateMachine(jsonFilePath string) ([]byte, IStateMachine, error) {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, nil, err
	}
	var stateMachine StateMachine
	err = json.Unmarshal(byteValue, &stateMachine)
	if err != nil {
		return nil, nil, err
	}
	return byteValue, &stateMachine, nil
}
