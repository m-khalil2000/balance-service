CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    balance NUMERIC(20,2) NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    amount NUMERIC(20,2) NOT NULL,
    state TEXT NOT NULL,
    source_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance indexes for high concurrency
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_state ON transactions(state);

-- Composite index for common queries
CREATE INDEX idx_transactions_user_state ON transactions(user_id, state);

-- Prepopulate users
INSERT INTO users (id,balance) VALUES
(1, 10000.00),
(2, 20000.00),
(3, 30000.00);