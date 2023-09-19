package models

import "time"

const (
	ParcelProcessedEvent        = "ParcelProcessed"
	ParcelProcessingFailedEvent = "ParcelProcessingFailed"
	ParcelUnloadedEvent         = "ParcelUnloaded"
	TransportStartedEvent       = "TransportStarted"
)

type ParcelProcessed struct {
	ParcelID              string    `json:"parcel_id"`
	DestinationFacilityID string    `json:"destination_facility_id"`
	Timestamp             time.Time `json:"timestamp"`
	FacilityID            string    `json:"facility_id"`
	SortingMachineID      string    `json:"sorting_machine_id"`
}

type ParcelUnloaded struct {
	ParcelID            string    `json:"parcel_id"`
	FacilityID          string    `json:"facility_id"`
	Timestamp           time.Time `json:"timestamp"`
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	TransportID         string    `json:"transport_id"`
}

type ParcelProcessingFailed struct {
	FacilityID       string    `json:"facility_id"`
	SortingMachineID string    `json:"sorting_machine_id"`
	ParcelID         string    `json:"parcel_id"`
	ErrMsg           string    `json:"err_msg"`
	Timestamp        time.Time `json:"timestamp"`
}

type TransportStarted struct {
	VehicleLicensePlate   string    `json:"vehicle_license_plate"`
	SourceFacilityID      string    `json:"source_facility_id"`
	DestinationFacilityID string    `json:"destination_facility_id"`
	Timestamp             time.Time `json:"timestamp"`
	TransportID           string    `json:"transport_id"`
}
