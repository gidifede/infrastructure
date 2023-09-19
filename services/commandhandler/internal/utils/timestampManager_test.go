package utils

import (
	"testing"
	"time"
)

func TestTimestampManager_StringToTimestamp(t *testing.T) {
	dateString1 := "2006-01-02T15:04:05.999+01:00"
	t1, _ := time.Parse(time.RFC3339Nano, dateString1)
	expTimestamp1 := t1.UnixNano()

	dateString2 := "2006-01-02T15:04:05.990Z"
	t2, _ := time.Parse(time.RFC3339Nano, dateString2)
	expTimestamp2 := t2.UnixNano()

	type args struct {
		dateStr string
	}
	tests := []struct {
		name    string
		tm      *TimestampManager
		args    args
		want    int64
		wantErr bool
	}{
		{name: "Parsing OK.", tm: &TimestampManager{}, args: args{dateStr: dateString1}, want: expTimestamp1, wantErr: false},
		{name: "Parsing OK.", tm: &TimestampManager{}, args: args{dateStr: dateString2}, want: expTimestamp2, wantErr: false},
		{name: "Parsing Failed.", tm: &TimestampManager{}, args: args{dateStr: "2006/01/02T15:04:05.990Z"}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TimestampManager{}
			got, err := tm.StringToTimestamp(tt.args.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimestampManager.StringToTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("TimestampManager.StringToTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestampManager_TimestampToString(t *testing.T) {
	now := time.Now().UTC()

	type args struct {
		timestampNano int64
	}
	tests := []struct {
		name string
		tm   *TimestampManager
		args args
		want string
	}{
		{name: "Conversion OK.", tm: &TimestampManager{}, args: args{timestampNano: 1136214245999000000}, want: "2006-01-02T15:04:05.999Z"},
		{name: "Conversion OK.", tm: &TimestampManager{}, args: args{timestampNano: now.UnixNano()}, want: now.Format(time.RFC3339Nano)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TimestampManager{}
			if got := tm.TimestampToString(tt.args.timestampNano); got != tt.want {
				t.Errorf("TimestampManager.TimestampToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
