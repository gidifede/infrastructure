package repository

import "time"

const (
	FacilityVehicleCollection    = "facility_vehicle"
	FacilityVehicleCollectionKey = "facility_id"

	FacilityCollection    = "facility"
	FacilityCollectionKey = "facility_id"

	VehicleCollection    = "vehicle"
	VehicleCollectionKey = "vehicle_id"
)

type FacilityVehicle struct {
	FacilityID string              `bson:"facility_id"`
	Vehicles   []VehicleInFacility `bson:"vehicles"`
}

type VehicleInFacility struct {
	VehicleID   string    `bson:"vehicle_id"`
	Status      string    `bson:"status"`
	ArrivedTime time.Time `bson:"arrived_time"`
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

type Vehicle struct {
	VehicleID string `bson:"vehicle_id"`
	Capacity  int    `bson:"capacity"`
	Type      string `bson:"type"`
}
