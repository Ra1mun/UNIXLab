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

func createProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	return sarama.NewSyncProducer(brokerList, config)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	brokerList := []string{"kafka:29092"}

	var producer sarama.SyncProducer
	var err error
	for i := 0; i < 30; i++ {
		producer, err = createProducer(brokerList)
		if err == nil {
			break
		}
		log.Printf("Failed to create producer, retrying in 1 second... (attempt %d/30)", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to create producer after 30 attempts: %s", err)
	}
	defer producer.Close()

	// Handle OS signals for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Producer started. Generating tasks...")

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
				producer.Close()
				producer, err = createProducer(brokerList)
				if err != nil {
					log.Printf("Failed to recreate producer: %s", err)
					time.Sleep(time.Second)
					continue
				}
			} else {
				fmt.Printf("Sent task: %s (partition=%d, offset=%d)\n", taskJSON, partition, offset)
			}

			time.Sleep(2 * time.Second)
		}
	}
}
