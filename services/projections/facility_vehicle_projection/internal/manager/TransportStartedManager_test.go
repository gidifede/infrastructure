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

func TestTransportStartedManager_ManageEvent_OK_VehicleRemovedFromFacility(t *testing.T) {
	vehicleID := uuid.New().String()
	facilityID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityVehicleCollection
	collection := client.Database("logistic").Collection(collectionName)
	vehicleAlreadyInFacility1 := repository.VehicleInFacility{VehicleID: "fddfaafd"}
	vehicleAlreadyInFacility2 := repository.VehicleInFacility{VehicleID: vehicleID}
	documents := []interface{}{
		repository.FacilityVehicle{
			FacilityID: facilityID,
			Vehicles:   []repository.VehicleInFacility{vehicleAlreadyInFacility1, vehicleAlreadyInFacility2},
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportStarted{
		VehicleLicensePlate:   vehicleID,
		SourceFacilityID:      facilityID,
		DestinationFacilityID: "agfadgda",
		Timestamp:             time.Now().UTC(),
	}
	expFacilityVehicle := repository.FacilityVehicle{
		FacilityID: facilityID,
		Vehicles: []repository.VehicleInFacility{
			vehicleAlreadyInFacility1,
		},
	}

	type fields struct {
		logisticEvent models.TransportStarted
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
		{name: "Event OK. Vehicle removed from facility", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportStartedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportStartedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
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

func TestTransportStartedManager_ManageEvent_KO_FacilityVehicleDoesntExist(t *testing.T) {
	vehicleID := uuid.New().String()
	facilityID := uuid.New().String()

	// Initialize collections
	// nothing

	// Test
	event := models.TransportStarted{
		VehicleLicensePlate:   vehicleID,
		SourceFacilityID:      facilityID,
		DestinationFacilityID: "agfadgda",
		Timestamp:             time.Now().UTC(),
	}

	type fields struct {
		logisticEvent models.TransportStarted
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
		{name: "Event KO. Facility vehicle doesn't exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportStartedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportStartedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportStartedManager_ManageEvent_KO_VehicleIsNotInFacility(t *testing.T) {
	vehicleID := uuid.New().String()
	facilityID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityVehicleCollection
	collection := client.Database("logistic").Collection(collectionName)
	vehicleAlreadyInFacility1 := repository.VehicleInFacility{VehicleID: "fddfaafd"}
	documents := []interface{}{
		repository.FacilityVehicle{
			FacilityID: facilityID,
			Vehicles:   []repository.VehicleInFacility{vehicleAlreadyInFacility1},
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.TransportStarted{
		VehicleLicensePlate:   vehicleID,
		SourceFacilityID:      facilityID,
		DestinationFacilityID: "agfadgda",
		Timestamp:             time.Now().UTC(),
	}
	expFacilityVehicle := repository.FacilityVehicle{
		FacilityID: facilityID,
		Vehicles: []repository.VehicleInFacility{
			vehicleAlreadyInFacility1,
		},
	}

	type fields struct {
		logisticEvent models.TransportStarted
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
		{name: "Event KO. Vehicle is not in facility", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportStartedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportStartedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
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
