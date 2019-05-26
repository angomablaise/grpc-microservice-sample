package item

import (
	"context"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/service/item/repository"
)

type server struct {
	repo repo.ItemRepository
}

// NewItemServiceServer : Inject ItemService
func NewItemServiceServer(repo repo.ItemRepository) api.ItemServiceServer {
	return &server{repo: repo}
}

func (s *server) Create(ctx context.Context, req *api.CreateItemRequest) (*api.CreateItemResponse, error) {
	id, err := s.repo.Insert(ctx, req.Item)
	if err != nil {
		return nil, err
	}

	return &api.CreateItemResponse{Id: id}, nil
}

func (s *server) Get(ctx context.Context, req *api.GetItemRequest) (*api.GetItemResponse, error) {
	item, err := s.repo.SelectByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &api.GetItemResponse{Item: item}, nil
}

func (s *server) Update(ctx context.Context, req *api.UpdateItemRequest) (*api.UpdateItemResponse, error) {
	updated, err := s.repo.Update(ctx, req.Item)
	if err != nil {
		return nil, err
	}

	return &api.UpdateItemResponse{Updated: updated}, nil
}

func (s *server) Delete(ctx context.Context, req *api.DeleteItemRequest) (*api.DeleteItemResponse, error) {
	deleted, err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &api.DeleteItemResponse{Deleted: deleted}, nil
}

func (s *server) GetAll(ctx context.Context, req *api.GetAllItemRequest) (*api.GetAllItemResponse, error) {
	items, err := s.repo.SelectAll(ctx)
	if err != nil {
		return nil, err
	}

	return &api.GetAllItemResponse{Items: items}, nil
}
