package converter

import (
	"reflect"
	"testing"

	"encoding/json"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func TestStringToCloudEvent(t *testing.T) {
	expCloudEvent := cloudevents.NewEvent()
	expCloudEvent.SetSubject("subject")
	expCloudEvent.SetID("Id")
	expCloudEvent.SetSource("Source")
	expCloudEvent.SetType("type")
	body, err := json.Marshal(expCloudEvent)
	if err != nil {
		t.Errorf("Cannot run test because wrong body was provided. Body: %v", expCloudEvent)
	}

	type args struct {
		body string
	}
	tests := []struct {
		name    string
		args    args
		want    *cloudevents.Event
		wantErr bool
	}{
		{name: "Conversion OK", args: args{body: string(body)}, want: &expCloudEvent, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToCloudEvent(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToCloudEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToCloudEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
