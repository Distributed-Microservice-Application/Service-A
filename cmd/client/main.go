package main

import (
	"context"
	"flag"
	"log"
	"time"

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

	// The server will handle saving the result to outbox and publishing to Kafka
	log.Println("Client completed successfully. The server has processed the request and will publish the result to Kafka via the outbox pattern.")
}