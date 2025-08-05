package main

import (
	"context"
	"log"

	"net/http"
	"os"
	API "service-a/cmd/api"
	"service-a/cmd/api/connection"
	DB "service-a/internal/database"
	kafkaStructure "service-a/internal/kafka"
	"service-a/internal/metrics"
	"service-a/internal/outbox"
	"service-a/internal/server"

	// "time"
)

// Get specific partition based on hostname
func ExtractPartition(serviceHostname string) int {
	var partition int = 0;
	switch serviceHostname {
	case "service-a-1":
		partition = 0
	case "service-a-2":
		partition = 1
	case "service-a-3":
		partition = 2
	}
	log.Printf("Extracted partition %d for hostname %s", partition, serviceHostname)
	return partition
}

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

	// ------------ Initialize Kafka writer with specific partition based on hostname ------------

	serviceHostname, _ := os.Hostname()
	var partition int = ExtractPartition(serviceHostname)

	log.Printf("Service hostname: %s, assigned to partition: %d", serviceHostname, partition)
	writer := kafkaStructure.NewKafkaWriterWithPartition("user-events", partition)
	if writer == nil {
		log.Fatal("Failed to create Kafka writer")
	}
	defer writer.Publisher.Close()

	// Initialize OutboxPublisher with 3-second check interval
	// publisher := outbox.NewOutboxPublisher(repo, writer, 3*time.Second)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the OutboxPublisher in a goroutine
	// go publisher.Start(ctx)

	// Start the HTTP server for Prometheus metrics
	metrics.StartMetricsServer(ctx, 9091)

	// -------- Open the GRPC Connection and define the API --------

	client, conn, err := connection.GRPC_Connection()
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// ------ API ------
	http.HandleFunc("/sum", API.SummationRequest(client))

	// Start the gRPC server in a goroutine so it doesn't block
	go func() {
		if err := server.StartServerWithOutbox(50051, repo); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Start HTTP server for the API
	log.Printf("Starting HTTP server on port 8080, hostname: %s", serviceHostname)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
