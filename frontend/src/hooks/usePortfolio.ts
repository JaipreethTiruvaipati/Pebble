import { useQuery } from "@tanstack/react-query";
import * as portfolioApi from "@/api/portfolio.api";

export function usePortfolio() {
  return useQuery({
    queryKey: ["portfolio"],
    queryFn: portfolioApi.getPortfolio,
  });
}

export function useInvestments(limit = 20) {
  return useQuery({
    queryKey: ["investments", limit],
    queryFn: () => portfolioApi.listInvestments(undefined, limit),
  });
}

export function useMarketSignal() {
  return useQuery({
    queryKey: ["market", "signal"],
    queryFn: portfolioApi.getMarketSignal,
    refetchInterval: 60_000,
  });
}

export function useWeeklyInsights() {
  return useQuery({
    queryKey: ["insights", "weekly"],
    queryFn: portfolioApi.getWeeklyInsights,
  });
}

export function useBenchmark() {
  return useQuery({
    queryKey: ["insights", "benchmark"],
    queryFn: portfolioApi.getBenchmarkInsights,
  });
}
