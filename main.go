package main

import (
	"context"
	"log"
	DB "service-a/internal/database"
	kafkaStructure "service-a/internal/kafka"
	"service-a/internal/outbox"
	"service-a/internal/server"
	"time"
)

func main() {
	log.Println("Starting Summation Service")

	// Initialize database connection
	db := DB.INIT_DB()
	if db == nil {
		log.Fatal("Failed to initialize database connection")
	}
	defer db.Close()

	// Initialize outbox repository
	repo := outbox.NewRepository(db)

	// Initialize Kafka writer
	writer := kafkaStructure.NewKafkaWriter("user-events")
	if writer == nil {
		log.Fatal("Failed to create Kafka writer")
	}
	defer writer.Publisher.Close()

	// Initialize OutboxPublisher with 3-second check interval
	publisher := outbox.NewOutboxPublisher(repo, writer, 3*time.Second)

	// Start the OutboxPublisher in a goroutine
	publisherCtx, publisherCancel := context.WithCancel(context.Background())
	defer publisherCancel()
	go publisher.Start(publisherCtx)

	// Create a new instance of the SummationServer with outbox repository
	if err := server.StartServerWithOutbox(50051, repo); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
