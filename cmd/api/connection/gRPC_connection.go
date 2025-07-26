package connection

import (
	"flag"
	"log"
	pb "service-a/internal/server/summation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func GRPC_Connection() (pb.SummationServiceClient, *grpc.ClientConn, error) {
	// ---------------------- Set up gRPC connection ----------------------
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("GRPC Connection: Did not connect: %v", err)
		return nil, nil, err
	}

	// Create a new gRPC client
	client := pb.NewSummationServiceClient(conn)

	return client, conn, nil
}
