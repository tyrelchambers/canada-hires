CREATE TABLE boycotts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    business_name VARCHAR(500) NOT NULL,
    business_address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, business_name)
);

CREATE INDEX idx_boycotts_user_id ON boycotts(user_id);
CREATE INDEX idx_boycotts_business_name ON boycotts(business_name);
CREATE INDEX idx_boycotts_created_at ON boycotts(created_at);