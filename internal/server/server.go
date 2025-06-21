package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"service-a/internal/outbox"

	pb "service-a/internal/server/summation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// SummationServer is the server implementation of the SummationService
type SummationServer struct {
	// Embed the unimplemented server
	pb.UnimplementedSummationServiceServer
	outboxRepo outbox.Repository
}

// NewSummationServer creates a new instance of SummationServer
func NewSummationServer() *SummationServer {
	return &SummationServer{}
}

// NewSummationServerWithOutbox creates a new instance of SummationServer with outbox repository
func NewSummationServerWithOutbox(repo outbox.Repository) *SummationServer {
	return &SummationServer{
		outboxRepo: repo,
	}
}

// CalculateSum implements the CalculateSum RPC method
func (s *SummationServer) CalculateSum(ctx context.Context, req *pb.SummationRequest) (*pb.SummationResponse, error) {
	log.Printf("Received request: a=%d, b=%d", req.GetA(), req.GetB())
	result := req.GetA() + req.GetB()

	// Save result to outbox if repository is available
	if s.outboxRepo != nil {
		err := s.outboxRepo.SaveOutbox(ctx, outbox.NewOutbox(result))
		if err != nil {
			log.Printf("Warning: Failed to save to outbox: %v", err)
			// Continue despite error to ensure the RPC call still works
		} else {
			log.Println("Successfully saved result to outbox table for async processing")
		}
	}

	return &pb.SummationResponse{Result: result}, nil
}

// StartServer starts the gRPC server on the specified port
func StartServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	// Register reflection service on gRPC server
	reflection.Register(server)

	// Register the SummationService with the gRPC server
	pb.RegisterSummationServiceServer(server, &SummationServer{})

	log.Printf("Starting gRPC server on port %d", port)
	return server.Serve(lis)
}

// StartServerWithOutbox starts the gRPC server on the specified port with outbox support
func StartServerWithOutbox(port int, repo outbox.Repository) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	// Register reflection service on gRPC server
	reflection.Register(server)

	// Register the SummationService with the gRPC server
	pb.RegisterSummationServiceServer(server, NewSummationServerWithOutbox(repo))

	log.Printf("Starting gRPC server with outbox support on port %d", port)
	return server.Serve(lis)
}
