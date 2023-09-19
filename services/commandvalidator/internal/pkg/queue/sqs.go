package queue

import (
	"commandvalidator/internal/utils"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SqsQueue struct {
	session sqs.SQS
}

func NewSqsQueue(session sqs.SQS) IQeueu {
	return &SqsQueue{
		session: session,
	}
}

func (queue *SqsQueue) SendMsg(ctx context.Context, queueName string, body string, messageGroup string) error {
	utils.AddClassAndMethodToMDC(queue)
	result, errGetURL := queue.session.GetQueueUrlWithContext(ctx, &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if errGetURL != nil {
		return errGetURL
	}

	input := &sqs.SendMessageInput{
		MessageGroupId:         &messageGroup,
		MessageBody:            aws.String(body),
		QueueUrl:               result.QueueUrl,
		MessageDeduplicationId: aws.String(fmt.Sprint(time.Now().Local().UnixNano())),
	}

	_, errSendMessage := queue.session.SendMessageWithContext(ctx, input)

	if errSendMessage != nil {
		return errSendMessage
	}

	return nil
}
