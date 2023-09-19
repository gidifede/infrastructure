package models

import (
	"time"
)

const (
	AcceptedEvent             = "Accepted"
	DeliveryCompletedEvent    = "DeliveryCompleted"
	RouteEventEnum            = "Route"
	RouteCollectionEventEnum  = "Route"
	ParcelCollectionEventEnum = "Parcel"
	AcceptedEventStatus       = "Accepted"
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
