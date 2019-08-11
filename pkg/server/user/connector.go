package server

import (
	"database/sql"
	"fmt"

	"contrib.go.opencensus.io/integrations/ocsql"
	_ "github.com/go-sql-driver/mysql" // Register MySQL Driver
	"github.com/jmoiron/sqlx"
	config "github.com/smockoro/grpc-microservice-sample/pkg/config/user"
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

	driverName, err := ocsql.Register("mysql", ocsql.WithAllTraceOptions())
	if err != nil {
		return nil, fmt.Errorf("driver not created")
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("can't open database")
	}

	return sqlx.NewDb(db, "mysql"), nil
}
