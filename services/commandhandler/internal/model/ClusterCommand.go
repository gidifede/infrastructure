package model

type ClusterCommand struct {
	Cluster Cluster `json:"cluster"`
}

type Cluster struct {
	ID string `json:"id"`
}
