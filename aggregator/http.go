package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kkboranbay/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type ApiError struct {
	Code int
	Err  error
}

// Error implements the Error interface.
func (e ApiError) Error() string {
	return e.Err.Error()
}

type HttpFunc func(http.ResponseWriter, *http.Request) error

func makeHttpFuncHandler(fn HttpFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(ApiError); ok {
				writeJSON(w, apiErr.Code, map[string]string{"error": apiErr.Error()})
			}
		}
	}
}

type HttpMetricHandler struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func NewHttpMetricHandler(reqName string) *HttpMetricHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})

	errCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "err_counter"),
		Name:      "aggregator",
	})

	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &HttpMetricHandler{
		reqCounter: reqCounter,
		errCounter: errCounter,
		reqLatency: reqLatency,
	}
}

func (h *HttpMetricHandler) instrument(next HttpFunc) HttpFunc {
	// instrument method which returns an anonymous function
	// that serves as an HTTP handler. This anonymous function is the one that gets
	// executed when a request is made to your routes.
	// This anonymous function is what you register as the HTTP handler for your routes.

	// If you were to place h.reqCounter.Inc() outside of the returned anonymous function,
	// it would only be executed once when the instrument method is called, and
	// it would not increment with each request.
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
				"err":     err,
			}).Info()
			h.reqLatency.Observe(latency)
			h.reqCounter.Inc()
			if err != nil {
				h.errCounter.Inc()
			}
		}(time.Now())

		err = next(w, r)
		return err
	}
}

func handleAggregate(srv Aggregator) HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "POST" {
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid HTTP method %s", r.Method),
			}
		}

		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  err,
			}
		}

		fmt.Printf("%+v", distance)

		if err := srv.AggregateDistance(distance); err != nil {
			return ApiError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}

		return writeJSON(w, http.StatusOK, map[string]string{"msg": "ok"})
	}
}

func handleGetInvoice(srv Aggregator) HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" {
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid HTTP method %s", r.Method),
			}
		}

		obuIDParam := r.URL.Query().Get("obu")
		if obuIDParam == "" {
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("missing OBU parameter"),
			}
		}
		obuID, err := strconv.Atoi(obuIDParam)
		if err != nil {
			return ApiError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid OBU parameter"),
			}
		}

		inv, err := srv.CalculateInvoice(obuID)
		if err != nil {
			return ApiError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}

		return writeJSON(w, http.StatusOK, inv)
	}
}
