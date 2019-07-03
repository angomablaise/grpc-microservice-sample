package server_test

import (
	"testing"

	config "github.com/smockoro/grpc-microservice-sample/pkg/config/user"
	server "github.com/smockoro/grpc-microservice-sample/pkg/server/user"
)

func TestConnectDB(t *testing.T) {
	cases := []struct {
		name       string
		cfg        *config.Config
		errorIsNil bool
	}{
		{name: "not loss data", cfg: &config.Config{
			DBUser:     "dbuser",
			DBPassword: "password",
			DBHost:     "host.com",
			DBSchema:   "schema"}, errorIsNil: true},
		{name: "loss db user", cfg: &config.Config{
			DBUser:     "",
			DBPassword: "password",
			DBHost:     "host.com",
			DBSchema:   "schema"}, errorIsNil: false},
		{name: "loss db password", cfg: &config.Config{
			DBUser:     "dbuser",
			DBPassword: "",
			DBHost:     "host.com",
			DBSchema:   "schema"}, errorIsNil: false},
		{name: "loss db host", cfg: &config.Config{
			DBUser:     "dbuser",
			DBPassword: "password",
			DBHost:     "",
			DBSchema:   "schema"}, errorIsNil: false},
		{name: "loss db schema", cfg: &config.Config{
			DBUser:     "dbuser",
			DBPassword: "password",
			DBHost:     "host.com",
			DBSchema:   ""}, errorIsNil: false},
	}

	for _, c := range cases {
		c := c // cascading
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if _, err := server.ConnectDB(c.cfg); (err != nil) == c.errorIsNil {
				if c.errorIsNil {
					t.Errorf("wanted CoonectDB(%s) is nil. but %s", c.cfg, err)
				} else {
					t.Errorf("wanted CoonectDB(%s) is not nil. but %s", c.cfg, err)
				}
			}
		})
	}
}
