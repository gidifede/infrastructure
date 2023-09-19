package validator

//go:generate mockgen -source=./IValidator.go -destination=./IValidatorMock.go -package=validator

import (
	"context"
)

type IValidator interface {
	GetMainSchemaPath(commandType string) (string, error)
	GetReferenceSchemaPath() string
	GetJSONSchema(ctx context.Context, filename string) ([]byte, error)
	ValidateCommand(body, bodyJSONSchema []byte, bodyJSONSchemaRef []byte) error
	CommandMatchAPIPath(apiPath string, command string) error
}
