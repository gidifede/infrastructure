package models

type Network struct {
	Items []Item `json:"items"`
}

type Item struct {
	Node      Node   `json:"node"`
	Timestamp string `json:"timestamp"`
}

type Node struct {
	NodeID      string      `json:"node_id"`
	NodeType    string      `json:"node_type"`
	NodeCompany string      `json:"node_company"`
	Address     Address     `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"long"`
}

type Address struct {
	Address string `json:"address"`
	ZipCode string `json:"zipcode"`
	City    string `json:"city"`
	Nation  string `json:"nation"`
}
