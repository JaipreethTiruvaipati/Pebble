export const ROUTES = {
  landing: "/",
  signup: "/signup",
  onboardingRisk: "/onboarding/risk",
  onboardingWallet: "/onboarding/wallet",
  dashboard: "/dashboard",
  logTransaction: "/log",
  portfolio: "/portfolio",
  insights: "/insights",
  history: "/history",
  settings: "/settings",
} as const;

export type RouteKey = keyof typeof ROUTES;
