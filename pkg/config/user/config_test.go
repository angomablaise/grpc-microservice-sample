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
			"GRPC_PORT":  "9000",
			"DBHost":     "localhost:9000",
			"DBUser":     "connect_user",
			"DBPassword": "password",
			"DBSchema":   "shema"}, errorIsNil: false},
		{name: "GRPC_PORT is lost", values: map[string]string{
			"GRPC_PORT":  "",
			"DBHost":     "localhost:9000",
			"DBUser":     "connect_user",
			"DBPassword": "password",
			"DBSchema":   "shema"}, errorIsNil: false},
		{name: "DBHost is lost", values: map[string]string{
			"GRPC_PORT":  "9000",
			"DBHost":     "",
			"DBUser":     "connect_user",
			"DBPassword": "password",
			"DBSchema":   "shema"}, errorIsNil: false},
		{name: "DBUser is lost", values: map[string]string{
			"GRPC_PORT":  "9000",
			"DBHost":     "localhost:9000",
			"DBUser":     "",
			"DBPassword": "password",
			"DBSchema":   "shema"}, errorIsNil: false},
		{name: "DBPassword is lost", values: map[string]string{
			"GRPC_PORT":  "9000",
			"DBHost":     "localhost:9000",
			"DBUser":     "connect_user",
			"DBPassword": "",
			"DBSchema":   "shema"}, errorIsNil: false},
		{name: "DBSchema is lost", values: map[string]string{
			"GRPC_PORT":  "9000",
			"DBHost":     "localhost:9000",
			"DBUser":     "connect_user",
			"DBPassword": "password",
			"DBSchema":   ""}, errorIsNil: false},
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
