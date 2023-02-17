package server

import (
	"context"
	"log"
	"storage/modules/database"
	"storage/modules/utils"
	"storage/proto/pb"
)

type Server struct {
	es *database.Elastic
	pb.UnimplementedStatsServer
}

func (s *Server) SaveStats(ctx context.Context, in *pb.SaveStatsRequest) (*pb.SaveStatsResponse, error) {
	log.Println("Save Stats REQ received: ", len(in.Dataset))

	err := s.es.BulkWrite(in.Dataset)
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

	esres, err := s.es.QueryStats(int(in.FromDate), int(in.ToDate), int(in.Interval))
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
