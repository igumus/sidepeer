package grpc

import (
	"context"
	"io"
	"log"

	"net"

	g "google.golang.org/grpc"
)

type GrpcServer interface {
	Start() error
	Server() *g.Server
	io.Closer
}

type grpcServer struct {
	app      *g.Server
	addr     string
	listener net.Listener
}

func NewGrpcServer(ctx context.Context, opts ...GrpcServerOption) (GrpcServer, error) {
	cfg, cfgErr := applyOptions(opts...)
	if cfgErr != nil {
		return nil, cfgErr
	}

	server := g.NewServer()

	ret := &grpcServer{
		app:  server,
		addr: cfg.String(),
	}

	return ret, nil
}

func (s *grpcServer) Server() *g.Server {
	return s.app
}

func (gs *grpcServer) Start() error {
	listen, listenErr := net.Listen("tcp", gs.addr)
	if listenErr != nil {
		return listenErr
	}
	gs.listener = listen
	go func() {
		log.Printf("info: starting grpc server at: %s\n", gs.addr)
		gs.app.Serve(gs.listener)
	}()
	return nil
}

func (gs *grpcServer) Close() error {
	log.Printf("info: graceful shutdown of grpc server")
	gs.app.GracefulStop()
	return nil
}
