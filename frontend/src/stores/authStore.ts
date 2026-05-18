import { create } from "zustand";

export type RiskProfile = "Conservative" | "Moderate" | "Aggressive";

export interface User {
  id: string;
  name: string;
  email: string;
  avatar: string;
  riskProfile: RiskProfile;
  joinedAt: string;
  onboarded: boolean;
}

interface AuthStore {
  user: User | null;
  isAuthenticated: boolean;
  login: (user: User) => void;
  logout: () => void;
  setRiskProfile: (r: RiskProfile) => void;
}

const mockUser: User = {
  id: "u_arjun",
  name: "Arjun Sharma",
  email: "arjun.sharma@pebble.in",
  avatar: "AS",
  riskProfile: "Moderate",
  joinedAt: "2025-01-12",
  onboarded: true,
};

export const useAuthStore = create<AuthStore>((set) => ({
  user: mockUser,
  isAuthenticated: true,
  login: (user) => set({ user, isAuthenticated: true }),
  logout: () => set({ user: null, isAuthenticated: false }),
  setRiskProfile: (r) => set((s) => (s.user ? { user: { ...s.user, riskProfile: r } } : s)),
}));
