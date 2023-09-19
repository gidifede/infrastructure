package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDatabase interface {
	InsertNewDocument(ctx context.Context, collectionName string, document any) error
	UpdateFieldsDocument(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToUpdate interface{}) error
	UpdateFieldsDocuments(ctx context.Context, collectionName string, documentIDKey string, documentsID []string, fieldsToUpdate interface{}) error
	InsertFieldDocument(ctx context.Context, collectionName string, documentIDKey string, documentID string, attribute string, newValue any) error
	DeleteDocument(ctx context.Context, collectionName string, documentID string) error
	DeleteFieldDocument(ctx context.Context, collectionName string, documentID string) error
	CountDocument(ctx context.Context, collectionName string, documentIDKey string, documentID string) (int64, error)
	RetrieveDocument(ctx context.Context, collectionName string, documentType string, filters map[string]interface{}) (interface{}, error)
	RetrieveFacilityDocument(ctx context.Context, collectionName string, filters map[string]interface{}) ([]Facility, error)
	RetrieveRouteDocument(ctx context.Context, collectionName string, filters map[string]interface{}) ([]Route, error)
	RetrieveVehicleTransportDocument(ctx context.Context, collectionName string, filters map[string]interface{}, options *options.FindOptions) ([]VehicleTransport, error)
	RetrieveParcelTransportDocument(ctx context.Context, collectionName string, filters map[string]interface{}, options *options.FindOptions) ([]Parcel, error)
}
