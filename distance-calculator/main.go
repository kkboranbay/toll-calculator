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

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, srv, client.NewHttpClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
