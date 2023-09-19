package manager

import (
	"context"
	"facility-parcel-projection/internal"
	"facility-parcel-projection/internal/models"
	"facility-parcel-projection/internal/repository"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type ParcelProcessingFailedManager struct {
	logisticEvent models.ParcelProcessingFailed
}

func newParcelProcessingFailedManager(logisticEvent models.ParcelProcessingFailed) IEventManager {
	return &ParcelProcessingFailedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *ParcelProcessingFailedManager) ManageEvent(ctx context.Context) error {
	// 1. check if parcel is in facility
	//		- yes: update status to processingFailed
	//		- no: return error
	filters := bson.M{repository.FacilityParcelCollectionKey: a.logisticEvent.FacilityID, "parcels.parcel_id": a.logisticEvent.ParcelID}
	update := bson.M{"$set": bson.M{"parcels.$.status": "processingFailed"}}
	err := internal.Repo.UpdateDocument(ctx, repository.FacilityParcelCollection, filters, update)
	if err != nil {
		log.Err(err).Msgf("cannot update facility_parcel, facility_id: %s, parcel_id: %s, update: %s", a.logisticEvent.FacilityID, a.logisticEvent.ParcelID, update)
		return err
	}
	return nil
}
