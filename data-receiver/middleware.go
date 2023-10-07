package main

import (
	"time"

	"github.com/kkboranbay/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

// Decorator Pattern
type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

// go get github.com/sirupsen/logrus
func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(start),
		}).Info("producing to Kafka")
	}(time.Now())
	return l.next.ProduceData(data)
}
