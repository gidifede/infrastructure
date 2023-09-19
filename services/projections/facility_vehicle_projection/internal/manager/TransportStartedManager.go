package manager

import (
	"context"
	"facility-vehicle-projection/internal"
	"facility-vehicle-projection/internal/models"
	"facility-vehicle-projection/internal/repository"

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
	// 1. check if vehicle is in facility
	//		- yes: remove vehicle from facility
	//		- no: return error
	update := bson.M{
		"$pull": bson.M{
			"vehicles": bson.M{"vehicle_id": a.logisticEvent.VehicleLicensePlate},
		},
	}
	err := internal.Repo.DeleteDocumentFields(ctx, repository.FacilityVehicleCollection, repository.FacilityVehicleCollectionKey, a.logisticEvent.SourceFacilityID, update)
	if err != nil {
		log.Err(err).Msgf("cannot delete vehicle from facility_vehicle, vehicle_id: %s, facility_id: %s", a.logisticEvent.VehicleLicensePlate, a.logisticEvent.SourceFacilityID)
		return err
	}
	return nil
}
