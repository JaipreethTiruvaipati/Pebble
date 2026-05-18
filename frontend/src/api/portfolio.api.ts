import { apiRequest } from "./client";
import type {
  PortfolioResponse,
  InvestmentListResponse,
  Investment,
  MarketSignalResponse,
  WeeklyDigest,
  BenchmarkResult,
  ReferralStats,
} from "@/types/api.types";

export function getPortfolio() {
  return apiRequest<PortfolioResponse>("/portfolio");
}

export function listInvestments(triggerType?: string, limit = 20) {
  const params = new URLSearchParams({ limit: String(limit) });
  if (triggerType) params.set("trigger_type", triggerType);
  return apiRequest<InvestmentListResponse>(`/investments?${params}`);
}

export function getInvestment(id: string) {
  return apiRequest<Investment>(`/investments/${id}`);
}

export function getMarketSignal() {
  return apiRequest<MarketSignalResponse>("/market/signal");
}

export function getWeeklyInsights() {
  return apiRequest<WeeklyDigest>("/insights/weekly");
}

export function getBenchmarkInsights() {
  return apiRequest<BenchmarkResult>("/insights/benchmark");
}

export function getReferralStats() {
  return apiRequest<ReferralStats>("/referrals/me");
}

export function redeemReferral(code: string) {
  return apiRequest<{ message: string }>("/referrals/redeem", {
    method: "POST",
    body: JSON.stringify({ code }),
  });
}
