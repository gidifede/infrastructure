package service

import (
	"context"
	"facility_queryhandler/internal"
	"facility_queryhandler/internal/models"
	"facility_queryhandler/internal/repository"
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
	internal.Repo = repository.NewMongoDB(*client.Database("logistic"))
}

func TestGetFacilityShortStats_OK(t *testing.T) {
	facilityID := uuid.New().String()
	timestampNow := time.Now().UTC()
	layout := "2006-01-02"
	currentDate := timestampNow.Format(layout)
	temp := time.Unix(timestampNow.Add(5*time.Minute).Unix(), 0)
	nextCutoffTime := temp.Format("15:04")
	numExpectedParcels := 50

	// Initialize collections
	collectionName := repository.FacilityExpectedParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.FacilityExpectedParcel{
			FacilityID: facilityID,
			Status:     repository.FacilityExpectedParcelStatusHealthy,
			Counter:    numExpectedParcels,
		},
		repository.FacilityExpectedParcel{
			FacilityID: "fadfadfdagdag",
			Status:     repository.FacilityExpectedParcelStatusHealthy,
			Counter:    0,
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.SortingMachineCollection
	collection = client.Database("logistic").Collection(collectionName)
	hourlyRate1 := repository.HourlyRate{Hour: 10, Rate: 5}
	hourlyRate2 := repository.HourlyRate{Hour: 12, Rate: 30}
	documents = []interface{}{
		repository.SortingMachine{
			SortingMachineID: "afdagfdaga",
			FacilityID:       facilityID,
			ItemProcessingRates: []repository.ItemProcessingRate{
				{
					Date: "2020-09-20",
					HourlyRates: []repository.HourlyRate{
						{Hour: 10, Rate: 5},
						{Hour: 12, Rate: 8},
					},
				},
				{
					Date: currentDate,
					HourlyRates: []repository.HourlyRate{
						hourlyRate1,
						hourlyRate2,
					},
				},
			},
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.RouteCollection
	collection = client.Database("logistic").Collection(collectionName)
	documents = []interface{}{
		repository.Route{
			SourceFacilityID: facilityID,
			DestFacilityID:   "fafdafsa",
			Cutoff:           []string{"10:30", nextCutoffTime, "17:00"},
		},
		repository.Route{
			SourceFacilityID: facilityID,
			DestFacilityID:   "flalfhol",
			Cutoff:           []string{"10:30", "11:30", "17:00"},
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	expFacilityShortStats := models.FacilityShortStats{
		FacilityHealth:  repository.FacilityExpectedParcelStatusHealthy,
		ParcelWaiting:   numExpectedParcels,
		ParcelProcessed: hourlyRate1.Rate + hourlyRate2.Rate,
		NextCutOffTime:  nextCutoffTime,
	}
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    models.FacilityShortStats
		wantErr bool
	}{
		{name: "OK", args: args{ctx: context.TODO(), facilityID: facilityID}, want: expFacilityShortStats, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFacilityShortStats(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFacilityShortStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFacilityShortStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFacilityParcelDetails_OK(t *testing.T) {
	facilityID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityParcelCollection
	collection := client.Database("logistic").Collection(collectionName)
	parcel1 := repository.ParcelInFacility{ParcelID: "fddfafd", NextHop: "afdad", DeliveryDestination: "fafdad", Status: "fafda"}
	parcel2 := repository.ParcelInFacility{ParcelID: "dsgskgs", NextHop: "453", DeliveryDestination: "3553", Status: "35535"}
	documents := []interface{}{
		repository.FacilityParcel{
			FacilityID: facilityID,
			Parcels:    []repository.ParcelInFacility{parcel1, parcel2},
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	expFacilityParcelDetails := []models.FacilityParcelDetails{
		{ParcelID: parcel1.ParcelID, NextHop: parcel1.NextHop, Status: parcel1.Status, Destination: parcel1.DeliveryDestination},
		{ParcelID: parcel2.ParcelID, NextHop: parcel2.NextHop, Status: parcel2.Status, Destination: parcel2.DeliveryDestination},
	}
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.FacilityParcelDetails
		wantErr bool
	}{
		{name: "OK", args: args{ctx: context.TODO(), facilityID: facilityID}, want: expFacilityParcelDetails, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFacilityParcelDetails(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFacilityParcelDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFacilityParcelDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFacilityParcelStats_OK(t *testing.T) {
	facilityID := uuid.New().String()
	facilityCapacity := 1000
	numExpectedParcels := 100
	timestampNow := time.Now().UTC()
	layout := "2006-01-02"
	currentDate := timestampNow.Format(layout)

	// Initialize collections
	collectionName := repository.FacilityCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.Facility{
			FacilityID: facilityID,
			Capacity:   facilityCapacity,
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
	documents = []interface{}{
		repository.FacilityExpectedParcel{
			FacilityID: facilityID,
			Status:     repository.FacilityExpectedParcelStatusHealthy,
			Counter:    numExpectedParcels,
		},
		repository.FacilityExpectedParcel{
			FacilityID: "fadfadfdagdag",
			Status:     repository.FacilityExpectedParcelStatusHealthy,
			Counter:    0,
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.SortingMachineCollection
	collection = client.Database("logistic").Collection(collectionName)
	hourlyRate1 := repository.HourlyRate{Hour: 10, Rate: 5}
	hourlyRate2 := repository.HourlyRate{Hour: 12, Rate: 19}
	documents = []interface{}{
		repository.SortingMachine{
			SortingMachineID: "afdagfdaga",
			FacilityID:       facilityID,
			ItemProcessingRates: []repository.ItemProcessingRate{
				{
					Date: "2020-09-20",
					HourlyRates: []repository.HourlyRate{
						{Hour: 10, Rate: 5},
						{Hour: 12, Rate: 8},
					},
				},
				{
					Date: currentDate,
					HourlyRates: []repository.HourlyRate{
						hourlyRate1,
						hourlyRate2,
					},
				},
			},
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	collectionName = repository.FacilityParcelCollection
	collection = client.Database("logistic").Collection(collectionName)
	parcel1 := repository.ParcelInFacility{ParcelID: "fddfafd", ArrivingTime: timestampNow, ExitTime: timestampNow.Add(time.Duration(5 * time.Second)), NextHop: "afdad", DeliveryDestination: "fafdad", Status: "fafda"}
	parcel2 := repository.ParcelInFacility{ParcelID: "dsgskgs", ArrivingTime: timestampNow, ExitTime: timestampNow.Add(time.Duration(2 * time.Second)), NextHop: "453", DeliveryDestination: "3553", Status: "35535"}
	documents = []interface{}{
		repository.FacilityParcel{
			FacilityID: facilityID,
			Parcels:    []repository.ParcelInFacility{parcel1, parcel2},
		},
	}
	_, err = collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	expFacilityParcelStats := models.FacilityParcelStats{
		Capacity:          facilityCapacity,
		ParcelWaiting:     numExpectedParcels,
		ParcelProcessed:   hourlyRate1.Rate + hourlyRate2.Rate,
		AvgProcessingTime: (parcel1.ExitTime.Sub(parcel1.ArrivingTime).Seconds() + parcel2.ExitTime.Sub(parcel2.ArrivingTime).Seconds()) / 2,
	}
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    models.FacilityParcelStats
		wantErr bool
	}{
		{name: "OK", args: args{ctx: context.TODO(), facilityID: facilityID}, want: expFacilityParcelStats, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFacilityParcelStats(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFacilityParcelStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFacilityParcelStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFacilitySortingMachineStats_OK(t *testing.T) {
	facilityID := uuid.New().String()
	capacity := 100

	// Initialize collections
	collectionName := repository.SortingMachineCollection
	collection := client.Database("logistic").Collection(collectionName)
	timestampNow := time.Now().UTC()
	timestampNowCopy := timestampNow
	layout := "2006-01-02"
	rate1 := repository.HourlyRate{Hour: 10, Rate: 5}
	rate2 := repository.HourlyRate{Hour: 12, Rate: 8}
	itemProcessingRates := []repository.ItemProcessingRate{
		{
			Date: "2020-09-20",
			HourlyRates: []repository.HourlyRate{
				rate1,
				rate2,
			},
		},
		{
			Date:        timestampNow.Format(layout),
			HourlyRates: []repository.HourlyRate{},
		},
		{
			Date:        timestampNow.Add(-time.Duration(24 * time.Hour)).Format(layout),
			HourlyRates: []repository.HourlyRate{},
		},
	}
	for i := 0; i < 5; i++ {
		currentDate := timestampNow.Format(layout)
		currentTime := timestampNow.Hour()
		for i := range itemProcessingRates {
			if itemProcessingRates[i].Date == currentDate {
				itemProcessingRates[i].HourlyRates = append(itemProcessingRates[i].HourlyRates, repository.HourlyRate{Hour: currentTime, Rate: 3})
			}
		}
		timestampNow = timestampNow.Add(-time.Duration(2 * time.Hour))
	}
	documents := []interface{}{
		repository.SortingMachine{
			SortingMachineID:    uuid.New().String(),
			FacilityID:          facilityID,
			ItemProcessingRates: itemProcessingRates,
			Capacity:            capacity,
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	expFacilitySMStats := models.FacilitySortingMachineStats{
		Capacity:               capacity,
		WorkingCapacityAverage: float64(rate1.Rate+rate2.Rate+15) / 3,
		ParcelProcessed: []models.ParcelProcessedItem{
			{Day: timestampNowCopy.Format(layout), Hour: timestampNowCopy.Hour(), Parcels: 3},
			{Day: timestampNowCopy.Add(-time.Duration(time.Hour)).Format(layout), Hour: timestampNowCopy.Add(-time.Duration(time.Hour)).Hour(), Parcels: 0},
			{Day: timestampNowCopy.Add(-time.Duration(2 * time.Hour)).Format(layout), Hour: timestampNowCopy.Add(-time.Duration(2 * time.Hour)).Hour(), Parcels: 3},
			{Day: timestampNowCopy.Add(-time.Duration(3 * time.Hour)).Format(layout), Hour: timestampNowCopy.Add(-time.Duration(3 * time.Hour)).Hour(), Parcels: 0},
			{Day: timestampNowCopy.Add(-time.Duration(4 * time.Hour)).Format(layout), Hour: timestampNowCopy.Add(-time.Duration(4 * time.Hour)).Hour(), Parcels: 3},
			{Day: timestampNowCopy.Add(-time.Duration(5 * time.Hour)).Format(layout), Hour: timestampNowCopy.Add(-time.Duration(5 * time.Hour)).Hour(), Parcels: 0},
			{Day: timestampNowCopy.Add(-time.Duration(6 * time.Hour)).Format(layout), Hour: timestampNowCopy.Add(-time.Duration(6 * time.Hour)).Hour(), Parcels: 3},
			{Day: timestampNowCopy.Add(-time.Duration(7 * time.Hour)).Format(layout), Hour: timestampNowCopy.Add(-time.Duration(7 * time.Hour)).Hour(), Parcels: 0},
			{Day: timestampNowCopy.Add(-time.Duration(8 * time.Hour)).Format(layout), Hour: timestampNowCopy.Add(-time.Duration(8 * time.Hour)).Hour(), Parcels: 3},
			{Day: timestampNowCopy.Add(-time.Duration(9 * time.Hour)).Format(layout), Hour: timestampNowCopy.Add(-time.Duration(9 * time.Hour)).Hour(), Parcels: 0},
		},
	}
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    models.FacilitySortingMachineStats
		wantErr bool
	}{
		{name: "OK", args: args{ctx: context.TODO(), facilityID: facilityID}, want: expFacilitySMStats, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFacilitySortingMachineStats(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFacilitySortingMachineStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFacilitySortingMachineStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFacilityVehicleStats(t *testing.T) {
	facilityID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityVehicleCollection
	collection := client.Database("logistic").Collection(collectionName)
	documents := []interface{}{
		repository.FacilityVehicle{
			FacilityID: facilityID,
			Vehicles: []repository.VehicleInFacility{
				{VehicleID: uuid.New().String(), Status: repository.FacilityVehicleStatusLoading},
				{VehicleID: uuid.New().String(), Status: repository.FacilityVehicleStatusLoading},
				{VehicleID: uuid.New().String(), Status: repository.FacilityVehicleStatusUnloading},
				{VehicleID: uuid.New().String(), Status: repository.FacilityVehicleStatusLoading},
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
	expFacilityVehicleStats := models.FacilityVehicleStats{
		VehiclesUnloading: 1,
		VehiclesLoading:   3,
	}
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    models.FacilityVehicleStats
		wantErr bool
	}{
		{name: "OK", args: args{ctx: context.TODO(), facilityID: facilityID}, want: expFacilityVehicleStats, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFacilityVehicleStats(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFacilityVehicleStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFacilityVehicleStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFacilityVehicleDetails(t *testing.T) {
	facilityID := uuid.New().String()

	// Initialize collections
	collectionName := repository.FacilityVehicleCollection
	collection := client.Database("logistic").Collection(collectionName)
	v1 := repository.VehicleInFacility{VehicleID: uuid.New().String(), Status: repository.FacilityVehicleStatusLoading}
	v2 := repository.VehicleInFacility{VehicleID: uuid.New().String(), Status: repository.FacilityVehicleStatusUnloading}
	documents := []interface{}{
		repository.FacilityVehicle{
			FacilityID: facilityID,
			Vehicles:   []repository.VehicleInFacility{v1, v2},
		},
	}
	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("Collection '%s' initialized.", collectionName)
	}

	// Test
	expFacilityVehicleDetails := []models.FacilityVehicleDetails{
		{
			VehicleLicencePlate: v1.VehicleID,
			Status:              v1.Status,
		},
		{
			VehicleLicencePlate: v2.VehicleID,
			Status:              v2.Status,
		},
	}
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.FacilityVehicleDetails
		wantErr bool
	}{
		{name: "OK", args: args{ctx: context.TODO(), facilityID: facilityID}, want: expFacilityVehicleDetails, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFacilityVehicleDetails(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFacilityVehicleDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFacilityVehicleDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}
