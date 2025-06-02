package main

import (
	_ "context"
	"fmt"
	"log"
	"os"
	"os/signal"
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

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Worker %s processing task: %s\n", hostname, string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Fatalf("Failed to create partition consumer: %s", err)
		case <-signals:
			fmt.Println("Shutting down worker...")
			return
		}
	}
}
