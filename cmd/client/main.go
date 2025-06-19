package main

import (
	"context"
	"flag"
	"log"
	"time"

	db "service-a/internal/database"
	Kafka "service-a/internal/kafka"
	"service-a/internal/outbox"
	pb "service-a/internal/server/summation"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultFirstNumber  = 10
	defaultSecondNumber = 20
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	numA = flag.Int("a", defaultFirstNumber, "First number to add")
	numB = flag.Int("b", defaultSecondNumber, "Second number to add")
)

func main() {
	flag.Parse()

	// ---------------------- Set up Database connection ----------------------
	database := db.INIT_DB()
	if database == nil {
		log.Fatal("Failed to initialize database connection")
	}
	defer database.Close()

	// ---------------------- Set up Outbox components ----------------------

	// Initialize repository
	repo := outbox.NewRepository(database)

	// Initialize Kafka writer (will be used by OutboxPublisher)
	writer := Kafka.NewKafkaWriter("user-events")
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

	// ---------------------- Set up gRPC connection ----------------------
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewSummationServiceClient(conn)

	// Create context with timeout for the gRPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// ---------------------- Make the gRPC call ----------------------
	log.Printf("Sending gRPC request with numbers: %d and %d", *numA, *numB)
	result, err := client.CalculateSum(ctx, &pb.SummationRequest{A: int32(*numA), B: int32(*numB)})
	if err != nil {
		log.Fatalf("Could not add numbers: %v", err)
	}
	log.Printf("Received sum from gRPC server: %d", result.Result)

	// ---------------------- Save result to Outbox ----------------------

	err = repo.SaveOutbox(ctx, outbox.NewOutbox(result.Result))
	if err != nil {
		log.Fatalf("Failed to save to outbox: %v", err)
	}
	log.Println("Successfully saved result to outbox table")

	// Sleep for a short duration to allow the outbox publisher to process the message
	time.Sleep(2 * time.Second)
}
