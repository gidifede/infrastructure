package statemachine

import (
	"reflect"
	"testing"
	"fmt"
)

var smJSONPathProduct = "./testdata/statemachine/product/stateMachine.json"

func TestLoadStateMachineProduct(t *testing.T) {
	body, expStateMachine, err := loadStateMachine(smJSONPathProduct)
	if err != nil {
		t.Errorf("Unable to run test because state machine %v is unloadable. Load error: %v", smJSONPathProduct, err.Error())
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

func TestStateMachineProduct_nextState(t *testing.T) {
	_, stateMachine, err := loadStateMachine(smJSONPathProduct)
	if err != nil {
		t.Errorf("Unable to run test because state machine %v is unloadable. Load error: %v", smJSONPathProduct, err.Error())
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
			{name: "Transition OK", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "ReadyForDelivery", action: "StartDelivery"}, want: "LastMileTransit", wantErr: false, wantErrString: ""},
			{name: "Transition OK", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "WorkInProgress", action: "FailProcessing"}, want: "Unworkable", wantErr: false, wantErrString: ""},
			{name: "Transition OK", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "WithdrawalPending", action: "Withdraw"}, want: "Withdrawn", wantErr: false, wantErrString: ""},
			{name: "Transition KO. Invalid transition", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "Accepted", action: "Accept"}, want: "", wantErr: true, wantErrString: "Accept is not a valid action for state Accepted in the loaded state machine"},
			{name: "Transition KO. Source state not exists", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "Anyone", action: "StartProcessing"}, want: "", wantErr: true, wantErrString: "Anyone is not a valid source state for the loaded state machine"},
			{name: "Transition KO. Transition not exists", fields: fields{States: stateMachine.GetStates(), Actions: stateMachine.GetActions(), TransitionMatrix: stateMachine.GetTransitionMatrix()},
				args: args{sourceState: "LastMileTransit", action: "Anyone"}, want: "", wantErr: true, wantErrString: "Anyone is not a valid action for the loaded state machine"},
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


func TestStateMachineProduct_validate(t *testing.T) {
	_, stateMachine, err := loadStateMachine(smJSONPathProduct)

	if err != nil {
		t.Errorf("Unable to run test because state machine %v is unloadable. Load error: %v", smJSONPathProduct, err.Error())
	}else{

		isValid, errorValidation := stateMachine.Validate()

		if errorValidation!= nil{
			t.Errorf("StateMachine.validate() valid= %v errorString = %v",isValid, errorValidation.Error())
			return
		}

		fmt.Printf("Validation State Machine Product: %v\n", isValid)

	}
	
}