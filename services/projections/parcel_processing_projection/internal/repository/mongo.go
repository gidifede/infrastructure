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

// Inserisco un documento che non esiste
func (m *MongoDB) InsertNewDocument(ctx context.Context, collectionName string, document any) error {
	collection := m.DB.Collection(collectionName)
	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

// Aggiorno piu attributi di un documento
func (m *MongoDB) UpdateFieldsDocument(ctx context.Context, collectionName string, documentIDKey string, documentID string, fieldsToUpdate interface{}) error {
	collection := m.DB.Collection(collectionName)
	filter := bson.M{documentIDKey: documentID}
	res, err := collection.UpdateOne(ctx, filter, fieldsToUpdate)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("cannot update collection %s, doc_id %s, with filter: %s", collectionName, documentID, fieldsToUpdate)
	}
	return nil
}

// Aggiorno piu attributi di piu documenti
func (m *MongoDB) UpdateFieldsDocuments(ctx context.Context, collectionName string, documentIDKey string, documentsID []string, fieldsToUpdate interface{}) error {
	collection := m.DB.Collection(collectionName)
	filter := bson.M{documentIDKey: bson.M{"$in": documentsID}}
	res, err := collection.UpdateMany(ctx, filter, fieldsToUpdate)
	if err != nil {
		return err
	}
	if (int)(res.ModifiedCount) != len(documentsID) {
		return fmt.Errorf("cannot update collection %s, doc_ids %s, with filter: %s", collectionName, documentsID, fieldsToUpdate)
	}
	return nil
}

// inserisco un attributo di documento
func (m *MongoDB) InsertFieldDocument(ctx context.Context, collectionName string, documentIDKey string, documentID string, attribute string, newValue any) error {
	collection := m.DB.Collection(collectionName)
	filter := bson.M{documentIDKey: documentID}
	update := bson.M{"$set": bson.M{attribute: newValue}}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("cannot update collection %s, doc_id %s, attribut %s, value %s", collectionName, documentID, attribute, newValue)
	}
	return nil
}

// Elimino un documento
func (m *MongoDB) DeleteDocument(ctx context.Context, collectionName string, documentID string) error {
	//do not implement
	return nil
}

// Elimino un attributo di un document
func (m *MongoDB) DeleteFieldDocument(ctx context.Context, collectionName string, documentID string) error {
	//do not implement
	return nil
}

// Faccio la count per documentID
func (m *MongoDB) CountDocument(ctx context.Context, collectionName string, documentIDKey string, documentID string) (int64, error) {
	collection := m.DB.Collection(collectionName)
	filter := bson.M{documentIDKey: documentID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Select per recuperare tutti i document in una collection (Tentativo di generalizzazione delal select)
func (m *MongoDB) RetrieveDocument(ctx context.Context, collectionName string, documentType string, filters map[string]interface{}) (interface{}, error) {
	collection := m.DB.Collection(collectionName)

	//Creo i filtri in maniera dinamica in base alla valorizzazione della mappa in input
	var filter interface{}
	if filters == nil {
		filter = bson.M{}
	} else {
		filter = bson.M(filters)
	}

	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documentList []any

	for cursor.Next(ctx) {
		result := GetDocumentStruct(documentType)
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		documentList = append(documentList, result)
	}

	return documentList, nil

}

func (m *MongoDB) RetrieveFacilityDocument(ctx context.Context, collectionName string, filters map[string]interface{}) ([]Facility, error) {
	collection := m.DB.Collection(collectionName)

	//Creo i filtri in maniera dinamica in base alla valorizzazione della mappa in input
	var filter interface{}
	if filters == nil {
		filter = bson.M{}
	} else {
		filter = bson.M(filters)
	}

	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documentList []Facility

	for cursor.Next(ctx) {
		result := &Facility{}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		documentList = append(documentList, *result)
	}

	return documentList, nil

}

func (m *MongoDB) RetrieveRouteDocument(ctx context.Context, collectionName string, filters map[string]interface{}) ([]Route, error) {
	collection := m.DB.Collection(collectionName)

	//Creo i filtri in maniera dinamica in base alla valorizzazione della mappa in input
	var filter interface{}
	if filters == nil {
		filter = bson.M{}
	} else {
		filter = bson.M(filters)
	}

	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documentList []Route

	for cursor.Next(ctx) {
		result := &Route{}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		documentList = append(documentList, *result)
	}

	return documentList, nil

}

func (m *MongoDB) RetrieveVehicleTransportDocument(ctx context.Context, collectionName string, filters map[string]interface{}, options *options.FindOptions) ([]VehicleTransport, error) {
	collection := m.DB.Collection(collectionName)

	//Creo i filtri in maniera dinamica in base alla valorizzazione della mappa in input
	var filter interface{}
	if filters == nil {
		filter = bson.M{}
	} else {
		filter = bson.M(filters)
	}

	cursor, err := collection.Find(ctx, filter, options)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documentList []VehicleTransport

	for cursor.Next(ctx) {
		result := &VehicleTransport{}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		documentList = append(documentList, *result)
	}

	return documentList, nil

}

func (m *MongoDB) RetrieveParcelTransportDocument(ctx context.Context, collectionName string, filters map[string]interface{}, options *options.FindOptions) ([]Parcel, error) {
	collection := m.DB.Collection(collectionName)

	//Creo i filtri in maniera dinamica in base alla valorizzazione della mappa in input
	var filter interface{}
	if filters == nil {
		filter = bson.M{}
	} else {
		filter = bson.M(filters)
	}

	cursor, err := collection.Find(ctx, filter, options)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var documentList []Parcel

	for cursor.Next(ctx) {
		result := &Parcel{}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		documentList = append(documentList, *result)
	}

	return documentList, nil

}
