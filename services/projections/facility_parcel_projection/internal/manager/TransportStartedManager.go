package manager

import (
	"context"
	"facility-parcel-projection/internal"
	"facility-parcel-projection/internal/models"
	"facility-parcel-projection/internal/repository"
	"fmt"

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
	// 1. get all parcels on the transport.
	// 2. For each parcel:
	// 		2.1. check if parcel is in facility
	//			- yes: remove parcel from facility
	//			- no: return error
	vehicleTransport, err := internal.Repo.GetVehicleTransportDocumentByTransportID(ctx, a.logisticEvent.TransportID)
	if err != nil {
		log.Err(err).Msgf("cannot get vehicle transport document, trasnport_id: %s, vehicle_id: %s", a.logisticEvent.TransportID, a.logisticEvent.VehicleLicensePlate)
		return err
	}
	errorList := []error{}
	for _, parcelID := range vehicleTransport.Parcels {
		update := bson.M{
			"$pull": bson.M{
				"parcels": bson.M{"parcel_id": parcelID},
			},
		}
		err := internal.Repo.DeleteDocumentFields(ctx, repository.FacilityParcelCollection, repository.FacilityParcelCollectionKey, a.logisticEvent.SourceFacilityID, update)
		if err != nil {
			errorList = append(errorList, err)
			log.Err(err).Msgf("cannot delete parcel from facility_parcel, parcel_id: %s, facility_id: %s", parcelID, a.logisticEvent.SourceFacilityID)
		}
	}
	if len(errorList) > 0 {
		return fmt.Errorf("%s", errorList)
	}
	return nil
}
