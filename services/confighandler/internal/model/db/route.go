package db

type Route struct {
	SourceFacilityID string `json:"source_facility_id" bson:"source_facility_id"`
	DestFacilityID string `json:"dest_facility_id" bson:"dest_facility_id"`
	Cutoff []string `json:"cutoff" bson:"cutoff"`
	
}