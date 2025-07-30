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
	Result    int32  `json:"result"`
	ServiceID string `json:"service_id"`
	Timestamp string `json:"timestamp"`
}

const (
	ServiceName = "SummationService"
)

func SummationRequest(client pb.SummationServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ServiceID string = r.RemoteAddr

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Service-ID", ServiceID)

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
			log.Printf("[%s] Using default values due to JSON decode error: %v", ServiceID, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// ---------------------- Make the gRPC call ----------------------
		log.Printf("[%s] Sending gRPC request with numbers: %d and %d", ServiceID, Data.A, Data.B)
		result, err := client.CalculateSum(ctx, &pb.SummationRequest{A: int32(Data.A), B: int32(Data.B)})
		if err != nil {
			log.Printf("[%s] gRPC call failed: %v", ServiceID, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "gRPC call failed", "service_id": ServiceID})
			return
		}

		log.Printf("[%s] Received sum from gRPC server: %d", ServiceID, result.Result)

		// Send JSON response with service identification
		response := ResponseData{
			Result:    result.Result,
			ServiceID: ServiceID,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("[%s] Failed to encode response: %v", ServiceID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("[%s] HTTP API request completed successfully from remote address %s", ServiceName, ServiceID)
	}
}
