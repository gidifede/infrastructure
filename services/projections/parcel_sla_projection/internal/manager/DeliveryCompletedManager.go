package manager

import (
	"context"
	"parcel-sla-projection/internal"
	"parcel-sla-projection/internal/models"
	"parcel-sla-projection/internal/repository"

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
		Alla ricezione dell'evento di delivery va aggiornato il campo delivery_date con la data effettiva di delivery
	*/
	parcelID := a.logisticEvent.ParcelID
	deliveryDate := a.logisticEvent.Timestamp

	//preparo il bson di update
	update := bson.M{
		"$set": bson.M{
			"delivery_date": deliveryDate,
		},
	}

	err := internal.Repo.UpdateFieldsDocument(ctx, repository.ParcelSLACollectionEnum, repository.ParcelCollectionDocumentKeyEnum, parcelID, update)
	if err != nil {
		log.Err(err).Msgf("cannot update parcel_sla doc, parcel_id %s, update %s", parcelID, update)
		return err
	}
	return nil
}
