package statemachineconfig

import (
	"command-handler/internal/utils"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3StateMachineConfig struct {
	connection s3iface.S3API
	bucket     string
	key        string
}

func NewS3StateMachineConfig(conn s3iface.S3API, bucket string, key string) IStateMachineConfig {
	return &S3StateMachineConfig{connection: conn, bucket: bucket, key: key}
}

func (s *S3StateMachineConfig) Get(ctx context.Context) ([]byte, error) {
	utils.AddClassAndMethodToMDC(s)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
	}
	result, err := s.connection.GetObjectWithContext(ctx, input)
	if err != nil {

		return nil, fmt.Errorf("Get state machine config failed. %v", err)
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("read state machine config failed. %v", err)
	}
	return body, nil
}
