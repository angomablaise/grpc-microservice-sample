package main

import (
	"context"
	"log"
	"time"

	userpb "github.com/smockoro/grpc-microservice-sample/pkg/api"
	config "github.com/smockoro/grpc-microservice-sample/pkg/config/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	cfg := config.NewConfig()
	conn, err := grpc.Dial("localhost:"+cfg.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to server: %v", err)
	}
	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)

	err = CheckGetAll(client)
	if err != nil {
		log.Printf("GetAll Error: %v", err)
	}
}

func CheckGetAll(client *userpb.UserServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "bearer sample_token")

	resp, err := client.GetAll(ctx, &userpb.GetAllUserRequest{})
	if err != nil {
		return err
	}
	log.Printf("GetAll: %v", resp)

	return nil
}
