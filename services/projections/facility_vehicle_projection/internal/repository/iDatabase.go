package repository

import "context"

type IDatabase interface {
	GetFacilityDocumentByID(ctx context.Context, facilityID string) (Facility, error)
	GetVehicleDocumentByID(ctx context.Context, vehicleID string) (Vehicle, error)
	GetFacilityVehicleDocumentByID(ctx context.Context, facilityID string) (FacilityVehicle, error)
	InsertNewDocument(ctx context.Context, collectionName string, document any) error
	UpdateDocumentFields(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToUpdate interface{}) error
	UpdateDocument(ctx context.Context, collectionName string, filters interface{}, fieldsToUpdate interface{}) error
	DeleteDocumentFields(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToRemove interface{}) error
}
