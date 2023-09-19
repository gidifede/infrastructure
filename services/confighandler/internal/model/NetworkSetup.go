package model

type NetworkSetup struct {
	Items []NetworkItem `bson:"items"`
}

type NetworkItem struct {
	Node            Node              `bson:"node"`
	ConnectTo       []ConnectTo       `bson:"connectTo"`
	SortingMachines []SortingMachines `bson:"sortingMachines"`
	Timestamp       string            `bson:"timestamp"`
}

type Node struct {
	NodeID      string      `bson:"nodeId"`
	NodeType    string      `bson:"nodeType"`
	NodeCompany string      `bson:"nodeCompany"`
	Address     Address     `bson:"address"`
	Coordinates Coordinates `bson:"coordinates"`
}

type Address struct {
	Address string `bson:"address"`
	ZipCode string `bson:"zipcode"`
	City    string `bson:"city"`
	Nation  string `bson:"nation"`
}

type Coordinates struct {
	Latitude  string `bson:"latitude"`
	Longitude string `bson:"longitude"`
}

type ConnectTo struct {
	NodeId   string `bson:"nodeId"`
	Distance int    `bson:"distance"`
}

type SortingMachines struct {
	Serial   string `bson:"serial"`
	Capacity int    `bson:"capacity"`
}
