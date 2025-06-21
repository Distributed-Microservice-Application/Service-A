package outbox

const (
	SaveOutbox = `INSERT INTO outbox (id, sum) VALUES ($1, $2)`

	GetOutboxs = `SELECT id, sum, sent_at, created_at FROM outbox WHERE sent_at IS NULL`

	MarkAsSent = `UPDATE outbox SET sent_at = $1 WHERE id = $2`
)
