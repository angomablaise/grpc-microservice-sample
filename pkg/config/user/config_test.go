package config_test

import (
	"os"
	"testing"

	config "github.com/smockoro/grpc-microservice-sample/pkg/config/user"
)

func TestNewConfig(t *testing.T) {
	cases := []struct {
		name       string
		values     map[string]string
		errorIsNil bool
	}{
		{name: "env value not loss", values: map[string]string{
			"GRPC_PORT":   "9000",
			"DB_HOST":     "localhost:9000",
			"DB_USER":     "connect_user",
			"DB_PASSWORD": "password",
			"DB_SCHEMA":   "shema"}, errorIsNil: false},
		{name: "GRPC_PORT is lost", values: map[string]string{
			"GRPC_PORT":   "",
			"DB_HOST":     "localhost:9000",
			"DB_USER":     "connect_user",
			"DB_PASSWORD": "password",
			"DB_SCHEMA":   "shema"}, errorIsNil: false},
		{name: "DB_HOST is lost", values: map[string]string{
			"GRPC_PORT":   "9000",
			"DB_HOST":     "",
			"DB_USER":     "connect_user",
			"DB_PASSWORD": "password",
			"DB_SCHEMA":   "shema"}, errorIsNil: false},
		{name: "DB_USER is lost", values: map[string]string{
			"GRPC_PORT":   "9000",
			"DB_HOST":     "localhost:9000",
			"DB_USER":     "",
			"DB_PASSWORD": "password",
			"DB_SCHEMA":   "shema"}, errorIsNil: false},
		{name: "DB_PASSWORD is lost", values: map[string]string{
			"GRPC_PORT":   "9000",
			"DB_HOST":     "localhost:9000",
			"DB_USER":     "connect_user",
			"DB_PASSWORD": "",
			"DB_SCHEMA":   "shema"}, errorIsNil: false},
		{name: "DB_SCHEMA is lost", values: map[string]string{
			"GRPC_PORT":   "9000",
			"DB_HOST":     "localhost:9000",
			"DB_USER":     "connect_user",
			"DB_PASSWORD": "password",
			"DB_SCHEMA":   ""}, errorIsNil: false},
	}

	for _, c := range cases {
		testSetEnvs(t, c.values) // don't Parallel because Enviroment Value is vibration
		t.Run(c.name, func(t *testing.T) {
			cfg := config.NewConfig()
			if cfg.Port != c.values["GRPC_PORT"] {
				t.Errorf("want %s but actual %s", c.values["GRPC_PORT"], cfg.Port)
			}
			if cfg.DBHost != c.values["DB_HOST"] {
				t.Errorf("want %s but actual %s", c.values["DB_HOST"], cfg.DBHost)
			}
			if cfg.DBUser != c.values["DB_USER"] {
				t.Errorf("want %s but actual %s", c.values["DB_USER"], cfg.DBUser)
			}
			if cfg.DBPassword != c.values["DB_PASSWORD"] {
				t.Errorf("want %s but actual %s", c.values["DB_PASSWORD"], cfg.DBPassword)
			}
			if cfg.DBSchema != c.values["DB_SCHEMA"] {
				t.Errorf("want %s but actual %s", c.values["DB_SCHEMA"], cfg.DBSchema)
			}
		})
		testClearEnvs(t, c.values)
	}

}

func testSetEnvs(t *testing.T, envmap map[string]string) {
	t.Helper()
	for key, value := range envmap {
		err := os.Setenv(key, value)
		if err != nil {
			t.Fatalf("err %s", err)
		}
	}
}

func testClearEnvs(t *testing.T, envmap map[string]string) {
	t.Helper()
	for key := range envmap {
		err := os.Setenv(key, "")
		if err != nil {
			t.Fatalf("err %s", err)
		}
	}
}
