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

-- Prepopulate users
INSERT INTO users (id,balance) VALUES
(1, 100.00),
(2, 200.00),
(3, 300.00);