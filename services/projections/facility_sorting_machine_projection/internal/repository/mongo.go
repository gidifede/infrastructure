package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	DB mongo.Database
}

func NewMongo(DB mongo.Database) IDatabase {
	return &MongoDB{DB: DB}
}

func (m *MongoDB) GetSortingMachineDocumentByID(ctx context.Context, sortingMachineID string) (SortingMachine, error) {
	collection := m.DB.Collection(SortingMachineCollection)
	filter := bson.M{SortingMachineCollectionKey: sortingMachineID}
	var doc SortingMachine
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return SortingMachine{}, err
	}
	return doc, nil
}

func (m *MongoDB) UpdateDocument(ctx context.Context, collectionName string, filters interface{}, fieldsToUpdate interface{}, arrayFilters *options.UpdateOptions) error {
	collection := m.DB.Collection(collectionName)
	res, err := collection.UpdateOne(ctx, filters, fieldsToUpdate, arrayFilters)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("cannot update collection %s with filter: %s", collectionName, fieldsToUpdate)
	}
	return nil
}

func (m *MongoDB) InsertNewDocument(ctx context.Context, collectionName string, document any) error {
	collection := m.DB.Collection(collectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}
