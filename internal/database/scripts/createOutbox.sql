CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    SUM INT NOT NULL,

    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_outbox_sent_at ON outbox (sent_at);

-- Create the publication for Debezium
CREATE PUBLICATION dbz_outbox_publication FOR TABLE public.outbox;

-- Grant necessary permissions
GRANT SELECT ON public.outbox TO postgres;
GRANT USAGE ON SCHEMA public TO postgres;