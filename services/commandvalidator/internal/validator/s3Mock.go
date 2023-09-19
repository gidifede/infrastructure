package validator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type MockS3 struct {
	s3iface.S3API
}

func (m MockS3) GetObjectWithContext(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	if strings.HasSuffix(*input.Key, "keyOK") {
		data, err := LoadMap("./testData/commandAPIPathMap.json")
		if err != nil {
			return nil, err
		}
		output := s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader(data))}
		return &output, nil
	} 
	return nil, fmt.Errorf("an error")
	

}

func LoadMap(jsonFilePath string) ([]byte, error) {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}
