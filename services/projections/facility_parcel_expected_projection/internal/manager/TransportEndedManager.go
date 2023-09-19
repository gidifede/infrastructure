package manager

import (
	"context"
	"facility-parcel-expected-projection/internal"
	"facility-parcel-expected-projection/internal/models"
	"facility-parcel-expected-projection/internal/repository"
	"facility-parcel-expected-projection/internal/utils"

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
	// 2. get number_of_parcels on transport from vehicle_transport collection
	// 3. check if facility exists in facility_expected_parcel
	//			- yes: update facility_expected_parcel doc with counter+=number_of_parcels, status=(compute it)
	//			- no: insert new facility_expected_parcel doc with counter=number_of_parcels, status=(compute it)

	facility, err := internal.Repo.GetFacilityDocumentByID(ctx, a.logisticEvent.FacilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility document, facility_id: %s", a.logisticEvent.FacilityID)
		return err
	}

	vehicleTransport, err := internal.Repo.GetVehicleTransportDocumentByTransportID(ctx, a.logisticEvent.TransportID)
	if err != nil {
		log.Err(err).Msgf("cannot get vehicle transport document, trasnport_id: %s, vehicle_id: %s", a.logisticEvent.TransportID, a.logisticEvent.VehicleLicensePlate)
		return err
	}

	facilityExpectedParcel, err := internal.Repo.GetFacilityExpectedParcelDocumentByID(ctx, a.logisticEvent.FacilityID)
	if err != nil {
		// insert new doc
		counter := len(vehicleTransport.Parcels)
		status := utils.ComputeStatus(facility.Capacity, counter)
		newFacilityExpectedParcel := repository.FacilityExpectedParcel{
			FacilityID: a.logisticEvent.FacilityID,
			Status:     status,
			Counter:    counter,
		}
		err = internal.Repo.InsertNewDocument(ctx, repository.FacilityExpectedParcelCollection, newFacilityExpectedParcel)
		if err != nil {
			log.Err(err).Msgf("cannot insert facility_expected_parcel document, doc: %v", newFacilityExpectedParcel)
			return err
		}
		return nil
	}
	// update doc
	counter := facilityExpectedParcel.Counter + len(vehicleTransport.Parcels)
	status := utils.ComputeStatus(facility.Capacity, counter)
	update := bson.M{"$set": bson.M{"status": status, "counter": counter}}
	err = internal.Repo.UpdateDocumentFields(ctx, repository.FacilityExpectedParcelCollection, repository.FacilityExpectedParcelCollectionKey, a.logisticEvent.FacilityID, update)
	if err != nil {
		log.Err(err).Msgf("cannot update facility_parcel document, facility_id: %s, update: %s", a.logisticEvent.FacilityID, update)
		return err
	}
	return nil
}
