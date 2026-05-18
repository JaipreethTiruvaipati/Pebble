DROP INDEX IF EXISTS idx_investments_trigger_type;
ALTER TABLE investments
    DROP COLUMN IF EXISTS broker_ref,
    DROP COLUMN IF EXISTS trigger_type;
