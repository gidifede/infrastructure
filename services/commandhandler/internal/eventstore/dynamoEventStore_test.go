package eventstore

import (
	dynamoMock "command-handler/internal/eventstore/mock"
	"context"
	"os"
	"reflect"
	"testing"

	overlog "github.com/Trendyol/overlog"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zlogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	overlog.New(zlogger)
	overlog.SetGlobalFields([]string{"trace id", "method", "class"})

	overlog.MDC().Set("trace id", "0123456789")
}

func TestDynamoEventStore_Append(t *testing.T) {

	dbMock := dynamoMock.MockDynamoDB{}
	tableName := "tablename"

	type fields struct {
		connection dynamodbiface.DynamoDBAPI
		tableName  string
	}
	type args struct {
		ctx   context.Context
		event Event
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		wantErrString string
	}{
		{name: "Append OK.", fields: fields{connection: &dbMock, tableName: tableName},
			args: args{ctx: context.TODO(), event: Event{AggregateID: "aggregateOK"}}, wantErr: false, wantErrString: ""},
		{name: "Append KO. An error", fields: fields{connection: &dbMock, tableName: tableName},
			args: args{ctx: context.TODO(), event: Event{AggregateID: "something"}}, wantErr: true, wantErrString: "Got error calling PutItem: An error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &DynamoEventStore{
				connection: tt.fields.connection,
				tableName:  tt.fields.tableName,
			}
			if err := es.Append(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("DynamoEventStore.Append() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestDynamoEventStore_Get(t *testing.T) {

	dbMock := dynamoMock.MockDynamoDB{}
	aggregateID := "C234-1234-1248"
	data1 := "A json"
	data2 := "Another json"
	expItem1 := Event{AggregateID: aggregateID, Timestamp: 1673885460,
		Type: "Logistic.PCL.Product.Accept.AcceptProduct", Data: data1,
		Source: "Logistic.PCL.UP.OMP", EventID: "eventId1",
		TimestampSent: 1673885460}
	expItem2 := Event{AggregateID: aggregateID, Timestamp: 1673996280,
		Type: "Logistic.PCL.Product.Withdraw.WithdrawProduct", Data: data2,
		Source: "Logistic.PCL.UP.ANOTHER", EventID: "eventId2",
		TimestampSent: 1673996280}
	expItems := []Event{expItem1, expItem2}

	aggregateIDNotFound := "idNotFound"

	tableName := "tablename"

	type fields struct {
		connection dynamodbiface.DynamoDBAPI
		tableName  string
	}
	type args struct {
		ctx         context.Context
		aggregateID string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          []Event
		wantErr       bool
		wantErrString string
	}{
		{name: "Get items OK.", fields: fields{connection: &dbMock, tableName: tableName}, args: args{ctx: context.TODO(), aggregateID: aggregateID},
			want: expItems, wantErr: false, wantErrString: ""},
		{name: "Get items KO. Aggregate id not found", fields: fields{connection: &dbMock, tableName: tableName},
			args: args{ctx: context.TODO(), aggregateID: aggregateIDNotFound}, want: []Event{}, wantErr: false,
			wantErrString: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &DynamoEventStore{
				connection: tt.fields.connection,
				tableName:  tt.fields.tableName,
			}
			got, err := es.Get(tt.args.ctx, tt.args.aggregateID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DynamoEventStore.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoEventStore.Get() = %v, want %v", got, tt.want)
			}
		})
	}

}
