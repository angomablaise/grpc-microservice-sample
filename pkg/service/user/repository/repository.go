package repository

import (
	"context"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
)

type UserRepository interface {
	Insert(context.Context, *api.User) (int64, error)
	SelectByID(context.Context, int64) (*api.User, error)
	SelectAll(context.Context) ([]*api.User, error)
	Update(context.Context, *api.User) (int64, error)
	Delete(context.Context, int64) (int64, error)
}
