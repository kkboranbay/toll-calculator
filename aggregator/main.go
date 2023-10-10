package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/kkboranbay/toll-calculator/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpAddr", ":3000", "the listen address of the HTTP server")
	grpcListenAddr := flag.String("grpcAddr", ":3001", "the listen address of the GRPC server")

	flag.Parse()

	store := NewMemoryStore()
	srv := NewInvoiceAggregator(store)
	srv = NewLogMiddleware(srv)
	go makeGRPCTransport(*grpcListenAddr, srv)
	makeHTTPTransport(*httpListenAddr, srv)
}

func makeGRPCTransport(listenAddr string, srv Aggregator) error {
	fmt.Println("GRPC Transport running on port ", listenAddr)
	// Make a TCP listeners
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer lis.Close()

	// Make a new GRPC native server with (options)
	server := grpc.NewServer()
	// Register (OUR) GRPC server implementation to the GRPC package.
	types.RegisterAggregatorServer(server, NewGRPCServer(srv))
	return server.Serve(lis)
}

func makeHTTPTransport(listenAddr string, srv Aggregator) {
	fmt.Println("HTTP Transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(srv))
	http.HandleFunc("/invoice", handleGetInvoice(srv))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(srv Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

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
