package main

import (
	"confighandler/internal/database"
	"confighandler/internal/model"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"

	"log"

	"confighandler/internal/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	overlog "github.com/Trendyol/overlog"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

type EventInput struct {
	Limit int64 `json:"limit"`
}

type SecretData struct {
	MongoUser string `json:"username"`
	MongoPass string `json:"password"`
	MongoHost string `json:"host"`
	MongoPort int    `json:"port"`
}

var (
	db          database.IDatabase
	ginLambda   *ginadapter.GinLambda
	xRayTraceID string
	c           *context.Context
	client      *mongo.Client
)

func getSecret() SecretData {
	secretName := os.Getenv("DB_PROJECTION_CREDENTIAL_SECRET_NAME")
	region := os.Getenv("REGION")

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err.Error())
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

func init() {
	var filename = "/opt/global-bundle.pem"
	rootPEM, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)

	}
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

	//clientOpts.SetWriteConcern(writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(1000)))
	clientOpts.SetTLSConfig(tlsConfig)
	hosts := []string{secrets.MongoHost}
	clientOpts.SetHosts(hosts)
	clientOpts.SetAuth(options.Credential{Username: secrets.MongoUser, Password: secrets.MongoPass})
	clientOpts.SetReadPreference(readpref.Primary())
	client, err = mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		fmt.Println("connect failed!")
		log.Fatal(err)
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
		log.Fatal(err)
		return
	}
	fmt.Println("ping successful!")
	db = database.NewMongo(*client.Database("logistic"))

	setupAPI()

}

func setupAPI() {
	log.Println("Gin cold start")
	r := gin.Default()

	r.POST("/v1/network/setup", NetworkSetup)
	r.POST("/v1/product/setup", ProductSetup)
	r.POST("/v1/fleet/setup", TransportSetup)
	r.POST("/v1/route/setup", RouteSetup)
	ginLambda = ginadapter.New(r)
}

func NetworkSetup(c *gin.Context) {
	var network model.NetworkSetup

	err := c.BindJSON(&network)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid network body"})
		return
	}

	err = db.InsertNetwork(c, network)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "Network inserted!")

}

func ProductSetup(c *gin.Context) {
	var product model.ProductSetup

	err := c.BindJSON(&product)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid product body"})
		return
	}

	err = db.InsertProduct(c, product)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "Product inserted!")
}

func TransportSetup(c *gin.Context) {
	var transport model.TransportSetup

	err := c.BindJSON(&transport)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid transport body"})
		return
	}

	err = db.InsertTransport(c, transport)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "transport inserted!")
}

func RouteSetup(c *gin.Context) {
	var route model.RouteSetup

	err := c.BindJSON(&route)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid route body"})
		return
	}

	err = db.InsertRoute(c, route)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "route inserted!")
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	xRayTraceID = utils.GetXRayTraceID(req)
	overlog.MDC().Set("trace id", xRayTraceID)

	c = &ctx

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

// Method for setup logger
func setUpLogging() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Default level is info, unless LOG_LEVEL env var is present
	if logLevel, isPresent := os.LookupEnv("LOG_LEVEL"); isPresent {
		l, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			log.Fatal(err)
		}
		zerolog.SetGlobalLevel(l)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zlogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	overlog.New(zlogger)
	overlog.SetGlobalFields([]string{"trace id", "method", "class"})

}

