package service

import (
	"context"
	"facility_queryhandler/internal"
	"facility_queryhandler/internal/models"
	"facility_queryhandler/internal/repository"
	"time"

	"github.com/rs/zerolog/log"
)

func GetNetwork(ctx context.Context) (models.Network, error) {
	filter := ""
	res, err := internal.Repo.SelectNetworkNodes(ctx, filter)
	if err != nil {
		return models.Network{}, err
	}
	return res, nil
}

func GetFacilityShortStats(ctx context.Context, facilityID string) (models.FacilityShortStats, error) {
	facilityShortStats := models.FacilityShortStats{}
	facilityExpParcelDoc, err := internal.Repo.GetFacilityExpectedParcelDocumentByID(ctx, facilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility short stats. Facility %s not found in %s", facilityID, repository.FacilityExpectedParcelCollection)
		return facilityShortStats, err
	}
	facilityShortStats.FacilityHealth = facilityExpParcelDoc.Status
	facilityShortStats.ParcelWaiting = facilityExpParcelDoc.Counter

	// Compute daily processed parcel
	layout := "2006-01-02"
	currentDate := time.Now().UTC().Format(layout)
	facilityShortStats.ParcelProcessed, err = internal.Repo.CountSortingMachineProcessedItemsByFacilityIDAndDate(ctx, facilityID, currentDate)
	if err != nil {
		log.Err(err).Msgf("cannot get facility short stats. Cannot compute daily processed parcel, facility %s, date %s", facilityID, currentDate)
		return facilityShortStats, err
	}

	// Compute next cutoff time
	t := time.Unix(time.Now().UTC().Unix(), 0)
	currentTime := t.Format("15:04")
	facilityShortStats.NextCutOffTime, err = internal.Repo.ComputeFacilityNextCutoffTimeByFacilityIDAndTime(ctx, facilityID, currentTime)
	if err != nil {
		log.Err(err).Msgf("cannot get facility short stats. Cannot compute next cutoff time, facility %s, time %s", facilityID, currentTime)
		return facilityShortStats, err
	}

	return facilityShortStats, nil
}

func GetFacilityParcelDetails(ctx context.Context, facilityID string) ([]models.FacilityParcelDetails, error) {
	facilityParcelDetails := []models.FacilityParcelDetails{}
	facilityParcelDoc, err := internal.Repo.GetFacilityParcelDocumentByID(ctx, facilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility parcel details. Facility %s not found in %s", facilityID, repository.FacilityParcelCollection)
		return facilityParcelDetails, err
	}
	for _, item := range facilityParcelDoc.Parcels {
		facilityParcelDetails = append(facilityParcelDetails,
			models.FacilityParcelDetails{
				ParcelID:    item.ParcelID,
				TimeIn:      item.ArrivingTime,
				TimeOut:     item.ExitTime,
				NextHop:     item.NextHop,
				Destination: item.DeliveryDestination,
				Status:      item.Status,
			},
		)
	}
	return facilityParcelDetails, nil
}

func GetFacilityParcelStats(ctx context.Context, facilityID string) (models.FacilityParcelStats, error) {
	// Get facility capacity
	facilityParcelStats := models.FacilityParcelStats{}
	facilityDoc, err := internal.Repo.GetFacilityDocumentByID(ctx, facilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility parcel stats. Facility %s not found in %s", facilityID, repository.FacilityCollection)
		return facilityParcelStats, err
	}
	facilityParcelStats.Capacity = facilityDoc.Capacity

	// Get facility expected parcels
	facilityExpParcelDoc, err := internal.Repo.GetFacilityExpectedParcelDocumentByID(ctx, facilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility parcel stats. Facility %s not found in %s", facilityID, repository.FacilityExpectedParcelCollection)
		return facilityParcelStats, err
	}
	facilityParcelStats.ParcelWaiting = facilityExpParcelDoc.Counter

	// Compute facility daily processed parcels
	layout := "2006-01-02"
	currentDate := time.Now().UTC().Format(layout)
	facilityParcelStats.ParcelProcessed, err = internal.Repo.CountSortingMachineProcessedItemsByFacilityIDAndDate(ctx, facilityID, currentDate)
	if err != nil {
		log.Err(err).Msgf("cannot get facility parcel stats. Cannot compute daily processed parcel, facility %s, date %s", facilityID, currentDate)
		return facilityParcelStats, err
	}

	// Compute facility daily avg processing time
	avgTime, err := internal.Repo.ComputeFacilityAverageProcessingTimeByFacilityID(ctx, facilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility parcel stats. Cannot compute current average processing time, facility %s", facilityID)
		return facilityParcelStats, err
	}
	facilityParcelStats.AvgProcessingTime = avgTime
	return facilityParcelStats, nil
}

func GetFacilitySortingMachineStats(ctx context.Context, facilityID string) (models.FacilitySortingMachineStats, error) {
	// Get sorting machine caoacity and next maintenance
	facilitySMStats := models.FacilitySortingMachineStats{}
	sortingMachineDoc, err := internal.Repo.GetSortingMachineDocumentByFacilityID(ctx, facilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get facility sorting machine stats. Facility %s not found in %s", facilityID, repository.SortingMachineCollection)
		return facilitySMStats, err
	}
	facilitySMStats.Capacity = sortingMachineDoc.Capacity
	facilitySMStats.NextMaintenance = sortingMachineDoc.NextMaintenance

	// compute daily average working capacity
	days := 0
	dailyCapacitySum := 0
	for _, itemDate := range sortingMachineDoc.ItemProcessingRates {
		days++
		for _, item := range itemDate.HourlyRates {
			dailyCapacitySum += item.Rate
		}
	}
	facilitySMStats.WorkingCapacityAverage = float64(dailyCapacitySum) / float64(days)

	// compute hourly parcel processed during last N hours
	parcelProcessed := []models.ParcelProcessedItem{}
	// Suppose we return last N=10 hours
	timestampNow := time.Now().UTC()
	for i := 0; i < 10; i++ {
		currentDate := timestampNow.Format("2006-01-02")
		currentHour := timestampNow.Hour()
		numParcels := 0
		for _, itemDate := range sortingMachineDoc.ItemProcessingRates {
			if itemDate.Date == currentDate {
				for _, itemHour := range itemDate.HourlyRates {
					if itemHour.Hour == currentHour {
						numParcels = itemHour.Rate
						break
					}
				}
				break
			}
		}
		parcelProcessed = append(parcelProcessed, models.ParcelProcessedItem{
			Day:     currentDate,
			Hour:    currentHour,
			Parcels: numParcels,
		})
		timestampNow = timestampNow.Add(-time.Duration(1 * time.Hour))
	}
	facilitySMStats.ParcelProcessed = parcelProcessed
	return facilitySMStats, nil
}

func GetFacilityVehicleStats(ctx context.Context, facilityID string) (models.FacilityVehicleStats, error) {
	facilityVehicleStats := models.FacilityVehicleStats{}
	var err error
	facilityVehicleStats.VehiclesLoading, err = internal.Repo.CountFacilityVehicleByFacilityIDAndVehicleStatus(ctx, facilityID, repository.FacilityVehicleStatusLoading)
	if err != nil {
		log.Err(err).Msgf("cannot count vehicles with status %s in facility %s", repository.FacilityVehicleStatusLoading, facilityID)
		return facilityVehicleStats, err
	}
	facilityVehicleStats.VehiclesUnloading, err = internal.Repo.CountFacilityVehicleByFacilityIDAndVehicleStatus(ctx, facilityID, repository.FacilityVehicleStatusUnloading)
	if err != nil {
		log.Err(err).Msgf("cannot count vehicles with status %s in facility %s", repository.FacilityVehicleStatusUnloading, facilityID)
		return facilityVehicleStats, err
	}
	return facilityVehicleStats, nil
}

func GetFacilityVehicleDetails(ctx context.Context, facilityID string) ([]models.FacilityVehicleDetails, error) {
	facilityVehicleDetails := []models.FacilityVehicleDetails{}
	facilityVehicleDoc, err := internal.Repo.GetFacilityVehicleDocByFacilityID(ctx, facilityID)
	if err != nil {
		log.Err(err).Msgf("cannot get vehicles in facility %s from doc %s", facilityID, repository.FacilityVehicleCollection)
		return facilityVehicleDetails, err
	}
	for _, item := range facilityVehicleDoc.Vehicles {
		facilityVehicleDetails = append(facilityVehicleDetails, models.FacilityVehicleDetails{
			VehicleLicencePlate: item.VehicleID,
			Status:              item.Status,
			ArrivedTime:         item.ArrivedTime,
			// NextStartTime:       ??, // fake for now
		})
	}
	return facilityVehicleDetails, nil
}
