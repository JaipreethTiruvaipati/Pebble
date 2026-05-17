CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    merchant VARCHAR(255) NOT NULL,
    total_amount DECIMAL(12, 2) NOT NULL,
    bill_s3_key VARCHAR(512),
    status VARCHAR(50) DEFAULT 'pending',
    scored_at TIMESTAMP WITH TIME ZONE,
    logged_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE line_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    quantity INT DEFAULT 1,
    impulse_score INT DEFAULT 0,
    category VARCHAR(100),
    reasoning TEXT,
    user_overridden BOOLEAN DEFAULT FALSE,
    override_score INT,
    penalty_amount DECIMAL(12, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
