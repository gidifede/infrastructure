package model

type RouteSetup []RouteItem

type CutoffTime struct {
	Timestamp []string `json:"timestamp" validate:"required,datetime"`
}

type RouteItem struct {
	SourceNodeID      string     `json:"source_node_id" validate:"required"`
	DestinationNodeID string     `json:"destination_node_id" validate:"required"`
	CutoffTime        CutoffTime `json:"cutoff_time" validate:"required"`
	Timestamp         string     `json:"timestamp" validate:"required,datetime"`
}