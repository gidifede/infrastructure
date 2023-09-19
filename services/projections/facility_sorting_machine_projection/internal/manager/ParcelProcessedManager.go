package manager

import (
	"context"
	"facility-sorting-machine-projection/internal"
	"facility-sorting-machine-projection/internal/models"
	"facility-sorting-machine-projection/internal/repository"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ParcelProcessedManager struct {
	logisticEvent models.ParcelProcessed
}

func newParcelProcessedManager(logisticEvent models.ParcelProcessed) IEventManager {
	return &ParcelProcessedManager{
		logisticEvent: logisticEvent,
	}
}

func (a *ParcelProcessedManager) ManageEvent(ctx context.Context) error {
	// 1. check if sorting machine exists:
	//		- no: return error
	//		- yes: go to 2
	// 2. extract date and hour by timestamp of the event
	// 3. check if an itemProcessingRate record exists in sorting_machine collection for the date
	// 		- yes: update the itemProcessingRate record in sorting_machine collection for the extracted date and hour & return
	// 		- no: go to 4
	// 4. (TODO) compute nextMaintenance based on last n-days itemProcessingRate records
	// 5. update the sorting_machine collection with a new itemProcessingRate record, for the extracted date and hour, and (TODO) update nextMaintenance

	sortingMachine, err := internal.Repo.GetSortingMachineDocumentByID(ctx, a.logisticEvent.SortingMachineID)
	if err != nil {
		log.Err(err).Msgf("cannot get sorting machine document, sorting_machine_id: %s", a.logisticEvent.SortingMachineID)
		return err
	}
	// year, month, day := a.logisticEvent.Timestamp.Date()
	eventHour := a.logisticEvent.Timestamp.Hour()
	layout := "2006-01-02"
	eventDate := a.logisticEvent.Timestamp.Format(layout)
	foundDate := false
	foundHour := false
	for _, itemProcessingRate := range sortingMachine.ItemProcessingRates {
		if itemProcessingRate.Date == eventDate {
			foundDate = true
			for _, hour := range itemProcessingRate.HourlyRates {
				if hour.Hour == eventHour {
					foundHour = true
					break
				}
			}
			break
		}
	}
	if foundDate {
		var update bson.M
		var filters bson.M
		var opts *options.UpdateOptions
		if foundHour {
			filters = bson.M{
				"sorting_machine_id": a.logisticEvent.SortingMachineID,
			}
			update = bson.M{
				"$inc": bson.M{
					"item_processing_rates.$[elem].hourly_rates.$[elem2].rate": 1,
				},
			}
			arrayFilters := options.ArrayFilters{
				Filters: []interface{}{
					bson.M{"elem.date": eventDate},
					bson.M{"elem2.hour": eventHour},
				},
			}
			opts = options.Update().SetArrayFilters(arrayFilters)
		} else {
			filters = bson.M{
				"sorting_machine_id":         a.logisticEvent.SortingMachineID,
				"item_processing_rates.date": eventDate,
			}
			newHourlyRate := repository.HourlyRate{
				Hour: eventHour,
				Rate: 1,
			}
			update = bson.M{
				"$push": bson.M{
					"item_processing_rates.$.hourly_rates": newHourlyRate,
				},
			}
		}
		err = internal.Repo.UpdateDocument(ctx, repository.SortingMachineCollection, filters, update, opts)
		if err != nil {
			log.Err(err).Msgf("cannot update sorting_machine doc, update %s", update)
			return err
		}
		return nil
	}

	// compute maintenance
	// ---------------TO-DO: compute next maintenance
	filters := bson.M{
		"sorting_machine_id": a.logisticEvent.SortingMachineID,
	}
	newItemProcessingRate := repository.ItemProcessingRate{
		Date: eventDate,
		HourlyRates: []repository.HourlyRate{
			{
				Hour: eventHour,
				Rate: 1,
			},
		},
	}

	update := bson.M{
		"$push": bson.M{
			"item_processing_rates": newItemProcessingRate,
		},
		// "$set": bson.M{
		// 	"next_maintenance": nextMaintenance,
		// },
	}
	var opts *options.UpdateOptions
	err = internal.Repo.UpdateDocument(ctx, repository.SortingMachineCollection, filters, update, opts)
	if err != nil {
		log.Err(err).Msgf("cannot update sorting_machine doc, update %s", update)
		return err
	}
	return nil
}
