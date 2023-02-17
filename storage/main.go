package main

import (
	"fmt"
	"log"
	"net"
	"storage/modules/server"
	"storage/proto/pb"

	"google.golang.org/grpc"
)

const port string = ":8090"

func main() {

	// start grpc server
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	fmt.Println("Storage main func start. gRPC server start on port: ", port)

	pb.RegisterStatsServer(s, &server.Server{})
	if err := s.Serve(listener); err != nil {
		log.Printf("failed to serve: %v/n", err)
	}
}
