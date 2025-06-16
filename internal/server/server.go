package server

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "service-a/internal/server/summation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// SummationServer is the server implementation of the SummationService
type SummationServer struct {
	// Embed the unimplemented server
	pb.UnimplementedSummationServiceServer
}

// NewSummationServer creates a new instance of SummationServer
func NewSummationServer() *SummationServer {
	return &SummationServer{}
}

// CalculateSum implements the CalculateSum RPC method
func (s *SummationServer) CalculateSum(ctx context.Context, req *pb.SummationRequest) (*pb.SummationResponse, error) {
	log.Printf("Received request: a=%d, b=%d", req.GetA(), req.GetB())
	result := req.GetA() + req.GetB()
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
