package example

import (
	"fmt"

	"github.com/samber/do"

	"github.com/holydocs/servicefile/examples/do/client/firebase"
	"github.com/holydocs/servicefile/examples/do/database/postgres"
)

/*
service:name Example
description: Example service for exampling stuff.
*/
type Service struct {
	db *postgres.Connection
	c  *firebase.Client
}

func NewService(i *do.Injector) *Service {
	return &Service{
		db: do.MustInvoke[*postgres.Connection](i),
		c:  do.MustInvoke[*firebase.Client](i),
	}
}

func (svc *Service) Do() error {
	err := svc.db.Query()
	if err != nil {
		return fmt.Errorf("querying db: %w", err)
	}

	err = svc.c.Request()
	if err != nil {
		return fmt.Errorf("requesting client: %w", err)
	}

	return nil
}
