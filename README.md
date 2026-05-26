# 🪨 Pebble

**Pebble** is a fintech microservices platform that helps users curb impulse spending by scoring purchases, applying penalties to impulsive transactions, and automatically investing the penalty funds into diversified portfolios.

---

## Architecture

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  React/Vite │────▶│   API Gateway    │────▶│   PostgreSQL    │
│  Frontend   │     │   (Chi :8080)    │────▶│   Redis         │
└─────────────┘     └──────┬───────────┘     └─────────────────┘
                           │ RabbitMQ
              ┌────────────┼────────────────────────┐
              ▼            ▼            ▼            ▼
        ┌───────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
        │   Bill    │ │ Scoring  │ │ Penalty  │ │Investment│
        │  Service  │ │ Service  │ │ Service  │ │ Service  │
        └───────────┘ └──────────┘ └──────────┘ └──────────┘
                                                      ▲
                      ┌──────────┐ ┌──────────┐       │
                      │  Market  │ │  Notify  │       │
                      │  Poller  │ │ Service  │───────┘
                      └──────────┘ └──────────┘
```

| Service | Purpose | Dependencies |
|---------|---------|-------------|
| **api-gateway** | REST API, JWT auth, rate limiting | Postgres, Redis, RabbitMQ |
| **bill-service** | Receipt upload → S3 → publish event | RabbitMQ, AWS S3 |
| **scoring-service** | OCR (Google Vision) + LLM scoring (Gemini) | Postgres, RabbitMQ, Redis |
| **penalty-service** | Calculate penalties, consent windows, expiry sweep | Postgres, RabbitMQ |
| **investment-service** | Pool penalty funds, execute trades via Smallcase | Postgres, Redis, RabbitMQ |
| **market-poller** | Fetch NSE/MCX/CCIL/AMFI data, cache signals | Redis |
| **notification-service** | FCM push + SES email dispatch | RabbitMQ |

---

## Prerequisites

| Tool | Version |
|------|---------|
| Go | 1.25+ |
| Docker & Docker Compose | Latest |
| Node.js | 20+ |
| OpenSSL | Any (for key generation) |

---

## Quick Start

### 1. Clone & Setup

```bash
git clone https://github.com/jaipreeth/pebble.git
cd pebble
make setup-dev    # generates JWT keys + copies .env.example → .env.local
```

### 2. Configure Environment

Edit `.env.local` with your actual API keys and secrets:

```bash
# Required for basic operation:
DATABASE_URL=postgres://ss:ss@localhost:5433/pebble?sslmode=disable
REDIS_URL=redis://localhost:6379
RABBITMQ_URL=amqp://ss:ss@localhost:5672/

# External APIs (fill in your keys):
GEMINI_API_KEY=your-key-here
RAZORPAY_KEY_ID=your-key-here
# ... see .env.example for the full list
```

For the frontend, edit `frontend/.env.local`:
```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_RAZORPAY_KEY_ID=your-publishable-key
# ... see frontend/.env.example for the full list
```

### 3. Start Infrastructure & Services

```bash
# Start Postgres, Redis, RabbitMQ + all 7 microservices
docker compose up -d --build

# Or run just the infra containers and the gateway locally:
docker compose up -d postgres redis rabbitmq
make run-gateway
```

### 4. Start Frontend

```bash
cd frontend
npm install
npm run dev    # → http://localhost:5173
```

### 5. Verify Everything Works

```bash
./verify_pebble.sh    # Automated E2E health check
```

---

## Project Structure

```
Pebble/
├── backend/
│   ├── Dockerfile              # Unified multi-stage (all 7 services)
│   ├── cmd/                    # Service entry points (main.go per service)
│   ├── internal/               # Shared internals (auth, cache, config, db, queue)
│   ├── pkg/                    # Domain packages (allocate, broker, llm, market, notify, ocr)
│   ├── migrations/             # PostgreSQL migrations (golang-migrate)
│   └── tests/                  # Integration tests + helpers
├── frontend/                   # React 19 + Vite + TanStack Router
├── contracts/                  # OpenAPI spec + shared TypeScript types
├── infra/
│   ├── amazon-mq/              # Amazon MQ configuration guide
│   └── terraform/              # AWS infrastructure (RDS, ElastiCache, ECS, IAM)
├── tests/load/                 # k6 load testing scripts
├── .github/
│   ├── workflows/              # CI/CD pipelines
│   ├── CODEOWNERS              # PR review ownership
│   └── pull_request_template.md
├── docker-compose.yml          # Local development stack
├── docker-compose.amazonmq.yml # Amazon MQ overlay
├── Makefile                    # Build, test, deploy commands
├── .env.example                # Backend env template (safe to commit)
└── .dockerignore               # Keeps secrets out of images
```

---

## Environment Variables

### Backend (`.env.example` → `.env.local`)

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | Yes* | PostgreSQL connection string |
| `REDIS_URL` | Yes* | Redis connection string |
| `RABBITMQ_URL` | Yes* | RabbitMQ AMQP URL |
| `JWT_PRIVATE_KEY_PATH` | Yes | Path to RS256 private key |
| `JWT_PUBLIC_KEY_PATH` | Yes | Path to RS256 public key |
| `GEMINI_API_KEY` | No | Google Gemini for LLM scoring |
| `RAZORPAY_KEY_ID` | No | Razorpay for wallet topup |
| `RAZORPAY_KEY_SECRET` | No | Razorpay server secret |
| `AWS_S3_BUCKET` | No | S3 bucket for bill images |
| `FIREBASE_CREDENTIALS_PATH` | No | FCM push notifications |
| `SMALLCASE_API_KEY` | No | Broker execution |

*Required by services that use them. Services gracefully skip unused dependencies.

### Frontend (`frontend/.env.example` → `frontend/.env.local`)

| Variable | Description |
|----------|-------------|
| `VITE_API_BASE_URL` | Backend API endpoint |
| `VITE_RAZORPAY_KEY_ID` | Razorpay publishable key |
| `VITE_FIREBASE_*` | Firebase client config (6 vars) |

> **Security:** `.env.local` files are in `.gitignore` and `.dockerignore`. Never commit real secrets.

---

## Docker

All 7 backend services share a single unified `backend/Dockerfile` using build args:

```bash
# Build a single service
docker build --build-arg SERVICE=api-gateway -t pebble-api-gateway -f backend/Dockerfile .

# Build all services
make docker-build

# Push to ECR
make docker-push ECR_REGISTRY=123456789.dkr.ecr.ap-south-1.amazonaws.com GIT_SHA=$(git rev-parse --short HEAD)
```

The runtime image uses **Google Distroless** (`gcr.io/distroless/static-debian12:nonroot`) — no shell, no package manager, non-root by default.

---

## Testing

```bash
make test               # Unit tests with race detector
make test-coverage      # + HTML coverage report
make test-integration   # Integration tests (requires running infra)
make load-test          # k6 load test (set K6_JWT env var first)
make lint               # golangci-lint
make vuln               # govulncheck
```

---

## CI/CD

| Workflow | Trigger | Steps |
|----------|---------|-------|
| `backend-ci.yml` | Push/PR to `main`/`develop` (backend changes) | Test → Lint → Vuln → Docker Build |
| `frontend-ci.yml` | Push/PR to `main`/`develop` (frontend changes) | Type Check → Lint → Build → Deploy to Cloudflare Pages |
| `deploy.yml` | Push to `main` | Build all 7 images → Push to ECR → Deploy to ECS |
| `security-scan.yml` | Push/PR (backend changes) | govulncheck → Trivy image scan (all 7 services) |

### Required GitHub Secrets

| Secret | Used By |
|--------|---------|
| `AWS_ACCESS_KEY_ID` | deploy.yml |
| `AWS_SECRET_ACCESS_KEY` | deploy.yml |
| `CLOUDFLARE_API_TOKEN` | frontend-ci.yml |
| `CLOUDFLARE_ACCOUNT_ID` | frontend-ci.yml |

---

## API Documentation

Full OpenAPI 3.0 specification: [`contracts/openapi.yaml`](contracts/openapi.yaml)

Key endpoints:
- `POST /api/v1/auth/login` — JWT authentication
- `POST /api/v1/transactions` — Record a purchase
- `POST /api/v1/transactions/bill` — Upload receipt for OCR
- `GET /api/v1/portfolio` — Investment portfolio view
- `GET /api/v1/market/signal` — Live market opportunity signals

---

IMPORTANT

The following things need to be filled by user

These items require your real credentials and cannot be automated:

Fill .env.local with real API keys (Gemini, Razorpay, Smallcase, etc.)
Fill frontend/.env.local with Firebase and Razorpay client config
Add GitHub Secrets in repo settings:
AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY
CLOUDFLARE_API_TOKEN / CLOUDFLARE_ACCOUNT_ID
Implement stub handlers (flagged as TODO in router.go):
handleRazorpayWebhook — add signature verification
handleRefresh — implement real refresh token rotation
handleLogout — clear httpOnly cookies
handleSignup — persist user to DB
## License

Proprietary — Pebble Fintech © 2026
