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
)

type HttpMetricHandler struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func NewHttpMetricHandler(reqName string) *HttpMetricHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})

	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &HttpMetricHandler{
		reqCounter: reqCounter,
		reqLatency: reqLatency,
	}
}

func (h *HttpMetricHandler) instrument(next http.HandlerFunc) http.HandlerFunc {
	// instrument method which returns an anonymous function
	// that serves as an HTTP handler. This anonymous function is the one that gets
	// executed when a request is made to your routes.
	// This anonymous function is what you register as the HTTP handler for your routes.

	// If you were to place h.reqCounter.Inc() outside of the returned anonymous function,
	// it would only be executed once when the instrument method is called, and
	// it would not increment with each request.
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			h.reqLatency.Observe(time.Since(start).Seconds())
		}(time.Now())

		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleAggregate(srv Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		fmt.Printf("%+v", distance)

		if err := srv.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func handleGetInvoice(srv Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuIDParam := r.URL.Query().Get("obu")
		if obuIDParam == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU parameter"})
			return
		}
		obuID, err := strconv.Atoi(obuIDParam)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU parameter"})
			return
		}

		inv, err := srv.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, inv)
	}
}
