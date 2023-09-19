package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	DB mongo.Database
}

func NewMongo(DB mongo.Database) IDatabase {
	return &MongoDB{DB: DB}
}

func (m *MongoDB) GetFacilityDocumentByID(ctx context.Context, facilityID string) (Facility, error) {
	collection := m.DB.Collection(FacilityCollection)
	filter := bson.M{FacilityCollectionKey: facilityID}
	var doc Facility
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return Facility{}, err
	}
	return doc, nil
}

func (m *MongoDB) GetVehicleDocumentByID(ctx context.Context, vehicleID string) (Vehicle, error) {
	collection := m.DB.Collection(VehicleCollection)
	filter := bson.M{VehicleCollectionKey: vehicleID}
	var doc Vehicle
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return Vehicle{}, err
	}
	return doc, nil
}

func (m *MongoDB) GetFacilityVehicleDocumentByID(ctx context.Context, facilityID string) (FacilityVehicle, error) {
	collection := m.DB.Collection(FacilityVehicleCollection)
	filter := bson.M{FacilityVehicleCollectionKey: facilityID}
	var doc FacilityVehicle
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return FacilityVehicle{}, err
	}
	return doc, nil
}

func (m *MongoDB) InsertNewDocument(ctx context.Context, collectionName string, document any) error {
	collection := m.DB.Collection(collectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) UpdateDocumentFields(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToUpdate interface{}) error {
	collection := m.DB.Collection(collectionName)
	filter := bson.M{documentIDKey: documentID}
	res, err := collection.UpdateOne(ctx, filter, fieldsToUpdate)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("cannot update collection %s, doc %s = %s with filter: %s", collectionName, documentIDKey, documentID, fieldsToUpdate)
	}
	return nil
}

func (m *MongoDB) UpdateDocument(ctx context.Context, collectionName string, filters interface{}, fieldsToUpdate interface{}) error {
	collection := m.DB.Collection(collectionName)
	res, err := collection.UpdateOne(ctx, filters, fieldsToUpdate)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("cannot update collection %s with filter: %s", collectionName, fieldsToUpdate)
	}
	return nil
}

func (m *MongoDB) DeleteDocumentFields(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToRemove interface{}) error {
	collection := m.DB.Collection(collectionName)
	filter := bson.M{documentIDKey: documentID}
	res, err := collection.UpdateOne(ctx, filter, fieldsToRemove)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("cannot remove element from collection %s, doc %s = %s with filter: %s", collectionName, documentIDKey, documentID, fieldsToRemove)
	}
	return nil
}
