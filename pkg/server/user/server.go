package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql" // Register MySQL Driver
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	"github.com/smockoro/grpc-microservice-sample/pkg/service/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Port       string
	DBHost     string
	DBUser     string
	DBPassword string
	DBSchema   string
}

func RunServer() error {
	cfg := getConfig()

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBSchema)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	server := user.NewUserServiceServer(db)
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

func getConfig() Config {

	var cfg Config
	cfg.Port = os.Getenv("GRPC_PORT")
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.DBSchema = os.Getenv("DB_SCHEMA")

	return cfg
}
