package repository

import "time"

const (
	FacilityExpectedParcelCollection      = "facility_expected_parcel"
	FacilityExpectedParcelCollectionKey   = "facility_id"
	FacilityExpectedParcelStatusHealthy   = "healthy"
	FacilityExpectedParcelStatusUnhealthy = "unhealthy"
	FacilityExpectedParcelStatusWarning   = "warning"

	FacilityCollection    = "facility"
	FacilityCollectionKey = "facility_id"

	VehicleTransportCollection    = "vehicle_transport"
	VehicleTransportCollectionKey = "transport_id"
)

type FacilityExpectedParcel struct {
	FacilityID string `bson:"facility_id"`
	Status     string `bson:"status"`
	Counter    int    `bson:"counter"`
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

type VehicleTransportHistory []struct {
	Status   string `bson:"status"`
	Data     string `bson:"data"`
	Edificio string `bson:"edificio,omitempty"`
}

type VehicleTransportPosition []struct {
	Latitudine  float64   `bson:"latitudine"`
	Longitudine float64   `bson:"longitudine"`
	Timestamp   time.Time `bson:"timestamp"`
}

type VehicleTransport struct {
	TransportID    string                   `bson:"transport_id"`
	VehicleID      string                   `bson:"vehicle_id"`
	Parcels        []string                 `bson:"parcels"`
	History        VehicleTransportHistory  `bson:"history"`
	Position       VehicleTransportPosition `bson:"position"`
	StartTimestamp time.Time                `bson:"start_timestamp"`
	IsActive       bool                     `bson:"is_active"`
}
