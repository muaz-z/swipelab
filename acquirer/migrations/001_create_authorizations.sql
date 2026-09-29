CREATE TABLE authorizations (
    id UUID PRIMARY KEY,

    merchant_id VARCHAR(100) NOT NULL,

    amount BIGINT NOT NULL CHECK (amount > 0),

    currency VARCHAR(3) NOT NULL,

    status VARCHAR(20) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)