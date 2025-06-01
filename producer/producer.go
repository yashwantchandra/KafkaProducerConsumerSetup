package main

import (
	"log"
	"strconv"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Create a Kafka writer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Replace with your broker address
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// Delivery report callback
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Message delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					log.Printf("Message delivered to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := "my-topic"

	// Produce messages
	for i := 0; i < 10; i++ {
		value := []byte("Hello Kafka " + strconv.Itoa(i))
		// Produce the message to the topic
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value,
		}, nil)
	}

	// Wait for message deliveries
	producer.Flush(15 * 1000) // Wait up to 15 seconds for deliveries
	log.Println("Messages produced")
}