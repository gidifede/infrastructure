package mockstatemachineconfig

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type MockS3 struct {
	s3iface.S3API
}

func (m MockS3) GetObjectWithContext(ctx context.Context, input *s3.GetObjectInput, options ...request.Option) (*s3.GetObjectOutput, error) {

	fmt.Printf("ctx :%v , request: %v",ctx,options)
	if *input.Key == "keyOK" {
		data, err := LoadStateMachine("../testdata/statemachine/cluster/stateMachine.json")
		if err != nil {
			return nil, err
		}
		output := s3.GetObjectOutput{Body: io.NopCloser(bytes.NewReader(data))}
		return &output, nil
	}

	return nil, fmt.Errorf("an error")
}

func LoadStateMachine(jsonFilePath string) ([]byte, error) {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}
