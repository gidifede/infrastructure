package model

type Network struct {
	Items []Item `bson:"items"`
}

type Item struct {
	Node      Node   `bson:"node"`
	Timestamp string `bson:"timestamp"`
}

type Node struct {
	NodeID      string      `bson:"nodeId"`
	NodeType    string      `bson:"nodeType"`
	NodeCompany string      `bson:"nodeCompany"`
	Address     Address     `bson:"address"`
	Coordinates Coordinates `bson:"coordinates"`
}

type Coordinates struct {
	Latitude  string `bson:"latitude"`
	Longitude string `bson:"longitude"`
}

type Address struct {
	Address string `bson:"address"`
	ZipCode string `bson:"zipcode"`
	City    string `bson:"city"`
	Nation  string `bson:"nation"`
}
