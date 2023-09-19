package manager

import (
	"context"
	"parcel-processing-projection/internal"
	"parcel-processing-projection/internal/models"
	"parcel-processing-projection/internal/repository"

	"github.com/rs/zerolog/log"
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
		Va preparato il bson con i campi
			- last_status
			- position
			- aggiungo lo stato alla history

	*/
	lastStatus := "TransportStarted"
	status := &repository.Status{Status: lastStatus, Date: a.logisticEvent.Timestamp}

	filters := bson.M{
		"transport_id": a.logisticEvent.TransportID,
		"is_active":    true,
	}
	vehicleTransport, err := internal.Repo.RetrieveVehicleTransportDocument(ctx, repository.VehicleTransportCollectionEnum, filters, nil)
	if err != nil {
		log.Err(err).Msgf("cannot get vehicle_transport docs, filters %s", filters)
		return err
	}

	//pacchi contenuti nel camion da aggiornare
	idsDaAggiornare := []string{}
	if len(vehicleTransport) != 0 {
		idsDaAggiornare = vehicleTransport[0].Parcels
	}

	update := bson.M{
		"$push": bson.M{
			"history": status,
		},
		"$set": bson.M{
			"last_status": lastStatus,
			"position": bson.M{"type": "Vehicle",
				"id": a.logisticEvent.VehicleLicensePlate},
		},
	}

	err = internal.Repo.UpdateFieldsDocuments(ctx, repository.ParcelCollectionEnum, repository.ParcelCollectionIndexEnum, idsDaAggiornare, update)
	if err != nil {
		log.Err(err).Msgf("cannot update parcel docs, update %s", update)
		return err
	}
	return nil
}
