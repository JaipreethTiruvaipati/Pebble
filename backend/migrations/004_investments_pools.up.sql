CREATE TABLE investments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_class VARCHAR(50) NOT NULL, -- e.g., 'equity', 'gold', 'bonds', 'mutual_funds'
    amount DECIMAL(12, 2) NOT NULL,
    units DECIMAL(18, 6) NOT NULL,
    nav_at_purchase DECIMAL(12, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_investments_user_id ON investments(user_id);
CREATE INDEX idx_investments_asset_class ON investments(asset_class);

CREATE TABLE pool_contributions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    penalty_id UUID REFERENCES penalties(id) ON DELETE SET NULL,
    amount DECIMAL(12, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pooled', -- 'pooled', 'invested'
    invested_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pool_contributions_status ON pool_contributions(status);
