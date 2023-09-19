package db

type Product struct {
	Name string `json:"name" bson:"name"`
	SLA  int    `json:"SLA" bson:"SLA"`
}