package converter

import (
	sm "command-handler/internal/statemachine"
	smMock "command-handler/internal/statemachine/mock"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/mock/gomock"
)

func loadCloudEventFromJSON(jsonFilePath string) (*cloudevents.Event, error) {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	event := cloudevents.NewEvent()
	err = json.Unmarshal(byteValue, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func Test_commandToEvent(t *testing.T) {
	commandOK, err := loadCloudEventFromJSON("./testdata/commandOK.json")
	if err != nil {
		t.Errorf("Unable to run test because a wrong command cloud event json was provided. Error: %v", err.Error())
	}
	expEvent, err := loadCloudEventFromJSON("./testdata/event.json")
	if err != nil {
		t.Errorf("Unable to run test because a wrong event cloud event json was provided. Error: %v", err.Error())
	}

	// commandKO, err := loadCloudEventFromJSON("./testdata/commandKO.json")
	// if err != nil {
	// 	t.Errorf("Unable to run test because a wrong command cloud event json was provided. Error: %v", err.Error())
	// }

	// commandWrongType, err := loadCloudEventFromJSON("./testdata/commandWrongType.json")
	// if err != nil {
	// 	t.Errorf("Unable to run test because a wrong command cloud event json was provided. Error: %v", err.Error())
	// }

	ctrl := gomock.NewController(t)
	smMock := smMock.NewMockIStateMachine(ctrl)
	// // Mapping OK
	// smMock.EXPECT().CommandToEvent(strings.Split(commandOK.Type(), ".")[4]).Return(strings.Split(expEvent.Type(), ".")[4], nil)
	// // Mapping KO. Command doesn't map to an event
	// smMock.EXPECT().CommandToEvent(strings.Split(commandKO.Type(), ".")[4]).Return("", fmt.Errorf("mapping failed"))
	// // Mapping KO. Command has wrong type field
	// // No call

	type args struct {
		command      *cloudevents.Event
		stateMachine sm.IStateMachine
	}
	tests := []struct {
		name          string
		args          args
		want          *cloudevents.Event
		wantErr       bool
		wantErrString string
	}{
		{name: "Mapping OK", args: args{command: commandOK, stateMachine: smMock}, want: expEvent, wantErr: false, wantErrString: ""},
		// {name: "Mapping KO. Command doesn't map to an event", args: args{command: commandKO, stateMachine: smMock}, want: nil, wantErr: true, wantErrString: "Mapping Failed."},
		// {name: "Mapping KO. Command has wrong type field", args: args{command: commandWrongType, stateMachine: smMock}, want: nil, wantErr: true, wantErrString: "Command cloud event has wrong type field"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CommandToEvent(tt.args.command, tt.args.stateMachine)
			if (err != nil) != tt.wantErr {
				t.Errorf("commandToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commandToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
