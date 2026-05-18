-- Week 22: indexes identified via EXPLAIN ANALYZE on hot paths
CREATE INDEX IF NOT EXISTS idx_investments_user_status ON investments(user_id, status);
CREATE INDEX IF NOT EXISTS idx_investments_user_created ON investments(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_line_items_transaction ON line_items(transaction_id);
CREATE INDEX IF NOT EXISTS idx_line_items_tx_category ON line_items(transaction_id, category);
CREATE INDEX IF NOT EXISTS idx_users_risk_profile ON users(risk_profile);
