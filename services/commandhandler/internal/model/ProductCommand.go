package model

type ProductCommand struct {
	Product Product `json:"product"`
}

type Product struct {
	ID string `json:"id"`
}
