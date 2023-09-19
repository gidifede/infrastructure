package repository

import "time"

const (
	SortingMachineCollection    = "sorting_machine"
	SortingMachineCollectionKey = "sorting_machine_id"
)

type SortingMachine struct {
	SortingMachineID    string               `bson:"sorting_machine_id"`
	FacilityID          string               `bson:"facility_id"`
	Capacity            int                  `bson:"capacity"`
	ItemProcessingRates []ItemProcessingRate `bson:"item_processing_rates"`
	InsertTimestamp     time.Time            `bson:"insert_timestamp"`
	NextMaintenance     time.Time            `bson:"next_maintenance"`
}

type ItemProcessingRate struct {
	Date        string       `bson:"date"`
	HourlyRates []HourlyRate `bson:"hourly_rates"`
}

type HourlyRate struct {
	Hour int `bson:"hour"`
	Rate int `bson:"rate"`
}
