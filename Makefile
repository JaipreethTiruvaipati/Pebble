.PHONY: run-gateway run-all migrate migrate-down migrate-test test lint docker-build docker-push keys setup-dev

# ── Local dev ─────────────────────────────────────────────────────────────────
run-gateway:
	go run ./backend/cmd/api-gateway/...

run-bill:
	go run ./backend/cmd/bill-service/...

run-scoring:
	go run ./backend/cmd/scoring-service/...

run-penalty:
	go run ./backend/cmd/penalty-service/...

run-investment:
	go run ./backend/cmd/investment-service/...

run-poller:
	go run ./backend/cmd/market-poller/...

run-notify:
	go run ./backend/cmd/notification-service/...

# ── Database ──────────────────────────────────────────────────────────────────
migrate:
	migrate -path backend/migrations -database "$${DATABASE_URL}" up

migrate-down:
	migrate -path backend/migrations -database "$${DATABASE_URL}" down 1

migrate-test:
	migrate -path backend/migrations -database "$${DATABASE_TEST_URL}" up

migrate-reset:
	migrate -path backend/migrations -database "$${DATABASE_URL}" drop -f
	migrate -path backend/migrations -database "$${DATABASE_URL}" up

# ── Testing ───────────────────────────────────────────────────────────────────
test:
	go test ./... -race -coverprofile=coverage.out

test-coverage:
	go test ./... -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-integration:
	go test ./backend/tests/... -race -v -tags integration

# ── Code quality ─────────────────────────────────────────────────────────────
lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

load-test:
	@test -n "$$K6_JWT" || (echo "Set K6_JWT to a valid access token" && exit 1)
	k6 run tests/load/k6-portfolio.js

# ── Docker ───────────────────────────────────────────────────────────────────
SERVICES = api-gateway bill-service scoring-service penalty-service \
           investment-service market-poller notification-service

docker-build:
	@for svc in $(SERVICES); do \
		echo "▶ Building $$svc..."; \
		docker build -t pebble-$$svc:latest -f backend/cmd/$$svc/Dockerfile .; \
	done

docker-push:
	@for svc in $(SERVICES); do \
		docker tag pebble-$$svc:latest $(ECR_REGISTRY)/pebble-$$svc:$(GIT_SHA); \
		docker push $(ECR_REGISTRY)/pebble-$$svc:$(GIT_SHA); \
	done

# ── One-time local setup ─────────────────────────────────────────────────────
setup-dev: keys
	@test -f .env.local || cp .env.example .env.local
	@echo "✅ Copy .env.local and ensure Docker (postgres, redis, rabbitmq) is running"

# ── Key generation (one-time local setup) ────────────────────────────────────
keys:
	@mkdir -p keys
	openssl genrsa -out keys/private.pem 2048
	openssl rsa -in keys/private.pem -pubout -out keys/public.pem
	@echo "✅ RS256 key pair generated in keys/"
	@echo "⚠️  Never commit keys/ to git!"
