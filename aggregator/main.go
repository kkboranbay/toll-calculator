package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/kkboranbay/toll-calculator/types"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "the listen address of the HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	srv := NewInvoiceAggregator(store)
	makeHTTPTransport(*listenAddr, srv)
}

func makeHTTPTransport(listenAddr string, srv Aggregator) {
	fmt.Println("HTTP Transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(srv))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(srv Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
