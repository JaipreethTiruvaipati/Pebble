import { apiRequest } from "./client";
import type { PenaltyRow, PendingPenaltyBanner } from "@/types/api.types";

export function listPenalties(status?: string) {
  const q = status ? `?status=${status}` : "";
  return apiRequest<PenaltyRow[]>(`/penalties${q}`);
}

export function contestPenalty(id: string) {
  return apiRequest(`/penalties/${id}/contest`, { method: "POST" });
}

export function confirmPenaltyEarly(id: string) {
  return apiRequest(`/penalties/${id}/confirm`, { method: "POST" });
}
