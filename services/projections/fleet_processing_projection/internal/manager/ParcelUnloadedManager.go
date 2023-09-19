package manager

import (
	"context"
	"fleet-processing-projection/internal"
	"fleet-processing-projection/internal/models"
	"fleet-processing-projection/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type ParcelUnloadedManager struct {
	logisticEvent models.ParcelUnloaded
}

func newParcelUnloadedManager(logisticEvent models.ParcelUnloaded) IEventManager {
	return &ParcelUnloadedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *ParcelUnloadedManager) ManageEvent(ctx context.Context) error {

	/*
		se è il primo unloading viene aggiunto lo status unloading altrimenti no
	*/

	lastStatus := "Transport Unloading"

	//Cerco se esiste gia lo stato nella history
	filter := bson.M{
		"transport_id":   a.logisticEvent.TransportID,
		"history.status": lastStatus,
	}

	exist := internal.Repo.FindDocument(ctx, repository.VehicleTransportCollectionEnum, filter)

	if !exist {
		filter := bson.M{
			repository.TransportIDEnum:       a.logisticEvent.TransportID,
			repository.TransportIsActiveEnum: false,
		}
		//Aggiornamento record - aggiungo il record nella history
		update := bson.M{
			"$push": bson.M{
				"history": &repository.History{Status: lastStatus, Timestamp: a.logisticEvent.Timestamp},
			},
		}
		internal.Repo.UpdateFieldsDocument(ctx, repository.VehicleTransportCollectionEnum, filter, update, nil)
	}

	return nil
}
