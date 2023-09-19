package manager

import (
	"context"
	"fleet-processing-projection/internal"
	"fleet-processing-projection/internal/models"
	"fleet-processing-projection/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type PositionManager struct {
	logisticEvent models.Position
}

func newPositionManager(logisticEvent models.Position) IEventManager {
	return &PositionManager{
		logisticEvent: logisticEvent,
	}
}

func (a *PositionManager) ManageEvent(ctx context.Context) error {

	/*
		Controllo se è stato inserito il campo o è la prima position inserita
	*/
	filterExist := bson.M{
		repository.TransportPositionEnum: bson.M{"$exists": true},
		repository.VehicleIDEnum:         a.logisticEvent.VehicleLicensePlate,
		repository.TransportIsActiveEnum: true,
	}
	exist, err := internal.Repo.ExistDocumentField(ctx, repository.VehicleTransportCollectionEnum, filterExist, repository.TransportPositionEnum)

	if err != nil {
		return err
	}

	var operation string
	if !exist {
		operation = "$set"
	} else {
		operation = "$push"
	}

	update := bson.M{
		operation: bson.M{
			repository.TransportPositionEnum: &[]repository.Position{
				repository.Position{
					Latitudine:  a.logisticEvent.Latitude,
					Longitudine: a.logisticEvent.Longitude,
					Timestamp:   a.logisticEvent.Timestamp,
				},
			},
		},
	}

	filter := bson.M{
		repository.VehicleIDEnum:         a.logisticEvent.VehicleLicensePlate,
		repository.TransportIsActiveEnum: true,
	}

	internal.Repo.UpdateFieldsDocument(ctx, repository.VehicleTransportCollectionEnum, filter, update, nil)

	return nil
}
