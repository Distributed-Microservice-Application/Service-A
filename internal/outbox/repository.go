package outbox

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	// SaveOutbox saves an outbox record to the database
	SaveOutbox(ctx context.Context, outbox Outbox) error

	// GetOutboxs retrieves all outbox records from the database
	GetOutboxs(ctx context.Context) ([]Outbox, error)

	// MarkAsSent marks an outbox record as sent by updating its SentAt timestamp
	MarkAsSent(ctx context.Context, id uuid.UUID) error
}

type DB struct {
	RepositoryDB *sql.DB
}

// SaveOutbox saves an outbox record to the database (stub implementation)
func (db *DB) SaveOutbox(ctx context.Context, outbox Outbox) error {
	_, err := db.RepositoryDB.ExecContext(ctx, SaveOutbox, outbox.ID, outbox.Sum)
	if err != nil {
		log.Println("Error saving outbox:", err)
		return err
	}
	return nil
}

// GetOutboxs retrieves all outbox records from the database (stub implementation)
func (db *DB) GetOutboxs(ctx context.Context) ([]Outbox, error) {
	rows, err := db.RepositoryDB.QueryContext(ctx, GetOutboxs)
	if err != nil {
		log.Println("Error retrieving outbox records:", err)
		return nil, err
	}
	defer rows.Close()

	var outboxs []Outbox
	for rows.Next() {
		var outbox Outbox
		if err := rows.Scan(&outbox.ID, &outbox.Sum, &outbox.SentAt, &outbox.CreatedAt); err != nil {
			log.Println("Error scanning outbox record:", err)
			return nil, err
		}
		outboxs = append(outboxs, outbox)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}

	return outboxs, nil
}

// MarkAsSent marks an outbox record as sent by updating its SentAt timestamp (stub implementation)
func (db *DB) MarkAsSent(ctx context.Context, id uuid.UUID) error {
	now := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	_, err := db.RepositoryDB.ExecContext(ctx, MarkAsSent, now, id)
	if err != nil {
		log.Println("Error marking outbox as sent:", err)
		return err
	}
	return nil
}

// NewRepository creates a new instance of the outbox repository
func NewRepository(db *sql.DB) Repository {
	return &DB{
		RepositoryDB: db,
	}
}
