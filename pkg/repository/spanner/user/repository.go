package repository

import (
	"context"

	"cloud.google.com/go/spanner"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/service/user/repository"
)

type userRepository struct {
	client *spanner.Client
}

func NewUserRepository(client *spanner.Client) repo.UserRepository {
	return &userRepository{client: client}
}

func (u *userRepository) Insert(ctx context.Context, user *api.User) (int64, error) {
	// 単発書き込み
	return 0, nil
}

func (u *userRepository) SelectByID(ctx context.Context, id int64) (*api.User, error) {
	// 単発の読み取り
	return &api.User{}, nil
}

func (u *userRepository) SelectAll(ctx context.Context) ([]*api.User, error) {
	// 多数の読み取り
	list := []*api.User{}
	return list, nil
}

func (u *userRepository) Update(ctx context.Context, user *api.User) (int64, error) {
	// 参照して書き込み
	return 0, nil
}

func (u *userRepository) Delete(ctx context.Context, id int64) (int64, error) {
	// 単発書き込み
	return 0, nil
}
