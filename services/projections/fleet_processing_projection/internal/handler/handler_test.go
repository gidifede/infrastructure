package handler

import (
	"context"
	"encoding/json"
	"fleet-processing-projection/internal"
	"fleet-processing-projection/internal/models"
	"fleet-processing-projection/internal/repository"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	parcelID = "FA1231ADERF1234"
	parcel2ID = "FA4321ADERF2346"
)

func init() {
	clientOpts := options.Client().ApplyURI("mongodb://root:rootpassword@localhost:27017") //.ApplyURI("mongodb://test:****@dds-bp*******1.mongodb.rds.aliyuncs.com:3717,dds-bp*******2.mongodb.rds.aliyuncs.com:3717/admin?replicaSet=mgset-XXXXX&ssl=true")
	var err error
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal().Msgf("connection failed. %s", err.Error())
		return
	}
	log.Debug().Msg("connection success")
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal().Msgf("ping failed. %s", err.Error())
		return
	}
	log.Debug().Msg("ping success")
	internal.Repo = repository.NewMongo(*client.Database("logistic"))

	// Create and initialize collections required for test
	collectionName := "facility"
	opts := options.CreateCollection()
	opts.SetValidator(bson.M{})
	err = client.Database("logistic").CreateCollection(context.Background(), collectionName, opts)
	if err != nil {
		log.Err(err)
	}
}

func TestHandleEvents(t *testing.T) {
	type args struct {
		ctx      context.Context
		sqsEvent events.SQSEvent
	}
	type testType struct {
		name    string
		args    args
		want    events.SQSEventResponse
		wantErr bool
	}
	tests := []testType{}

	// Add test cases
	logisticEvent3 := models.ParcelLoaded{
		TransportID: "4ADSF-0900-020923",
		VehicleLicensePlate: "RD312FT",
		ParcelID:   parcelID,
		FacilityID: "FI1",
		Timestamp:  time.Now()}
	b3, _ := json.Marshal(&logisticEvent3)
	var m3 map[string]interface{}
	_ = json.Unmarshal(b3, &m3)
	cloudEvent3 := buildCloudEvent("id", "Logistic.PCL.Fleet.ParcelLoaded", "source", "subject", time.Now(), m3)
	body3, _ := json.Marshal(cloudEvent3)
	snsMsg3 := MessageReceived{Message: string(body3)}
	bodySNS3, _ := json.Marshal(&snsMsg3)
	test3 := testType{name: "Caricamento Pacco Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS3)}}}}}
	tests = append(tests, test3)

	// Add test cases
	logisticEventLoaded3 := models.ParcelLoaded{
		TransportID: "4ADSF-0900-020923",
		VehicleLicensePlate: "RD312FT",
		ParcelID:   parcel2ID,
		FacilityID: "FI1",
		Timestamp:  time.Now()}
	bloaded3, _ := json.Marshal(&logisticEventLoaded3)
	var mloaded3 map[string]interface{}
	_ = json.Unmarshal(bloaded3, &mloaded3)
	cloudEventloaded3 := buildCloudEvent("id", "Logistic.PCL.Fleet.ParcelLoaded", "source", "subject", time.Now(), mloaded3)
	bodyloaded3, _ := json.Marshal(cloudEventloaded3)
	snsMsgloaded3 := MessageReceived{Message: string(bodyloaded3)}
	bodySNSloaded3, _ := json.Marshal(&snsMsgloaded3)
	testloaded3 := testType{name: "Caricamento Pacco Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNSloaded3)}}}}}
	tests = append(tests, testloaded3)

	logisticEvent4 := models.TransportStarted{
		TransportID: "4ADSF-0900-020923",
		VehicleLicensePlate:   "RD312FT",
		SourceFacilityID:      "RM1",
		DestinationFacilityID: "FI1",
		Timestamp:             time.Now()}
	b4, _ := json.Marshal(&logisticEvent4)
	var m4 map[string]interface{}
	_ = json.Unmarshal(b4, &m4)
	cloudEvent4 := buildCloudEvent("id", "Logistic.PCL.Fleet.TransportStarted", "source", "subject", time.Now(), m4)
	body4, _ := json.Marshal(cloudEvent4)
	snsMsg4 := MessageReceived{Message: string(body4)}
	bodySNS4, _ := json.Marshal(&snsMsg4)
	test4 := testType{name: "Inizio trasporto Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS4)}}}}}
	tests = append(tests, test4)


	logisticEventPosition1 := models.Position{
		VehicleLicensePlate: "RD312FT",
		Latitude:            40.7128,
		Longitude:           -74.0060,
		Timestamp:           time.Now(),
	}
	bPosition1, _ := json.Marshal(&logisticEventPosition1)
	var mPosition1 map[string]interface{}
	_ = json.Unmarshal(bPosition1, &mPosition1)
	cloudEventPosition1 := buildCloudEvent("id", "Logistic.PCL.Fleet.Position", "source", "subject", time.Now(), mPosition1)
	bodyPosition1, _ := json.Marshal(cloudEventPosition1)
	snsMsgPosition1 := MessageReceived{Message: string(bodyPosition1)}
	bodySNSPosition1, _ := json.Marshal(&snsMsgPosition1)
	testPosition1 := testType{name: "Position1 Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNSPosition1)}}}}}
	tests = append(tests, testPosition1)

	logisticEventPosition2 := models.Position{
		VehicleLicensePlate: "RD312FT",
		Latitude:            40.7138,
		Longitude:           -74.0160,
		Timestamp:           time.Now(),
	}
	bPosition2, _ := json.Marshal(&logisticEventPosition2)
	var mPosition2 map[string]interface{}
	_ = json.Unmarshal(bPosition2, &mPosition2)
	cloudEventPosition2 := buildCloudEvent("id", "Logistic.PCL.Fleet.Position", "source", "subject", time.Now(), mPosition2)
	bodyPosition2, _ := json.Marshal(cloudEventPosition2)
	snsMsgPosition2 := MessageReceived{Message: string(bodyPosition2)}
	bodySNSPosition2, _ := json.Marshal(&snsMsgPosition2)
	testPosition2 := testType{name: "Position2 Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNSPosition2)}}}}}
	tests = append(tests, testPosition2)

	logisticEvent5 := models.TransportEnded{
		TransportID: "4ADSF-0900-020923",
		VehicleLicensePlate: "RD312FT",
		FacilityID:          "FI1",
		Timestamp:           time.Now()}
	b5, _ := json.Marshal(&logisticEvent5)
	var m5 map[string]interface{}
	_ = json.Unmarshal(b5, &m5)
	cloudEvent5 := buildCloudEvent("id", "Logistic.PCL.Fleet.TransportEnded", "source", "subject", time.Now(), m5)
	body5, _ := json.Marshal(cloudEvent5)
	snsMsg5 := MessageReceived{Message: string(body5)}
	bodySNS5, _ := json.Marshal(&snsMsg5)
	test5 := testType{name: "Trasporto Concluso Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS5)}}}}}
	tests = append(tests, test5)

	logisticEvent6 := models.ParcelUnloaded{
		TransportID: "4ADSF-0900-020923",
		ParcelID:   parcelID,
		FacilityID: "FI1",
		Timestamp:  time.Now()}
	b6, _ := json.Marshal(&logisticEvent6)
	var m6 map[string]interface{}
	_ = json.Unmarshal(b6, &m6)
	cloudEvent6 := buildCloudEvent("id", "Logistic.PCL.Fleet.ParcelUnloaded", "source", "subject", time.Now(), m6)
	body6, _ := json.Marshal(cloudEvent6)
	snsMsg6 := MessageReceived{Message: string(body6)}
	bodySNS6, _ := json.Marshal(&snsMsg6)
	test6 := testType{name: "Scaricamento Veicolo Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS6)}}}}}
	tests = append(tests, test6)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleEvents(tt.args.ctx, tt.args.sqsEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleEvents() = %v, want %v", got, tt.want)
			}
		})
	}
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
