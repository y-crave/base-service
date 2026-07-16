-- Subscription plans and per-user subscription state.

CREATE TABLE plans (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code           TEXT NOT NULL UNIQUE,
    name           TEXT NOT NULL,
    duration_days  INT NOT NULL,
    price_cents    INT NOT NULL DEFAULT 0,
    currency       TEXT NOT NULL DEFAULT 'USD',
    is_active      BOOLEAN NOT NULL DEFAULT true,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO plans (code, name, duration_days, price_cents, currency)
VALUES ('premium_monthly', 'Premium Monthly', 30, 0, 'USD')
ON CONFLICT (code) DO NOTHING;

CREATE TABLE user_subscriptions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL,
    plan_id      UUID NOT NULL REFERENCES plans (id),
    status       TEXT NOT NULL,
    store        TEXT NOT NULL,
    receipt_ref  TEXT,
    expires_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT user_subscriptions_status_check CHECK (status IN ('active', 'expired', 'cancelled')),
    CONSTRAINT user_subscriptions_store_check CHECK (store IN ('ios', 'android', 'dev'))
);

CREATE INDEX idx_user_subscriptions_user_id ON user_subscriptions (user_id);
