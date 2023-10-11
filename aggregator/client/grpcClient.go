package client

import (
	"github.com/kkboranbay/toll-calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGrpcClient(endpoint string) (*GrpcClient, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GrpcClient{
		Endpoint:         endpoint,
		AggregatorClient: c,
	}, nil
}
