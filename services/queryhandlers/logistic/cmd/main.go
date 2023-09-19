package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logistic_queryhandler/internal/database"
	"logistic_queryhandler/internal/utils"
	"os"

	overlog "github.com/Trendyol/overlog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"

	// "github.com/aws/aws-xray-sdk-go/xray"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	_ "github.com/go-sql-driver/mysql"
)

type Credential struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	Engine               string `json:"engine"`
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	DbName               string `json:"dbname"`
	DbInstanceIdentifier string `json:"dbInstanceIdentifier"`
}

type SecretData struct {
	MongoUser string `json:"username"`
	MongoPass string `json:"password"`
	MongoHost string `json:"host"`
	MongoPort int    `json:"port"`
}

var client *mongo.Client

var (
	dataStore   database.IDatabase
	ginLambda   *ginadapter.GinLambda
	xRayTraceID string
)

func init() {

	setUpLogging()

	setupDbConnection()

	setupAPI()
}

func setupAPI() {

	log.Debug().Msg("Gin cold start")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/v1/network/setup", GetNetwork)

	ginLambda = ginadapter.New(r)
}

func GetNetwork(c *gin.Context) {
	filter := ""
	ret, err := dataStore.SelectNetworkNodes(c, filter)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, ret)

}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	xRayTraceID = utils.GetXRayTraceID(req)

	overlog.MDC().Set("trace id", xRayTraceID)

	return ginLambda.ProxyWithContext(ctx, req)
}

// Method for setup logger
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

func getSecret() SecretData {
	secretName := "docdb-secret"
	region := "eu-central-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal().Err(err)
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
		log.Fatal().Err(err)
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString

	var secretData SecretData
	err = json.Unmarshal([]byte(secretString), &secretData)
	if err != nil {
		panic(err.Error())
	}

	return secretData
}

func setupDbConnection() {
	var filename = "/opt/global-bundle.pem"
	rootPEM, err := ioutil.ReadFile(filename)
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM([]byte(rootPEM)); !ok {
		fmt.Printf("get certs from %s fail!\n", filename)
		return
	}
	tlsConfig := &tls.Config{
		RootCAs:            roots,
		InsecureSkipVerify: true,
	}

	secrets := getSecret()

	clientOpts := options.Client() //.ApplyURI("mongodb://test:****@dds-bp*******1.mongodb.rds.aliyuncs.com:3717,dds-bp*******2.mongodb.rds.aliyuncs.com:3717/admin?replicaSet=mgset-XXXXX&ssl=true")
	// clientOpts.SetReadPreference(readpref.Secondary())
	clientOpts.SetWriteConcern(writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(1000)))
	clientOpts.SetTLSConfig(tlsConfig)
	hosts := []string{secrets.MongoHost}
	clientOpts.SetHosts(hosts)
	clientOpts.SetAuth(options.Credential{Username: secrets.MongoUser, Password: secrets.MongoPass})
	client, err = mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		fmt.Println("connect failed!")
		log.Fatal().Err(err)
		return
	}
	fmt.Println("connect successful!")

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
		fmt.Println("ping failed!")
		log.Fatal().Err(err)
		return
	}
	fmt.Println("ping successful!")

	dataStore = &database.MongoDB{
		DB: *client.Database("logistic"),
	}
}
