package manager

import (
	"context"
	"facility-parcel-projection/internal"
	"facility-parcel-projection/internal/models"
	"facility-parcel-projection/internal/repository"
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

func TestTransportStartedManager_ManageEvent_OK_ParcelRemovedFromFacility(t *testing.T) {
	vehicleID := uuid.New().String()
	transportID := uuid.New().String()
	parcelID := uuid.New().String()
	facilityID := uuid.New().String()
	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility1 := repository.ParcelInFacility{ParcelID: "fddfaafd"}
	parcelAlreadyInFacility2 := repository.ParcelInFacility{ParcelID: parcelID}
	documents := []interface{}{
		repository.FacilityParcel{
			FacilityID: facilityID,
			Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility1, parcelAlreadyInFacility2},
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
			Parcels:     []string{parcelID},
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

	// Test
	event := models.TransportStarted{
		VehicleLicensePlate:   vehicleID,
		SourceFacilityID:      facilityID,
		DestinationFacilityID: "agfadgda",
		Timestamp:             time.Now().UTC(),
		TransportID:           transportID,
	}
	expFacilityParcel := repository.FacilityParcel{
		FacilityID: facilityID,
		Parcels: []repository.ParcelInFacility{
			parcelAlreadyInFacility1,
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
		{name: "Event OK. Parcel removed from facility", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &TransportStartedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("TransportStartedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetFacilityParcelDocumentByID(context.TODO(), facilityID)
			if err != nil {
				t.Errorf("TransportStartedManager.ManageEvent(). Cannot check test on doc id  %v", facilityID)
			}
			if !reflect.DeepEqual(got, expFacilityParcel) {
				t.Errorf("TransportStartedManager.ManageEvent(). Read %v, want %v", got, expFacilityParcel)
			}
		})
	}
}

func TestTransportStartedManager_ManageEvent_KO_CannotFindVehicleTransport(t *testing.T) {
	vehicleID := uuid.New().String()
	transportID := uuid.New().String()

	// Test
	event := models.TransportStarted{
		VehicleLicensePlate:   vehicleID,
		SourceFacilityID:      "gsgfs",
		DestinationFacilityID: "agfadgda",
		Timestamp:             time.Now().UTC(),
		TransportID:           transportID,
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
		{name: "Event KO. Cannot find vehicle transport", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
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

func TestTransportStartedManager_ManageEvent_KO_CannotFindFacility(t *testing.T) {
	vehicleID := uuid.New().String()
	parcelID := uuid.New().String()
	facilityID := uuid.New().String()
	transportID := uuid.New().String()

	// Initialize collections
	collectionName := repository.VehicleTransportCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.VehicleTransport{
			VehicleID:   "adfgafa",
			TransportID: "fsgfd",
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     []string{parcelID},
			TransportID: transportID,
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     []string{"fsfdafa"},
			TransportID: "fsagda",
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
		TransportID:           transportID,
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
		{name: "Event KO. Cannot find facility", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
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

func TestTransportStartedManager_ManageEvent_KO_CannotFindParcelInFacility(t *testing.T) {
	vehicleID := uuid.New().String()
	parcelID := uuid.New().String()
	facilityID := uuid.New().String()
	transportID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility1 := repository.ParcelInFacility{ParcelID: "fddfaafd"}
	documents := []interface{}{
		repository.FacilityParcel{
			FacilityID: facilityID,
			Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility1},
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
			TransportID: "fsgfs",
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     []string{parcelID},
			TransportID: transportID,
		},
		repository.VehicleTransport{
			VehicleID:   vehicleID,
			Parcels:     []string{"fsfdafa"},
			TransportID: "fagdag",
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
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
		TransportID:           transportID,
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
		{name: "Event KO. Cannot find parcel in facility", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
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
