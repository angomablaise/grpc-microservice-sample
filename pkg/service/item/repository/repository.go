package repository

import (
	"context"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
)

type ItemRepository interface {
	Insert(context.Context, *api.Item) (int64, error)
	SelectByID(context.Context, int64) (*api.Item, error)
	SelectAll(context.Context) ([]*api.Item, error)
	Update(context.Context, *api.Item) (int64, error)
	Delete(context.Context, int64) (int64, error)
}
