package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartMetricsServer starts the Prometheus metrics HTTP server in a separate goroutine
func StartMetricsServer(ctx context.Context, port int) {
	go func() {
		// Create HTTP server for metrics
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		}

		log.Printf("Starting Prometheus metrics server on port %d", port)

		// Start server in a goroutine
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("Metrics server error: %v", err)
			}
		}()

		// Wait for context cancellation
		<-ctx.Done()
		log.Println("Shutting down Prometheus metrics server...")

		// Gracefully shutdown the server
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down metrics server: %v", err)
		}
	}()
}
