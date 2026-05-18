ALTER TABLE investments
    ADD COLUMN IF NOT EXISTS trigger_type VARCHAR(50),
    ADD COLUMN IF NOT EXISTS broker_ref TEXT;

CREATE INDEX IF NOT EXISTS idx_investments_trigger_type ON investments(trigger_type);
