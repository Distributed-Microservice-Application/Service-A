package main

import (
	"context"
	"log"
	DB "service-a/internal/database"
	kafkaStructure "service-a/internal/kafka"
	"service-a/internal/metrics"
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

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the OutboxPublisher in a goroutine
	go publisher.Start(ctx)

	// Start the HTTP server for Prometheus metrics
	metrics.StartMetricsServer(ctx, 9091)

	// Create a new instance of the SummationServer with outbox repository
	if err := server.StartServerWithOutbox(50051, repo); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
