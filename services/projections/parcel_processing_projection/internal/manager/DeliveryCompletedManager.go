package manager

import (
	"context"
	"parcel-processing-projection/internal"
	"parcel-processing-projection/internal/models"
	"parcel-processing-projection/internal/repository"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type DeliveryCompletedManager struct {
	logisticEvent models.DeliveryCompleted
}

func newDeliveryCompletedManager(logisticEvent models.DeliveryCompleted) IEventManager {
	return &DeliveryCompletedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *DeliveryCompletedManager) ManageEvent(ctx context.Context) error {
	/*
		Va preparato il bson con i campi
			- last_status
			- position
			- aggiorno la percentuale di avanzamento
			- aggiungo lo stato alla history

	*/
	lastStatus := "Delivered"
	status := &repository.Status{Status: lastStatus, Date: a.logisticEvent.Timestamp}

	update := bson.M{
		"$push": bson.M{
			"history": status,
		},
		"$set": bson.M{
			"last_status":                lastStatus,
			"parcel_path.path_completed": 100, //Se viene deliverato la percentuale del completamento è all 100%
		},
	}

	err := internal.Repo.UpdateFieldsDocument(ctx, repository.ParcelCollectionEnum, repository.ParcelCollectionIndexEnum, a.logisticEvent.ParcelID, update)
	if err != nil {
		log.Err(err).Msgf("cannot update parcel doc with id %s, update %s", a.logisticEvent.ParcelID, update)
		return err
	}
	return nil
}
