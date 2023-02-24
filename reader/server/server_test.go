package server_test

import (
	"log"
	"os"
	"reader/modules/file_generator"
	"reader/server"
	"testing"

	"github.com/CusterJ/data-aggr/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func BenchmarkSaveFileStream(b *testing.B) {
	filename := "big.json"
	file_generator.GenerateFileBySize(1024*1024*51, filename)
	defer os.Remove(filename)

	var rpcPort string = ":8090"

	tls := insecure.NewBundle().TransportCredentials()
	conn, err := grpc.Dial(rpcPort, grpc.WithTransportCredentials(tls))
	// conn, err := grpc.Dial(rpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewStatsClient(conn)

	srv := server.NewServer(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := srv.SaveFileStream(filename)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkSaveFullFile(b *testing.B) {
	filename := "big.json"
	file_generator.GenerateFileBySize(1024*1024*51, filename)
	defer os.Remove(filename)

	var rpcPort string = ":8090"

	conn, err := grpc.Dial(rpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.Dial(rpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// defer conn.Close()

	client := pb.NewStatsClient(conn)

	srv := server.NewServer(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := srv.SaveFile(filename)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkSaveFileVsSaveFileStream(b *testing.B) {
	filename := "big.json"
	file_generator.GenerateFileBySize(1024*1024*50, filename)
	// defer os.Remove("data.json")

	var rpcPort string = ":8090"

	conn, err := grpc.Dial(rpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.Dial(rpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewStatsClient(conn)

	srv := server.NewServer(client)

	b.Run("SaveFile", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := srv.SaveFile(filename)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
		}
	})

	b.Run("SaveFileStream", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := srv.SaveFileStream(filename)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
		}
	})
}
