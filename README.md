# Web Push SaaS Platform & PWA Demo Client

Высоконагруженная SaaS-платформа рассылки и аналитики Web Push уведомлений для Progressive Web Applications (PWA).

## Стек технологий

- Backend: Go (Golang 1.22+)
- Broker: Apache Kafka (KRaft mode)
- Database: PostgreSQL 16 (с партиционированием по неделям)
- Cache/Dedup: Redis 7
- Frontend: Vanilla JS, Service Workers, Pico CSS v2
- Migrations: golang-migrate

## Структура проекта

- `cmd/api` — Точка входа основного бэкенда SaaS платформы (порт 8080)
- `cmd/pizza` — Автономный тестовый сайт доставки пиццы PWA (порт 8085)
- `internal/` — Бизнес-логика (Clean Architecture: delivery, usecase, repository, domain, connectors)
- `migrations/` — SQL-миграции схемы PostgreSQL
- `frontend/` — Панель управления SaaS (Single Page Application)

## Требования

- Go 1.22 или выше
- Docker и Docker Compose
- golang-migrate CLI (необязательно, миграции накатываются автоматически при старте приложения)

## Быстрый запуск

### 1. Настройка окружения

Проверьте файл `.env` в корне проекта.

### 2. Запуск инфраструктурных контейнеров

Запустите PostgreSQL, Redis, Kafka и Kafka UI:

```bash
docker-compose up -d
```

Проверить статус контейнеров:

```bash
docker-compose ps
```

### 3. Запуск основного SaaS приложения

В первом терминале выполните:

```bash
make run-api
```

Или вручную:

```bash
go run cmd/api/main.go
```

При старте автоматически применятся SQL-миграции и проверится создание недельных партиций в PostgreSQL.

Панель управления SaaS будет доступна по адресу: `http://localhost:8080`
Kafka UI доступен по адресу: `http://localhost:8082`

### 4. Запуск автономного сайта пиццерии (PWA)

Во втором терминале выполните:

```bash
make run-pizza
```

Или вручную:

```bash
go run cmd/pizza/main.go
```

Сайт пиццерии будет доступен по адресу: `http://localhost:8085`

## Порядок проверки функционала

1. Перейдите в панель SaaS (`http://localhost:8080`).
2. Зарегистрируйте новую компанию (автоматически создастся аккаунт и компания).
3. Перейдите во вкладку **Сайты**, нажмите **+ Добавить сайт**:
   - Название: `Luigi's Pizza`
   - Код сайта: `pizza_app`
4. Перейдите во вкладку **API-Ключи** и выпишите новый ключ для внешней интеграции.
5. Откройте сайт пиццерии (`http://localhost:8085`) в браузерах Chrome или Firefox.
6. Нажмите **Включить Push-Уведомления** и разрешите отправку уведомлений в браузере.
7. Вернитесь в панель SaaS (`http://localhost:8080`), перейдите во вкладку **Рассылки**, нажмите **+ Создать рассылку** и нажмите **Запустить**.
8. Уведомление придет на устройство через фоновый Service Worker.
9. Проверьте клик по уведомлению и переход во вкладку **Аналитика** для просмотра CTR.

## Тестирование на мобильных устройствах через Cloudflare Tunnel

Так как Web Push API требует защищенное соединение (HTTPS), для проверки на мобильных устройствах или из внешней сети используйте Cloudflare Tunnel с протоколом HTTP/2:

1. Проброс основного SaaS API:

```bash
cloudflared tunnel --url http://localhost:8080
```

Скопируйте полученную ссылку (например, `https://your-saas.trycloudflare.com`).

2. Укажите эту ссылку в `.env`:

```env
PIZZA_SAAS_API_URL=https://your-saas.trycloudflare.com/api/v1
```

3. Проброс сайта пиццерии:

```bash
cloudflared tunnel --url http://localhost:8085
```

Откройте полученную HTTPS-ссылку пиццерии на телефоне для оформления подписки.

## Команды Makefile

- `make run-api` — запуск основного SaaS бэкенда.
- `make run-pizza` — запуск сайта пиццерии.
- `make migrate-up` — ручной запуск миграций PostgreSQL через Docker.
- `make migrate-down` — откат всех миграций PostgreSQL.
