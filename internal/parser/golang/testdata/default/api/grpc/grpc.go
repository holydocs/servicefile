package grpc

import (
	"github.com/samber/do"

	"github.com/holydocs/servicefile/examples/do/service/example"
)

/*
service:replies
description: Provides user management APIs to other services
technology:grpc-server
proto:grpc
*/
type Server struct {
	svc *example.Service
}

func NewServer(i *do.Injector) *Server {
	return &Server{
		svc: do.MustInvoke[*example.Service](i),
	}
}

func (srv *Server) Serve() error {
	return nil
}
