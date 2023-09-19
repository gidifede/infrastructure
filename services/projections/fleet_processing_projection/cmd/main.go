package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fleet-processing-projection/internal"
	"fleet-processing-projection/internal/handler"
	"fleet-processing-projection/internal/repository"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	overlog "github.com/Trendyol/overlog"
)

type SecretData struct {
	MongoUser string `json:"username"`
	MongoPass string `json:"password"`
	MongoHost string `json:"host"`
	MongoPort int    `json:"port"`
}

func getSecret() SecretData {
	secretName := os.Getenv("DB_PROJECTION_CREDENTIAL_SECRET_NAME")
	region := "eu-central-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal().Msgf("cannot load default config. %s", err.Error())
		// No need to return error. log.Fatal() will exit
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal().Msgf("cannot get secret. %s", err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString

	var secretData SecretData
	err = json.Unmarshal([]byte(secretString), &secretData)
	if err != nil {
		log.Fatal().Msgf("cannot unmarshal secret data: %s", secretData)
	}

	return secretData
}

func init() {
	setUpLogging()
	var filename = "/opt/global-bundle.pem"
	rootPEM, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal().Msgf("cannot read cert file. %s", err.Error())
	}
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM([]byte(rootPEM)); !ok {
		log.Fatal().Msgf("get certs from %s fail!\n", filename)
	}
	tlsConfig := &tls.Config{
		RootCAs:            roots,
		InsecureSkipVerify: true,
	}

	secrets := getSecret()

	clientOpts := options.Client() //.ApplyURI("mongodb://test:****@dds-bp*******1.mongodb.rds.aliyuncs.com:3717,dds-bp*******2.mongodb.rds.aliyuncs.com:3717/admin?replicaSet=mgset-XXXXX&ssl=true")
	clientOpts.SetTLSConfig(tlsConfig)
	hosts := []string{secrets.MongoHost}
	clientOpts.SetHosts(hosts)
	clientOpts.SetAuth(options.Credential{Username: secrets.MongoUser, Password: secrets.MongoPass})
	clientOpts.SetReadPreference(readpref.Primary())
	clientOpts.SetRetryWrites(false)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal().Msgf("connection failed. %s", err.Error())
	}
	log.Debug().Msg("connection success")

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		fmt.Println("disconnect failed!")
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("disconnect successful!")
	// }()

	// Call Ping to verify that the deployment is up and the Client was
	// configured successfully. As mentioned in the Ping documentation, this
	// reduces application resiliency as the server may be temporarily
	// unavailable when Ping is called.
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal().Msgf("ping failed. %s", err.Error())
	}
	log.Debug().Msg("ping success")
	internal.Repo = repository.NewMongo(*client.Database("test"))
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
