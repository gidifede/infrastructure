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

func (m *MongoDB) GetFacilityExpectedParcelDocumentByID(ctx context.Context, facilityID string) (FacilityExpectedParcel, error) {
	collection := m.DB.Collection(FacilityExpectedParcelCollection)
	filter := bson.M{FacilityExpectedParcelCollectionKey: facilityID}
	var doc FacilityExpectedParcel
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return FacilityExpectedParcel{}, err
	}
	return doc, nil
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

func (m *MongoDB) InsertNewDocument(ctx context.Context, collectionName string, document any) error {
	collection := m.DB.Collection(collectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) GetVehicleTransportDocumentByTransportID(ctx context.Context, transportID string) (VehicleTransport, error) {
	collection := m.DB.Collection(VehicleTransportCollection)
	filter := bson.M{VehicleTransportCollectionKey: transportID}
	var doc VehicleTransport
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return VehicleTransport{}, err
	}
	return doc, nil
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
