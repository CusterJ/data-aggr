package server

import (
	"context"
	"fmt"
	"log"
	fr "reader/modules/file_reader"
	"reader/modules/utils"
	"reader/proto/pb"
	"time"
)

type Server struct {
	sc pb.StatsClient
}

func NewServer(sc pb.StatsClient) *Server {
	return &Server{
		sc: sc,
	}
}

func (s *Server) SaveFile(filename string) error {

	// read strings one by one
	// fr.ReadFileByLines(filename)

	fmt.Println("Reading file: ", filename)
	data, err := fr.ReadFullFile(filename)
	utils.Check(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	rec := new(pb.SaveStatsRequest)

	var dataset []*pb.Dataset

	for _, v := range data {
		dataset = append(dataset, &pb.Dataset{
			Id:     v.ID,
			Time:   int32(v.Time),
			Signal: v.Signal,
			Data:   v.Data,
		})
	}

	rec.Dataset = dataset

	// grpc req
	res, err := s.sc.SaveStats(ctx, rec)
	utils.Check(err)
	log.Println("SaveStats result: ", res)
	return nil
}
