package manager

import (
	"context"
	"fleet-processing-projection/internal"
	"fleet-processing-projection/internal/models"
	"fleet-processing-projection/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type ParcelLoadedManager struct {
	logisticEvent models.ParcelLoaded
}

func newParcelLoadedManager(logisticEvent models.ParcelLoaded) IEventManager {
	return &ParcelLoadedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *ParcelLoadedManager) ManageEvent(ctx context.Context) error {

	/*
		Per ogni pacco caricato verifico se esiste il record, se esiste aggiungo il pacco altrimenti faccio il primo insert
	*/

	//Cerco se esiste gia lo stato nella history
	filter := bson.M{
		"transport_id": a.logisticEvent.TransportID,
	}

	exist := internal.Repo.FindDocument(ctx, repository.VehicleTransportCollectionEnum, filter)

	if !exist {
		lastStatus := "Transport Loading"
		//Insert primo record
		vehicleTransport := &repository.VehicleTransport{
			TransportID:    a.logisticEvent.TransportID,
			VehicleID:      a.logisticEvent.VehicleLicensePlate,
			Parcels:        []string{a.logisticEvent.ParcelID},
			History:        []repository.History{{Status: lastStatus, Timestamp: a.logisticEvent.Timestamp}},
			StartTimestamp: a.logisticEvent.Timestamp,
			IsActive:       true,
		}
		internal.Repo.InsertNewDocument(ctx, repository.VehicleTransportCollectionEnum, vehicleTransport)
	} else {

		filter := bson.M{
			repository.TransportIDEnum:       a.logisticEvent.TransportID,
			repository.TransportIsActiveEnum: true,
		}
		//Aggiornamento record - aggiungo il pacco
		update := bson.M{
			"$push": bson.M{
				"parcels": a.logisticEvent.ParcelID,
			},
		}

		internal.Repo.UpdateFieldsDocument(ctx, repository.VehicleTransportCollectionEnum, filter, update, nil)
	}

	return nil
}
