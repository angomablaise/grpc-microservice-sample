package server_test

import (
	"context"
	"os"
	"testing"
	"time"

	userpb "github.com/smockoro/grpc-microservice-sample/pkg/api"
	"github.com/smockoro/grpc-microservice-sample/pkg/server/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestRunServer(t *testing.T) {
	go server.RunServer()
	time.Sleep(1 * time.Second) // Server Start uping

	conn, err := grpc.Dial("localhost:"+os.Getenv("GRPC_PORT"), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Did not connect to server: %v", err)
	}
	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)

	cases := []struct {
		name string
		f    func(t *testing.T)
	}{
		{name: "Create_NotAuthrizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			_, err := client.Create(ctx, &userpb.CreateUserRequest{
				User: &userpb.User{
					Name:    "Bob",
					Age:     11,
					Mail:    "bob@sample.com",
					Address: "Tokyo",
				},
			})

			if err == nil {
				t.Errorf("It is expected that err is not nil but err is nil")
			}
		}},
		{name: "CreateUpdateDelete_AuthrizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "bearer sample_token")
			_, err := client.Create(ctx, &userpb.CreateUserRequest{
				User: &userpb.User{
					Name:    "Bob",
					Age:     11,
					Mail:    "bob@sample.com",
					Address: "Tokyo",
				},
			})

			if err != nil {
				t.Errorf("It is expected that err is nil but err is not nil: %v", err)
			}
		}},
		{name: "GetAll_NotAuthorizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			_, err := client.GetAll(ctx, &userpb.GetAllUserRequest{})

			if err == nil {
				t.Errorf("It is expected that err is not nil but err is nil")
			}
		}},
		{name: "GetAll_AuthorizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "bearer sample_token")
			_, err := client.GetAll(ctx, &userpb.GetAllUserRequest{})

			if err != nil {
				t.Errorf("It is expected that err is nil but err is not nil: %v", err)
			}
		}},
		{name: "Update_NotAuthrizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			id := testGetId(t, client)
			_, err := client.Update(ctx, &userpb.UpdateUserRequest{
				User: &userpb.User{
					Id:      id,
					Name:    "Bob",
					Age:     88,
					Mail:    "bobbobbob@sample.com",
					Address: "Athena",
				},
			})

			if err == nil {
				t.Errorf("It is expected that err is not nil but err is nil")
			}
		}},
		{name: "Update_AuthrizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "bearer sample_token")
			id := testGetId(t, client)
			_, err := client.Update(ctx, &userpb.UpdateUserRequest{
				User: &userpb.User{
					Id:      id,
					Name:    "Bob",
					Age:     88,
					Mail:    "bobbobbob@sample.com",
					Address: "Athena",
				},
			})

			if err != nil {
				t.Errorf("It is expected that err is nil but err is not nil: %v", err)
			}
		}},
		{name: "Get_NotAuthrizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			id := testGetId(t, client)
			_, err := client.Get(ctx, &userpb.GetUserRequest{Id: id})

			if err == nil {
				t.Errorf("It is expected that err is not nil but err is nil")
			}
		}},
		{name: "Get_AuthrizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "bearer sample_token")
			id := testGetId(t, client)
			_, err := client.Get(ctx, &userpb.GetUserRequest{Id: id})

			if err != nil {
				t.Errorf("It is expected that err is nil but err is not nil: %v", err)
			}
		}},
		{name: "Delete_NotAuthrizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			id := testGetId(t, client)
			_, err := client.Delete(ctx, &userpb.DeleteUserRequest{Id: id})

			if err == nil {
				t.Errorf("It is expected that err is not nil but err is nil")
			}
		}},
		{name: "Delete_AuthrizationHeader", f: func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "bearer sample_token")
			id := testGetId(t, client)
			_, err := client.Delete(ctx, &userpb.DeleteUserRequest{Id: id})

			if err != nil {
				t.Errorf("It is expected that err is nil but err is not nil: %v", err)
			}
		}},
	}

	for _, c := range cases {
		t.Run(c.name, c.f)
	}
}

func testGetId(t *testing.T, client userpb.UserServiceClient) int64 {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "bearer sample_token")
	resp, err := client.GetAll(ctx, &userpb.GetAllUserRequest{})

	if err != nil {
		t.Fatal(err)
	}

	return resp.Users[len(resp.Users)-1].Id
}
