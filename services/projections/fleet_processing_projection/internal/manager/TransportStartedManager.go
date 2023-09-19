package manager

import (
	"context"
	"fleet-processing-projection/internal"
	"fleet-processing-projection/internal/models"
	"fleet-processing-projection/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type TransportStartedManager struct {
	logisticEvent models.TransportStarted
}

func newTransportStartedManager(logisticEvent models.TransportStarted) IEventManager {
	return &TransportStartedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *TransportStartedManager) ManageEvent(ctx context.Context) error {
	/*
		Aggiungo il record alla history
	*/

	filter := bson.M{
		repository.TransportIDEnum:       a.logisticEvent.TransportID,
		repository.TransportIsActiveEnum: true,
	}

	lastStatus := "Transport Started"
	update := bson.M{
		"$push": bson.M{
			"history": &repository.History{Status: lastStatus, Timestamp: a.logisticEvent.Timestamp},
		},
	}
	internal.Repo.UpdateFieldsDocument(ctx, repository.VehicleTransportCollectionEnum, filter, update, nil)
	return nil
}
