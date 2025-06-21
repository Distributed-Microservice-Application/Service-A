package outbox

import (
	"context"
	"github.com/google/uuid"
	"log"
	"service-a/internal/kafka"
	"sync"
	"time"
)

type OutboxPublisher struct {
	Repository  Repository
	KafkaWriter *kafkaStructure.KafkaPublisher
	// Interval defines how often the outbox publisher checks for new messages to send
	Interval time.Duration
}

// NewOutboxPublisher creates a new OutboxPublisher with the given repository and Kafka writer
func NewOutboxPublisher(repo Repository, Publisher *kafkaStructure.KafkaPublisher, interval time.Duration) *OutboxPublisher {
	return &OutboxPublisher{
		Repository:  repo,
		KafkaWriter: Publisher,
		Interval:    interval,
	}
}

// Start begins the outbox publishing process, checking for new messages at the defined interval
func (p *OutboxPublisher) Start(ctx context.Context) {
	// Ensure the interval is not zero or negative
	if p.Interval <= 0 {
		log.Println("Invalid interval, using default of 3 seconds")
		p.Interval = 3 * time.Second
	}

	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	log.Println("OutboxPublisher started, checking for messages every", p.Interval)

	for {
		select {
		case <-ticker.C:
			p.publishOutboxMessages(ctx)
		case <-ctx.Done():
			log.Println("OutboxPublisher stopped")
			return
		}
	}
}

// publishOutboxMessages retrieves outbox messages and sends them to Kafka
func (p *OutboxPublisher) publishOutboxMessages(ctx context.Context) {
	outboxs, err := p.Repository.GetOutboxs(ctx)
	if err != nil {
		log.Println("Error retrieving outbox records:", err)
		return
	}

	// Check if there are no outbox messages to send
	if len(outboxs) == 0 {
		return
	}

	log.Printf("Found %d outbox messages to send", len(outboxs))

	// send the outbox messages to Kafka with goroutines
	var wg sync.WaitGroup
	for _, outbox := range outboxs {
		wg.Add(1)
		go func(outbox Outbox) {
			defer wg.Done()
			err := p.KafkaWriter.SendMessage(outbox.Sum, ctx)
			if err != nil {
				log.Println("Error sending message to Kafka:", err)
				return
			}
			err = p.Repository.MarkAsSent(ctx, outbox.ID)
			if err != nil {
				log.Println("Error marking outbox as sent:", err)
			}
		}(outbox)
	}

	go func() {
		wg.Wait()
		log.Println("All outbox messages processed")
	}()

	log.Println("All outbox messages processed and sent to Kafka")
}

// NewOutbox creates a new Outbox instance with the current time as CreatedAt
func NewOutbox(sum int32) Outbox {
	return Outbox{
		ID:        uuid.New(),
		Sum:       sum,
		CreatedAt: time.Now(),
	}
}
