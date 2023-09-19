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

func (m *MongoDB) FindDocument(ctx context.Context, collectionName string, filters interface{}) bool {
	transport := &VehicleTransport{}
	collection := m.DB.Collection(collectionName)
	err := collection.FindOne(ctx, filters).Decode(&transport)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false //Nessun document trovato
		} else {
			return false //Nessun document trovato
		}
	}
	return true //Document trovato
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
func (m *MongoDB) UpdateFieldsDocument(ctx context.Context, collectionName string, filters interface{}, fieldsToUpdate interface{}, options *options.UpdateOptions) error {
	collection := m.DB.Collection(collectionName)
	res, err := collection.UpdateOne(ctx, filters, fieldsToUpdate, options)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("cannot update collection %s, with filter: %s", collectionName, filters)
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

func (m *MongoDB) ExistDocumentField(ctx context.Context, collectionName string, filters interface{}, attributeName string) (bool, error) {

	collection := m.DB.Collection(collectionName)

	result := collection.FindOne(ctx, filters)

	if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else if result.Err() != nil {
		return false, result.Err()
	}

	return true, nil
}
