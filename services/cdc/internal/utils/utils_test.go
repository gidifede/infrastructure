package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestConvertToCloudEvent(t *testing.T) {

	cloudeventMapping := make(map[string]string)
	cloudeventMapping["eventId"] = "testId"
	cloudeventMapping["timestampSent"] = "1683270621006"
	cloudeventMapping["data"] = "{}"
	cloudeventMappingWrongTimestamp := make(map[string]string)
	cloudeventMappingWrongTimestamp["eventId"] = "testId"
	cloudeventMappingWrongTimestamp["timestampSent"] = "wrongtimestampformat"

	tests := []struct {
		name    string
		arg     map[string]string
		wantErr bool
	}{
		{name: "1. OK.", arg: cloudeventMapping, wantErr: false},
		{name: "2. KO wrong timestamp format.", arg: cloudeventMappingWrongTimestamp, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err := ConvertToCloudEvent(tt.arg)
			fmt.Println(ret, err)
			if (err != nil) && tt.wantErr == false {
				t.Errorf("ConvertToCloudEvent() errorString = %v", err.Error())
			}
			if (err == nil) && tt.wantErr == true {
				t.Errorf("ConvertToCloudEvent() errorString = %v", err.Error())
			}
			if ret != nil {
				if ret.ID() != tt.arg["eventId"] {
					t.Errorf("ConvertToCloudEvent() errorString = %v", err.Error())
				}
			}
		})
	}

}

func getDummyStream() (*events.DynamoDBEvent, error) {
	var configFileName = "testData/dummystream.json"
	var streamEvent events.DynamoDBEvent
	content, err := os.ReadFile(configFileName)

	if err != nil {
		return nil, err
	}

	text := string(content)
	err = json.Unmarshal([]byte(text), &streamEvent)

	if err != nil {
		return nil, err
	}

	return &streamEvent, nil
}

func TestParseDynamoStream(t *testing.T) {
	streamEvent, errConfig := getDummyStream()

	if errConfig != nil {
		t.Fatal(errConfig)
	}
	for _, record := range streamEvent.Records {
		jsonm := make(map[string]string)
		msgJSON, err := ParseDynamoStream(context.TODO(), record, jsonm)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(msgJSON)
		t.Log("Test with dummy stream succesfull!")
	}

}
