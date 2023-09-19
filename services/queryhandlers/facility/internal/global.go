package internal

import (
	"facility_queryhandler/internal/repository"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var (
	Repo        repository.IDatabase
	GinLambda   *ginadapter.GinLambda
	XRayTraceID string
)
