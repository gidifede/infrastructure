package eventstore

type Event struct {
	AggregateID   string `json:"aggregate_id"`
	Type          string `json:"type"`
	Data          string `json:"data"`
	Timestamp     int64  `json:"timestamp"`
	Source        string `json:"source"`
	EventID       string `json:"eventId"`
	TimestampSent int64  `json:"timestampSent"`
}
