package handler

import (
	"context"
	"facility_queryhandler/internal"
	"facility_queryhandler/internal/utils"

	overlog "github.com/Trendyol/overlog"
	"github.com/aws/aws-lambda-go/events"
)

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	internal.XRayTraceID = utils.GetXRayTraceID(req)

	overlog.MDC().Set("trace id", internal.XRayTraceID)

	return internal.GinLambda.ProxyWithContext(ctx, req)
}
