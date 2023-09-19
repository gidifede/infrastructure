package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDatabase interface {
	InsertNewDocument(ctx context.Context, collectionName string, document any) error
	FindDocument(ctx context.Context, collectionName string, filters interface{}) bool
	UpdateFieldsDocument(ctx context.Context, collectionName string, filters interface{}, fieldsToUpdate interface{}, options *options.UpdateOptions) error
	InsertFieldDocument(ctx context.Context, collectionName string, documentIDKey string, documentID string, attribute string, newValue any) error
	ExistDocumentField(ctx context.Context, collectionName string, filters interface{}, attributeName string) (bool, error)
}
