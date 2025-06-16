package main

import (
	"log"
	"service-a/internal/server"
)

func main() {
	log.Println("Starting Summation Service")

	// Create a new instance of the SummationServer
	if err := server.StartServer(50051); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
