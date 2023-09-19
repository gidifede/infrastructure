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

type TransportStartedManager struct {
	logisticEvent models.TransportStarted
}

func newTransportStartedManager(logisticEvent models.TransportStarted) IEventManager {
	return &TransportStartedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *TransportStartedManager) ManageEvent(ctx context.Context) error {
	// 1. check if source_facility exists in facility collection
	//		- yes: go to 2
	//		- no: return error
	// 2. get number_of_parcels on transport from vehicle_transport collection
	// 3. check if facility exists in facility_expected_parcel
	//			- yes: update facility_expected_parcel doc with counter-=number_of_parcels, status=(compute it)
	//			- no: do nothing! If facility doesn't exist, maybe it has only accepted parcels. Log it

	sourceFacility, err := internal.Repo.GetFacilityDocumentByID(ctx, a.logisticEvent.SourceFacilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility document, source_facility_id: %s", a.logisticEvent.SourceFacilityID)
		return err
	}

	vehicleTransport, err := internal.Repo.GetVehicleTransportDocumentByTransportID(ctx, a.logisticEvent.TransportID)
	if err != nil {
		log.Err(err).Msgf("cannot get vehicle transport document, trasnport_id: %s, vehicle_id: %s", a.logisticEvent.TransportID, a.logisticEvent.VehicleLicensePlate)
		return err
	}

	facilityExpectedParcel, err := internal.Repo.GetFacilityExpectedParcelDocumentByID(ctx, a.logisticEvent.SourceFacilityID)
	if err != nil {
		log.Debug().Msgf("facility %s hasn't expected parcel (TransportStarted due to parcel accepted)", a.logisticEvent.SourceFacilityID)
		return nil
	}
	// update doc
	counter := facilityExpectedParcel.Counter - len(vehicleTransport.Parcels)
	status := utils.ComputeStatus(sourceFacility.Capacity, counter)
	update := bson.M{"$set": bson.M{"status": status, "counter": counter}}
	err = internal.Repo.UpdateDocumentFields(ctx, repository.FacilityExpectedParcelCollection, repository.FacilityExpectedParcelCollectionKey, a.logisticEvent.SourceFacilityID, update)
	if err != nil {
		log.Err(err).Msgf("cannot update facility_parcel document, facility_id: %s, update: %s", a.logisticEvent.SourceFacilityID, update)
		return err
	}
	return nil
}
