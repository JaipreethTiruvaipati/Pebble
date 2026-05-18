import { create } from "zustand";
import { getToken, clearAuth } from "@/api/client";
import * as authApi from "@/api/auth.api";
import type { UserProfile } from "@/types/api.types";

export type RiskProfile = "Conservative" | "Moderate" | "Aggressive";

export interface User {
  id: string;
  name: string;
  email: string;
  avatar: string;
  riskProfile: RiskProfile;
  joinedAt: string;
  onboarded: boolean;
  streakCount: number;
  effectivePenaltyRate: number;
}

function profileToUser(p: UserProfile): User {
  const riskMap: Record<string, RiskProfile> = {
    conservative: "Conservative",
    moderate: "Moderate",
    aggressive: "Aggressive",
  };
  const email = p.email || "user@pebble.in";
  const name = email.split("@")[0].replace(/[._]/g, " ");
  return {
    id: p.id,
    name: name.charAt(0).toUpperCase() + name.slice(1),
    email: p.email,
    avatar: name.slice(0, 2).toUpperCase(),
    riskProfile: riskMap[p.risk_profile?.toLowerCase()] ?? "Moderate",
    joinedAt: p.streak_last_updated ?? new Date().toISOString(),
    onboarded: true,
    streakCount: p.streak_count ?? 0,
    effectivePenaltyRate: p.effective_penalty_rate ?? 0.1,
  };
}

interface AuthStore {
  user: User | null;
  profile: UserProfile | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string, referralCode?: string) => Promise<void>;
  loadProfile: () => Promise<void>;
  logout: () => void;
  setRiskProfile: (r: RiskProfile) => void;
  hydrate: () => void;
}

export const useAuthStore = create<AuthStore>((set, get) => ({
  user: null,
  profile: null,
  isAuthenticated: !!getToken(),
  isLoading: false,

  hydrate: () => {
    set({ isAuthenticated: !!getToken() });
  },

  login: async (email, password, referralCode) => {
    set({ isLoading: true });
    try {
      await authApi.login(email, password, referralCode);
      await get().loadProfile();
    } finally {
      set({ isLoading: false });
    }
  },

  loadProfile: async () => {
    if (!getToken()) {
      set({ user: null, profile: null, isAuthenticated: false });
      return;
    }
    set({ isLoading: true });
    try {
      const profile = await authApi.getMe();
      set({
        profile,
        user: profileToUser(profile),
        isAuthenticated: true,
      });
    } catch {
      clearAuth();
      set({ user: null, profile: null, isAuthenticated: false });
    } finally {
      set({ isLoading: false });
    }
  },

  logout: () => {
    clearAuth();
    set({ user: null, profile: null, isAuthenticated: false });
  },

  setRiskProfile: (r) =>
    set((s) => (s.user ? { user: { ...s.user, riskProfile: r } } : s)),
}));
