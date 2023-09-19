package repository

import (
	"time"
)

const (
	ProductEnum                     = "Product"
	ProductCollectionEnum           = "product"
	ParcelCollectionDocumentKeyEnum = "parcel_id"
	ParcelSLACollectionEnum         = "parcel_sla"
)

type Product struct {
	Name string `bson:"name"`
	SLA  string `bson:"SLA"`
}

type ParcelSLA struct {
	ParcelID             string    `bson:"parcel_id"`
	ExpectedDeliveryDate time.Time `bson:"expected_delivery_date"`
	DeliveryDate         time.Time `bson:"delivery_date,omitempty"`
}
