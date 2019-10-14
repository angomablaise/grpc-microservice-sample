package repository

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/service/store/repository"
)

// NewStoreRepository : Use go-redis client
func NewStoreRepository(client *redis.Client) repo.StoreRepository {
	return &storeRepository{client: client}
}

type storeRepository struct {
	client *redis.Client
}

func (rr *storeRepository) Set(ctx context.Context, store *api.Store) (string, error) {
	key := strconv.FormatInt(store.Id, 10)
	bytes, err := json.Marshal(store)
	if err != nil {
		return "", err
	}
	err = rr.client.WithContext(ctx).Set(key, bytes, 0).Err()
	if err != nil {
		return "", err
	}
	return key, nil
}

func (rr *storeRepository) Get(ctx context.Context, key string) (*api.Store, error) {
	value, err := rr.client.WithContext(ctx).Get(key).Bytes()
	if err != nil {
		return nil, err
	}
	var store *api.Store
	if err := json.Unmarshal(value, &store); err != nil {
		return nil, err
	}
	return store, nil
}

func (rr *storeRepository) Del(ctx context.Context, key string) (string, error) {
	err := rr.client.WithContext(ctx).Del(key).Err()
	if err != nil {
		return "", err
	}
	return key, nil
}
