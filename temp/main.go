package main

import (
	"context"
	"log"
	"time"

	"github.com/kkboranbay/toll-calculator/aggregator/client"
	"github.com/kkboranbay/toll-calculator/types"
)

func main() {
	c, err := client.NewGrpcClient(":3001")
	if err != nil {
		log.Fatal(err)
	}

	req := &types.AggregatorRequest{
		ObuID: 111,
		Value: 77.77,
		Unix:  time.Now().UnixNano(),
	}

	if _, err := c.Aggregate(context.Background(), req); err != nil {
		log.Fatal(err)
	}
}
