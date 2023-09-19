package utils

//go:generate mockgen -source=./ITimestampManager.go -destination=./mock/ITimestampManagerMock.go -package=mock

type ITimestampManager interface {
	GenerateTimestamp() int64
	StringToTimestamp(string) (int64, error)
	TimestampToString(int64) string
}
