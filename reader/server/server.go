package server

import (
	"context"
	"fmt"
	"log"
	fr "reader/modules/file_reader"
	"reader/modules/utils"

	"time"

	"github.com/CusterJ/data-aggr/proto/pb"
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

	var dataset []*pb.FooData

	for _, v := range data {
		dataset = append(dataset, &pb.FooData{
			Id:     v.ID,
			Time:   int32(v.Time),
			Signal: v.Signal,
			Data:   v.Data,
		})
	}

	rec.FooData = dataset

	// grpc req
	res, err := s.sc.SaveStats(ctx, rec)
	utils.Check(err)
	log.Println("SaveStats result: ", res)
	return nil
}

func (s *Server) SaveFileStream(filename string) error {
	fmt.Println("Reading file: ", filename)

	rpc, err := s.sc.SaveStatsStream(context.Background())
	utils.Check(err)

	stream := fr.NewJsonStream()
	go func() {
		for data := range stream.Watch() {
			if data.Error != nil {
				log.Println(data.Error)
			}
			// log.Println(data.FooData.ID, ": ", data.FooData.Signal)
			err := rpc.Send(&pb.FooData{
				Id:     data.FooData.ID,
				Time:   int32(data.FooData.Time),
				Signal: data.FooData.Signal,
				Data:   data.FooData.Data,
			})
			utils.Check(err)
		}
	}()

	stream.Start(filename)

	res, err := rpc.CloseAndRecv()
	log.Println("rpc.CloseAndRecv response: ", res)
	utils.Check(err)

	return nil
}
