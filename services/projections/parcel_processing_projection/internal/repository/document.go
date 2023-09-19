package repository

import (
	"time"
)

const (
	ParcelEnum                     = "Parcel"
	ProcessedEnum                  = "Processed"
	FailProcessedEnum              = "FailProcessed"
	LoadedEnum                     = "Loaded"
	UnloadedEnum                   = "Unloaded"
	TransportEndedEnum             = "TransportEnded"
	TransportStartedEnum           = "TransportStarted"
	DeliveryCompletedEnum          = "DeliveryCompleted"
	RouteEnum                      = "Route"
	FacilityEnum                   = "Facility"
	RouteCollectionEnum            = "routes"
	ParcelCollectionEnum           = "parcel"
	FacilityCollectionEnum         = "facility"
	VehicleTransportCollectionEnum = "vehicle_transport"
	ParcelCollectionIndexEnum      = "parcel_id"
)

type Parcel struct {
	Name       string     `bson:"parcel_name"`
	ID         string     `bson:"parcel_id"`
	Type       string     `bson:"parcel_type"`
	LastStatus string     `bson:"last_status"`
	Position   Position   `bson:"position"`
	ParcelPath ParcelPath `bson:"parcel_path"`
	History    []Status   `bson:"history"`
	Sender     struct {
		Name    string `bson:"name"`
		Address string `bson:"address"`
		Zipcode string `bson:"zipcode"`
		City    string `bson:"city"`
		Nation  string `bson:"nation"`
	} `bson:"sender"`
	Receiver struct {
		Name    string `bson:"name"`
		Address string `bson:"address"`
		Zipcode string `bson:"zipcode"`
		City    string `bson:"city"`
		Nation  string `bson:"nation"`
		Number  string `bson:"number"`
		Email   string `bson:"email"`
	} `bson:"receiver"`
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

type Route struct {
	SourceFacilityID string   `bson:"source_facility_id"`
	DestFacilityID   string   `bson:"dest_facility_id"`
	Cutoff           []string `bson:"cutoff"`
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

type VehicleTransport struct {
	TransportID string   `bson:"transport_id"`
	VehicleID   string   `bson:"vehicle_id"`
	Parcels     []string `bson:"parcels"`
	History     []struct {
		Status   string    `bson:"status"`
		Data     time.Time `bson:"data"`
		Edificio string    `bson:"edificio,omitempty"`
	} `bson:"history"`
	Position []struct {
		Latitudine  float64   `bson:"latitudine"`
		Longitudine float64   `bson:"longitudine"`
		Timestamp   time.Time `bson:"timestamp"`
	} `bson:"position"`
	StartTimestamp time.Time `bson:"start_timestamp"`
	IsActive       bool      `bson:"is_active"`
}

// Veniva usato nel tentativo gi generalizzazione delle select
func GetDocumentStruct(documentType string) any {
	if documentType == ParcelEnum {
		return &Parcel{}
	}
	if documentType == RouteEnum {
		return &Route{}
	}
	if documentType == FacilityEnum {
		return &Facility{}
	}
	return nil
}
