package models

import (
	"time"
)

const (
	ParcelProcessedEvent        = "ParcelProcessed"
	ParcelProcessingFailedEvent = "ParcelProcessingFailed"
)

type ParcelProcessed struct {
	FacilityID            string    `json:"facility_id"`
	SortingMachineID      string    `json:"sorting_machine_id"`
	ParcelID              string    `json:"parcel_id"`
	DestinationFacilityID string    `json:"destination_facility_id"`
	Timestamp             time.Time `json:"timestamp"`
}

type ParcelProcessingFailed struct {
	FacilityID       string    `json:"facility_id"`
	SortingMachineID string    `json:"sorting_machine_id"`
	ParcelID         string    `json:"parcel_id"`
	ErrMsg           string    `json:"err_msg"`
	Timestamp        time.Time `json:"timestamp"`
}
