CREATE UNLOGGED TABLE IF NOT EXISTS payment_history (
    correlation_id UUID PRIMARY KEY,
    amount DECIMAL(10, 2) NOT NULL,
    requested_at TIMESTAMP NOT NULL,
    type VARCHAR(20) NOT NULL
);