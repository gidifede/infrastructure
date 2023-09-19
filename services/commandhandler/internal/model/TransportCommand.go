package model

type TransportCommand struct {
	Transport Transport `json:"transport"`
}

type Transport struct {
	ID string `json:"id"`
}
