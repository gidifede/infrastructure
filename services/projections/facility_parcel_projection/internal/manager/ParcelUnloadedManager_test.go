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

func TestParcelUnloadedManager_ManageEvent_OK_FacilityParcelDocExists(t *testing.T) {
	facilityID := uuid.New().String()
	parcelID := uuid.New().String()
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

	collectionName = repository.ParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	parcel := repository.Parcel{
		ID:         parcelID,
		Name:       "parcel name",
		Type:       "parcel type",
		LastStatus: "a status",
		Position:   repository.Position{},
		ParcelPath: repository.ParcelPath{},
		History:    []repository.Status{},
		Sender:     repository.Sender{},
		Receiver:   repository.Receiver{Address: "Receiver's address"},
	}
	documents = []interface{}{parcel}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.FacilityParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility := repository.ParcelInFacility{ParcelID: "fsfdsfs"}
	documents = []interface{}{
		repository.FacilityParcel{
			FacilityID: facilityID,
			Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility},
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.ParcelUnloaded{
		FacilityID: facilityID,
		ParcelID:   parcelID,
		Timestamp:  time.Now().UTC(),
	}
	expFacilityParcel := repository.FacilityParcel{
		FacilityID: facilityID,
		Parcels: []repository.ParcelInFacility{
			parcelAlreadyInFacility,
			{
				ParcelID:            parcelID,
				ArrivingTime:        event.Timestamp,
				DeliveryDestination: parcel.Receiver.Address,
				Status:              "sorting",
			}},
	}

	type fields struct {
		logisticEvent models.ParcelUnloaded
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
		{name: "Event OK. Facility parcel doc exists", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelUnloadedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelUnloadedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
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

func TestParcelUnloadedManager_ManageEvent_OK_FacilityParcelDocDoesntExist(t *testing.T) {
	facilityID := uuid.New().String()
	parcelID := uuid.New().String()
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

	collectionName = repository.ParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	parcel := repository.Parcel{
		ID:         parcelID,
		Name:       "parcel name",
		Type:       "parcel type",
		LastStatus: "a status",
		Position:   repository.Position{},
		ParcelPath: repository.ParcelPath{},
		History:    []repository.Status{},
		Sender:     repository.Sender{},
		Receiver:   repository.Receiver{Address: "Receiver's address"},
	}
	documents = []interface{}{parcel}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.ParcelUnloaded{
		FacilityID: facilityID,
		ParcelID:   parcelID,
		Timestamp:  time.Now().UTC(),
	}
	expFacilityParcel := repository.FacilityParcel{
		FacilityID: facilityID,
		Parcels: []repository.ParcelInFacility{
			{
				ParcelID:            parcelID,
				ArrivingTime:        event.Timestamp,
				DeliveryDestination: parcel.Receiver.Address,
				Status:              "sorting",
			}},
	}

	type fields struct {
		logisticEvent models.ParcelUnloaded
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
		{name: "Event OK. Facility parcel doc doesn't exist", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelUnloadedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelUnloadedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
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

func TestParcelUnloadedManager_ManageEvent_KO_ParcelDoesntExist(t *testing.T) {
	facilityID := uuid.New().String()
	parcelID := uuid.New().String()
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
	event := models.ParcelUnloaded{
		FacilityID: facilityID,
		ParcelID:   parcelID,
		Timestamp:  time.Now().UTC(),
	}

	type fields struct {
		logisticEvent models.ParcelUnloaded
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
			a := &ParcelUnloadedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelUnloadedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParcelUnloadedManager_ManageEvent_KO_FacilityDoesntExist(t *testing.T) {
	facilityID := uuid.New().String()
	parcelID := uuid.New().String()
	// Initialize collections
	collectionName := repository.ParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcel := repository.Parcel{
		ID:         parcelID,
		Name:       "parcel name",
		Type:       "parcel type",
		LastStatus: "a status",
		Position:   repository.Position{},
		ParcelPath: repository.ParcelPath{},
		History:    []repository.Status{},
		Sender:     repository.Sender{},
		Receiver:   repository.Receiver{Address: "Receiver's address"},
	}
	documents := []interface{}{parcel}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.ParcelUnloaded{
		FacilityID: facilityID,
		ParcelID:   parcelID,
		Timestamp:  time.Now().UTC(),
	}

	type fields struct {
		logisticEvent models.ParcelUnloaded
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
			a := &ParcelUnloadedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelUnloadedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParcelUnloadedManager_ManageEvent_KO_ParcelInFacilityDocAlreadyExists(t *testing.T) {
	facilityID := uuid.New().String()
	parcelID := uuid.New().String()
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

	collectionName = repository.ParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	parcel := repository.Parcel{
		ID:         parcelID,
		Name:       "parcel name",
		Type:       "parcel type",
		LastStatus: "a status",
		Position:   repository.Position{},
		ParcelPath: repository.ParcelPath{},
		History:    []repository.Status{},
		Sender:     repository.Sender{},
		Receiver:   repository.Receiver{Address: "Receiver's address"},
	}
	documents = []interface{}{parcel}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.FacilityParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	parcelAlreadyInFacility := repository.ParcelInFacility{ParcelID: parcelID}
	documents = []interface{}{
		repository.FacilityParcel{
			FacilityID: facilityID,
			Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility},
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	event := models.ParcelUnloaded{
		FacilityID: facilityID,
		ParcelID:   parcelID,
		Timestamp:  time.Now().UTC(),
	}
	expFacilityParcel := repository.FacilityParcel{
		FacilityID: facilityID,
		Parcels:    []repository.ParcelInFacility{parcelAlreadyInFacility},
	}

	type fields struct {
		logisticEvent models.ParcelUnloaded
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
		{name: "Event KO. Parcel in facility doc already exists", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelUnloadedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelUnloadedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
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
