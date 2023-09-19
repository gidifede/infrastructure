package database

import (
	"confighandler/internal/model"
	"confighandler/internal/model/db"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type MyMongo struct {
	DB mongo.Database
}

func NewMongo(DB mongo.Database) IDatabase {
	return &MyMongo{DB: DB}
}

func (m *MyMongo) InsertProduct(c context.Context, product model.ProductSetup) error {

	coll := m.DB.Collection("product")

	toInsert := db.Product{
		Name: product.Product.Name,
		SLA: product.Product.SLA,
	}

	result, err := coll.InsertOne(c, toInsert)
	if err != nil {
		log.Error().Msgf("Error: %s", err)
		return err
	}
	fmt.Printf("Document inserted with ID: %s\n", result.InsertedID)

	return nil
}

func (m *MyMongo) InsertRoute(c context.Context, route model.RouteSetup) error {
	coll := m.DB.Collection("route")

	for _, r := range route {
		
		toInsert := db. Route{
			SourceFacilityID: r.SourceNodeID,
			DestFacilityID: r.DestinationNodeID,
			Cutoff: r.CutoffTime.Timestamp,
		}

		result, err := coll.InsertOne(c, toInsert)
		
		if err != nil {
			log.Error().Msgf("Error: %s", err)
			return err
		}
		fmt.Printf("Document inserted with ID: %s\n", result.InsertedID)
	}
	

	return nil
}

func (m *MyMongo) InsertTransport(c context.Context, transport model.TransportSetup) error {

	coll := m.DB.Collection("vehicle")

	for _, i := range transport {
		for _, t := range i.Vehicles {
			toInsert := db.Vehicle{
				VehicleID: t.ID,
				Type: t.Type,
				Capacity: t.Capacity,
				LicensePlate: t.LicensePlate,
			}

	        result, err := coll.InsertOne(c, toInsert)
	        if err != nil {
		        log.Error().Msgf("Error: %s", err)
		        return err
	        }
	        fmt.Printf("Document inserted with ID: %s\n", result.InsertedID)
		}
	}

	return nil
}

func (m *MyMongo) InsertNetwork(c context.Context, network model.NetworkSetup) error {

	coll := m.DB.Collection("network")

	result, err := coll.InsertOne(c, network)
	if err != nil {
		log.Error().Msgf("Error: %s", err)
		return err
	}
	fmt.Printf("Document inserted with ID: %s\n", result.InsertedID)

	return nil
}
