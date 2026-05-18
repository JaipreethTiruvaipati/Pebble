# Amazon MQ (RabbitMQ) — Production Messaging

Week 17 replaces local Docker RabbitMQ with **Amazon MQ for RabbitMQ** in staging/production.

## Connection URL format

```bash
# TLS (recommended)
RABBITMQ_URL=amqps://USERNAME:PASSWORD@b-xxxxxxxx.mq.ap-south-1.amazonaws.com:5671/

# Non-TLS (dev only)
RABBITMQ_URL=amqp://USERNAME:PASSWORD@b-xxxxxxxx.mq.ap-south-1.amazonaws.com:5672/
```

Set this in ECS task definitions (Secrets Manager) for every service that uses `internal/queue`:

- api-gateway (if publishing)
- bill-service
- scoring-service
- penalty-service
- **investment-service**
- **notification-service**

## Provisioning checklist

1. Create Amazon MQ broker in **ap-south-1** (RabbitMQ 3.12, single-instance dev / active/standby prod).
2. Create application user with configure/write/read on `pebble.events` and `pebble.dlx`.
3. Open security group: ECS tasks → broker on **5671** (TLS).
4. Update `RABBITMQ_URL` in `.env.local` for local integration tests against a dev broker (optional).

## Local development

Keep using `docker-compose.yml` RabbitMQ for laptops. Use Amazon MQ only when `APP_ENV=production` or `RABBITMQ_USE_AMAZON_MQ=true`.

No code changes are required — the same `amqp091-go` client reads `RABBITMQ_URL`.
