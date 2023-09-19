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

func TestParcelProcessingFailedManager_ManageEvent_OK_FacilityParcelDocUpdated(t *testing.T) {
	facilityID := uuid.New().String()
	parcelID := uuid.New().String()
	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility1 := repository.ParcelInFacility{ParcelID: "fsfdsfs"}
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

	// Test
	event := models.ParcelProcessingFailed{
		ParcelID:         parcelID,
		Timestamp:        time.Now().UTC(),
		FacilityID:       facilityID,
		SortingMachineID: "",
		ErrMsg:           "an error msg",
	}

	parcelAlreadyInFacility2.Status = "processingFailed"
	expFacilityParcel := repository.FacilityParcel{
		FacilityID: facilityID,
		Parcels: []repository.ParcelInFacility{
			parcelAlreadyInFacility1,
			parcelAlreadyInFacility2},
	}

	type fields struct {
		logisticEvent models.ParcelProcessingFailed
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
		{name: "Event OK. Facility parcel doc updated", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelProcessingFailedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelProcessingFailedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetFacilityParcelDocumentByID(context.TODO(), facilityID)

			if err != nil {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Cannot check test on doc id  %v", facilityID)
			}
			// Avoid comparing timestamp
			for i, item := range got.Parcels {
				expFacilityParcel.Parcels[i].ArrivingTime = item.ArrivingTime
				expFacilityParcel.Parcels[i].ExitTime = item.ExitTime
			}
			if !reflect.DeepEqual(got, expFacilityParcel) {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Read %v, want %v", got, expFacilityParcel)
			}
		})
	}
}

func TestParcelProcessingFailedManager_ManageEvent_KO_FacilityParcelDocNotFound(t *testing.T) {
	facilityID := uuid.New().String()
	parcelID := uuid.New().String()
	// Initialize collections
	// --Nothing

	// Test
	event := models.ParcelProcessingFailed{
		ParcelID:         parcelID,
		Timestamp:        time.Now().UTC(),
		FacilityID:       facilityID,
		SortingMachineID: "",
		ErrMsg:           "an error msg",
	}

	type fields struct {
		logisticEvent models.ParcelProcessingFailed
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
		{name: "Event KO. Facility parcel doc not found", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelProcessingFailedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelProcessingFailedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParcelProcessingFailedManager_ManageEvent_KO_ParcelNotInFacility(t *testing.T) {
	facilityID := uuid.New().String()
	parcelID := uuid.New().String()
	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility1 := repository.ParcelInFacility{ParcelID: "fsfdsfs"}
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

	// Test
	event := models.ParcelProcessingFailed{
		ParcelID:         parcelID,
		Timestamp:        time.Now().UTC(),
		FacilityID:       facilityID,
		SortingMachineID: "",
		ErrMsg:           "an error msg",
	}

	type fields struct {
		logisticEvent models.ParcelProcessingFailed
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
		{name: "Event KO. Parcel not in facility", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelProcessingFailedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelProcessingFailedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
