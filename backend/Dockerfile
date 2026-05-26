# ── Pebble Backend — Unified Multi-Stage Dockerfile ──────────────────────────
# Usage:
#   docker build --build-arg SERVICE=api-gateway -t pebble-api-gateway .
#   docker build --build-arg SERVICE=bill-service -t pebble-bill-service .
#
# All 7 microservices share this single Dockerfile.
# The SERVICE arg selects the entrypoint from backend/cmd/${SERVICE}/.

ARG SERVICE=api-gateway

# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Cache dependency downloads (only re-fetched when go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy full source for compilation
COPY backend/ backend/

# Build the target service with optimizations:
#   -s: strip symbol table  -w: strip DWARF debug info  → ~30% smaller binary
ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/service ./backend/cmd/${SERVICE}/

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
# Distroless: no shell, no package manager, non-root by default.
# Attack surface is minimal — only the Go binary and CA certs.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /out/service /app/service

EXPOSE 8080

ENTRYPOINT ["/app/service"]
