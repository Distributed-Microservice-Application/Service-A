package outbox

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Outbox represents the structure of the outbox table
// It contains the ID, sum of messages, sent timestamp, and creation timestamp.
type Outbox struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	Sum       int32        `json:"sum" db:"sum"`
	SentAt    sql.NullTime `json:"sent_at" db:"sent_at"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
}
