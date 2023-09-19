package db

type Vehicle struct {
	VehicleID    string `json:"vehicle_id" bson:"vehicle_id"`
	Type         string `json:"type" bson:"type"`
	Capacity     int    `json:"capacity" bson:"capacity"`
	LicensePlate string `json:"license_plate" bson:"license_plate"`
}
