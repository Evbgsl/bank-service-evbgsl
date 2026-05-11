CREATE EXTENSION IF NOT EXISTS pgcrypto;

DROP TABLE IF EXISTS cards;

CREATE TABLE cards (
                       id BIGSERIAL PRIMARY KEY,
                       user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
                       card_number_encrypted BYTEA NOT NULL,
                       expiry_encrypted BYTEA NOT NULL,
                       card_number_hmac TEXT NOT NULL,
                       masked_number VARCHAR(32) NOT NULL,
                       cvv_hash TEXT NOT NULL,
                       status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);