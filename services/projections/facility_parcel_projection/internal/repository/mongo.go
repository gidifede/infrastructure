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

func (m *MongoDB) GetFacilityParcelDocumentByID(ctx context.Context, facilityID string) (FacilityParcel, error) {
	collection := m.DB.Collection(FacilityParcelCollection)
	filter := bson.M{FacilityParcelCollectionKey: facilityID}
	var doc FacilityParcel
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return FacilityParcel{}, err
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

func (m *MongoDB) GetParcelDocumentByID(ctx context.Context, parcelID string) (Parcel, error) {
	collection := m.DB.Collection(ParcelCollection)
	filter := bson.M{ParcelCollectionKey: parcelID}
	var doc Parcel
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return Parcel{}, err
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

func (m *MongoDB) InsertNewDocument(ctx context.Context, collectionName string, document any) error {
	collection := m.DB.Collection(collectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
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

func (m *MongoDB) GetRouteDocumentBySourceAndDest(ctx context.Context, sourceFacilityID string, destFacilityID string) (Route, error) {
	collection := m.DB.Collection(RouteCollection)
	filter := bson.M{"source_facility_id": sourceFacilityID, "dest_facility_id": destFacilityID}
	var doc Route
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return Route{}, err
	}
	return doc, nil
}
