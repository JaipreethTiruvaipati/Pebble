/**
 * Week 22 load test — portfolio endpoint at 500 concurrent users.
 *
 * Prerequisites:
 *   docker compose up -d postgres redis rabbitmq
 *   go run ./backend/cmd/api-gateway/...
 *
 * Run:
 *   export K6_BASE_URL=http://localhost:8080
 *   export K6_JWT=<access_token from POST /api/v1/auth/login>
 *   k6 run tests/load/k6-portfolio.js
 */
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const errorRate = new Rate('errors');
const portfolioLatency = new Trend('portfolio_latency', true);

export const options = {
  stages: [
    { duration: '30s', target: 100 },
    { duration: '1m', target: 500 },
    { duration: '2m', target: 500 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],
    portfolio_latency: ['p(95)<200'],
    errors: ['rate<0.01'],
  },
};

const baseURL = __ENV.K6_BASE_URL || 'http://localhost:8080';
const token = __ENV.K6_JWT || '';

export default function () {
  if (!token) {
    console.error('Set K6_JWT to a valid Bearer token');
    return;
  }
  const res = http.get(`${baseURL}/api/v1/portfolio`, {
    headers: {
      Authorization: `Bearer ${token}`,
      Accept: 'application/json',
    },
    tags: { name: 'portfolio' },
  });
  portfolioLatency.add(res.timings.duration);
  const ok = check(res, {
    'status is 200': (r) => r.status === 200,
    'has allocation': (r) => r.json('allocation_pct') !== undefined,
  });
  errorRate.add(!ok);
  sleep(0.1);
}
