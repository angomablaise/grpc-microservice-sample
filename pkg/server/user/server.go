package server

import (
	"context"
	"fmt"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	config "github.com/smockoro/grpc-microservice-sample/pkg/config/user"
	"github.com/smockoro/grpc-microservice-sample/pkg/lib"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/repository/mysql/user"
	"github.com/smockoro/grpc-microservice-sample/pkg/service/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

	stackTracer := lib.NewStackTracer()
	repo := repo.NewUserRepository(db)
	server := user.NewUserServiceServer(repo, stackTracer)

	opts := []grpc_zap.Option{}
	zapLogger, _ := zap.NewProduction()
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(
				grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(zapLogger, opts...),
			grpc_auth.UnaryServerInterceptor(tokenAuthentication),
		),
	)

	api.RegisterUserServiceServer(s, server)
	reflection.Register(s)

	log.Println("starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}

func tokenAuthentication(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	if token != "sample_token" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	newCtx := context.WithValue(ctx, "authentication", "ok")
	return newCtx, nil
}
