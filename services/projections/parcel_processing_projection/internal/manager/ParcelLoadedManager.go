package manager

import (
	"context"
	"parcel-processing-projection/internal"
	"parcel-processing-projection/internal/models"
	"parcel-processing-projection/internal/repository"

	"github.com/rs/zerolog/log"
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
		Va preparato il bson con i campi
			- last_status
			- aggiungo lo stato alla history

	*/
	lastStatus := "Loaded"
	status := &repository.Status{Status: lastStatus, Date: a.logisticEvent.Timestamp}

	update := bson.M{
		"$push": bson.M{
			"history": status,
		},
		"$set": bson.M{
			"last_status": lastStatus,
		},
	}

	err := internal.Repo.UpdateFieldsDocument(ctx, repository.ParcelCollectionEnum, repository.ParcelCollectionIndexEnum, a.logisticEvent.ParcelID, update)
	if err != nil {
		log.Err(err).Msgf("cannot update parcel doc with id %s, update %s", a.logisticEvent.ParcelID, update)
		return err
	}
	return nil
}
