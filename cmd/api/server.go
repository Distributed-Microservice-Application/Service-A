package API

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "service-a/internal/server/summation"
)

type RequestData struct {
	A int32 `json:"a"`
	B int32 `json:"b"`
}

type ResponseData struct {
	Result int32 `json:"result"`
}

func SummationRequest(client pb.SummationServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Handle different HTTP methods
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Only POST method allowed"})
			return
		}

		var Data RequestData

		// Try to decode JSON from request body
		if err := json.NewDecoder(r.Body).Decode(&Data); err != nil {
			// If JSON decode fails, use default values for testing
			Data.A = 10
			Data.B = 20
			log.Printf("Using default values due to JSON decode error: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// ---------------------- Make the gRPC call ----------------------
		log.Printf("Sending gRPC request with numbers: %d and %d", Data.A, Data.B)
		result, err := client.CalculateSum(ctx, &pb.SummationRequest{A: int32(Data.A), B: int32(Data.B)})
		if err != nil {
			log.Printf("gRPC call failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "gRPC call failed"})
			return
		}

		log.Printf("Received sum from gRPC server: %d", result.Result)

		// Send JSON response
		response := ResponseData{Result: result.Result}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println("HTTP API request completed successfully")
	}
}
