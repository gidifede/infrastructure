package database

import (
	"confighandler/internal/model"
	"context"
)

type IDatabase interface {
	InsertNetwork(ctx context.Context, network model.NetworkSetup) (error)
	InsertProduct(ctx context.Context, network model.ProductSetup) (error)
	InsertTransport(ctx context.Context, network model.TransportSetup) (error)
	InsertRoute(ctx context.Context, network model.RouteSetup) (error)
}
