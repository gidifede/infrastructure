package database

import (
	"context"
	"logistic_queryhandler/internal/model"
)

type IDatabase interface {
	SelectNetworkNodes(c context.Context, filter string) (model.Network, error)
}
