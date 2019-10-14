package repository

import (
	"context"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
)

type StoreRepository interface {
	Set(context.Context, *api.Store) (string, error)
	Get(context.Context, string) (*api.Store, error)
	Del(context.Context, string) (string, error)
}
