package main

import (
	"context"
	"flag"
	"log"
	"time"

	Kafka "service-a/internal/kafka"
	pb "service-a/internal/server/summation"

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

	// ---------------------- Set up Kafka writer ----------------------
	writer := Kafka.NewKafkaWriter("user-events")
	if writer == nil {
		log.Fatal("Failed to create Kafka writer")
	}
	defer writer.Close()

	// ---------------------- Set up gRPC connection ----------------------
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewSummationServiceClient(conn)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// ---------------------- Make the gRPC call ----------------------

	log.Printf("Sending gRPC request with numbers: %d and %d", *numA, *numB)
	result, err := client.CalculateSum(ctx, &pb.SummationRequest{A: int32(*numA), B: int32(*numB)})
	if err != nil {
		log.Fatalf("Could not add numbers: %v", err)
	}
	log.Printf("Received sum from gRPC server: %d", result.Result)

	// ---------------------- Send result to Kafka ----------------------
	err = Kafka.SendMessage(writer, result.Result, ctx)
	if err != nil {
		log.Fatalf("Failed to send message to Kafka: %v", err)
	}
}
