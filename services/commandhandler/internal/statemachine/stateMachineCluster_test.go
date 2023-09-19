package statemachine

import (
	"fmt"
	"reflect"
	"testing"
)

var smJSONPathCluster = "./testdata/statemachine/cluster/stateMachine.json"

func TestLoadStateMachineCluster(t *testing.T) {
	body, expStateMachine, err := loadStateMachine(smJSONPathCluster)
	if err != nil {
		t.Errorf("Unable to run test because state machine %v is unloadable. Load error: %v", smJSONPathCluster, err.Error())
	} else {
		type args struct {
			jsonfileContent []byte
		}
		tests := []struct {
			name    string
			args    args
			want    IStateMachine
			wantErr bool
		}{
			{name: "Load OK.", args: args{jsonfileContent: body}, want: expStateMachine, wantErr: false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := LoadStateMachine(tt.args.jsonfileContent)
				if (err != nil) != tt.wantErr {
					t.Errorf("LoadStateMachine() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("LoadStateMachine() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

func TestStateMachineCluster_nextState(t *testing.T) {
	_, stateMachine, err := loadStateMachine(smJSONPathCluster)
	if err != nil {
		t.Errorf("Unable to run test because state machine %v is unloadable. Load error: %v", smJSONPathCluster, err.Error())
	} else {
		type fields struct {
			States           []string
			Actions          []string
			TransitionMatrix [][]string
		}
		type args struct {
			sourceState string
			action      string
		}
		tests := []struct {
			name          string
			fields        fields
			args          args
			want          string
			wantErr       bool
			wantErrString string
		}{
			{name: "Transition OK .1", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "ReadyToFilling", action: "AddProduct"}, want: "Filling", wantErr: false, wantErrString: ""},
			{name: "Transition OK .2", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "Filling", action: "Close"}, want: "FilledAndCLose", wantErr: false, wantErrString: ""},
			{name: "Transition OK .3", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "FilledAndCLose", action: "StartTransport"}, want: "NetworkTransit", wantErr: false, wantErrString: ""},
			{name: "Transition KO. Invalid transition", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "ReadyToFilling", action: "Create"}, want: "", wantErr: true, wantErrString: "Create is not a valid action for state ReadyToFilling in the loaded state machine"},
			{name: "Transition KO. Source state not exists", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "Anyone", action: "StartProductProcessing"}, want: "", wantErr: true, wantErrString: "Anyone is not a valid source state for the loaded state machine"},
			{name: "Transition KO. Transition not exists", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "NetworkTransit", action: "Anyone"}, want: "", wantErr: true, wantErrString: "Anyone is not a valid action for the loaded state machine"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := &StateMachine{
					States:           tt.fields.States,
					Actions:          tt.fields.Actions,
					TransitionMatrix: tt.fields.TransitionMatrix,
				}
				got, err := s.NextState(tt.args.sourceState, tt.args.action)
				if (err != nil) != tt.wantErr {
					t.Errorf("StateMachine.nextState() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if (err != nil) && err.Error() != tt.wantErrString {
					t.Errorf("StateMachine.nextState() errorString = %v, wantErrString = %v", err.Error(), tt.wantErrString)
				}
				if got != tt.want {
					t.Errorf("StateMachine.nextState() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}


func TestStateMachineCluster_validate(t *testing.T) {
	_, stateMachine, err := loadStateMachine(smJSONPathCluster)

	if err != nil {
		t.Errorf("Unable to run test because state machine %v is unloadable. Load error: %v", smJSONPathCluster, err.Error())
	}else{

		isValid, errorValidation := stateMachine.Validate()

		if errorValidation!= nil{
			t.Errorf("StateMachine.validate() valid= %v errorString = %v",isValid, errorValidation.Error())
			return
		}

		fmt.Printf("Validation State Machine Cluster: %v\n", isValid)

	}
	
}