#!/bin/bash

# ==============================================================================
# Pebble Fintech - Automated End-to-End Sanity Check & Production Verification
# ==============================================================================

set -e # Exit immediately if a command exits with a non-zero status

# Ensure jq is installed for JSON parsing
if ! command -v jq &> /dev/null; then
    echo "❌ Error: 'jq' is required but not installed. Please install it (e.g., sudo apt install jq or brew install jq)."
    exit 1
fi

echo "🚀 Starting Pebble Automated Verification Pipeline..."

# 1. Infrastructure Spin-Up
echo "📦 Spinning up Docker infrastructure..."
# We use --build to ensure we're testing the latest code
docker compose up -d --build

# Define cleanup function to run on exit or failure
cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        echo "🚨 TEST FAILED! Dumping API Gateway logs:"
        docker logs pebble-api-gateway --tail 50
        echo "🚨 Dumping Postgres logs:"
        docker logs pebble-postgres --tail 50
    fi
    
    echo "🧹 Cleaning up Docker environment..."
    docker compose down
    
    if [ $exit_code -eq 0 ]; then
        echo ""
        echo "==========================================================="
        echo "✅ SYSTEM HEALTHY - READY FOR PRODUCTION"
        echo "==========================================================="
    else
        echo ""
        echo "==========================================================="
        echo "❌ SYSTEM VERIFICATION FAILED (Exit Code: $exit_code)"
        echo "==========================================================="
    fi
    exit $exit_code
}

# Trap EXIT to always run cleanup
trap cleanup EXIT

# 2. Health Checks
echo "⏳ Waiting for API Gateway to become healthy..."
max_attempts=30
attempt=1
gateway_ready=false

while [ $attempt -le $max_attempts ]; do
    # API Gateway exposes /health which returns {"status":"ok"}
    status=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health || true)
    
    if [ "$status" = "200" ]; then
        gateway_ready=true
        echo "✅ API Gateway is healthy!"
        break
    fi
    
    echo "   Attempt $attempt/$max_attempts: Gateway not ready yet (Status: $status). Waiting 2s..."
    sleep 2
    attempt=$((attempt + 1))
done

if [ "$gateway_ready" = false ]; then
    echo "❌ API Gateway failed to become healthy within 60 seconds."
    exit 1
fi

# Migrations are automatically run by the API Gateway on startup (in main.go).
# The database seed is also handled automatically when we hit the dev login endpoint.

# 3. Automated API Flow Testing
echo "🔐 Step 1: Authenticating dummy user..."
LOGIN_RES=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@pebble.in","password":"demo"}')

# Extract JWT token
TOKEN=$(echo "$LOGIN_RES" | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo "❌ Failed to extract JWT token from login response."
    echo "Response: $LOGIN_RES"
    exit 1
fi
echo "✅ Authentication successful. JWT Token acquired."

echo "👤 Step 2: Fetching user profile..."
PROFILE_RES=$(curl -s -w "\n%{http_code}" -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN")

PROFILE_BODY=$(echo "$PROFILE_RES" | head -n -1)
PROFILE_STATUS=$(echo "$PROFILE_RES" | tail -n 1)

if [ "$PROFILE_STATUS" != "200" ]; then
    echo "❌ Failed to fetch user profile. HTTP Status: $PROFILE_STATUS"
    echo "Response: $PROFILE_BODY"
    exit 1
fi

# Assert risk_profile and penalty_rate exist
RISK_PROFILE=$(echo "$PROFILE_BODY" | jq -r '.risk_profile')
PENALTY_RATE=$(echo "$PROFILE_BODY" | jq -r '.effective_penalty_rate')

if [ "$RISK_PROFILE" == "null" ] || [ "$PENALTY_RATE" == "null" ]; then
    echo "❌ User profile is missing required fields (risk_profile or effective_penalty_rate)."
    echo "Response: $PROFILE_BODY"
    exit 1
fi
echo "✅ User profile verified. Risk Profile: $RISK_PROFILE, Penalty Rate: $PENALTY_RATE"

echo "💰 Step 3: Fetching wallet balance..."
WALLET_RES=$(curl -s -w "\n%{http_code}" -X GET http://localhost:8080/api/v1/wallet/balance \
  -H "Authorization: Bearer $TOKEN")

WALLET_BODY=$(echo "$WALLET_RES" | head -n -1)
WALLET_STATUS=$(echo "$WALLET_RES" | tail -n 1)

if [ "$WALLET_STATUS" != "200" ]; then
    echo "❌ Failed to fetch wallet balance. HTTP Status: $WALLET_STATUS"
    echo "Response: $WALLET_BODY"
    exit 1
fi

BALANCE=$(echo "$WALLET_BODY" | jq -r '.balance')
echo "✅ Wallet balance fetched successfully. Current Balance: $BALANCE"

# If we reached here, everything succeeded. The trap will handle the success message and cleanup.
