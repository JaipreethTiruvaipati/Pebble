CREATE INDEX IF NOT EXISTS idx_penalties_user_status ON penalties(user_id, status);
CREATE INDEX IF NOT EXISTS idx_penalties_expires ON penalties(expires_at) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_transactions_user_logged ON transactions(user_id, logged_at DESC);
CREATE INDEX IF NOT EXISTS idx_wallet_tx_user ON wallet_transactions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_pool_status ON pool_contributions(status);
