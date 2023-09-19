package handler

import (
	"context"
	"encoding/json"
	"parcel-sla-projection/internal"
	"parcel-sla-projection/internal/models"
	"parcel-sla-projection/internal/repository"
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
	collectionName := "product"
	opts := options.CreateCollection()
	opts.SetValidator(bson.M{})
	err = client.Database("logistic").CreateCollection(context.Background(), collectionName, opts)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' created.", collectionName)
		collection := client.Database("logistic").Collection(collectionName)
		documents := []interface{}{
			repository.Product{Name: "PosteDeliveryWeb", SLA: "3"},
			repository.Product{Name: "PosteDeliveryWebExpress", SLA: "1"},
		}
		_, err := collection.InsertMany(context.TODO(), documents)
		if err != nil {
			log.Err(err)
		} else {
			log.Debug().Msgf("Collection '%s' initialized.", collectionName)
		}
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
	logisticEvent := models.Accepted{
		Parcel: struct {
			Name string "json:\"name\""
			ID   string "json:\"id\""
			Type string "json:\"type\""
		}{
			Name: "PosteDeliveryWeb", ID: parcelID, Type: "BOX"},
		FacilityID: "RM1",
		Sender: struct {
			Name    string "json:\"name\""
			Address string "json:\"address\""
			Zipcode string "json:\"zipcode\""
			City    string "json:\"city\""
			Nation  string "json:\"nation\""
		}{
			Name:    "Pippo Baudo",
			Nation:  "Italia",
			City:    "Roma",
			Address: "Viale Europa 190",
			Zipcode: "000123",
		},
		Receiver: struct {
			Name    string "json:\"name\""
			Address string "json:\"address\""
			Zipcode string "json:\"zipcode\""
			City    string "json:\"city\""
			Nation  string "json:\"nation\""
			Number  string "json:\"number\""
			Email   string "json:\"email\""
		}{
			Name:    "Silvio Berlusconi",
			Nation:  "Italia",
			City:    "Milano",
			Address: "Viale delle Colline 90",
			Zipcode: "000111",
			Number:  "3333333333",
			Email:   "silvio.berlusconi@mediaset.it",
		},
		Timestamp: time.Now()}
	b, _ := json.Marshal(&logisticEvent)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	cloudEvent := buildCloudEvent("id", "Logistic.PCL.Parcel.Accepted", "source", "subject", time.Now(), m)
	body, _ := json.Marshal(cloudEvent)
	snsMsg := MessageReceived{Message: string(body)}
	bodySNS, _ := json.Marshal(&snsMsg)
	test := testType{name: "Accettazione Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS)}}}}}
	tests = append(tests, test)

	logisticEvent2 := models.DeliveryCompleted{
		ParcelID:  parcelID,
		Timestamp: time.Now()}
	b2, _ := json.Marshal(&logisticEvent2)
	var m2 map[string]interface{}
	_ = json.Unmarshal(b2, &m2)
	cloudEvent2 := buildCloudEvent("id", "Logistic.PCL.Parcel.DeliveryCompleted", "source", "subject", time.Now(), m2)
	body2, _ := json.Marshal(cloudEvent2)
	snsMsg2 := MessageReceived{Message: string(body2)}
	bodySNS2, _ := json.Marshal(&snsMsg2)
	test2 := testType{name: "Consegna pacco Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS2)}}}}}
	tests = append(tests, test2)

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
