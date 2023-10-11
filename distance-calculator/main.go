package main

import (
	"log"

	"github.com/kkboranbay/toll-calculator/aggregator/client"
)

const kafkaTopic = "obudata"
const aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"

// Transport could be HTTP, GRPC, Kafka -> attach business logic to this transport

func main() {
	srv := NewCalculateService()
	srv = NewLogMiddleware(srv)

	httpClient := client.NewHttpClient(aggregatorEndpoint)
	// grpcClient, err := client.NewGrpcClient(":3001")
	// if err != nil {
	// log.Fatal(err)
	// }

	// kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, srv, grpcClient)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, srv, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
