package main

import (
	_ "context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("tasks", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to create partition consumer: %s", err)
	}
	defer partitionConsumer.Close()

	// Handle OS signals for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Get hostname for worker identification
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	fmt.Printf("Worker %s started and waiting for tasks...\n", hostname)

	// Process messages
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Worker %s processing task: %s\n", hostname, string(msg.Value))
			// Simulate task processing
			// In a real application, you would implement actual task processing logic here
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error: %s\n", err.Error())
		case <-signals:
			fmt.Println("Shutting down worker...")
			return
		}
	}
}
