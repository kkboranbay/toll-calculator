package main

import (
	"time"

	"github.com/kkboranbay/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) CalculateDistance(data types.OBUData) (distance float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"error":    err,
			"distance": distance,
		}).Info("Calculation")
	}(time.Now())

	distance, err = l.next.CalculateDistance(data)
	return
}
