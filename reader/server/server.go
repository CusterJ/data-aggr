package server

import (
	"context"
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

// Read full file and send it like slice of domain.FooData
func (s *Server) SaveFile(filename string) error {
	// log.Println("Reading file: ", filename)
	data, err := fr.ReadFullFile(filename)
	utils.Check(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	rec := new(pb.SaveStatsRequest)

	var dataset []*pb.FooData

	for i, v := range data {
		dataset = append(dataset, &pb.FooData{
			Id:     v.ID,
			Time:   int32(v.Time),
			Signal: v.Signal,
			Data:   v.Data,
		})

		if i%5001 == 0 {
			rec.FooData = dataset
			_, err := s.sc.SaveStats(ctx, rec)
			utils.Check(err)
			dataset = nil
		}
	}

	rec.FooData = dataset

	// grpc req
	_, err = s.sc.SaveStats(ctx, rec)
	utils.Check(err)
	// log.Println("SaveStats result: ", res)
	return nil
}

// Read file by tokens and send them with stream to save
func (s *Server) SaveFileStream(filename string) error {
	// log.Println("Reading file: ", filename)

	rpc, err := s.sc.SaveStatsStream(context.TODO())
	utils.Check(err)

	defer rpc.CloseSend()

	stream := fr.NewJsonStream()
	go func() {
		for data := range stream.Watch() {
			if data.Error != nil {
				log.Println(data.Error)
			}

			err := rpc.Send(&pb.FooData{
				Id:     data.FooData.ID,
				Time:   int32(data.FooData.Time),
				Signal: data.FooData.Signal,
				Data:   data.FooData.Data,
			})
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	stream.Start(filename)

	// _, err = rpc.CloseAndRecv()
	// log.Println("rpc.CloseAndRecv response: ", res)
	// if err != nil {
	// 	return err
	// }
	return nil
}
