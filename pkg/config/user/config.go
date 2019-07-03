package config

import (
	"os"
)

type Config struct {
	Port       string
	DBHost     string
	DBUser     string
	DBPassword string
	DBSchema   string
}

func NewConfig() *Config {
	var cfg Config
	cfg.Port = os.Getenv("GRPC_PORT")
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.DBSchema = os.Getenv("DB_SCHEMA")
	return &cfg
}
