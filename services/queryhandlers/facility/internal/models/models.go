package models

import "time"

type FacilityShortStats struct {
	FacilityHealth  string `json:"status"`
	ParcelWaiting   int    `json:"parcel_waiting"`
	ParcelProcessed int    `json:"parcel_processed"`
	NextCutOffTime  string `json:"next_cutoff_time"`
}

type FacilityParcelDetails struct {
	ParcelID    string    `json:"parcel_id"`
	TimeIn      time.Time `json:"time_in"`
	TimeOut     time.Time `json:"time_out"`
	NextHop     string    `json:"next_hop"`
	Destination string    `json:"destination"`
	Status      string    `json:"status"`
}

type FacilityParcelStats struct {
	Capacity          int     `json:"capacity"`
	ParcelWaiting     int     `json:"parcel_waiting"`
	ParcelProcessed   int     `json:"parcel_processed"`
	AvgProcessingTime float64 `json:"avg_processing_time"`
}

type ParcelProcessedItem struct {
	Day     string `json:"day"`
	Hour    int    `json:"hour"`
	Parcels int    `json:"parcels"`
}

type FacilitySortingMachineStats struct { // fix based on json schema
	// return n hours (put zero if no parcel in that hour)
	Capacity               int                   `json:"capacity"`
	ParcelProcessed        []ParcelProcessedItem `json:"parcel_processed"`
	WorkingCapacityAverage float64               `json:"working_capacity_average"` //-> avg by day
	NextMaintenance        time.Time             `json:"next_maintenance"`
}

type FacilityVehicleDetails struct {
	VehicleLicencePlate string    `json:"vehicle_licence_plate"`
	Status              string    `json:"status"`
	ArrivedTime         time.Time `json:"arrived_time"`
	NextStartTime       time.Time `json:"next_start_time"` // by vehicle id -> transport id -> understand how
}

type FacilityVehicleStats struct {
	VehiclesUnloading int `json:"vehicles_unloading"`
	VehiclesLoading   int `json:"vehicles_loading"`
}
