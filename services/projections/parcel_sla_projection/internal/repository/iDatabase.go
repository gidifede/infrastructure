package repository

import "context"

type IDatabase interface {
	InsertNewDocument(ctx context.Context, collectionName string, document any) error
	UpdateFieldsDocument(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToUpdate interface{}) error
	RetrieveProductDocument(ctx context.Context, collectionName string, filters map[string]interface{}) ([]Product, error)
}
