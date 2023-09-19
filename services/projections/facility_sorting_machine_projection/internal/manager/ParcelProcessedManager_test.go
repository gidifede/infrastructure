package manager

import (
	"context"
	"facility-sorting-machine-projection/internal"
	"facility-sorting-machine-projection/internal/models"
	"facility-sorting-machine-projection/internal/repository"
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

func TestParcelProcessedManager_ManageEvent_OK_SortingMachineUpdateItemProcessingRateDateFoundHourFound(t *testing.T) {
	sortingMachineID := uuid.New().String()

	// Initialize collections
	collectionName := repository.SortingMachineCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.SortingMachine{
			SortingMachineID: sortingMachineID,
			ItemProcessingRates: []repository.ItemProcessingRate{
				{
					Date: "2023-09-20",
					HourlyRates: []repository.HourlyRate{
						{Hour: 10, Rate: 5},
						{Hour: 12, Rate: 8},
					},
				},
				{
					Date: "2023-09-21",
					HourlyRates: []repository.HourlyRate{
						{Hour: 10, Rate: 5},
						{Hour: 12, Rate: 8},
					},
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
	eventTimestamp, _ := time.Parse(time.RFC3339, "2023-09-20T10:10:14Z")
	event := models.ParcelProcessed{
		Timestamp:        eventTimestamp,
		SortingMachineID: sortingMachineID,
	}

	expSortingMachine := repository.SortingMachine{
		SortingMachineID: sortingMachineID,
		ItemProcessingRates: []repository.ItemProcessingRate{
			{
				Date: "2023-09-20",
				HourlyRates: []repository.HourlyRate{
					{Hour: 10, Rate: 6},
					{Hour: 12, Rate: 8},
				},
			},
			{
				Date: "2023-09-21",
				HourlyRates: []repository.HourlyRate{
					{Hour: 10, Rate: 5},
					{Hour: 12, Rate: 8},
				},
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
		{name: "Event OK. Sorting Machine Update Item Processing Rate Date Found Hour Found", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelProcessedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelProcessedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetSortingMachineDocumentByID(context.TODO(), sortingMachineID)

			if err != nil {
				t.Errorf("ParcelProcessedManager.ManageEvent(). Cannot check test on doc id  %v", sortingMachineID)
			}
			if !reflect.DeepEqual(got, expSortingMachine) {
				t.Errorf("ParcelProcessedManager.ManageEvent(). Read %v, want %v", got, expSortingMachine)
			}
		})
	}
}

func TestParcelProcessedManager_ManageEvent_OK_SortingMachineUpdateItemProcessingRateDateFoundHourNotFound(t *testing.T) {
	sortingMachineID := uuid.New().String()

	// Initialize collections
	collectionName := repository.SortingMachineCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.SortingMachine{
			SortingMachineID: sortingMachineID,
			ItemProcessingRates: []repository.ItemProcessingRate{
				{
					Date: "2023-09-20",
					HourlyRates: []repository.HourlyRate{
						{Hour: 12, Rate: 8},
					},
				},
				{
					Date: "2023-09-21",
					HourlyRates: []repository.HourlyRate{
						{Hour: 10, Rate: 5},
						{Hour: 12, Rate: 8},
					},
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
	eventTimestamp, _ := time.Parse(time.RFC3339, "2023-09-20T10:10:14Z")
	event := models.ParcelProcessed{
		Timestamp:        eventTimestamp,
		SortingMachineID: sortingMachineID,
	}

	expSortingMachine := repository.SortingMachine{
		SortingMachineID: sortingMachineID,
		ItemProcessingRates: []repository.ItemProcessingRate{
			{
				Date: "2023-09-20",
				HourlyRates: []repository.HourlyRate{
					{Hour: 12, Rate: 8},
					{Hour: 10, Rate: 1},
				},
			},
			{
				Date: "2023-09-21",
				HourlyRates: []repository.HourlyRate{
					{Hour: 10, Rate: 5},
					{Hour: 12, Rate: 8},
				},
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
		{name: "Event OK. Sorting Machine Update Item Processing Rate Date Found Hour Not Found", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelProcessedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelProcessedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetSortingMachineDocumentByID(context.TODO(), sortingMachineID)

			if err != nil {
				t.Errorf("ParcelProcessedManager.ManageEvent(). Cannot check test on doc id  %v", sortingMachineID)
			}
			if !reflect.DeepEqual(got, expSortingMachine) {
				t.Errorf("ParcelProcessedManager.ManageEvent(). Read %v, want %v", got, expSortingMachine)
			}
		})
	}
}

func TestParcelProcessedManager_ManageEvent_OK_SortingMachineUpdateItemProcessingRateDateNotFound(t *testing.T) {
	sortingMachineID := uuid.New().String()

	// Initialize collections
	collectionName := repository.SortingMachineCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.SortingMachine{
			SortingMachineID: sortingMachineID,
			ItemProcessingRates: []repository.ItemProcessingRate{
				{
					Date: "2023-09-21",
					HourlyRates: []repository.HourlyRate{
						{Hour: 10, Rate: 5},
						{Hour: 12, Rate: 8},
					},
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
	eventTimestamp, _ := time.Parse(time.RFC3339, "2023-09-20T10:10:14Z")
	event := models.ParcelProcessed{
		Timestamp:        eventTimestamp,
		SortingMachineID: sortingMachineID,
	}

	expSortingMachine := repository.SortingMachine{
		SortingMachineID: sortingMachineID,
		ItemProcessingRates: []repository.ItemProcessingRate{
			{
				Date: "2023-09-21",
				HourlyRates: []repository.HourlyRate{
					{Hour: 10, Rate: 5},
					{Hour: 12, Rate: 8},
				},
			},
			{
				Date: "2023-09-20",
				HourlyRates: []repository.HourlyRate{
					{Hour: 10, Rate: 1},
				},
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
		{name: "Event OK. Sorting Machine Update Item Processing Rate Date Not Found", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelProcessedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelProcessedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := internal.Repo.GetSortingMachineDocumentByID(context.TODO(), sortingMachineID)

			if err != nil {
				t.Errorf("ParcelProcessedManager.ManageEvent(). Cannot check test on doc id  %v", sortingMachineID)
			}
			if !reflect.DeepEqual(got, expSortingMachine) {
				t.Errorf("ParcelProcessedManager.ManageEvent(). Read %v, want %v", got, expSortingMachine)
			}
		})
	}
}

func TestParcelProcessedManager_ManageEvent_KO_SortingMachineNotFound(t *testing.T) {
	sortingMachineID := uuid.New().String()

	// Initialize collections
	collectionName := repository.SortingMachineCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.SortingMachine{
			SortingMachineID: "dkfjdas",
			ItemProcessingRates: []repository.ItemProcessingRate{
				{
					Date: "2023-09-20",
					HourlyRates: []repository.HourlyRate{
						{Hour: 10, Rate: 5},
						{Hour: 12, Rate: 8},
					},
				},
				{
					Date: "2023-09-21",
					HourlyRates: []repository.HourlyRate{
						{Hour: 10, Rate: 5},
						{Hour: 12, Rate: 8},
					},
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
	eventTimestamp, _ := time.Parse(time.RFC3339, "2023-09-20T10:10:14Z")
	event := models.ParcelProcessed{
		Timestamp:        eventTimestamp,
		SortingMachineID: sortingMachineID,
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
		{name: "Event KO. Sorting Machine Not Found", args: args{ctx: context.TODO()}, fields: fields{logisticEvent: event}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ParcelProcessedManager{
				logisticEvent: tt.fields.logisticEvent,
			}
			if err := a.ManageEvent(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ParcelProcessedManager.ManageEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
