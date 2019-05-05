package main

import (
	"fmt"
	"os"

	"github.com/smockoro/grpc-microservice-sample/pkg/server/user"
)

func main() {
	if err := server.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
