package repository

import (
	"context"
	"facility_queryhandler/internal/models"
)

type IDatabase interface {
	SelectNetworkNodes(c context.Context, filter string) (models.Network, error)
	GetFacilityExpectedParcelDocumentByID(ctx context.Context, facilityID string) (FacilityExpectedParcel, error)
	CountSortingMachineProcessedItemsByFacilityIDAndDate(ctx context.Context, facilityID string, date string) (int, error)
	ComputeFacilityNextCutoffTimeByFacilityIDAndTime(ctx context.Context, facilityID string, currentTime string) (string, error)
	GetFacilityParcelDocumentByID(ctx context.Context, facilityID string) (FacilityParcel, error)
	GetFacilityDocumentByID(ctx context.Context, facilityID string) (Facility, error)
	ComputeFacilityAverageProcessingTimeByFacilityID(ctx context.Context, facilityID string) (float64, error)
	GetSortingMachineDocumentByFacilityID(ctx context.Context, facilityID string) (SortingMachine, error)
	CountFacilityVehicleByFacilityIDAndVehicleStatus(ctx context.Context, facilityID string, vehicleStatus string) (int, error)
	GetFacilityVehicleDocByFacilityID(ctx context.Context, facilityID string) (FacilityVehicle, error)
}
