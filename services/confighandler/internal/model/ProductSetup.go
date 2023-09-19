package model

type ProductSetup struct {
	Product   Product `json:"product" validate:"required"`
	Timestamp string  `json:"timestamp" validate:"required,datetime"`
}

type Product struct {
	Name string `json:"name" validate:"required"`
	SLA  int    `json:"SLA" validate:"required"`
}