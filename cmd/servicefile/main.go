package main

import (
	"log"

	"github.com/holydocs/servicefile/internal/api/cli"
)

func main() {
	cmd := cli.Command()

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
