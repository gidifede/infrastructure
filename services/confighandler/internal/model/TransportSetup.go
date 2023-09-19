package model

type TransportSetup []TransportItem

type Vehicle struct {
	ID           string `json:"id" validate:"required"`
	Type         string `json:"type" validate:"required"`
	Capacity     int    `json:"capacity" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required"`
}

type TransportItem struct {
	Vehicles  []Vehicle `json:"vehicles"`
	Timestamp string    `json:"timestamp" validate:"required,datetime"`
}

