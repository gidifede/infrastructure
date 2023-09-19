package manager

import (
	"context"
	"facility-vehicle-projection/internal"
	"facility-vehicle-projection/internal/models"
	"facility-vehicle-projection/internal/repository"
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

func TestTransportEndedManager_ManageEvent_OK_FacilityVehicleDocExists(t *testing.T) {
	facilityID := uuid.New().String()
	vehicleID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.Facility{
			FacilityID:   facilityID,
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
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.VehicleCollection
	collection = client.Database("logistic").Collection(collectionName)
	vehicle := repository.Vehicle{
		VehicleID: vehicleID,
		Capacity:  12345,
		Type:      "a vehicle type",
	}
	documents = []interface{}{vehicle}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.FacilityVehicleCollection
	collection = client.Database("logistic").Collection(collectionName)
	vehicleAlreadyInFacility := repository.VehicleInFacility{VehicleID: "fsfdsfs"}
	documents = []interface{}{
		repository.FacilityVehicle{
			FacilityID: facilityID,
			Vehicles:   []repository.VehicleInFacility{vehicleAlreadyInFacility},
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportEnded{
		FacilityID:          facilityID,
		VehicleLicensePlate: vehicleID,
		Timestamp:           time.Now().UTC(),
	}
	expFacilityVehicle := repository.FacilityVehicle{
		FacilityID: facilityID,
		Vehicles: []repository.VehicleInFacility{
			vehicleAlreadyInFacility,
			{
				VehicleID:   vehicleID,
				Status:      "arrived",
				ArrivedTime: event.Timestamp,
			},
		},
	}
	type fields struct {
		logisticEvent models.TransportEnded
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{name: "Event OK. Facility vehicle doc exists", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportEndedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportEndedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetFacilityVehicleDocumentByID(context.TODO(), facilityID)

			if err != nil {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Cannot check test on doc id  %v", facilityID)
			}
			// Avoid comparing timestamp
			for i, item := range got.Vehicles {
				expFacilityVehicle.Vehicles[i].ArrivedTime = item.ArrivedTime
			}
			if !reflect.DeepEqual(got, expFacilityVehicle) {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Read %v, want %v", got, expFacilityVehicle)
			}
		})
	}
}

func TestTransportEndedManager_ManageEvent_OK_FacilityVehicleDocDoesntExist(t *testing.T) {
	facilityID := uuid.New().String()
	vehicleID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.Facility{
			FacilityID:   facilityID,
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
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.VehicleCollection
	collection = client.Database("logistic").Collection(collectionName)
	vehicle := repository.Vehicle{
		VehicleID: vehicleID,
		Capacity:  12345,
		Type:      "a vehicle type",
	}
	documents = []interface{}{vehicle}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportEnded{
		FacilityID:          facilityID,
		VehicleLicensePlate: vehicleID,
		Timestamp:           time.Now().UTC(),
	}
	expFacilityVehicle := repository.FacilityVehicle{
		FacilityID: facilityID,
		Vehicles: []repository.VehicleInFacility{
			{
				VehicleID:   vehicleID,
				Status:      "arrived",
				ArrivedTime: event.Timestamp,
			},
		},
	}
	type fields struct {
		logisticEvent models.TransportEnded
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{name: "Event OK. Facility vehicle doc doesn't exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportEndedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportEndedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetFacilityVehicleDocumentByID(context.TODO(), facilityID)

			if err != nil {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Cannot check test on doc id  %v", facilityID)
			}
			// Avoid comparing timestamp
			for i, item := range got.Vehicles {
				expFacilityVehicle.Vehicles[i].ArrivedTime = item.ArrivedTime
			}
			if !reflect.DeepEqual(got, expFacilityVehicle) {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Read %v, want %v", got, expFacilityVehicle)
			}
		})
	}
}

func TestTransportEndedManager_ManageEvent_KO_VehicleDoesntExist(t *testing.T) {
	facilityID := uuid.New().String()
	vehicleID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.Facility{
			FacilityID:   facilityID,
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
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportEnded{
		FacilityID:          facilityID,
		VehicleLicensePlate: vehicleID,
		Timestamp:           time.Now().UTC(),
	}
	type fields struct {
		logisticEvent models.TransportEnded
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{name: "Event KO. Vehicle doesn't exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
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

func TestTransportEndedManager_ManageEvent_KO_FacilityDoesntExist(t *testing.T) {
	facilityID := uuid.New().String()
	vehicleID := uuid.New().String()

	// Initialize collections

	// Test
	event := models.TransportEnded{
		FacilityID:          facilityID,
		VehicleLicensePlate: vehicleID,
		Timestamp:           time.Now().UTC(),
	}
	type fields struct {
		logisticEvent models.TransportEnded
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{name: "Event KO. Facility doesn't exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
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

func TestTransportEndedManager_ManageEvent_KO_VehicleInFacilityDocAlreadyExists(t *testing.T) {
	facilityID := uuid.New().String()
	vehicleID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.Facility{
			FacilityID:   facilityID,
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
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.VehicleCollection
	collection = client.Database("logistic").Collection(collectionName)
	vehicle := repository.Vehicle{
		VehicleID: vehicleID,
		Capacity:  12345,
		Type:      "a vehicle type",
	}
	documents = []interface{}{vehicle}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.FacilityVehicleCollection
	collection = client.Database("logistic").Collection(collectionName)
	vehicleAlreadyInFacility := repository.VehicleInFacility{VehicleID: vehicleID}
	documents = []interface{}{
		repository.FacilityVehicle{
			FacilityID: facilityID,
			Vehicles:   []repository.VehicleInFacility{vehicleAlreadyInFacility},
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportEnded{
		FacilityID:          facilityID,
		VehicleLicensePlate: vehicleID,
		Timestamp:           time.Now().UTC(),
	}
	expFacilityVehicle := repository.FacilityVehicle{
		FacilityID: facilityID,
		Vehicles:   []repository.VehicleInFacility{vehicleAlreadyInFacility},
	}
	type fields struct {
		logisticEvent models.TransportEnded
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{name: "Event KO. Vehicle in Facility doc already exists", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportEndedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportEndedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetFacilityVehicleDocumentByID(context.TODO(), facilityID)

			if err != nil {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Cannot check test on doc id  %v", facilityID)
			}
			// Avoid comparing timestamp
			for i, item := range got.Vehicles {
				expFacilityVehicle.Vehicles[i].ArrivedTime = item.ArrivedTime
			}
			if !reflect.DeepEqual(got, expFacilityVehicle) {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Read %v, want %v", got, expFacilityVehicle)
			}
		})
	}
}
