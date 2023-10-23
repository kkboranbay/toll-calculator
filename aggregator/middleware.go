package main

import (
	"time"

	"github.com/kkboranbay/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter

	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter

	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram

	next Aggregator
}

func NewMetricsMiddleware(next Aggregator) Aggregator {
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})

	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})

	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})

	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})

	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &MetricsMiddleware{
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
		next:           next,
	}
}

func (l *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		l.reqLatencyAgg.Observe(time.Since(start).Seconds())
		l.reqCounterAgg.Inc()

		if err != nil {
			l.errCounterAgg.Inc()
		}
	}(time.Now())

	err = l.next.AggregateDistance(distance)
	return
}

func (l *MetricsMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		l.reqLatencyCalc.Observe(time.Since(start).Seconds())
		l.reqCounterCalc.Inc()

		if err != nil {
			l.errCounterCalc.Inc()
		}
	}(time.Now())

	invoice, err = l.next.CalculateInvoice(obuID)
	return
}

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
