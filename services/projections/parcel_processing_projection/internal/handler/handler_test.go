package handler

import (
	"context"
	"encoding/json"
	"parcel-processing-projection/internal"
	"parcel-processing-projection/internal/models"
	"parcel-processing-projection/internal/repository"
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
	collectionName := "facility"
	opts := options.CreateCollection()
	opts.SetValidator(bson.M{})
	err = client.Database("logistic").CreateCollection(context.Background(), collectionName, opts)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' created.", collectionName)
		collection := client.Database("logistic").Collection(collectionName)
		documents := []interface{}{
			repository.Facility{
				FacilityID:   "MI1",
				Capacity:     3000,
				FacilityType: "HUB",
				FacilityLocation: repository.FacilityLocation{
					Address: "Viale dei Soldi 112",
					Zipcode: "29065",
					City:    "Milano",
					Nation:  "Italia",
				},
				Latitudine: "123,123",
				Longitude:  "987.987",
				Company:    "SDA",
				Connections: []repository.Connection{
					{
						FacilityDestinationID: "BO1",
						Distance:              300,
					},
				},
			},
			repository.Facility{
				FacilityID:   "BO1",
				Capacity:     3000,
				FacilityType: "HUB",
				FacilityLocation: repository.FacilityLocation{
					Address: "Viale del piazzale 112",
					Zipcode: "29065",
					City:    "Bologna",
					Nation:  "Italia",
				},
				Latitudine: "123,123",
				Longitude:  "987.987",
				Company:    "SDA",
				Connections: []repository.Connection{
					{
						FacilityDestinationID: "MI1",
						Distance:              300,
					},
				},
			},
			repository.Facility{
				FacilityID:   "FI1",
				Capacity:     3000,
				FacilityType: "HUB",
				FacilityLocation: repository.FacilityLocation{
					Address: "Viale Dante 748",
					Zipcode: "29065",
					City:    "Firenze",
					Nation:  "Italia",
				},
				Latitudine: "123,123",
				Longitude:  "987.987",
				Company:    "SDA",
				Connections: []repository.Connection{
					{
						FacilityDestinationID: "BO1",
						Distance:              300,
					},
				},
			},
			repository.Facility{
				FacilityID:   "RM1",
				Capacity:     3000,
				FacilityType: "HUB",
				FacilityLocation: repository.FacilityLocation{
					Address: "Viale dei lupi 748",
					Zipcode: "29065",
					City:    "Roma",
					Nation:  "Italia",
				},
				Latitudine: "123,123",
				Longitude:  "987.987",
				Company:    "SDA",
				Connections: []repository.Connection{
					{
						FacilityDestinationID: "FI1",
						Distance:              300,
					},
				},
			},
		}
		_, err := collection.InsertMany(context.TODO(), documents)
		if err != nil {
			log.Err(err)
		} else {
			log.Debug().Msgf("Collection '%s' initialized.", collectionName)
		}
	}

	collectionName = "routes"
	opts = options.CreateCollection()
	opts.SetValidator(bson.M{})
	err = client.Database("logistic").CreateCollection(context.Background(), collectionName, opts)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' created.", collectionName)
		collection := client.Database("logistic").Collection(collectionName)
		documents := []interface{}{
			repository.Route{
				SourceFacilityID: "RM1",
				DestFacilityID:   "FI1",
				Cutoff: []string{
					"15:30",
					"17:30",
					"18:30",
				},
			},
			repository.Route{
				SourceFacilityID: "FI1",
				DestFacilityID:   "BO1",
				Cutoff: []string{
					"16:30",
					"17:30",
					"18:30",
				},
			},
			repository.Route{
				SourceFacilityID: "BO1",
				DestFacilityID:   "MI1",
				Cutoff: []string{
					"16:30",
					"17:30",
					"18:30",
				},
			},
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

	logisticEvent2 := models.ParcelProcessed{
		FacilityID:            "RM1",
		SortingMachineID:      "12ASDFE32",
		ParcelID:              parcelID,
		DestinationFacilityID: "FI1",
		Timestamp:             time.Now()}
	b2, _ := json.Marshal(&logisticEvent2)
	var m2 map[string]interface{}
	_ = json.Unmarshal(b2, &m2)
	cloudEvent2 := buildCloudEvent("id", "Logistic.PCL.Faility.ParcelProcessed", "source", "subject", time.Now(), m2)
	body2, _ := json.Marshal(cloudEvent2)
	snsMsg2 := MessageReceived{Message: string(body2)}
	bodySNS2, _ := json.Marshal(&snsMsg2)
	test2 := testType{name: "Processamento Pacco Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS2)}}}}}
	tests = append(tests, test2)

	logisticEvent3 := models.ParcelLoaded{
		TransportID: "1ADSF-1530-020123",
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

	logisticEvent4 := models.TransportStarted{
		TransportID: "1ADSF-1530-020123",
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

	logisticEvent5 := models.TransportEnded{
		TransportID: "1ADSF-1530-020123",
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
		TransportID: "1ADSF-1530-020123",
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

	logisticEvent7 := models.ParcelProcessed{
		FacilityID:            "FI1",
		SortingMachineID:      "543GFD543",
		ParcelID:              parcelID,
		DestinationFacilityID: "BO1",
		Timestamp:             time.Now()}
	b7, _ := json.Marshal(&logisticEvent7)
	var m7 map[string]interface{}
	_ = json.Unmarshal(b7, &m7)
	cloudEvent7 := buildCloudEvent("id", "Logistic.PCL.Facility.ParcelProcessed", "source", "subject", time.Now(), m7)
	body7, _ := json.Marshal(cloudEvent7)
	snsMsg7 := MessageReceived{Message: string(body7)}
	bodySNS7, _ := json.Marshal(&snsMsg7)
	test7 := testType{name: "Processamento Pacco Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS7)}}}}}
	tests = append(tests, test7)

	logisticEvent8 := models.ParcelLoaded{
		TransportID: "2ADSF-1530-020123",
		ParcelID:   parcelID,
		FacilityID: "FI1",
		Timestamp:  time.Now()}
	b8, _ := json.Marshal(&logisticEvent8)
	var m8 map[string]interface{}
	_ = json.Unmarshal(b8, &m8)
	cloudEvent8 := buildCloudEvent("id", "Logistic.PCL.Fleet.ParcelLoaded", "source", "subject", time.Now(), m8)
	body8, _ := json.Marshal(cloudEvent8)
	snsMsg8 := MessageReceived{Message: string(body8)}
	bodySNS8, _ := json.Marshal(&snsMsg8)
	test8 := testType{name: "Caricamento Pacco Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS8)}}}}}
	tests = append(tests, test8)

	logisticEvent9 := models.TransportStarted{
		TransportID: "2ADSF-1530-020123",
		VehicleLicensePlate:   "RD312FT",
		SourceFacilityID:      "FI1",
		DestinationFacilityID: "BO1",
		Timestamp:             time.Now()}
	b9, _ := json.Marshal(&logisticEvent9)
	var m9 map[string]interface{}
	_ = json.Unmarshal(b9, &m9)
	cloudEvent9 := buildCloudEvent("id", "Logistic.PCL.Fleet.TransportStarted", "source", "subject", time.Now(), m9)
	body9, _ := json.Marshal(cloudEvent9)
	snsMsg9 := MessageReceived{Message: string(body9)}
	bodySNS9, _ := json.Marshal(&snsMsg9)
	test9 := testType{name: "Inizio trasporto Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS9)}}}}}
	tests = append(tests, test9)

	logisticEvent10 := models.TransportEnded{
		TransportID: "2ADSF-1530-020123",
		VehicleLicensePlate: "RD312FT",
		FacilityID:          "BO1",
		Timestamp:           time.Now()}
	b10, _ := json.Marshal(&logisticEvent10)
	var m10 map[string]interface{}
	_ = json.Unmarshal(b10, &m10)
	cloudEvent10 := buildCloudEvent("id", "Logistic.PCL.Fleet.TransportEnded", "source", "subject", time.Now(), m10)
	body10, _ := json.Marshal(cloudEvent10)
	snsMsg10 := MessageReceived{Message: string(body10)}
	bodySNS10, _ := json.Marshal(&snsMsg10)
	test10 := testType{name: "Trasporto Concluso Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS10)}}}}}
	tests = append(tests, test10)

	logisticEvent11 := models.ParcelUnloaded{
		TransportID: "2ADSF-1530-020123",
		ParcelID:   parcelID,
		FacilityID: "BO1",
		Timestamp:  time.Now()}
	b11, _ := json.Marshal(&logisticEvent11)
	var m11 map[string]interface{}
	_ = json.Unmarshal(b11, &m11)
	cloudEvent11 := buildCloudEvent("id", "Logistic.PCL.Fleet.ParcelUnloaded", "source", "subject", time.Now(), m11)
	body11, _ := json.Marshal(cloudEvent11)
	snsMsg11 := MessageReceived{Message: string(body11)}
	bodySNS11, _ := json.Marshal(&snsMsg11)
	test11 := testType{name: "Scaricamento Veicolo Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS11)}}}}}
	tests = append(tests, test11)

	logisticEvent12 := models.ParcelProcessed{
		FacilityID:            "BO1",
		SortingMachineID:      "876TRE654",
		ParcelID:              parcelID,
		DestinationFacilityID: "MI1",
		Timestamp:             time.Now()}
	b12, _ := json.Marshal(&logisticEvent12)
	var m12 map[string]interface{}
	_ = json.Unmarshal(b12, &m12)
	cloudEvent12 := buildCloudEvent("id", "Logistic.PCL.Facility.ParcelProcessed", "source", "subject", time.Now(), m12)
	body12, _ := json.Marshal(cloudEvent12)
	snsMsg12 := MessageReceived{Message: string(body12)}
	bodySNS12, _ := json.Marshal(&snsMsg12)
	test12 := testType{name: "Processamento Pacco Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS12)}}}}}
	tests = append(tests, test12)

	logisticEvent13 := models.ParcelLoaded{
		TransportID: "3ADSF-1530-020123",
		ParcelID:   parcelID,
		FacilityID: "BO1",
		Timestamp:  time.Now()}
	b13, _ := json.Marshal(&logisticEvent13)
	var m13 map[string]interface{}
	_ = json.Unmarshal(b13, &m13)
	cloudEvent13 := buildCloudEvent("id", "Logistic.PCL.Fleet.ParcelLoaded", "source", "subject", time.Now(), m13)
	body13, _ := json.Marshal(cloudEvent13)
	snsMsg13 := MessageReceived{Message: string(body13)}
	bodySNS13, _ := json.Marshal(&snsMsg13)
	test13 := testType{name: "Caricamento Pacco Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS13)}}}}}
	tests = append(tests, test13)

	logisticEvent14 := models.TransportStarted{
		TransportID: "3ADSF-1530-020123",
		VehicleLicensePlate:   "RD312FT",
		SourceFacilityID:      "BO1",
		DestinationFacilityID: "MI1",
		Timestamp:             time.Now()}
	b14, _ := json.Marshal(&logisticEvent14)
	var m14 map[string]interface{}
	_ = json.Unmarshal(b14, &m14)
	cloudEvent14 := buildCloudEvent("id", "Logistic.PCL.Fleet.TransportStarted", "source", "subject", time.Now(), m14)
	body14, _ := json.Marshal(cloudEvent14)
	snsMsg14 := MessageReceived{Message: string(body14)}
	bodySNS14, _ := json.Marshal(&snsMsg14)
	test14 := testType{name: "Inizio trasporto Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS14)}}}}}
	tests = append(tests, test14)

	logisticEvent15 := models.TransportEnded{
		TransportID: "3ADSF-1530-020123",
		VehicleLicensePlate: "RD312FT",
		FacilityID:          "MI1",
		Timestamp:           time.Now()}
	b15, _ := json.Marshal(&logisticEvent15)
	var m15 map[string]interface{}
	_ = json.Unmarshal(b15, &m15)
	cloudEvent15 := buildCloudEvent("id", "Logistic.PCL.Fleet.TransportEnded", "source", "subject", time.Now(), m15)
	body15, _ := json.Marshal(cloudEvent15)
	snsMsg15 := MessageReceived{Message: string(body15)}
	bodySNS15, _ := json.Marshal(&snsMsg15)
	test15 := testType{name: "Trasporto Concluso Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS15)}}}}}
	tests = append(tests, test15)

	logisticEvent16 := models.ParcelUnloaded{
		TransportID: "3ADSF-1530-020123",
		ParcelID:   parcelID,
		FacilityID: "MI1",
		Timestamp:  time.Now()}
	b16, _ := json.Marshal(&logisticEvent16)
	var m16 map[string]interface{}
	_ = json.Unmarshal(b16, &m16)
	cloudEvent16 := buildCloudEvent("id", "Logistic.PCL.Fleet.ParcelUnloaded", "source", "subject", time.Now(), m16)
	body16, _ := json.Marshal(cloudEvent16)
	snsMsg16 := MessageReceived{Message: string(body16)}
	bodySNS16, _ := json.Marshal(&snsMsg16)
	test16 := testType{name: "Scaricamento Veicolo Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS16)}}}}}
	tests = append(tests, test16)

	logisticEvent17 := models.DeliveryCompleted{
		ParcelID:  parcelID,
		Timestamp: time.Now()}
	b17, _ := json.Marshal(&logisticEvent17)
	var m17 map[string]interface{}
	_ = json.Unmarshal(b17, &m17)
	cloudEvent17 := buildCloudEvent("id", "Logistic.PCL.Parcel.DeliveryCompleted", "source", "subject", time.Now(), m17)
	body17, _ := json.Marshal(cloudEvent17)
	snsMsg17 := MessageReceived{Message: string(body17)}
	bodySNS17, _ := json.Marshal(&snsMsg17)
	test17 := testType{name: "Consegna pacco Test", args: args{ctx: context.TODO(), sqsEvent: events.SQSEvent{Records: []events.SQSMessage{{Body: string(bodySNS17)}}}}}
	tests = append(tests, test17)

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
