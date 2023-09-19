package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	DB mongo.Database
}

func NewMongo(DB mongo.Database) IDatabase {
	return &MongoDB{DB: DB}
}
