package manager

import (
	"context"
	"fleet-processing-projection/internal"
	"fleet-processing-projection/internal/models"
	"fleet-processing-projection/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type TransportEndedManager struct {
	logisticEvent models.TransportEnded
}

func newTransportEndedManager(logisticEvent models.TransportEnded) IEventManager {
	return &TransportEndedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *TransportEndedManager) ManageEvent(ctx context.Context) error {
	/*
		Aggiungo il record alla history
	*/
	lastStatus := "Transport Ended"

	filter := bson.M{
		repository.TransportIDEnum:       a.logisticEvent.TransportID,
		repository.TransportIsActiveEnum: true,
	}

	update := bson.M{
		"$push": bson.M{
			"history": &repository.History{Status: lastStatus, Timestamp: a.logisticEvent.Timestamp},
		},
		"$set": bson.M{
			"is_active": false,
		},
	}
	internal.Repo.UpdateFieldsDocument(ctx, repository.VehicleTransportCollectionEnum, filter, update, nil)
	return nil
}
