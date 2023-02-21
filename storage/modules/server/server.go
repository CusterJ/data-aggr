package server

import (
	"context"
	"io"
	"log"
	"storage/modules/database"
	"storage/modules/utils"

	"github.com/CusterJ/data-aggr/proto/pb"
)

type Server struct {
	ES *database.Elastic
	pb.UnimplementedStatsServer
}

func (s *Server) SaveStatsStream(ps pb.Stats_SaveStatsStreamServer) error {
	log.Println("Save Stats Stream REQ received")

	count, recvCount := 0, 0
	var dataSet []*pb.FooData
	for {
		// Get a packet
		data, err := ps.Recv()
		if err == io.EOF {
			log.Printf("Received io.EOF: %#v", data)
			break
		}

		// check err while receiving stream data - ps.Recv
		if err != nil {
			log.Printf("EXIT => for loop stream recv error: %s \n ", err)

			ps.SendAndClose(&pb.SaveStatsResponse{
				Saved: false,
			})

			return err
		}

		// Concat packet to create final note
		dataSet = append(dataSet, &pb.FooData{
			Id:     data.Id,
			Time:   data.Time,
			Signal: data.Signal,
			Data:   data.Data,
		})
		count++
		recvCount++

		if count%1001 == 0 {
			err := s.ES.BulkWrite(dataSet)
			utils.Check(err)

			count = 0
			dataSet = nil
		}
	}

	err := s.ES.BulkWrite(dataSet)
	utils.Check(err)

	ps.SendAndClose(&pb.SaveStatsResponse{
		Saved: true,
	})

	log.Println("Recieved and saved stream data: ", recvCount)

	return nil
}

func (s *Server) SaveStats(ctx context.Context, in *pb.SaveStatsRequest) (*pb.SaveStatsResponse, error) {
	log.Println("Save Stats REQ received: ", len(in.FooData))

	err := s.ES.BulkWrite(in.FooData)
	utils.Check(err)
	if err != nil {
		return &pb.SaveStatsResponse{
			Saved: false,
		}, err
	}

	return &pb.SaveStatsResponse{
		Saved: true,
	}, nil
}

func (s *Server) GetStats(ctx context.Context, in *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	log.Println("GetStats REQ received: ", in)

	esres, err := s.ES.QueryStats(int(in.FromDate), int(in.ToDate), int(in.Interval))
	utils.Check(err)

	pbres := new(pb.GetStatsResponse)

	hb := []*pb.HistoBuckets{}
	sb := []*pb.Buckets{}

	for _, v := range esres.Aggregations.Histogram.HistoBuckets {
		for _, j := range v.SignalCount.Buckets {
			sb = append(sb, &pb.Buckets{
				Key:      j.Key,
				DocCount: uint32(j.DocCount),
			})
		}
		hb = append(hb, &pb.HistoBuckets{
			KeyAsString: v.KeyAsString,
			Key:         uint64(v.Key),
			DocCount:    uint32(v.DocCount),
			SignalCount: &pb.SignalCount{
				DocCountErrorUpperBound: uint32(v.SignalCount.DocCountErrorUpperBound),
				SumOtherDocCount:        uint32(v.SignalCount.SumOtherDocCount),
				Buckets:                 sb,
			},
		})
	}

	pbres.Aggrs = &pb.Aggrs{
		Took:     uint32(esres.Took),
		TimedOut: esres.TimedOut,
		Hits: &pb.Hits{Total: &pb.Total{
			Value:    uint32(esres.Hits.Total.Value),
			Relation: esres.Hits.Total.Relation,
		}},
		Aggregations: &pb.Aggregations{
			Histogram: &pb.Histogram{
				Buckets:  hb,
				Interval: esres.Aggregations.Histogram.Interval,
			},
		},
	}

	return pbres, nil
}
