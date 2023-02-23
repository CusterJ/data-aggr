package server_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reader/server"
	"strings"
	"testing"

	"github.com/CusterJ/data-aggr/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGenerateHandlerFullFile(t *testing.T) {
	defer os.Remove("data.json")

	var rpcPort string = ":8090"

	conn, err := grpc.Dial(rpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.Dial(grpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewStatsClient(conn)

	testCases := []struct {
		name     string
		length   string
		expected string
	}{
		{
			name:     "Valid_length_100",
			length:   "100",
			expected: "Generate and read file with lengh: 100 -- done",
		},
		{
			name:     "Valid_length_500",
			length:   "500",
			expected: "Generate and read file with lengh: 500 -- done",
		},
		{
			name:     "Invalid_length_5001",
			length:   "50001",
			expected: "Length must be greater than zero and not greater than 50.000",
		},
		{
			name:     "Invalid_length_-100",
			length:   "-100",
			expected: "Length must be greater than zero and not greater than 50.000",
		},
		{
			name:     "Invalid_length_empty",
			length:   "",
			expected: "To generate file add length parameter: /generate?length=100",
		},
		{
			name:     "Invalid_length_abc",
			length:   "abc",
			expected: "Generate and read file error! Add a number for parameter length",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			srv := server.NewServer(client)
			req, err := http.NewRequest("GET", "/generate?length="+tc.length, nil)
			if err != nil {
				log.Println("http.NewRequest error")
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(srv.GenerateHandler)
			handler.ServeHTTP(rr, req)

			// Status code check
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			// response check
			body := strings.TrimSpace(rr.Body.String())
			if !strings.Contains(body, tc.expected) {
				t.Errorf("handler returned unexpected body: got %q", body)
			}

		})
	}
}
