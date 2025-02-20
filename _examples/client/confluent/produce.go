package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	kafkaConfig *kafka.ConfigMap
	producer    *kafka.Producer
	admin       *kafka.AdminClient
	topic       = "test-topic"
)

func init() {
	kafkaConfig = &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "test-client",
		"acks":              "all"}

	var err error
	producer, err = kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatal(err)
	}

	admin, err = kafka.NewAdminClient(kafkaConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// get metadata
	defer admin.Close()

	metadata, err := admin.GetMetadata(&topic, false, 1000)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v", metadata)

	// produce message
	deliveryChan := make(chan kafka.Event, 1)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("message")},
		deliveryChan)

	if err != nil {
		panic(err)
	}
	<-deliveryChan
}
