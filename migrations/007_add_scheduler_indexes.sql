CREATE INDEX IF NOT EXISTS idx_payment_schedules_due
    ON payment_schedules (status, payment_date);

CREATE INDEX IF NOT EXISTS idx_credits_status
    ON credits (status);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id
    ON transactions (user_id);