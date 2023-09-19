package repository

import "context"

type IDatabase interface {
	GetFacilityParcelDocumentByID(ctx context.Context, facilityID string) (FacilityParcel, error)
	GetFacilityDocumentByID(ctx context.Context, facilityID string) (Facility, error)
	GetParcelDocumentByID(ctx context.Context, parcelID string) (Parcel, error)
	GetVehicleTransportDocumentByTransportID(ctx context.Context, transportID string) (VehicleTransport, error)
	GetRouteDocumentBySourceAndDest(ctx context.Context, sourceFacilityID string, destFacilityID string) (Route, error)
	InsertNewDocument(ctx context.Context, collectionName string, document any) error
	UpdateDocumentFields(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToUpdate interface{}) error
	UpdateDocument(ctx context.Context, collectionName string, filters interface{}, fieldsToUpdate interface{}) error
	DeleteDocumentFields(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToRemove interface{}) error
}
