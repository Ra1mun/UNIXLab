package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

type Task struct {
	TaskID    int     `json:"task_id"`
	TaskType  string  `json:"task_type"`
	Timestamp float64 `json:"timestamp"`
}

func generateTask() Task {
	taskTypes := []string{"process_data", "generate_report", "send_notification", "update_database"}
	return Task{
		TaskID:    rand.Intn(9000) + 1000,
		TaskType:  taskTypes[rand.Intn(len(taskTypes))],
		Timestamp: float64(time.Now().UnixNano()) / 1e9,
	}
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Kafka configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Create producer
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	// Handle OS signals for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Producer started. Generating tasks...")

	// Generate and send tasks
	for {
		select {
		case <-signals:
			fmt.Println("Shutting down producer...")
			return
		default:
			task := generateTask()
			taskJSON, err := json.Marshal(task)
			if err != nil {
				log.Printf("Error marshaling task: %s", err)
				continue
			}

			msg := &sarama.ProducerMessage{
				Topic: "tasks",
				Value: sarama.StringEncoder(taskJSON),
			}

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Failed to send message: %s", err)
			} else {
				fmt.Printf("Sent task: %s (partition=%d, offset=%d)\n", taskJSON, partition, offset)
			}

			time.Sleep(2 * time.Second)
		}
	}
} 