package main

import (
	"context"

	"github.com/kkboranbay/toll-calculator/types"
)

// GRPCServer is used to implement types.AggregatorServer.
type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	srv Aggregator
}

func NewGRPCServer(srv Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		srv: srv,
	}
}

// transport layer
// JSON -> types.Distance -> all done (same type)
// GRPC -> types.AggregateRequest -> type.Distance

// business layer ->business layer type (main type everyone needs to convert to)

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregatorRequest) error {
	distance := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return s.srv.AggregateDistance(distance)
}

// func (s *GRPCServer) AggregateDistance(distance types.Distance) error {
// 	return s.Aggregate()
// }
