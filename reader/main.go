package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	fg "reader/modules/file_generator"
	"reader/modules/utils"
	"reader/server"

	"github.com/CusterJ/data-aggr/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// const bufferSize = 1024 * 1024
const port string = ":8090"
const defaultFileLength int = 10

func main() {
	fileName := flag.String("r", "", "a string")
	generateFile := flag.Bool("g", false, "generate file or not")
	generateFilelength := flag.Int("l", defaultFileLength, "generate file with this length")

	flag.Parse()
	fmt.Println("Flags parsed")

	// generate new json file
	if *generateFile {
		fmt.Println("Generating file: ", *generateFilelength)
		err := fg.GenerateNewFile(*generateFilelength)
		utils.Check(err)
	}

	//grpc client connection
	rpcPort := os.Getenv("RPC_PORT")
	if rpcPort == "" {
		rpcPort = port
	}

	fmt.Println("GRPC connecting: ", rpcPort)

	conn, err := grpc.Dial(rpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.Dial(grpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewStatsClient(conn)

	// Server init
	srv := server.NewServer(client)

	// read json file
	if *fileName != "" {
		err := srv.SaveFile(*fileName)
		utils.Check(err)
	}

	httpServer := http.Server{
		Addr:              ":8002",
		Handler:           nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}

	fmt.Println("HTTP server starting on port: ", httpServer.Addr)

	http.HandleFunc("/generate", srv.GenerateHandler)
	http.HandleFunc("/stats", srv.GetStats)
	if err := httpServer.ListenAndServe(); err != nil {
		fmt.Println("httpServer error - exit program")
		os.Exit(1)
	}
}
