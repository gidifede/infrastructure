package repository

import "time"

const (
	NetworkCollection = "network"

	FacilityExpectedParcelCollection      = "facility_expected_parcel"
	FacilityExpectedParcelCollectionKey   = "facility_id"
	FacilityExpectedParcelStatusHealthy   = "healthy"
	FacilityExpectedParcelStatusUnhealthy = "unhealthy"
	FacilityExpectedParcelStatusWarning   = "warning"

	SortingMachineCollection            = "sorting_machine"
	SortingMachineCollectionFacilityKey = "facility_id"

	RouteCollection = "route"

	FacilityParcelCollection    = "facility_parcel"
	FacilityParcelCollectionKey = "facility_id"

	FacilityCollection    = "facility"
	FacilityCollectionKey = "facility_id"

	FacilityVehicleCollection      = "facility_vehicle"
	FacilityVehicleCollectionKey   = "facility_id"
	FacilityVehicleStatusLoading   = "loading"
	FacilityVehicleStatusUnloading = "unloading"
)

type FacilityExpectedParcel struct {
	FacilityID string `bson:"facility_id"`
	Status     string `bson:"status"`
	Counter    int    `bson:"counter"`
}

type SortingMachine struct {
	SortingMachineID    string               `bson:"sorting_machine_id"`
	FacilityID          string               `bson:"facility_id"`
	Capacity            int                  `bson:"capacity"`
	ItemProcessingRates []ItemProcessingRate `bson:"item_processing_rates"`
	InsertTimestamp     time.Time            `bson:"insert_timestamp"`
	NextMaintenance     time.Time            `bson:"next_maintenance"`
}

type ItemProcessingRate struct {
	Date        string       `bson:"date"`
	HourlyRates []HourlyRate `bson:"hourly_rates"`
}

type HourlyRate struct {
	Hour int `bson:"hour"`
	Rate int `bson:"rate"`
}

type Route struct {
	SourceFacilityID string   `bson:"source_facility_id"`
	DestFacilityID   string   `bson:"dest_facility_id"`
	Cutoff           []string `bson:"cutoff"`
}

type FacilityParcel struct {
	FacilityID string             `bson:"facility_id"`
	Parcels    []ParcelInFacility `bson:"parcels"`
}

type ParcelInFacility struct {
	ParcelID            string    `bson:"parcel_id"`
	ArrivingTime        time.Time `bson:"arriving_time"`
	ExitTime            time.Time `bson:"exit_time"`
	NextHop             string    `bson:"next_hop"`
	DeliveryDestination string    `bson:"delivery_destination"`
	Status              string    `bson:"status"`
}

type FacilityLocation struct {
	Address string `bson:"address"`
	Zipcode string `bson:"zipcode"`
	City    string `bson:"city"`
	Nation  string `bson:"nation"`
}

type Connection struct {
	FacilityDestinationID string `bson:"facility_destination_id"`
	Distance              int    `bson:"distance"`
}

type Facility struct {
	FacilityID       string           `bson:"facility_id"`
	Capacity         int              `bson:"capacity"`
	FacilityType     string           `bson:"facility_type"`
	FacilityLocation FacilityLocation `bson:"facility_location"`
	Latitudine       string           `bson:"latitudine"`
	Longitude        string           `bson:"longitude"`
	Company          string           `bson:"company"`
	Connections      []Connection     `bson:"connections"`
}

type VehicleInFacility struct {
	VehicleID   string    `bson:"vehicle_id"`
	Status      string    `bson:"status"`
	ArrivedTime time.Time `bson:"arrived_time"`
}

type FacilityVehicle struct {
	FacilityID string              `bson:"facility_id"`
	Vehicles   []VehicleInFacility `bson:"vehicles"`
}
