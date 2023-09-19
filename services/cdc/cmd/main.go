package main

import (
	"cdc/internal/publisher"
	"cdc/internal/utils"
	"context"

	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	overlog "github.com/Trendyol/overlog"
)

const SnsTopicArn = "SNS_TOPIC_ARN"

var snsService publisher.Publisher
var snsTopic string

func init() {
	setUpLogging()
	log.Debug().Msg("Fetching configurations...")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	snsAwsService := sns.New(sess)
	snsTopic = os.Getenv(SnsTopicArn)
	xray.AWS(snsAwsService.Client)

	snsService = publisher.NewSNSPublisher(snsAwsService)

	log.Debug().Msg("Fetch succesfully endend...")
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

func handleRequest(ctx context.Context, e events.DynamoDBEvent) error {

	for _, record := range e.Records {
		log.Debug().Msgf("Processing request data for event ID %s, type %s, message %v.\n", record.EventID, record.EventName, record.Change.NewImage)
		if record.EventName != "REMOVE" {

			jsonm := make(map[string]string)
			utils.ParseDynamoStream(ctx, record, jsonm)

			ce, errConvert := utils.ConvertToCloudEvent(jsonm)

			if errConvert != nil {
				log.Err(errConvert).Msg("Error converting cdc from Dynamo to CloudEvent format")
				return errConvert
			}

			err := sendXRaySegment(ctx, jsonm["eventId"])
			if err != nil {
				log.Err(err).Msg("Error sending Xray annotated segment")
			}

			xRayTraceID := xray.TraceID(ctx)
			overlog.MDC().Set("trace id", xRayTraceID)

			err = snsService.PublishEvent(ctx, ce, snsTopic)
			if err != nil {
				log.Err(err).Msg("Publishing failed")
				return err
			}
		}
	}

	return nil
}

func sendXRaySegment(ctx context.Context, eventID string) error {

	ctx, subsegment := xray.BeginSubsegment(ctx, "MySubsegment")
	xray.AddAnnotation(ctx, "eventId", eventID)

	defer subsegment.Close(nil)

	log.Debug().Msgf("Subsegment created with eventID: %s", eventID)

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
