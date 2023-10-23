package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kkboranbay/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	store := makeStore()
	srv := NewInvoiceAggregator(store)
	srv = NewMetricsMiddleware(srv)
	srv = NewLogMiddleware(srv)
	grpcListenAddr := os.Getenv("AGG_GRPC_PORT")
	httpListenAddr := os.Getenv("AGG_HTTP_PORT")

	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAddr, srv))
	}()
	log.Fatal(makeHTTPTransport(httpListenAddr, srv))
}

func makeGRPCTransport(listenAddr string, srv Aggregator) error {
	fmt.Println("GRPC Transport running on port ", listenAddr)
	// Make a TCP listeners
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer lis.Close()

	// Make a new GRPC native server
	server := grpc.NewServer()
	// Register (OUR) GRPC server implementation to the GRPC package.
	types.RegisterAggregatorServer(server, NewGRPCServer(srv))
	return server.Serve(lis)
}

func makeHTTPTransport(listenAddr string, srv Aggregator) error {
	fmt.Println("HTTP Transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(srv))
	http.HandleFunc("/invoice", handleGetInvoice(srv))
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(listenAddr, nil)
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

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type given %s", storeType)
		return nil
	}
}
