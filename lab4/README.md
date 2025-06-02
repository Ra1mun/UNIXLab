# Microservice Task Processing System

This is a microservice architecture that demonstrates horizontal scaling and message processing using Apache Kafka as a message broker.

## Architecture

The system consists of the following components:

1. **Producer Service** (Python)
   - Generates random tasks and sends them to Kafka
   - Tasks include different types: process_data, generate_report, send_notification, update_database

2. **Worker Service** (Go)
   - Consumes tasks from Kafka
   - Horizontally scalable (3 replicas by default)
   - Each worker processes tasks independently

3. **Message Broker** (Apache Kafka)
   - Handles message distribution between producer and workers
   - Ensures reliable message delivery

## Prerequisites

- Docker
- Docker Compose

## Running the System

1. Build and start all services:
   ```bash
   docker-compose up --build
   ```

2. To scale the number of workers:
   ```bash
   docker-compose up --scale worker=5
   ```

3. To stop all services:
   ```bash
   docker-compose down
   ```

## Monitoring

- The producer will output messages when it sends tasks
- Each worker will output messages when it receives and processes tasks
- You can see the worker hostnames in the logs to verify load balancing

## Architecture Details

- The producer generates tasks every 2 seconds
- Workers automatically consume tasks from Kafka
- Kafka ensures that each task is processed by only one worker
- The system is horizontally scalable - you can add more workers without modifying the code 