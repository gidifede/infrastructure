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

type ParcelUnloadedManager struct {
	logisticEvent models.ParcelUnloaded
}

func newParcelUnloadedManager(logisticEvent models.ParcelUnloaded) IEventManager {
	return &ParcelUnloadedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *ParcelUnloadedManager) ManageEvent(ctx context.Context) error {
	// 1. check if facility exists in facility collection
	//		- yes: go to 2
	//		- no: return error
	// 2. check if parcel exists in parcel collection && get parcel destination
	//		- yes: go to 3
	//		- no: return error
	// 3. check if facility exists in facility_parcel
	//		- yes: go to 4
	//		- no: insert facility_parcel doc with parcel && return
	// 4. check if parcel exists in facility of facility_parcel
	// 		- yes: return error
	//		- no: update facility_parcel doc with parcel

	_, err := internal.Repo.GetFacilityDocumentByID(ctx, a.logisticEvent.FacilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility document, facility_id: %s", a.logisticEvent.FacilityID)
		return err
	}
	parcel, err := internal.Repo.GetParcelDocumentByID(ctx, a.logisticEvent.ParcelID)
	if err != nil {
		log.Err(err).Msgf("cannot get parcel document, parcel_id: %s", a.logisticEvent.ParcelID)
		return err
	}
	readFacilityParcel, err := internal.Repo.GetFacilityParcelDocumentByID(ctx, a.logisticEvent.FacilityID)
	if err != nil {
		// insert new doc
		facilityParcel := repository.FacilityParcel{
			FacilityID: a.logisticEvent.FacilityID,
			Parcels: []repository.ParcelInFacility{
				{
					ParcelID:     a.logisticEvent.ParcelID,
					ArrivingTime: a.logisticEvent.Timestamp,
					//ExitTime:             unset,
					//NextHop:              unset,
					DeliveryDestination: parcel.Receiver.Address,
					Status:              "sorting",
				},
			},
		}
		err = internal.Repo.InsertNewDocument(ctx, repository.FacilityParcelCollection, facilityParcel)
		if err != nil {
			log.Err(err).Msgf("cannot insert facility_parcel document, doc: %s", facilityParcel)
			return err
		}
		return nil
	}
	// check if parcel already exists
	exists := false
	for _, parcel := range readFacilityParcel.Parcels {
		if parcel.ParcelID == a.logisticEvent.ParcelID {
			exists = true
		}
	}
	if exists {
		errMsg := fmt.Sprintf("parcel already exists in facility, facility_id: %s, parcel_id: %s", a.logisticEvent.FacilityID, a.logisticEvent.ParcelID)
		log.Error().Msgf(errMsg)
		return fmt.Errorf(errMsg)
	}
	parcelInFacility := repository.ParcelInFacility{
		ParcelID:     a.logisticEvent.ParcelID,
		ArrivingTime: a.logisticEvent.Timestamp,
		//ExitTime:             unset,
		//NextHop:              unset,
		DeliveryDestination: parcel.Receiver.Address,
		Status:              "sorting",
	}
	update := bson.M{"$addToSet": bson.M{"parcels": parcelInFacility}}
	err = internal.Repo.UpdateDocumentFields(ctx, repository.FacilityParcelCollection, repository.FacilityParcelCollectionKey, a.logisticEvent.FacilityID, update)
	if err != nil {
		log.Err(err).Msgf("cannot update facility_parcel document, facility_id: %s, update: %s", a.logisticEvent.FacilityID, update)
		return err
	}
	return nil
}
