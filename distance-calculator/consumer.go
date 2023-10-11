package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/kkboranbay/toll-calculator/aggregator/client"
	"github.com/kkboranbay/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggrClient  client.Client
}

func NewKafkaConsumer(topic string, srv CalculatorServicer, aggrClient client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:    c,
		calcService: srv,
		aggrClient:  aggrClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		// timeout may be set to -1 for indefinite wait.
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume errors %s", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization errors %s", err)
			continue
		}

		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation errors %s", err)
			continue
		}
		_ = distance

		req := &types.AggregatorRequest{
			ObuID: int32(data.OBUID),
			Value: distance,
			Unix:  time.Now().UnixNano(),
		}
		if err := c.aggrClient.Aggregate(context.Background(), req); err != nil {
			logrus.Errorf("aggregate error %s", err)
			continue
		}
	}
}
