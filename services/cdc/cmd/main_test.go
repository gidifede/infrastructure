package main

import (
	mock_publisher "cdc/internal/publisher/mock"
	"cdc/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
)

func Test_handle_events_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	// eventStore = es.NewMockIEventStore(ctrl)
	snsService = mock_publisher.NewMockPublisher(ctrl)
	// snsService = publisher.NewSNSPublisher(mock_publisher.MockSNSClient{})
	snsTopic = "testTopicArn"
	snsEvents := events.DynamoDBEvent{}
	file, _ := ioutil.ReadFile("../internal/publisher/testData/testStream.json")

	_ = json.Unmarshal([]byte(file), &snsEvents)

	jsonm := make(map[string]string)
	// jsonm["aggregate_id"] = snsEvents.Records[0].Change.NewImage["aggregate_id"]

	for name, value := range snsEvents.Records[0].Change.NewImage {
		if value.DataType() == events.DataTypeNumber {
			jsonm[name] = value.Number()
		} else {
			jsonm[name] = value.String()
		}
	}

	ce, _ := utils.ConvertToCloudEvent(jsonm)

	snsService.(*mock_publisher.MockPublisher).EXPECT().PublishEvent(context.TODO(), ce, snsTopic).Return(nil)

	type args struct {
		ctx      context.Context
		snsEvent events.DynamoDBEvent
	}

	tests := []struct {
		name           string
		args           args
		wantErr        bool
		startErrString string
	}{
		{name: "1. OK.", args: args{ctx: context.TODO(), snsEvent: snsEvents}, wantErr: false, startErrString: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleRequest(tt.args.ctx, tt.args.snsEvent)
			if (err != nil) && tt.wantErr == false {
				t.Errorf("handle_events() errorString = %v, wantErrString = %v", err.Error(), tt.startErrString)
			}
		})
	}
}

func Test_handle_events_KO(t *testing.T) {
	ctrl := gomock.NewController(t)
	snsService = mock_publisher.NewMockPublisher(ctrl)

	snsEvents := events.DynamoDBEvent{}
	dynamoStream, _ := ioutil.ReadFile("../internal/publisher/testData/testStream.json")
	_ = json.Unmarshal([]byte(dynamoStream), &snsEvents)

	snsEventsDelete := events.DynamoDBEvent{}
	dynamoStreamDelete, _ := ioutil.ReadFile("../internal/publisher/testData/testStream_delete.json")
	_ = json.Unmarshal([]byte(dynamoStreamDelete), &snsEventsDelete)

	snsEventsWrongTs := events.DynamoDBEvent{}
	dynamoStreamWrongTs, _ := ioutil.ReadFile("../internal/publisher/testData/testStream_wrongts.json")
	_ = json.Unmarshal([]byte(dynamoStreamWrongTs), &snsEventsWrongTs)

	jsonm := make(map[string]string)

	for name, value := range snsEvents.Records[0].Change.NewImage {
		if value.DataType() == events.DataTypeNumber {
			jsonm[name] = value.Number()
		} else {
			jsonm[name] = value.String()
		}
	}

	ce, _ := utils.ConvertToCloudEvent(jsonm)

	snsService.(*mock_publisher.MockPublisher).EXPECT().PublishEvent(context.TODO(), ce, "").Return(nil)
	snsService.(*mock_publisher.MockPublisher).EXPECT().PublishEvent(context.TODO(), ce, "unknownTopic").Return(errors.New(""))

	type args struct {
		ctx      context.Context
		snsEvent events.DynamoDBEvent
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		snsTopic string
	}{
		{name: "1. KO missing topicArn.", args: args{ctx: context.TODO(), snsEvent: snsEvents}, wantErr: true, snsTopic: ""},
		{name: "2. no Insert Events.", args: args{ctx: context.TODO(), snsEvent: snsEventsDelete}, wantErr: false, snsTopic: "testTopicArn"},
		{name: "3. wrong ts in dynamodb event.", args: args{ctx: context.TODO(), snsEvent: snsEventsWrongTs}, wantErr: true, snsTopic: "testTopicArn"},
		{name: "4. wrong topic arn.", args: args{ctx: context.TODO(), snsEvent: snsEvents}, wantErr: true, snsTopic: "unknownTopic"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snsTopic = tt.snsTopic
			err := handleRequest(tt.args.ctx, tt.args.snsEvent)
			if (err != nil) && tt.wantErr == false {
				t.Errorf("handle_events() errorString = %v", err.Error())
			}
		})
	}
}
