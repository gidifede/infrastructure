package utils

import (
	"time"
)

type TimestampManager struct{}

func NewTimestampManager() ITimestampManager {
	return &TimestampManager{}
}

func (tm *TimestampManager) GenerateTimestamp() int64 {
	return time.Now().UnixMilli()
}

func (tm *TimestampManager) StringToTimestamp(dateStr string) (int64, error) {
	t, err := time.Parse(time.RFC3339Nano, dateStr)
	return t.UnixNano(), err
}

func (tm *TimestampManager) TimestampToString(timestampNano int64) string {
	unixTimeUTC := time.Unix(0, timestampNano).UTC()         //gives unix time stamp in utc
	unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339Nano) // converts utc time to RFC3339 format
	return unitTimeInRFC3339
}
