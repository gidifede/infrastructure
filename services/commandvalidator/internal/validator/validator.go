package validator

import (
	"commandvalidator/internal/utils"

	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	overlog "github.com/Trendyol/overlog"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/rs/zerolog/log"

	"github.com/xeipuuv/gojsonschema"
)

type Validator struct {
	connection               s3iface.S3API
	commandAPIPathMapEnabled bool
	commandAPIPathMap        map[string]string
	bucketName               string
	bucketConfPrefix         string
	configVersion            string
	// appConfigRetriever       configs.AppConfigRetriver
}

func NewValidator(conn s3iface.S3API, bucketName string, bucketConfPrefix string, configVersion string, commandAPIPathMapEnabled bool, commandAPIPathMapFilename string) IValidator {

	class, method := "Validator", "NewValidator"

	overlog.MDC().Set("method", method)
	overlog.MDC().Set("class", class)

	if commandAPIPathMapEnabled {
		// Get map from s3
		input := &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(bucketConfPrefix + "/" + commandAPIPathMapFilename),
		}
		result, err := conn.GetObjectWithContext(context.Background(), input)
		if err == nil {
			defer result.Body.Close()
			body, err := ioutil.ReadAll(result.Body)
			if err == nil {
				var commandAPIPathMap = map[string]string{}
				err = json.Unmarshal(body, &commandAPIPathMap)
				if err == nil {
					log.Info().Msgf("using command - API path map")
					return &Validator{
						connection:               conn,
						commandAPIPathMapEnabled: commandAPIPathMapEnabled,
						commandAPIPathMap:        commandAPIPathMap,
						bucketName:               bucketName,
						bucketConfPrefix:         bucketConfPrefix,
						configVersion:            configVersion,
						// appConfigRetriever:       appConfigRetriever,
					}
				}
				log.Err(err).Msgf("")
			}
			log.Err(err).Msgf("")
		}
		log.Err(err).Msgf("")
	}
	log.Info().Msgf("not using command - API path map")
	return &Validator{
		connection:               conn,
		commandAPIPathMapEnabled: false,
		commandAPIPathMap:        nil,
		bucketName:               bucketName,
		bucketConfPrefix:         bucketConfPrefix,
		configVersion:            configVersion,
		// appConfigRetriever:       appConfigRetriever,
	}
}

func (v *Validator) GetMainSchemaPath(commandType string) (string, error) {
	utils.AddClassAndMethodToMDC(v)

	log.Debug().Msgf("Get Json Schema for command: %s\n", commandType)

	str := strings.Split(commandType, ".")
	if len(str) < 2 {
		return "", fmt.Errorf("invalid command type")
	}
	commandJSON := str[len(str)-1]
	filename := v.bucketConfPrefix + "/" + v.configVersion + "/JSONSCHEMA/" + str[len(str)-2] + "/Commands/" + commandJSON + ".json"

	log.Debug().Msgf("File json sche, bucketNamema: %s ", filename)
	return filename, nil
}

func (v *Validator) GetReferenceSchemaPath() string {
	utils.AddClassAndMethodToMDC(v)

	log.Debug().Msg("Get reference Json Schema")
	filename := v.bucketConfPrefix + "/" + v.configVersion + "/JSONSCHEMA/Common/common.json"

	log.Debug().Msgf("Reference file json schema: %s ", filename)
	return filename
}

func (v *Validator) GetJSONSchema(ctx context.Context, filename string) ([]byte, error) {
	utils.AddClassAndMethodToMDC(v)

	log.Debug().Msgf("Retriving json schema: %s from %s", filename, v.bucketName)

	input := &s3.GetObjectInput{
		Bucket: aws.String(v.bucketName),
		Key:    aws.String(filename),
	}

	result, err := v.connection.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("get json schema failed. %v", err)
	}
	defer result.Body.Close()
	bodyJSONSchema, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("read json schema failed. %v", err)
	}
	return bodyJSONSchema, nil
}

func (v *Validator) ValidateCommand(body, bodyJSONSchema []byte, bodyJSONSchemaRef []byte) error {
	utils.AddClassAndMethodToMDC(v)

	sl := gojsonschema.NewSchemaLoader()
	ref := gojsonschema.NewBytesLoader(bodyJSONSchemaRef)

	err := sl.AddSchema("https://example.com/schemas/common", ref)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewBytesLoader(bodyJSONSchema)
	schema, err := sl.Compile(schemaLoader)
	if err != nil {
		return err
	}

	documentLoader := gojsonschema.NewBytesLoader(body)
	result, err := schema.Validate(documentLoader)
	// result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if result.Valid() {
		log.Debug().Msg("Valid document")
		return nil
	}

	log.Debug().Msg("Invalid document for the following errors:\n")
	for _, desc := range result.Errors() {
		log.Debug().Msgf("%s", desc)
	}
	return fmt.Errorf("validation error")

}

func (v *Validator) CommandMatchAPIPath(apiPath string, command string) error {
	utils.AddClassAndMethodToMDC(v)

	commandTypeElements := strings.Split(command, ".")
	if len(commandTypeElements) < 1 {
		log.Debug().Msg(fmt.Sprintf("invalid command type: %v", command))
		return fmt.Errorf("command type does not match called API")
	}
	commandTypeName := commandTypeElements[len(commandTypeElements)-1]

	apiPathElements := strings.Split(apiPath, "/")
	if len(apiPathElements) < 1 {
		log.Debug().Msg(fmt.Sprintf("invalid API path: %v", apiPath))
		return fmt.Errorf("command type does not match called API")
	}
	pathCommandName := apiPathElements[len(apiPathElements)-1]

	if v.commandAPIPathMapEnabled {
		if v.commandAPIPathMap[pathCommandName] != commandTypeName {
			log.Debug().Msg(fmt.Sprintf("command type does not match called API. API path: %v, command type: %v (normalized: %v)", apiPath, command, commandTypeName))
			return fmt.Errorf("command type does not match called API")
		}
		return nil
	}

	commandTypeName = strings.ToLower(commandTypeName)
	pathCommandName = strings.ToLower(pathCommandName)
	pathCommandName = strings.Replace(pathCommandName, "_", "", -1)
	if pathCommandName != commandTypeName {
		log.Debug().Msg(fmt.Sprintf("command type does not match called API. API path: %v (normalized : %v), command type: %v (normalized: %v)", apiPath, pathCommandName, command, commandTypeName))
		return fmt.Errorf("command type does not match called API")
	}
	return nil

}
