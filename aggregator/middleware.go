package main

import (
	"time"

	"github.com/kkboranbay/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())

	err = l.next.AggregateDistance(distance)
	return
}

func (l *LogMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	var (
		totalDistance float64
		totalAmount   float64
	)

	defer func(start time.Time) {
		if invoice != nil {
			totalDistance = invoice.TotalDistance
			totalAmount = invoice.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":          time.Since(start),
			"err":           err,
			"obuID":         obuID,
			"totalDistance": totalDistance,
			"totalAmount":   totalAmount,
		}).Info("CalculateInvoice")
	}(time.Now())

	invoice, err = l.next.CalculateInvoice(obuID)
	return
}
