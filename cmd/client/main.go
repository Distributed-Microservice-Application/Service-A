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

	// Set up a connection to the server
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(*addr, opts...)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewSummationServiceClient(conn)

	// Contact the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sendRequest(ctx, client, *numA, *numB)
}

func sendRequest(ctx context.Context, client pb.SummationServiceClient, a, b int) int32 {
	resp, err := client.CalculateSum(ctx, &pb.SummationRequest{A: int32(a), B: int32(b)})
	if err != nil {
		log.Fatalf("Could not calculate sum: %v", err)
	}
	log.Printf("Sum result: %d + %d = %d", a, b, resp.GetResult())
	return resp.GetResult()
}
