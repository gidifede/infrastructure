package repository

import "context"

type IDatabase interface {
	GetFacilityExpectedParcelDocumentByID(ctx context.Context, facilityID string) (FacilityExpectedParcel, error)
	GetFacilityDocumentByID(ctx context.Context, facilityID string) (Facility, error)
	GetVehicleTransportDocumentByTransportID(ctx context.Context, transportID string) (VehicleTransport, error)
	InsertNewDocument(ctx context.Context, collectionName string, document any) error
	UpdateDocumentFields(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToUpdate interface{}) error
}
