package server

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Register MySQL Driver
	"github.com/jmoiron/sqlx"
	"github.com/smockoro/grpc-microservice-sample/pkg/config/user"
)

// ConnectDB : connect to mysql server
func ConnectDB(cfg *config.Config) (*sqlx.DB, error) {
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("DBUser is none")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DBPassword is none")
	}
	if cfg.DBHost == "" {
		return nil, fmt.Errorf("DBHost is none")
	}
	if cfg.DBSchema == "" {
		return nil, fmt.Errorf("DBSchema is none")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBSchema)
	return sqlx.Open("mysql", dsn)
}
