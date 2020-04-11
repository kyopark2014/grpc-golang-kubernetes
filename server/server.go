package main

import (
	"context"
	"grpc-golang-server/log"
	"grpc-golang-server/proto"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {
	log.I("Start the service...")
	listener, grpcErr := net.Listen("tcp", ":4040")
	if grpcErr != nil {
		log.E("Localisation function Failed to listen : %v", grpcErr)
		os.Exit(1)
	}

	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) Add(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	log.D("Server: add()...")
	a, b := request.GetA(), request.GetB()

	result := a + b

	log.D("%v + %v = %v", a, b, result)

	return &proto.Response{Result: result}, nil
}

func (s *server) Multiply(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	log.D("Server: multiply)...")
	a, b := request.GetA(), request.GetB()

	result := a * b

	log.D("%v x %v = %v", a, b, result)

	return &proto.Response{Result: result}, nil
}
