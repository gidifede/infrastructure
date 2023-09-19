package database

import (
	"context"
	"parcel_queryhandler/internal/model"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	DB mongo.Database
}

func NewMongoDB(DB mongo.Database) IDatabase {
	return &MongoDB{
		DB: DB,
	}
}

func (m *MongoDB) SelectNetworkNodes(c context.Context, filter string) (model.Network, error) {

	var results []model.Network

	cursor, err := m.DB.Collection("network").Find(c, bson.M{})
	for cursor.Next(c) {
		var n model.Network
		if err := cursor.Decode(&n); err != nil {
			log.Error().Msg(err.Error())
			return model.Network{}, err
		}
		results = append(results, n)
	}
	if err != nil {
		log.Error().Msg(err.Error())
		return model.Network{}, err
	}
	defer cursor.Close(c)

	return results[0], nil
}
