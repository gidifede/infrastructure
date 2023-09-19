package manager

import (
	"context"
	"facility-parcel-projection/internal"
	"facility-parcel-projection/internal/models"
	"facility-parcel-projection/internal/repository"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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
	// 1. check if parcel is in facility
	//		- yes: update next hop, status=processed, exit_time (retrieve route to compute exit_time)
	//		- no: return error

	// Get route
	route, err := internal.Repo.GetRouteDocumentBySourceAndDest(ctx, a.logisticEvent.FacilityID, a.logisticEvent.DestinationFacilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get route, source_facility_id:%s, dest_facility_id: %s", a.logisticEvent.FacilityID, a.logisticEvent.DestinationFacilityID)
		return err
	}

	// Compute exit_time = nearest cut off time compared to event timestamp
	todayAndTomorrowCutoffTimestamps := []time.Time{}
	for _, cutoff := range route.Cutoff {
		cutoffSplitted := strings.Split(cutoff, ":")
		if len(cutoffSplitted) != 2 {
			return fmt.Errorf("invalid cutoff time, cutoff: %s, source: %s, dest: %s", cutoff, a.logisticEvent.FacilityID, a.logisticEvent.DestinationFacilityID)
		}
		cutoffHour, err := strconv.ParseInt(cutoffSplitted[0], 10, 0)
		if err != nil {
			return fmt.Errorf("invalid cutoff time (invalid hour), cutoff: %s, source: %s, dest: %s", cutoff, a.logisticEvent.FacilityID, a.logisticEvent.DestinationFacilityID)
		}
		cutoffMinute, err := strconv.ParseInt(cutoffSplitted[1], 10, 0)
		if err != nil {
			return fmt.Errorf("invalid cutoff time (invalid minute), cutoff: %s, source: %s, dest: %s", cutoff, a.logisticEvent.FacilityID, a.logisticEvent.DestinationFacilityID)
		}

		// add today timestamp
		todayTimestamp := a.logisticEvent.Timestamp.UTC().Truncate(24 * time.Hour).Add(time.Duration(cutoffHour) * time.Hour).Add(time.Duration(cutoffMinute) * time.Minute)
		todayAndTomorrowCutoffTimestamps = append(todayAndTomorrowCutoffTimestamps, todayTimestamp)

		// add tomorrow timestamp
		tomorrowTimestamp := a.logisticEvent.Timestamp.UTC().Truncate(24 * time.Hour).Add(time.Duration(cutoffHour+24) * time.Hour).Add(time.Duration(cutoffMinute) * time.Minute)
		todayAndTomorrowCutoffTimestamps = append(todayAndTomorrowCutoffTimestamps, tomorrowTimestamp)
	}
	sort.Slice(todayAndTomorrowCutoffTimestamps, func(i, j int) bool {
		return todayAndTomorrowCutoffTimestamps[i].Before(todayAndTomorrowCutoffTimestamps[j])
	})

	var exitTime time.Time
	for _, cutoffTimestamp := range todayAndTomorrowCutoffTimestamps {
		if cutoffTimestamp.After(a.logisticEvent.Timestamp) {
			exitTime = cutoffTimestamp
			break
		}
	}

	// Update collection
	filters := bson.M{repository.FacilityParcelCollectionKey: a.logisticEvent.FacilityID, "parcels.parcel_id": a.logisticEvent.ParcelID}
	update := bson.M{"$set": bson.M{"parcels.$.status": "processed", "parcels.$.next_hop": a.logisticEvent.DestinationFacilityID, "parcels.$.exit_time": exitTime}}
	err = internal.Repo.UpdateDocument(ctx, repository.FacilityParcelCollection, filters, update)
	if err != nil {
		log.Err(err).Msgf("cannot update facility_parcel, facility_id: %s, parcel_id: %s, update: %s", a.logisticEvent.FacilityID, a.logisticEvent.ParcelID, update)
		return err
	}
	return nil
}
