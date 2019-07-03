package user_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	srv "github.com/smockoro/grpc-microservice-sample/pkg/service/user"
	mock "github.com/smockoro/grpc-microservice-sample/testdata"
)

func TestNewUserServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	s := srv.NewUserServiceServer(repo)

	if reflect.TypeOf(s).String() != "*user.server" {
		t.Errorf("want %s but actual %s", "*user.server", reflect.TypeOf(s))
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	s := srv.NewUserServiceServer(repo)

	users := map[string]*api.User{
		"no lost data": &api.User{
			Name: "Bob", Age: 16, Mail: "sample@sample.com", Address: "Tokyo",
		},
		"name is lost": &api.User{
			Age: 16, Mail: "sample@sample.com", Address: "Tokyo",
		},
		"age is lost": &api.User{
			Name: "Bob", Mail: "sample@sample.com", Address: "Tokyo",
		},
		"mail is lost": &api.User{
			Name: "Bob", Age: 16, Address: "Tokyo",
		},
		"address is lost": &api.User{
			Name: "Bob", Age: 16, Mail: "sample@sample.com",
		},
	}

	cases := []struct {
		name       string
		req        *api.CreateUserRequest
		user       *api.User
		res        *api.CreateUserResponse
		errorIsNil bool
	}{
		{
			name:       "no lost data",
			req:        &api.CreateUserRequest{User: users["no lost data"]},
			user:       users["no lost data"],
			res:        &api.CreateUserResponse{Id: 1},
			errorIsNil: true,
		},
		{
			name:       "name is lost",
			req:        &api.CreateUserRequest{User: users["name is lost"]},
			user:       users["name is lost"],
			res:        &api.CreateUserResponse{Id: 1},
			errorIsNil: true,
		},
		{
			name:       "age is lost",
			req:        &api.CreateUserRequest{User: users["age is lost"]},
			user:       users["age is lost"],
			res:        &api.CreateUserResponse{Id: 1},
			errorIsNil: true,
		},
		{
			name:       "address is lost",
			req:        &api.CreateUserRequest{User: users["address is lost"]},
			user:       users["address is lost"],
			res:        &api.CreateUserResponse{Id: 1},
			errorIsNil: true,
		},
		{
			name:       "*api.User is lost",
			req:        &api.CreateUserRequest{},
			errorIsNil: true,
		},
		{
			name:       "return err",
			req:        &api.CreateUserRequest{User: &api.User{}},
			user:       &api.User{},
			errorIsNil: false,
		},
	}

	ctx := context.Background()

	for _, c := range cases {
		c := c // cascading
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if c.errorIsNil {
				repo.EXPECT().Insert(ctx, c.user).Return(int64(1), nil)
			} else {
				repo.EXPECT().Insert(ctx, c.user).Return(int64(-1), fmt.Errorf("Error"))
			}
			res, err := s.Create(ctx, c.req)
			if err != nil && c.errorIsNil {
				t.Errorf("want %s actual %s", c.res, res)
			}
		})
	}

}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	_ = srv.NewUserServiceServer(repo)

}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	_ = srv.NewUserServiceServer(repo)

}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	_ = srv.NewUserServiceServer(repo)
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	_ = srv.NewUserServiceServer(repo)

}
