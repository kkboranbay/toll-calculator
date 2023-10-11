package client

import (
	"context"

	"github.com/kkboranbay/toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregatorRequest) error
}
