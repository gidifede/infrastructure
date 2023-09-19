package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDatabase interface {
	GetSortingMachineDocumentByID(ctx context.Context, sortingMachineID string) (SortingMachine, error)
	InsertNewDocument(ctx context.Context, collectionName string, document any) error
	UpdateDocument(ctx context.Context, collectionName string, filters interface{}, fieldsToUpdate interface{}, arrayFilters *options.UpdateOptions) error
}
