package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"github.com/kkboranbay/toll-calculator/types"
)

var kafkaTopic = "obudata"

// go get github.com/confluentinc/confluent-kafka-go/kafka
func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  *kafka.Producer
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// Start another goroutine to check if we have delivered the data.
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &DataReceiver{
		// channel always block when its full
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Produce messages to topic (asynchronously)
	dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)

	return nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	// In order to create a WebSocket endpoint, we effectively need to upgrade an incoming connection from a standard HTTP endpoint
	// to a long-lasting WebSocket connection. In order to do this, we are going to be using some of the functionality from the
	// very cool gorilla/websocket package!
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	dr.conn = conn
	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("OBU client connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}

		if err := dr.produceData(data); err != nil {
			log.Println("kafka product error:", err)
		}

		// so when we run first, make receiver, then open another terminal, there run make obu
		// in recerver terminal we get some messages but after when channel is full, then its not get a new data
		// because channel is full, nobody is consuming!!!
		// dr.msgch <- data
	}
}
