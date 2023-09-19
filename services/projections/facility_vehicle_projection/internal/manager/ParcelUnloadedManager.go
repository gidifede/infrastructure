package manager

import (
	"context"
	"facility-vehicle-projection/internal"
	"facility-vehicle-projection/internal/models"
	"facility-vehicle-projection/internal/repository"

	"github.com/rs/zerolog/log"
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
	// 2. check if vehicle is in facility
	//		- no: return error
	//		- yes: go to 3
	// 3. update vehicle status to unloading
	filters := bson.M{repository.FacilityVehicleCollectionKey: a.logisticEvent.FacilityID, "vehicles.vehicle_id": a.logisticEvent.VehicleLicensePlate}
	update := bson.M{"$set": bson.M{"vehicles.$.status": "unloading"}}
	err := internal.Repo.UpdateDocument(ctx, repository.FacilityVehicleCollection, filters, update)
	if err != nil {
		log.Err(err).Msgf("cannot update facility_vehicle, facility_id: %s, vehicle_id: %s, update: %s", a.logisticEvent.FacilityID, a.logisticEvent.VehicleLicensePlate, update)
		return err
	}
	return nil
}
