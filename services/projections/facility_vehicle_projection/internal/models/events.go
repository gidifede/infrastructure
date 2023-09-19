package models

import "time"

const (
	ParcelLoadedEvent     = "ParcelLoaded"
	ParcelUnloadedEvent   = "ParcelUnloaded"
	TransportStartedEvent = "TransportStarted"
	TransportEndedEvent   = "TransportEnded"
)

type ParcelLoaded struct {
	ParcelID            string    `json:"parcel_id"`
	FacilityID          string    `json:"facility_id"`
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	Timestamp           time.Time `json:"timestamp"`
	TransportID         string    `json:"transport_id"`
}

type ParcelUnloaded struct {
	ParcelID            string    `json:"parcel_id"`
	FacilityID          string    `json:"facility_id"`
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	Timestamp           time.Time `json:"timestamp"`
	TransportID         string    `json:"transport_id"`
}

type TransportStarted struct {
	VehicleLicensePlate   string    `json:"vehicle_license_plate"`
	SourceFacilityID      string    `json:"source_facility_id"`
	DestinationFacilityID string    `json:"destination_facility_id"`
	Timestamp             time.Time `json:"timestamp"`
	TransportID           string    `json:"transport_id"`
}

type TransportEnded struct {
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	FacilityID          string    `json:"facility_id"`
	Timestamp           time.Time `json:"timestamp"`
	TransportID         string    `json:"transport_id"`
}
