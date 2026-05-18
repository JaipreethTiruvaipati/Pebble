import { apiRequest } from "./client";
import type { ReferralStats } from "@/types/api.types";

export function getReferralMe() {
  return apiRequest<ReferralStats>("/referrals/me");
}
