package main

import (
    "log"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
    // Create a Kafka reader
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092", // Replace with your broker address
        "group.id":          "my-group",         // Consumer group ID
        "auto.offset.reset": "earliest",         // Start reading at the earliest message
    })
    if err != nil {
        log.Fatalf("Failed to create consumer: %s", err)
    }
    defer consumer.Close()

    topic := "my-topic"
    consumer.SubscribeTopics([]string{topic}, nil)

    // Read messages
    for {
        msg, err := consumer.ReadMessage(-1) // -1 means wait indefinitely for a message
        if err == nil {
            log.Printf("Received message: %s", string(msg.Value))
        } else {
            log.Printf("Error reading message: %v", err)
        }
    }
}