package user

import (
	"context"
	"fmt"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/service/user/repository"
	"go.opencensus.io/trace"
)

type server struct {
	repo repo.UserRepository
}

// NewUserServiceServer : return User Service
func NewUserServiceServer(repo repo.UserRepository) api.UserServiceServer {
	return &server{repo: repo}
}

func (s *server) Create(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	ctx, span := trace.StartSpan(ctx, "oc.user-service.service.create")
	defer span.End()

	id, err := s.repo.Insert(ctx, req.User)
	if err != nil {
		return nil, err
	}

	return &api.CreateUserResponse{Id: id}, nil
}

func (s *server) Get(ctx context.Context, req *api.GetUserRequest) (*api.GetUserResponse, error) {
	ctx, span := trace.StartSpan(ctx, "oc.user-service.service.get")
	defer span.End()

	user, err := s.repo.SelectByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &api.GetUserResponse{User: user}, nil
}

func (s *server) Update(ctx context.Context, req *api.UpdateUserRequest) (*api.UpdateUserResponse, error) {
	ctx, span := trace.StartSpan(ctx, "oc.user-service.service.update")
	defer span.End()

	updated, err := s.repo.Update(ctx, req.User)
	if err != nil {
		return nil, err
	}

	return &api.UpdateUserResponse{Updated: updated}, nil
}

func (s *server) Delete(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	ctx, span := trace.StartSpan(ctx, "oc.user-service.service.delete")
	defer span.End()

	deleted, err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &api.DeleteUserResponse{Deleted: deleted}, nil
}

func (s *server) GetAll(ctx context.Context, req *api.GetAllUserRequest) (*api.GetAllUserResponse, error) {
	ctx, span := trace.StartSpan(ctx, "oc.user-service.service.getall")
	defer span.End()
	fmt.Printf("%v\n", span)

	users, err := s.repo.SelectAll(ctx)
	if err != nil {
		return nil, err
	}

	return &api.GetAllUserResponse{Users: users}, nil
}
