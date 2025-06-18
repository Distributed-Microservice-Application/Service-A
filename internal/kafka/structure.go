package kafkaStructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

// Message represents the structure of messages we'll send
type Message struct {
	Sum       int32     `json:"sum"`
	Timestamp time.Time `json:"timestamp"`
}

// NewKafkaWriter creates a new Kafka writer with proper configuration
func NewKafkaWriter(topic string) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,       // Only wait for leader acknowledgment
		Async:        true,                   // Use asynchronous mode for better performance
		// BatchTimeout: 100 * time.Millisecond, // Increased batch timeout for better batching
		// BatchSize:    100,                    // Number of messages to batch before sending
		Logger:       kafka.LoggerFunc(log.Printf),
		ErrorLogger:  kafka.LoggerFunc(log.Printf), // Added error logger for async errors
	}
	return writer
}

// SendMessage sends a message to the Kafka topic with proper error handling
func SendMessage(writer *kafka.Writer, msg int32, ctx context.Context) error {
	message := Message{
		Sum:       msg,
		Timestamp: time.Now(),
	}

	// Convert the message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Create a Kafka message with a unique key
	kafkaMessage := kafka.Message{
		Key:   []byte(uuid.New().String()),
		Value: messageBytes,
		// Add headers for better debugging
		Headers: []kafka.Header{
			{
				Key:   "content-type",
				Value: []byte("application/json"),
			},
			{
				Key:   "timestamp",
				Value: []byte(time.Now().Format(time.RFC3339)),
			},
		},
	}

	// Send the message
	err = writer.WriteMessages(ctx, kafkaMessage)
	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %v", err)
	}

	log.Printf("Successfully sent message to Kafka topic %s: %+v", writer.Topic, message)
	return nil
}
