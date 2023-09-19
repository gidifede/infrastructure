package models

import (
	"time"
)

const (
	AcceptedEvent               = "Accepted"
	ParcelProcessedEvent        = "ParcelProcessed"
	ParcelProcessingFailedEvent = "ParcelProcessingFailed"
	ParcelLoadedEvent           = "ParcelLoaded"
	ParcelUnloadedEvent         = "ParcelUnloaded"
	TransportEndedEvent         = "TransportEnded"
	TransportStartedEvent       = "TransportStarted"
	DeliveryCompletedEvent      = "DeliveryCompleted"
	RouteEventEnum              = "Route"
	AcceptedEventStatus         = "Accepted"
)

type Accepted struct {
	Parcel struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"parcel"`
	FacilityID string `json:"facility_id"`
	Sender     struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Zipcode string `json:"zipcode"`
		City    string `json:"city"`
		Nation  string `json:"nation"`
	} `json:"sender"`
	Receiver struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Zipcode string `json:"zipcode"`
		City    string `json:"city"`
		Nation  string `json:"nation"`
		Number  string `json:"number"`
		Email   string `json:"email"`
	} `json:"receiver"`
	Timestamp time.Time `json:"timestamp"`
}

type DeliveryCompleted struct {
	ParcelID  string    `json:"parcel_id"`
	Timestamp time.Time `json:"timestamp"`
}

type ParcelLoaded struct {
	ParcelID            string    `json:"parcel_id"`
	FacilityID          string    `json:"facility_id"`
	TransportID			string	  `json:"transport_id"`
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	Timestamp           time.Time `json:"timestamp"`
}

type ParcelUnloaded struct {
	ParcelID            string    `json:"parcel_id"`
	FacilityID          string    `json:"facility_id"`
	TransportID			string	  `json:"transport_id"`
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	Timestamp           time.Time `json:"timestamp"`
}

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

// Per questo evento si deve recuperare i parcel che sta transportando
type TransportStarted struct {
	VehicleLicensePlate   string    `json:"vehicle_license_plate"`
	SourceFacilityID      string    `json:"source_facility_id"`
	TransportID			  string	`json:"transport_id"`
	DestinationFacilityID string    `json:"destination_facility_id"`
	Timestamp             time.Time `json:"timestamp"`
}

// Per questo evento si deve recuperare i parcel che sta transportando
type TransportEnded struct {
	VehicleLicensePlate string    `json:"vehicle_license_plate"`
	FacilityID          string    `json:"facility_id"`
	TransportID			  string	`json:"transport_id"`
	Timestamp           time.Time `json:"timestamp"`
}
