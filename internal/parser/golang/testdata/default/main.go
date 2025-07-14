package main

import (
	"github.com/denchenko/servicefile/examples/do/api/grpc"
	"github.com/denchenko/servicefile/examples/do/client/firebase"
	"github.com/denchenko/servicefile/examples/do/database/postgres"
	"github.com/denchenko/servicefile/examples/do/service/example"
	"github.com/samber/do"
)

func main() {
	i := do.New()

	do.ProvideValue(i, postgres.NewConnection)
	do.ProvideValue(i, firebase.NewClient)
	do.ProvideValue(i, example.NewService)

	server := grpc.NewServer(i)
	server.Serve()

	i.Shutdown()
}
