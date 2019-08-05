package user_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	srv "github.com/smockoro/grpc-microservice-sample/pkg/service/user"
	mock "github.com/smockoro/grpc-microservice-sample/testdata/mock/repository"
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
	s := srv.NewUserServiceServer(repo)

	cases := []struct {
		name string
		f    func(t *testing.T)
	}{
		{name: "Get OK", f: func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			reqID := 1
			user := &api.User{Name: "Bob", Age: 16, Mail: "sample@sample.com", Address: "Tokyo"}
			req := &api.GetUserRequest{Id: int64(reqID)}
			repo.EXPECT().SelectByID(ctx, int64(reqID)).Return(user, nil)
			_, err := s.Get(ctx, req)
			if err != nil {
				t.Errorf("want %s actual %s", "nil", err)
			}
		}},
		{name: "Get NG", f: func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := &api.GetUserRequest{}
			repo.EXPECT().SelectByID(ctx, int64(0)).Return(nil, fmt.Errorf("Error"))
			_, err := s.Get(ctx, req)
			if err == nil {
				t.Errorf("want %s actual %s", err, "nil")
			}
		}},
	}

	for _, c := range cases {
		t.Run(c.name, c.f)
	}

}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	s := srv.NewUserServiceServer(repo)

	cases := []struct {
		name string
		f    func(t *testing.T)
	}{
		{name: "Update OK", f: func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			user := &api.User{
				Id:      1,
				Name:    "Bob",
				Age:     16,
				Mail:    "sample@sample.com",
				Address: "Tokyo",
			}
			req := &api.UpdateUserRequest{User: user}
			repo.EXPECT().Update(ctx, user).Return(int64(1), nil)
			_, err := s.Update(ctx, req)
			if err != nil {
				t.Errorf("want %s actual %s", "nil", err)
			}
		}},
		{name: "Update NG", f: func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			user := &api.User{}
			req := &api.UpdateUserRequest{User: user}
			repo.EXPECT().Update(ctx, user).Return(int64(0), fmt.Errorf("Error"))
			_, err := s.Update(ctx, req)
			if err == nil {
				t.Errorf("want %s actual %s", err, "nil")
			}
		}},
	}

	for _, c := range cases {
		t.Run(c.name, c.f)
	}

}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	s := srv.NewUserServiceServer(repo)

	cases := []struct {
		name string
		f    func(t *testing.T)
	}{
		{name: "Delete OK", f: func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			reqID := 1
			req := &api.DeleteUserRequest{Id: int64(reqID)}
			repo.EXPECT().Delete(ctx, int64(reqID)).Return(int64(1), nil)
			_, err := s.Delete(ctx, req)
			if err != nil {
				t.Errorf("want %s actual %s", "nil", err)
			}
		}},
		{name: "Delete NG", f: func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			reqID := 0
			req := &api.DeleteUserRequest{Id: int64(reqID)}
			repo.EXPECT().Delete(ctx, int64(reqID)).Return(int64(0), fmt.Errorf("Error"))
			_, err := s.Delete(ctx, req)
			if err == nil {
				t.Errorf("want %s actual %s", err, "nil")
			}
		}},
	}

	for _, c := range cases {
		t.Run(c.name, c.f)
	}
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockUserRepository(ctrl)
	s := srv.NewUserServiceServer(repo)

	cases := []struct {
		name string
		f    func(t *testing.T)
	}{
		{name: "GetAll OK", f: func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			users := []*api.User{
				&api.User{
					Id:      1,
					Name:    "Bob",
					Age:     11,
					Mail:    "sample@sample.com",
					Address: "Tokyo",
				},
				&api.User{
					Id:      2,
					Name:    "Alice",
					Age:     13,
					Mail:    "example@sample.com",
					Address: "London",
				},
			}
			req := &api.GetAllUserRequest{}
			repo.EXPECT().SelectAll(ctx).Return(users, nil)
			_, err := s.GetAll(ctx, req)
			if err != nil {
				t.Errorf("want %s actual %s", "nil", err)
			}
		}},
		{name: "GetAll NG", f: func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			req := &api.GetAllUserRequest{}
			repo.EXPECT().SelectAll(ctx).Return(nil, fmt.Errorf("Error"))
			_, err := s.GetAll(ctx, req)
			if err == nil {
				t.Errorf("want %s actual %s", err, "nil")
			}
		}},
	}

	for _, c := range cases {
		t.Run(c.name, c.f)
	}

}
