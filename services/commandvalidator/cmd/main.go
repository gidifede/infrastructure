package main

import (
	"commandvalidator/internal/pkg/queue"
	"commandvalidator/internal/utils"
	validator "commandvalidator/internal/validator"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	overlog "github.com/Trendyol/overlog"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ENV VARIABLES
const S3BucketName = "S3_BUCKET_NAME"
const S3ConfigPrefix = "S3_CONFIG_PREFIX"
const SqsQueueName = "SQS_QUEUE_NAME"
const MessageGroup = "MESSAGE_GROUP"
const ConfigVersion = "CONFIG_VERSION"
const CommandOrigin = "COMMAND_ORIGIN"
const CommandAPIPathMapEnabled = "COMMAND_API_PATH_MAP_ENABLED"
const CommandAPIPathMapFilename = "COMMAND_API_PATH_MAP_FILENAME"

// RESPONSE
type ResponseBody struct {
	Message string `json:"message"`
}

// GLOBAL VARIABLES
var (
	s3Session               *s3.S3
	sessionQueue            *sqs.SQS
	queueManagement         queue.IQeueu
	bucketName              string
	bucketConfPrefix        string
	queueNameCommandHandler string
	messageGroup            string
	configVersion           string
	commandOrigin           string
	commandValidator        validator.IValidator
	xRayTraceID             string
)

// INIT METHOD
func init() {

	setUpLogging()

	//Initialize session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Open s3 session
	s3Session = s3.New(sess)
	// Open sqs session
	sessionQueue = sqs.New(sess)

	queueManagement = queue.NewSqsQueue(*sessionQueue)
	//Xray tracing
	xray.AWS(s3Session.Client)
	xray.AWS(sessionQueue.Client)

	//Retrieve variables from env
	bucketName = os.Getenv(S3BucketName)
	queueNameCommandHandler = os.Getenv(SqsQueueName)
	messageGroup = os.Getenv(MessageGroup)
	bucketConfPrefix = os.Getenv(S3ConfigPrefix)
	configVersion = os.Getenv(ConfigVersion)
	commandOrigin = os.Getenv(CommandOrigin)
	//xRayTraceID = os.Getenv("x-amzn-trace-id")
	commandAPIPathMapEnabled, err := strconv.ParseBool(os.Getenv(CommandAPIPathMapEnabled))
	commandAPIPathMapFilename := os.Getenv(CommandAPIPathMapFilename)
	if err != nil {
		commandAPIPathMapEnabled = false
	}
	commandValidator = validator.NewValidator(s3Session, bucketName, bucketConfPrefix, configVersion, commandAPIPathMapEnabled, commandAPIPathMapFilename)

}

// Metodo per la gestione della richiesta inviata alla lambda
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	xRayTraceID = utils.GetXRayTraceID(request)

	overlog.MDC().Set("trace id", xRayTraceID)

	//Conversion body from APIGatewayProxyRequest to Cloud Events
	commandReceived, errConvert := convertToCloudEvent(request.Body)
	if errConvert != nil {
		log.Err(errConvert).Msg("")
		//Generate BAD REQEUST
		return generateResponse(400, errConvert.Error()), nil
	}

	// Check API path and command type match
	errMatch := commandValidator.CommandMatchAPIPath(request.Path, commandReceived.Type())
	if errMatch != nil {
		log.Err(errMatch).Msg("")
		//Generate BAD REQEUST
		return generateResponse(400, errMatch.Error()), nil
	}

	//Retrieve JSONSchema
	uri, errGetPath := commandValidator.GetMainSchemaPath(commandReceived.Type())
	if errGetPath != nil {
		log.Err(errGetPath).Msg("")
		//Generate INTERNAL SERVER ERROR
		return generateResponse(500, errGetPath.Error()), nil
	}
	JSONSchema, configErr := commandValidator.GetJSONSchema(ctx, uri)
	if configErr != nil {
		log.Err(configErr).Msg("")
		//Generate INTERNAL SERVER ERROR
		return generateResponse(500, configErr.Error()), nil
	}

	jsonrefpath := commandValidator.GetReferenceSchemaPath()
	jsonRefSchema, configRefErr := commandValidator.GetJSONSchema(ctx, jsonrefpath)
	if configRefErr != nil {
		log.Err(configErr).Msg("")
		//Generate INTERNAL SERVER ERROR
		return generateResponse(500, configRefErr.Error()), nil
	}

	//Validation Command
	errValidation := commandValidator.ValidateCommand(commandReceived.Data(), JSONSchema, jsonRefSchema)
	if errValidation != nil {
		//Generate BAD REQUEST
		return generateResponse(400, errValidation.Error()), nil
	}

	//Send to FIFO QUEUE for Command Handler
	errSendFIFOMsg := queueManagement.SendMsg(ctx, queueNameCommandHandler, request.Body, messageGroup)
	if errSendFIFOMsg != nil {
		log.Err(errSendFIFOMsg)
		//Generate INTERNAL SERVER ERROR
		return generateResponse(500, errSendFIFOMsg.Error()), nil
	}

	// Response ACCEPTED
	return generateResponse(202, "Command Accepted"), nil
}

// Metodo con responsabilita di conversione
// del body del messaggio proveniente da un
// evento di ApiGateway in formalism Cloud Event
func convertToCloudEvent(body string) (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	err := json.Unmarshal([]byte(body), &event)
	if nil != err {
		return nil, err
	}
	return &event, nil
}

// Metodo per la generazione della risposta per l'api gateway
func generateResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	response := &ResponseBody{
		Message: fmt.Sprint(message),
	}
	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("build response error: %v", message), StatusCode: 500}
	}
	responseAPI := events.APIGatewayProxyResponse{Body: string(body), StatusCode: statusCode}
	return responseAPI
}

func setUpLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Default level is info, unless LOG_LEVEL env var is present
	if logLevel, isPresent := os.LookupEnv("LOG_LEVEL"); isPresent {
		l, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			log.Err(err).Msg("")
		}
		zerolog.SetGlobalLevel(l)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zlogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	overlog.New(zlogger)
	overlog.SetGlobalFields([]string{"trace id", "method", "class"})

}

// Main per lo start della lambda
func main() {
	lambda.Start(HandleRequest)
}
