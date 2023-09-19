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

func TestProcessedManager_ManageEvent_OK_FacilityParcelDocUpdated(t *testing.T) {
	sourceFacilityID := uuid.New().String()
	destFacilityID := uuid.New().String()
	parcelID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility := repository.ParcelInFacility{
		ParcelID:            parcelID,
		Status:              "processing",
		ArrivingTime:        time.Now().UTC(),
		DeliveryDestination: "a destination",
	}
	documents := []interface{}{
		repository.FacilityParcel{
			FacilityID: sourceFacilityID,
			Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility},
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.RouteCollection
	collection = client.Database("logistic").Collection(collectionName)
	route := repository.Route{
		SourceFacilityID: sourceFacilityID,
		DestFacilityID:   destFacilityID,
		Cutoff:           []string{"10:30", "11:30", "17:00"},
	}
	documents = []interface{}{route}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	timestamp := time.Now().UTC()
	timestamp = timestamp.Truncate(24 * time.Hour).Add(15 * time.Hour)
	event := models.ParcelProcessed{
		ParcelID:              parcelID,
		Timestamp:             timestamp,
		DestinationFacilityID: destFacilityID,
		FacilityID:            sourceFacilityID,
		SortingMachineID:      "gafdhgshjg",
	}

	exitTime := timestamp.Truncate(24 * time.Hour).Add(17 * time.Hour)
	expFacilityParcel := repository.FacilityParcel{
		FacilityID: sourceFacilityID,
		Parcels: []repository.ParcelInFacility{
			{
				ParcelID:            parcelID,
				ArrivingTime:        parcelAlreadyInFacility.ArrivingTime,
				ExitTime:            exitTime,
				NextHop:             destFacilityID,
				DeliveryDestination: parcelAlreadyInFacility.DeliveryDestination,
				Status:              "processed",
			},
		},
	}

	type fields struct {
		logisticEvent models.ParcelProcessed
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
			a := &ProcessedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ProcessedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetFacilityParcelDocumentByID(context.TODO(), sourceFacilityID)

			if err != nil {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Cannot check test on doc id  %v", sourceFacilityID)
			}
			// Compare exit time
			if expFacilityParcel.Parcels[0].ExitTime.Hour() != got.Parcels[0].ExitTime.Hour() {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Read exit_time %v, want exit_time %v", got.Parcels[0].ExitTime, expFacilityParcel.Parcels[0].ExitTime)
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

func TestProcessedManager_ManageEvent_OK_FacilityParcelDocUpdatedCutoffNextDay(t *testing.T) {
	sourceFacilityID := uuid.New().String()
	destFacilityID := uuid.New().String()
	parcelID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility := repository.ParcelInFacility{
		ParcelID:            parcelID,
		Status:              "processing",
		ArrivingTime:        time.Now().UTC(),
		DeliveryDestination: "a destination",
	}
	documents := []interface{}{
		repository.FacilityParcel{
			FacilityID: sourceFacilityID,
			Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility},
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.RouteCollection
	collection = client.Database("logistic").Collection(collectionName)
	route := repository.Route{
		SourceFacilityID: sourceFacilityID,
		DestFacilityID:   destFacilityID,
		Cutoff:           []string{"10:30", "11:30", "17:00"},
	}
	documents = []interface{}{route}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	timestamp := time.Now().UTC()
	timestamp = timestamp.Truncate(24 * time.Hour).Add(19 * time.Hour)
	event := models.ParcelProcessed{
		ParcelID:              parcelID,
		Timestamp:             timestamp,
		DestinationFacilityID: destFacilityID,
		FacilityID:            sourceFacilityID,
		SortingMachineID:      "gafdhgshjg",
	}

	exitTime := timestamp.Truncate(24 * time.Hour).Add(24 * time.Hour).Add(10 * time.Hour)
	expFacilityParcel := repository.FacilityParcel{
		FacilityID: sourceFacilityID,
		Parcels: []repository.ParcelInFacility{
			{
				ParcelID:            parcelID,
				ArrivingTime:        parcelAlreadyInFacility.ArrivingTime,
				ExitTime:            exitTime,
				NextHop:             destFacilityID,
				DeliveryDestination: parcelAlreadyInFacility.DeliveryDestination,
				Status:              "processed",
			},
		},
	}

	type fields struct {
		logisticEvent models.ParcelProcessed
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
		{name: "Event OK. Facility parcel doc updated, cut off next day", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ProcessedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ProcessedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetFacilityParcelDocumentByID(context.TODO(), sourceFacilityID)

			if err != nil {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Cannot check test on doc id  %v", sourceFacilityID)
			}
			// Compare exit time
			if expFacilityParcel.Parcels[0].ExitTime.Hour() != got.Parcels[0].ExitTime.Hour() || expFacilityParcel.Parcels[0].ExitTime.Day() != got.Parcels[0].ExitTime.Day() {
				t.Errorf("ParcelUnloadedManager.ManageEvent(). Read exit_time %v, want exit_time %v", got.Parcels[0].ExitTime, expFacilityParcel.Parcels[0].ExitTime)
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

func TestProcessedManager_ManageEvent_K0_RouteDoesntExist(t *testing.T) {
	sourceFacilityID := uuid.New().String()
	destFacilityID := uuid.New().String()
	parcelID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility := repository.ParcelInFacility{
		ParcelID:            parcelID,
		Status:              "processing",
		ArrivingTime:        time.Now().UTC(),
		DeliveryDestination: "a destination",
	}
	documents := []interface{}{
		repository.FacilityParcel{
			FacilityID: sourceFacilityID,
			Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility},
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	timestamp := time.Now().UTC()
	timestamp = timestamp.Truncate(24 * time.Hour).Add(15 * time.Hour)
	event := models.ParcelProcessed{
		ParcelID:              parcelID,
		Timestamp:             timestamp,
		DestinationFacilityID: destFacilityID,
		FacilityID:            sourceFacilityID,
		SortingMachineID:      "gafdhgshjg",
	}

	type fields struct {
		logisticEvent models.ParcelProcessed
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
		{name: "Event K0. Route doesn't exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ProcessedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ProcessedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessedManager_ManageEvent_KO_FacilityParcelDocDoesntExist(t *testing.T) {
	sourceFacilityID := uuid.New().String()
	destFacilityID := uuid.New().String()
	parcelID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility := repository.ParcelInFacility{
		ParcelID:            parcelID,
		Status:              "processing",
		ArrivingTime:        time.Now().UTC(),
		DeliveryDestination: "a destination",
	}
	documents := []interface{}{
		repository.FacilityParcel{
			FacilityID: "dfafadfad",
			Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility},
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.RouteCollection
	collection = client.Database("logistic").Collection(collectionName)
	route := repository.Route{
		SourceFacilityID: sourceFacilityID,
		DestFacilityID:   destFacilityID,
		Cutoff:           []string{"10:30", "11:30", "17:00"},
	}
	documents = []interface{}{route}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	timestamp := time.Now().UTC()
	timestamp = timestamp.Truncate(24 * time.Hour).Add(15 * time.Hour)
	event := models.ParcelProcessed{
		ParcelID:              parcelID,
		Timestamp:             timestamp,
		DestinationFacilityID: destFacilityID,
		FacilityID:            sourceFacilityID,
		SortingMachineID:      "gafdhgshjg",
	}

	type fields struct {
		logisticEvent models.ParcelProcessed
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
		{name: "Event KO. Facility parcel doc doesn't exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ProcessedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ProcessedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessedManager_ManageEvent_KO_ParcelDoesntExist(t *testing.T) {
	sourceFacilityID := uuid.New().String()
	destFacilityID := uuid.New().String()
	parcelID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility := repository.ParcelInFacility{
		ParcelID:            "sdgfsgsf",
		Status:              "processing",
		ArrivingTime:        time.Now().UTC(),
		DeliveryDestination: "a destination",
	}
	documents := []interface{}{
		repository.FacilityParcel{
			FacilityID: sourceFacilityID,
			Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility},
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.RouteCollection
	collection = client.Database("logistic").Collection(collectionName)
	route := repository.Route{
		SourceFacilityID: sourceFacilityID,
		DestFacilityID:   destFacilityID,
		Cutoff:           []string{"10:30", "11:30", "17:00"},
	}
	documents = []interface{}{route}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	timestamp := time.Now().UTC()
	timestamp = timestamp.Truncate(24 * time.Hour).Add(15 * time.Hour)
	event := models.ParcelProcessed{
		ParcelID:              parcelID,
		Timestamp:             timestamp,
		DestinationFacilityID: destFacilityID,
		FacilityID:            sourceFacilityID,
		SortingMachineID:      "gafdhgshjg",
	}

	type fields struct {
		logisticEvent models.ParcelProcessed
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
		{name: "Event KO. Parcel doesn't exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ProcessedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ProcessedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
