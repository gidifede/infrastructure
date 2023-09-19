package repository

import (
	"time"
)

const (
	VehicleTransportCollectionEnum 		= "vehicle_transport"
	VehicleIDEnum      					= "vehicle_id"
	TransportIDEnum						= "transport_id"
	TransportPositionEnum 				= "position"
	TransportHistoryEnum 				= "history"
	TransportIsActiveEnum 				= "is_active"
)


type VehicleTransport struct {
	TransportID 	string   	`bson:"transport_id"`
	VehicleID   	string   	`bson:"vehicle_id"`
	Parcels     	[]string 	`bson:"parcels"`
	History     	[]History 	`bson:"history"`
	Position 		[]Position 	`bson:"position,omitempty"`
	StartTimestamp 	time.Time 	`bson:"start_timestamp"`
	IsActive       	bool      	`bson:"is_active"`
}

type Position struct{
	Latitudine  float64   `bson:"latitudine"`
	Longitudine float64   `bson:"longitudine"`
	Timestamp   time.Time `bson:"timestamp"`
}

type History struct{
	Status   	  string    `bson:"status"`
	Timestamp     time.Time `bson:"timestamp"`
	Edificio 	  string    `bson:"edificio,omitempty"`
}