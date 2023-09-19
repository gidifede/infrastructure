package repository

import (
	"context"
	"facility_queryhandler/internal/models"
	"fmt"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	DB mongo.Database
}

func NewMongoDB(DB mongo.Database) IDatabase {
	return &MongoDB{
		DB: DB,
	}
}

func (m *MongoDB) SelectNetworkNodes(c context.Context, filter string) (models.Network, error) {

	var results []models.Network

	cursor, err := m.DB.Collection(NetworkCollection).Find(c, bson.M{})
	for cursor.Next(c) {
		var n models.Network
		if err := cursor.Decode(&n); err != nil {
			log.Err(err).Msg("")
			return models.Network{}, err
		}
		results = append(results, n)
	}
	if err != nil {
		log.Err(err).Msg("")
		return models.Network{}, err
	}
	defer cursor.Close(c)

	return results[0], nil
}

func (m *MongoDB) GetFacilityExpectedParcelDocumentByID(ctx context.Context, facilityID string) (FacilityExpectedParcel, error) {
	collection := m.DB.Collection(FacilityExpectedParcelCollection)
	filter := bson.M{FacilityExpectedParcelCollectionKey: facilityID}
	var doc FacilityExpectedParcel
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return FacilityExpectedParcel{}, err
	}
	return doc, nil
}

func (m *MongoDB) CountSortingMachineProcessedItemsByFacilityIDAndDate(ctx context.Context, facilityID string, date string) (int, error) {
	collection := m.DB.Collection(SortingMachineCollection)
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"facility_id": facilityID,
			},
		},
		{
			"$unwind": "$item_processing_rates",
		},
		{
			"$match": bson.M{
				"item_processing_rates.date": date,
			},
		},
		{
			"$unwind": "$item_processing_rates.hourly_rates",
		},
		{
			"$group": bson.M{
				"_id": nil,
				"totalRate": bson.M{
					"$sum": "$item_processing_rates.hourly_rates.rate",
				},
			},
		},
		{
			"$project": bson.M{
				"_id":       0,
				"totalRate": 1,
			},
		},
	}

	// Execute the aggregation query
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Err(err).Msgf("cannot execute query to sum hourly rates")
		return 0, err
	}
	defer cursor.Close(ctx)

	// Iterate over the results
	var result struct {
		TotalRate int `bson:"totalRate"`
	}
	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			log.Err(err).Msgf("cannot decode hourly rate sum result")
			return 0, err
		}
		return result.TotalRate, nil
	}
	return 0, nil
}

func (m *MongoDB) ComputeFacilityNextCutoffTimeByFacilityIDAndTime(ctx context.Context, facilityID string, currentTime string) (string, error) {
	collection := m.DB.Collection(RouteCollection)
	// Define a structure to hold the cutoff times and their differences from the target
	type CutoffWithDiff struct {
		Cutoff string
		Diff   time.Duration
	}

	// Find the document with the specified source_facility_id
	filter := bson.M{"source_facility_id": facilityID}
	var doc struct {
		Cutoff []string `bson:"cutoff"`
	}

	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		log.Err(err).Msgf("cannot execute query to get next cutoff time")
		return "", err
	}

	// Calculate the time differences and store them in a slice
	var cutoffWithDiffs []CutoffWithDiff
	targetTime, err := time.Parse("15:04", currentTime)
	if err != nil {
		log.Err(err).Msgf("cannot parse given time %s", currentTime)
		return "", err
	}

	for _, cutoff := range doc.Cutoff {
		cutoffTime, err := time.Parse("15:04", cutoff)
		if err != nil {
			log.Err(err).Msgf("cannot parse cutoff time %s", cutoff)
			return "", err
		}

		diff := cutoffTime.Sub(targetTime)
		cutoffWithDiffs = append(cutoffWithDiffs, CutoffWithDiff{Cutoff: cutoff, Diff: diff})
	}

	// Sort the slice by time difference
	sort.Slice(cutoffWithDiffs, func(i, j int) bool {
		return cutoffWithDiffs[i].Diff < cutoffWithDiffs[j].Diff
	})

	// Find the first positive element
	var firstPositiveElement CutoffWithDiff
	found := false
	for _, item := range cutoffWithDiffs {
		if item.Diff > 0 {
			firstPositiveElement = item
			found = true
			break
		}
	}
	if found {
		return firstPositiveElement.Cutoff, nil
	}
	errorMsg := fmt.Sprintf("cannot find a cutoff greather than given time %s", currentTime)
	log.Error().Msgf(errorMsg)
	return "", fmt.Errorf(errorMsg)
}

func (m *MongoDB) GetFacilityParcelDocumentByID(ctx context.Context, facilityID string) (FacilityParcel, error) {
	collection := m.DB.Collection(FacilityParcelCollection)
	filter := bson.M{FacilityParcelCollectionKey: facilityID}
	var doc FacilityParcel
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return FacilityParcel{}, err
	}
	return doc, nil
}

func (m *MongoDB) GetFacilityDocumentByID(ctx context.Context, facilityID string) (Facility, error) {
	collection := m.DB.Collection(FacilityCollection)
	filter := bson.M{FacilityCollectionKey: facilityID}
	var doc Facility
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return Facility{}, err
	}
	return doc, nil
}

func (m *MongoDB) ComputeFacilityAverageProcessingTimeByFacilityID(ctx context.Context, facilityID string) (float64, error) {
	collection := m.DB.Collection(FacilityParcelCollection)
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"facility_id": facilityID,
			},
		},
		{
			"$unwind": "$parcels",
		},
		{
			"$group": bson.M{
				"_id": nil,
				"averageDifference": bson.M{
					"$avg": bson.M{
						"$subtract": []interface{}{"$parcels.exit_time", "$parcels.arriving_time"},
					},
				},
			},
		},
	}

	// Execute the aggregation query
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Err(err).Msgf("cannot execute query to compute avg processing time")
		return 0, err
	}
	defer cursor.Close(ctx)

	// Read the result
	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		log.Err(err).Msgf("cannot decode avg processing time result")
		return 0, err
	}

	// the average difference (in milliseconds)
	if len(result) > 0 {
		averageDifference := result[0]["averageDifference"].(float64)
		// return the average difference (in seconds)
		return averageDifference / 1000, nil
	}
	return 0, nil
}

func (m *MongoDB) GetSortingMachineDocumentByFacilityID(ctx context.Context, facilityID string) (SortingMachine, error) {
	collection := m.DB.Collection(SortingMachineCollection)
	filter := bson.M{SortingMachineCollectionFacilityKey: facilityID}
	var doc SortingMachine
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return SortingMachine{}, err
	}
	return doc, nil
}

func (m *MongoDB) CountFacilityVehicleByFacilityIDAndVehicleStatus(ctx context.Context, facilityID string, vehicleStatus string) (int, error) {
	collection := m.DB.Collection(FacilityVehicleCollection)
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"facility_id": facilityID,
			},
		},
		{
			"$unwind": "$vehicles",
		},
		{
			"$match": bson.M{
				"vehicles.status": vehicleStatus,
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	// Execute the aggregation query
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Err(err).Msgf("cannot count vehicles with status %s in facility %s", vehicleStatus, facilityID)
		return 0, err
	}

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		log.Err(err).Msgf("cannot decode result of counting vehicles with status %s in facility %s", vehicleStatus, facilityID)
		return 0, err
	}

	if len(result) > 0 {
		count := result[0]["count"].(int32)
		return int(count), nil
	}
	return 0, nil
}

func (m *MongoDB) GetFacilityVehicleDocByFacilityID(ctx context.Context, facilityID string) (FacilityVehicle, error) {
	collection := m.DB.Collection(FacilityVehicleCollection)
	filter := bson.M{FacilityVehicleCollectionKey: facilityID}
	var doc FacilityVehicle
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return FacilityVehicle{}, err
	}
	return doc, nil
}
