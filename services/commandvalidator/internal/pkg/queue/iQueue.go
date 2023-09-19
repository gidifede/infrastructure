package queue

//go:generate mockgen -source=./iQueue.go -destination=./IQueueMock.go -package=queue

import (
	"context"
)

type IQeueu interface {
	SendMsg(ctx context.Context, queueURL string, body string, nmessageGroup string) error
}
