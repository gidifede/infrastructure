package main

import (
	"os"

	"command-handler/internal"
	es "command-handler/internal/eventstore"
	"command-handler/internal/handler"
	"command-handler/internal/utils"

	overlog "github.com/Trendyol/overlog"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	setUpLogging()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	//Open Dynamo session
	esConnection := dynamodb.New(sess)

	internal.EventStore = es.NewDynamoEventStore(esConnection, os.Getenv("DYNAMO_TABLE_NAME"))
	xray.AWS(esConnection.Client)

	// // Open s3 session
	// s3Connection := s3.New(sess)
	// internal.StateMachineConfig = smConfig.NewS3StateMachineConfig(s3Connection, os.Getenv("S3_BUCKET_NAME"), strings.Join([]string{os.Getenv("S3_CONFIG_PREFIX"), "stateMachine.json"}, "/"))
	// xray.AWS(s3Connection.Client)

	internal.ConfigLoaded = true
	// // Load state machine
	// smBody, err := internal.StateMachineConfig.Get(context.Background())
	// if err != nil {
	// 	internal.ConfigLoaded = false
	// 	log.Err(err).Msg("state machine config loading failed. ")
	// }
	// internal.StateMachine, err = sm.LoadStateMachine(smBody)
	// if err != nil {
	// 	internal.ConfigLoaded = false
	// 	log.Err(err).Msg("state machine loading failed.")
	// } else {
	// 	_, err = internal.StateMachine.Validate()
	// 	if err != nil {
	// 		internal.ConfigLoaded = false
	// 		log.Err(err).Msg("state machine validation failed.")
	// 	}
	// }
	internal.Timestamp = utils.NewTimestampManager()
}

func main() {
	lambda.Start(handler.HandleEvents)
}

func setUpLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Default level is info, unless LOG_LEVEL env var is present
	if logLevel, isPresent := os.LookupEnv("LOG_LEVEL"); isPresent {
		l, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			log.Err(err).Msgf("")
		}
		zerolog.SetGlobalLevel(l)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zlogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	overlog.New(zlogger)
	overlog.SetGlobalFields([]string{"trace id", "method", "class"})
}
