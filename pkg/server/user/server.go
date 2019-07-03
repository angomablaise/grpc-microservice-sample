package server

import (
	"fmt"
	"log"
	"net"

	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	"github.com/smockoro/grpc-microservice-sample/pkg/config/user"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/repository/mysql/user"
	"github.com/smockoro/grpc-microservice-sample/pkg/service/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// RunServer : Component Injected and Startup gRPC Server
func RunServer() error {
	cfg := config.NewConfig()

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := ConnectDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	repo := repo.NewUserRepository(db)
	server := user.NewUserServiceServer(repo)
	s := grpc.NewServer()

	api.RegisterUserServiceServer(s, server)
	reflection.Register(s)

	log.Println("starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}
