package manager

import (
	"context"
	"facility-vehicle-projection/internal"
	"facility-vehicle-projection/internal/models"
	"facility-vehicle-projection/internal/repository"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type TransportEndedManager struct {
	logisticEvent models.TransportEnded
}

func newTransportEndedManager(logisticEvent models.TransportEnded) IEventManager {
	return &TransportEndedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *TransportEndedManager) ManageEvent(ctx context.Context) error {
	// 1. check if facility exists in facility collection
	//		- yes: go to 2
	//		- no: return error
	// 2. check if vehicle exists in vehicle collection
	//		- yes: go to 3
	//		- no: return error
	// 3. check if facility exists in facility_vehicle
	//		- yes: go to 4
	//		- no: insert facility_vehicle doc with vehicle && return
	// 4. check if vehicle exists in facility of facility_vehicle
	// 		- yes: return error
	//		- no: update facility_vehicle doc with vehicle

	_, err := internal.Repo.GetFacilityDocumentByID(ctx, a.logisticEvent.FacilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility document, facility_id: %s", a.logisticEvent.FacilityID)
		return err
	}
	_, err = internal.Repo.GetVehicleDocumentByID(ctx, a.logisticEvent.VehicleLicensePlate)
	if err != nil {
		log.Err(err).Msgf("cannot get vehicle document, vehicle_id: %s", a.logisticEvent.VehicleLicensePlate)
		return err
	}
	facilityVehicle, err := internal.Repo.GetFacilityVehicleDocumentByID(ctx, a.logisticEvent.FacilityID)
	if err != nil {
		// insert new doc
		newFacilityVehicle := repository.FacilityVehicle{
			FacilityID: a.logisticEvent.FacilityID,
			Vehicles: []repository.VehicleInFacility{
				{
					VehicleID:   a.logisticEvent.VehicleLicensePlate,
					Status:      "arrived",
					ArrivedTime: a.logisticEvent.Timestamp,
				},
			},
		}
		err := internal.Repo.InsertNewDocument(ctx, repository.FacilityVehicleCollection, newFacilityVehicle)
		if err != nil {
			log.Err(err).Msgf("cannot insert facility_vehicle document, doc: %s", newFacilityVehicle)
			return err
		}
		return nil
	}
	exists := false
	for _, vehicle := range facilityVehicle.Vehicles {
		if vehicle.VehicleID == a.logisticEvent.VehicleLicensePlate {
			exists = true
			break
		}
	}
	if exists {
		errMsg := fmt.Sprintf("vehicle already exists in facility, facility_id: %s, vehicle_id: %s", a.logisticEvent.FacilityID, a.logisticEvent.VehicleLicensePlate)
		log.Error().Msgf(errMsg)
		return fmt.Errorf(errMsg)
	}
	// update
	vehicleInFacility := repository.VehicleInFacility{
		VehicleID:   a.logisticEvent.VehicleLicensePlate,
		Status:      "arrived",
		ArrivedTime: a.logisticEvent.Timestamp,
	}
	update := bson.M{"$addToSet": bson.M{"vehicles": vehicleInFacility}}
	err = internal.Repo.UpdateDocumentFields(ctx, repository.FacilityVehicleCollection, repository.FacilityVehicleCollectionKey, a.logisticEvent.FacilityID, update)
	if err != nil {
		log.Err(err).Msgf("cannot update facility_vehicle document, facility_id: %s, update: %s", a.logisticEvent.FacilityID, update)
		return err
	}
	return nil
}
