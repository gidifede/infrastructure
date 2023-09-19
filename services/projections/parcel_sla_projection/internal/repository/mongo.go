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

// Select per la collection di product
func (m *MongoDB) RetrieveProductDocument(ctx context.Context, collectionName string, filters map[string]interface{}) ([]Product, error) {
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

	var documentList []Product

	for cursor.Next(ctx) {
		result := &Product{}
		if err := cursor.Decode(&result); err != nil {
			fmt.Println(err)
			return nil, err
		}
		documentList = append(documentList, *result)
	}

	return documentList, nil
}
