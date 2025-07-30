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

type KafkaPublisher struct {
	Publisher *kafka.Writer
	Partition int // Specific partition for this publisher
}

// FixedPartitionBalancer always sends messages to a specific partition
type FixedPartitionBalancer struct {
	Partition int
}

func (f *FixedPartitionBalancer) Balance(msg kafka.Message, partitions ...int) int {
	// Always return the fixed partition
	return f.Partition
}

// // NewKafkaWriter creates a new Kafka writer with proper configuration
// func NewKafkaWriter(topic string) *KafkaPublisher {
// 	writer := &kafka.Writer{
// 		Addr:         kafka.TCP("kafka:29092"), // Use internal Kafka address
// 		Topic:        topic,
// 		Balancer:     &kafka.LeastBytes{},
// 		RequiredAcks: kafka.RequireOne, // Only wait for leader acknowledgment
// 		Async:        true,             // Use asynchronous mode for better performance
// 		Logger:      kafka.LoggerFunc(log.Printf),
// 		ErrorLogger: kafka.LoggerFunc(log.Printf), // Added error logger for async errors
// 	}
// 	return &KafkaPublisher{Publisher: writer, Partition: 0} // Default to partition 0
// }

// NewKafkaWriterWithPartition creates a new Kafka writer with a specific partition
func NewKafkaWriterWithPartition(topic string, partition int) *KafkaPublisher {
	writer := &kafka.Writer{
		Addr:         kafka.TCP("kafka:29092"), // Use internal Kafka address
		Topic:        topic,
		Balancer:     &FixedPartitionBalancer{Partition: partition}, // Use custom balancer for fixed partition
		RequiredAcks: kafka.RequireOne,                              // Only wait for leader acknowledgment
		Async:        true,                                          // Use asynchronous mode for better performance
		Logger:       kafka.LoggerFunc(log.Printf),
		ErrorLogger:  kafka.LoggerFunc(log.Printf), // Added error logger for async errors
	}
	return &KafkaPublisher{Publisher: writer, Partition: partition}
}

// SendMessage sends a message to the Kafka topic with proper error handling
func (p *KafkaPublisher) SendMessage(msg int32, ctx context.Context) error {
	message := Message{
		Sum:       msg,
		Timestamp: time.Now(),
	}

	// Convert the message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Create a Kafka message with a unique key - don't set partition here, let balancer handle it
	kafkaMessage := kafka.Message{
		Key:   []byte(uuid.New().String()),
		Value: messageBytes,
		// Remove Partition field - let the FixedPartitionBalancer handle it

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
			{
				Key:   "partition",
				Value: []byte(fmt.Sprintf("%d", p.Partition)),
			},
		},
	}

	// Send the message
	err = p.Publisher.WriteMessages(ctx, kafkaMessage)
	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %v", err)
	}

	log.Printf("Successfully sent message to Kafka topic %s (partition %d): %+v", p.Publisher.Topic, p.Partition, message)
	return nil
}
