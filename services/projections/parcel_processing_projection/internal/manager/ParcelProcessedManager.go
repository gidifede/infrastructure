package manager

import (
	"context"
	"parcel-processing-projection/internal"
	"parcel-processing-projection/internal/models"
	"parcel-processing-projection/internal/repository"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type ProcessedManager struct {
	logisticEvent models.ParcelProcessed
}

func newParcelProcessedManager(logisticEvent models.ParcelProcessed) IEventManager {
	return &ProcessedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *ProcessedManager) ManageEvent(ctx context.Context) error {
	/*
		Va preparato il bson con i campi
			- last_status
			- position
			- aggiorno la percentuale di avanzamento
			- aggiungo lo stato alla history
	*/

	lastStatus := "Processed"
	filtersParcel := bson.M{"parcel_id": a.logisticEvent.ParcelID}
	parcel, err := internal.Repo.RetrieveParcelTransportDocument(ctx, repository.ParcelCollectionEnum, filtersParcel, nil)
	if err != nil {
		log.Err(err).Msgf("cannot get parcel doc with id %s", a.logisticEvent.ParcelID)
		return err
	}
	pathCompleted := CalculatePathPercentage(parcel[0].ParcelPath.Path, a.logisticEvent.FacilityID)
	status := &repository.Status{Status: lastStatus, Date: a.logisticEvent.Timestamp}

	filtersFacilities := bson.M{"facility_id": a.logisticEvent.FacilityID}
	facilities, err := internal.Repo.RetrieveFacilityDocument(ctx, repository.FacilityCollectionEnum, filtersFacilities)
	if err != nil {
		log.Err(err).Msgf("cannot get facility doc with id %s, update %s", a.logisticEvent.FacilityID)
		return err
	}

	facilityType := facilities[0].FacilityType

	update := bson.M{
		"$push": bson.M{
			"history": status,
		},
		"$set": bson.M{
			"last_status": lastStatus,
			"position": bson.M{"type": facilityType,
				"id": a.logisticEvent.FacilityID},
			"parcel_path.path_completed": pathCompleted,
		},
	}

	err = internal.Repo.UpdateFieldsDocument(ctx, repository.ParcelCollectionEnum, repository.ParcelCollectionIndexEnum, a.logisticEvent.ParcelID, update)
	if err != nil {
		log.Err(err).Msgf("cannot update parcel doc with id %s, update %s", a.logisticEvent.ParcelID, update)
		return err
	}
	return nil

}
