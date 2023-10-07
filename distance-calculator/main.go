package main

import (
	"log"
)

const kafkaTopic = "obudata"

// Transport could be HTTP, GRPC, Kafka -> attach business logic to this transport

func main() {
	srv := NewCalculateService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, srv)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
