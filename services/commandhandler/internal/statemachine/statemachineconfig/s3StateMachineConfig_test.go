package statemachineconfig

import (
	"context"
	"reflect"
	"testing"

	smConfigMock "command-handler/internal/statemachine/statemachineconfig/mock"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

/*
	Test for state machine config retrieve using json on testData folder
*/

var smJSONPathProduct = "../testdata/statemachine/cluster/stateMachine.json"

func TestS3StateMachineConfig_Get(t *testing.T) {
	s3Mock := smConfigMock.MockS3{}
	expSmBody, err := smConfigMock.LoadStateMachine(smJSONPathProduct)
	if err != nil {
		t.Errorf("Cannot start test because load state machine failed: %v", err)
	}

	type fields struct {
		connection s3iface.S3API
		bucket     string
		key        string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name          string
		args          args
		fields        fields
		want          []byte
		wantErr       bool
		wantErrString string
	}{
		{name: "Get config OK.", args: args{ctx: context.TODO()}, fields: fields{connection: s3Mock, bucket: "aBucket", key: "keyOK"}, want: expSmBody, wantErr: false, wantErrString: ""},
		{name: "Get config KO.", args: args{ctx: context.TODO()}, fields: fields{connection: s3Mock, bucket: "aBucket", key: "keyKO"}, want: nil, wantErr: true, wantErrString: "Get state machine config failed. An Error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &S3StateMachineConfig{
				connection: tt.fields.connection,
				bucket:     tt.fields.bucket,
				key:        tt.fields.key,
			}
			got, err := s.Get(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("S3StateMachineConfig.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S3StateMachineConfig.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
