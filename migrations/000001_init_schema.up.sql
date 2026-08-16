-- 1. Таблица Клиентов (Tenants)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Таблица Пользователей (Users)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'TENANT_ADMIN',
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Таблица API-Ключей для интеграций
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. PWA Приложения (Apps)
CREATE TABLE IF NOT EXISTS apps (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    app_code VARCHAR(64) UNIQUE NOT NULL,
    vapid_public_key TEXT NOT NULL,
    vapid_private_key TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Подписки устройств (Subscriptions)
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    user_identifier VARCHAR(255),
    endpoint TEXT NOT NULL UNIQUE,
    p256dh_key TEXT NOT NULL,
    auth_key TEXT NOT NULL,
    browser VARCHAR(50),
    os VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    deactivated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 6. Кампании рассылок (Campaigns)
CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY,
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    icon_url TEXT,
    target_url TEXT,
    status VARCHAR(50) DEFAULT 'DRAFT',
    total_targeted INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 7. Логи доставок с партиционированием по месяцам (Delivery Logs)
CREATE TABLE IF NOT EXISTS delivery_logs (
    id UUID NOT NULL,
    campaign_id UUID NOT NULL,
    subscription_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    error_code INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Создаем базовые партиции на 2026 год (для диплома и разработки)
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m01 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m02 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m03 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m04 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m05 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m06 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m07 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m08 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m09 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m10 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m11 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE IF NOT EXISTS delivery_logs_y2026_m12 PARTITION OF delivery_logs
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

-- 8. Высокопроизводительные индексы
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_apps_app_code ON apps(app_code);
CREATE INDEX IF NOT EXISTS idx_subscriptions_app_active ON subscriptions(app_id, is_active);
CREATE INDEX IF NOT EXISTS idx_subscriptions_endpoint ON subscriptions(endpoint);
CREATE INDEX IF NOT EXISTS idx_campaigns_app_id ON campaigns(app_id);
CREATE INDEX IF NOT EXISTS idx_delivery_logs_campaign ON delivery_logs(campaign_id);