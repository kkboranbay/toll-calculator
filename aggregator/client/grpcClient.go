package client

import (
	"context"

	"github.com/kkboranbay/toll-calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NewGrpcClient(endpoint string) (*GrpcClient, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GrpcClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GrpcClient) Aggregate(ctx context.Context, req *types.AggregatorRequest) error {
	_, err := c.client.Aggregate(ctx, req)
	return err
}
