CREATE INDEX IF NOT EXISTS idx_transactions_created_at
    ON transactions (created_at);

CREATE INDEX IF NOT EXISTS idx_transactions_from_account_id
    ON transactions (from_account_id);

CREATE INDEX IF NOT EXISTS idx_transactions_to_account_id
    ON transactions (to_account_id);

CREATE INDEX IF NOT EXISTS idx_payment_schedules_payment_date
    ON payment_schedules (payment_date);