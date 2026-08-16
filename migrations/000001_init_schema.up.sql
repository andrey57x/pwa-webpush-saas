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

-- 3. Таблица API-Ключей для внешней интеграции
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
    browser VARCHAR(255),
    os VARCHAR(255),
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

-- 7. Логи доставок c еженедельным партиционированием (Delivery Logs)
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

-- 8. PL/pgSQL функция АВТОМАТИЧЕСКОГО создания недельной партиции для любой даты
CREATE OR REPLACE FUNCTION create_delivery_logs_partition_for_date(target_date DATE)
RETURNS VOID AS $$
DECLARE
    start_of_week DATE;
    end_of_week DATE;
    partition_name TEXT;
BEGIN
    -- Находим понедельник недели для указанной даты
    start_of_week := date_trunc('week', target_date)::DATE;
    -- Конец недели (следующий понедельник)
    end_of_week := start_of_week + INTERVAL '7 days';
    
    -- Имя партиции вида: delivery_logs_y2026_w33
    partition_name := 'delivery_logs_y' || to_char(start_of_week, 'YYYY') || '_w' || to_char(start_of_week, 'IW');

    -- Создаем партицию, если её еще нет
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF delivery_logs FOR VALUES FROM (%L) TO (%L);',
        partition_name, start_of_week, end_of_week
    );
END;
$$ LANGUAGE plpgsql;

-- Создаем первичные партиции на текущую дату
SELECT create_delivery_logs_partition_for_date(CURRENT_DATE);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_apps_app_code ON apps(app_code);
CREATE INDEX IF NOT EXISTS idx_subscriptions_app_active ON subscriptions(app_id, is_active);
CREATE INDEX IF NOT EXISTS idx_subscriptions_endpoint ON subscriptions(endpoint);
CREATE INDEX IF NOT EXISTS idx_campaigns_app_id ON campaigns(app_id);
CREATE INDEX IF NOT EXISTS idx_delivery_logs_campaign ON delivery_logs(campaign_id);