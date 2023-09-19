package models

import (
	"time"
)

const (
	ParcelLoadedEvent           = "ParcelLoaded"
	ParcelUnloadedEvent         = "ParcelUnloaded"
	TransportEndedEvent         = "TransportEnded"
	TransportStartedEvent       = "TransportStarted"
	PositionEvent				= "Position"
)

type Position struct {
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	Latitude			float64	  `json:"latitude"`
	Longitude			float64	  `json:"longitute"`
	Timestamp           time.Time `json:"timestamp"`
}

type ParcelLoaded struct {
	TransportID			string	  `json:"transport_id"`
	ParcelID            string    `json:"parcel_id"`
	FacilityID          string    `json:"facility_id"`
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	Timestamp           time.Time `json:"timestamp"`
}

type ParcelUnloaded struct {
	TransportID			string	  `json:"transport_id"`
	ParcelID            string    `json:"parcel_id"`
	FacilityID          string    `json:"facility_id"`
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	Timestamp           time.Time `json:"timestamp"`
}


type TransportStarted struct {
	TransportID			string	  `json:"transport_id"`
	VehicleLicensePlate   string    `json:"vehicle_license_plate"`
	SourceFacilityID      string    `json:"source_facility_id"`
	DestinationFacilityID string    `json:"destination_facility_id"`
	Timestamp             time.Time `json:"timestamp"`
}

type TransportEnded struct {
	TransportID			string	  `json:"transport_id"`
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	FacilityID          string    `json:"facility_id"`
	Timestamp           time.Time `json:"timestamp"`
}
