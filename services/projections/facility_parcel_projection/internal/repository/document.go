package repository

import "time"

const (
	FacilityParcelCollection    = "facility_parcel"
	FacilityParcelCollectionKey = "facility_id"

	FacilityCollection    = "facility"
	FacilityCollectionKey = "facility_id"

	ParcelCollection    = "parcel"
	ParcelCollectionKey = "parcel_id"

	VehicleTransportCollection    = "vehicle_transport"
	VehicleTransportCollectionKey = "transport_id"

	RouteCollection = "route"
)

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

type Parcel struct {
	Name       string     `bson:"parcel_name"`
	ID         string     `bson:"parcel_id"`
	Type       string     `bson:"parcel_type"`
	LastStatus string     `bson:"last_status"`
	Position   Position   `bson:"position"`
	ParcelPath ParcelPath `bson:"parcel_path"`
	History    []Status   `bson:"history"`
	Sender     Sender     `bson:"sender"`
	Receiver   Receiver   `bson:"receiver"`
}

type Sender struct {
	Name    string `bson:"name"`
	Address string `bson:"address"`
	Zipcode string `bson:"zipcode"`
	City    string `bson:"city"`
	Nation  string `bson:"nation"`
}

type Receiver struct {
	Name    string `bson:"name"`
	Address string `bson:"address"`
	Zipcode string `bson:"zipcode"`
	City    string `bson:"city"`
	Nation  string `bson:"nation"`
	Number  string `bson:"number"`
	Email   string `bson:"email"`
}

type ParcelPath struct {
	Path          []string `bson:"path"`
	PathCompleted int      `bson:"path_completed"`
}

type Position struct {
	PositionType string `bson:"type"`
	PositionID   string `bson:"id"`
}

type Status struct {
	Status string    `bson:"status"`
	Date   time.Time `bson:"date"`
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

type Route struct {
	SourceFacilityID string   `bson:"source_facility_id"`
	DestFacilityID   string   `bson:"dest_facility_id"`
	Cutoff           []string `bson:"cutoff"`
}
