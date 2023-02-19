package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"storage/modules/database"
	"storage/modules/server"

	"github.com/CusterJ/data-aggr/proto/pb"
	"google.golang.org/grpc"
)

const port string = ":8090"

func main() {
	rpcPort := os.Getenv("RPC_PORT")
	if rpcPort == "" {
		rpcPort = port
	}

	// start grpc server
	listener, err := net.Listen("tcp", rpcPort)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	fmt.Println("Storage main func start. gRPC server start on port: ", rpcPort)

	elastic := database.NewElastic()
	pb.RegisterStatsServer(s, &server.Server{
		ES:                       elastic,
		UnimplementedStatsServer: pb.UnimplementedStatsServer{},
	})

	if err := s.Serve(listener); err != nil {
		log.Printf("failed to serve: %v/n", err)
	}
}
