package manager

import (
	"context"
	"facility-parcel-expected-projection/internal"
	"facility-parcel-expected-projection/internal/models"
	"facility-parcel-expected-projection/internal/repository"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

func init() {
	clientOpts := options.Client().ApplyURI("mongodb://root:rootpassword@localhost:27017") //.ApplyURI("mongodb://test:****@dds-bp*******1.mongodb.rds.aliyuncs.com:3717,dds-bp*******2.mongodb.rds.aliyuncs.com:3717/admin?replicaSet=mgset-XXXXX&ssl=true")
	var err error
	client, err = mongo.Connect(context.TODO(), clientOpts)
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
}

func TestTransportEndedManager_ManageEvent_OK_FacilityExpectedParcelUpdated(t *testing.T) {
	facilityID := uuid.New().String()
	vehicleID := uuid.New().String()
	transportID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.Facility{
			FacilityID:   facilityID,
			Capacity:     10,
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
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.VehicleTransportCollection
	collection = client.Database("logistic").Collection(collectionName)
	parcelsOnVehicle := []string{"", "", "", "", ""}
	documents = []interface{}{
		repository.VehicleTransport{
			VehicleID:   "adfgafa",
			TransportID: "afdafad",
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     parcelsOnVehicle,
			TransportID: transportID,
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     []string{"fsfdafa"},
			TransportID: "fdfdsfg",
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.FacilityExpectedParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	facilityExpectedParcel1 := repository.FacilityExpectedParcel{
		FacilityID: facilityID,
		Status:     repository.FacilityExpectedParcelStatusHealthy,
		Counter:    4,
	}
	facilityExpectedParcel2 := repository.FacilityExpectedParcel{
		FacilityID: "fadfdaf",
		Status:     repository.FacilityExpectedParcelStatusHealthy,
		Counter:    4,
	}
	documents = []interface{}{facilityExpectedParcel1, facilityExpectedParcel2}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportEnded{
		VehicleLicensePlate: vehicleID,
		FacilityID:          facilityID,
		Timestamp:           time.Now().UTC(),
		TransportID:         transportID,
	}
	expFacilityExpectedParcel := repository.FacilityExpectedParcel{
		FacilityID: facilityID,
		Status:     repository.FacilityExpectedParcelStatusUnhealthy,
		Counter:    facilityExpectedParcel1.Counter + len(parcelsOnVehicle),
	}

	type fields struct {
		logisticEvent models.TransportEnded
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Event OK. Facility Expected Parcel Updated", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportEndedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportEndedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetFacilityExpectedParcelDocumentByID(context.TODO(), facilityID)
			if err != nil {
				t.Errorf("TransportEndedManager.ManageEvent(). Cannot check test on doc id  %v", facilityID)
			}
			if !reflect.DeepEqual(got, expFacilityExpectedParcel) {
				t.Errorf("TransportEndedManager.ManageEvent(). Read %v, want %v", got, expFacilityExpectedParcel)
			}
		})
	}
}

func TestTransportEndedManager_ManageEvent_OK_FacilityExpectedParcelInsertNewDoc(t *testing.T) {
	facilityID := uuid.New().String()
	vehicleID := uuid.New().String()
	transportID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.Facility{
			FacilityID:   facilityID,
			Capacity:     10,
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
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.VehicleTransportCollection
	collection = client.Database("logistic").Collection(collectionName)
	parcelsOnVehicle := []string{"", "", "", "", ""}
	documents = []interface{}{
		repository.VehicleTransport{
			VehicleID:   "adfgafa",
			TransportID: "afdafad",
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     parcelsOnVehicle,
			TransportID: transportID,
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     []string{"fsfdafa"},
			TransportID: "fdfdsfg",
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.FacilityExpectedParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	facilityExpectedparcel1 := repository.FacilityExpectedParcel{
		FacilityID: "fadfdaf",
		Status:     repository.FacilityExpectedParcelStatusHealthy,
		Counter:    4,
	}
	documents = []interface{}{facilityExpectedparcel1}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportEnded{
		VehicleLicensePlate: vehicleID,
		FacilityID:          facilityID,
		Timestamp:           time.Now().UTC(),
		TransportID:         transportID,
	}
	expFacilityExpectedParcel := repository.FacilityExpectedParcel{
		FacilityID: facilityID,
		Status:     repository.FacilityExpectedParcelStatusHealthy,
		Counter:    len(parcelsOnVehicle),
	}

	type fields struct {
		logisticEvent models.TransportEnded
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Event OK. Facility Expected Parcel Insert New Doc", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportEndedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportEndedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetFacilityExpectedParcelDocumentByID(context.TODO(), facilityID)
			if err != nil {
				t.Errorf("TransportEndedManager.ManageEvent(). Cannot check test on doc id  %v", facilityID)
			}
			if !reflect.DeepEqual(got, expFacilityExpectedParcel) {
				t.Errorf("TransportEndedManager.ManageEvent(). Read %v, want %v", got, expFacilityExpectedParcel)
			}
		})
	}
}

func TestTransportEndedManager_ManageEvent_KO_FacilityDoesntExist(t *testing.T) {
	facilityID := uuid.New().String()
	vehicleID := uuid.New().String()
	transportID := uuid.New().String()

	// Initialize collections
	collectionName := repository.VehicleTransportCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelsOnVehicle := []string{"", "", "", "", ""}
	documents := []interface{}{
		repository.VehicleTransport{
			VehicleID:   "adfgafa",
			TransportID: "afdafad",
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     parcelsOnVehicle,
			TransportID: transportID,
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     []string{"fsfdafa"},
			TransportID: "fdfdsfg",
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.FacilityExpectedParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	facilityExpectedParcel1 := repository.FacilityExpectedParcel{
		FacilityID: "fadfdaf",
		Status:     repository.FacilityExpectedParcelStatusHealthy,
		Counter:    4,
	}
	documents = []interface{}{facilityExpectedParcel1}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportEnded{
		VehicleLicensePlate: vehicleID,
		FacilityID:          facilityID,
		Timestamp:           time.Now().UTC(),
		TransportID:         transportID,
	}

	type fields struct {
		logisticEvent models.TransportEnded
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Event KO. Facility Doesn't Exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportEndedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportEndedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportEndedManager_ManageEvent_KO_VehicleTransportDoesntExist(t *testing.T) {
	facilityID := uuid.New().String()
	vehicleID := uuid.New().String()
	transportID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.Facility{
			FacilityID:   facilityID,
			Capacity:     10,
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
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.VehicleTransportCollection
	collection = client.Database("logistic").Collection(collectionName)
	documents = []interface{}{
		repository.VehicleTransport{
			VehicleID:   "adfgafa",
			TransportID: "afdafad",
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     []string{"fsfdafa"},
			TransportID: "fdfdsfg",
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.FacilityExpectedParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	facilityExpectedParcel1 := repository.FacilityExpectedParcel{
		FacilityID: facilityID,
		Status:     repository.FacilityExpectedParcelStatusHealthy,
		Counter:    4,
	}
	facilityExpectedParcel2 := repository.FacilityExpectedParcel{
		FacilityID: "fadfdaf",
		Status:     repository.FacilityExpectedParcelStatusHealthy,
		Counter:    4,
	}
	documents = []interface{}{facilityExpectedParcel1, facilityExpectedParcel2}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportEnded{
		VehicleLicensePlate: vehicleID,
		FacilityID:          facilityID,
		Timestamp:           time.Now().UTC(),
		TransportID:         transportID,
	}

	type fields struct {
		logisticEvent models.TransportEnded
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Event KO. Vehicle Transport Doesn't Exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportEndedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportEndedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
