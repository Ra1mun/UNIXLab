package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

func createConsumer(brokerList []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	return sarama.NewConsumer(brokerList, config)
}

func main() {
	brokerList := []string{"kafka:29092"}

	var consumer sarama.Consumer
	var err error
	for i := 0; i < 30; i++ {
		consumer, err = createConsumer(brokerList)
		if err == nil {
			break
		}
		log.Printf("Failed to create consumer, retrying in 1 second... (attempt %d/30)", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to create consumer after 30 attempts: %s", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("tasks", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to create partition consumer: %s", err)
	}
	defer partitionConsumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	fmt.Printf("Worker %s started and waiting for tasks...\n", hostname)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	isShuttingDown := false

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			if isShuttingDown {
				continue
			}
			wg.Add(1)
			go func(msg *sarama.ConsumerMessage) {
				defer wg.Done()
				fmt.Printf("Worker %s processing task: %s\n", hostname, string(msg.Value))
				// Simulate task processing
				time.Sleep(time.Second)
				fmt.Printf("Worker %s completed task: %s\n", hostname, string(msg.Value))
			}(msg)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error from partition consumer: %s", err)
		case <-signals:
			fmt.Println("Received shutdown signal. Stopping new task processing...")
			isShuttingDown = true

			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			select {
			case <-done:
				fmt.Println("All tasks completed, shutting down gracefully")
			case <-time.After(10 * time.Second):
				fmt.Println("Shutdown timeout reached, forcing exit")
			}
			return
		}
	}
}
