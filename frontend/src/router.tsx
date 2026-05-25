import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import {
  Outlet,
  Link,
  createRootRouteWithContext,
  createRoute,
  createRouter,
  HeadContent,
  Scripts,
  useRouter,
} from "@tanstack/react-router";
import * as React from "react";

// stylesheet import for TanStack Start
import appCss from "./styles.css?url";

// Page imports
import Landing from "./pages/Landing";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import { ProtectedRoute } from "./components/layout/ProtectedRoute";
import Insights from "./pages/Insights";
import LogTransaction from "./pages/LogTransaction";
import Portfolio from "./pages/Portfolio";
import Settings from "./pages/Settings";
import Signup from "./pages/Signup";
import TransactionHistory from "./pages/TransactionHistory";
import RiskProfile from "./pages/Onboarding/RiskProfile";
import WalletSetup from "./pages/Onboarding/WalletSetup";

import { ROUTES } from "./routes";

function NotFoundComponent() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="max-w-md text-center">
        <h1 className="text-7xl font-bold text-foreground">404</h1>
        <h2 className="mt-4 text-xl font-semibold text-foreground">Page not found</h2>
        <p className="mt-2 text-sm text-muted-foreground">
          The page you're looking for doesn't exist or has been moved.
        </p>
        <div className="mt-6">
          <Link
            to="/"
            className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
          >
            Go home
          </Link>
        </div>
      </div>
    </div>
  );
}

function ErrorComponent({ error, reset }: { error: Error; reset: () => void }) {
  console.error(error);
  const router = useRouter();

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="max-w-md text-center">
        <h1 className="text-xl font-semibold tracking-tight text-foreground">
          This page didn't load
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Something went wrong on our end. You can try refreshing or head back home.
        </p>
        <div className="mt-6 flex flex-wrap justify-center gap-2">
          <button
            onClick={() => {
              router.invalidate();
              reset();
            }}
            className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
          >
            Try again
          </button>
          <a
            href="/"
            className="inline-flex items-center justify-center rounded-md border border-input bg-background px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent"
          >
            Go home
          </a>
        </div>
      </div>
    </div>
  );
}

function RootShell({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <Scripts />
      </body>
    </html>
  );
}

import { ToastContainer } from "./components/ui/Toast";
import { ErrorBoundary } from "./components/layout/ErrorBoundary";
import { useGlobalErrorHandler } from "./hooks/useGlobalErrorHandler";
import { useNotifications } from "./hooks/useNotifications";

function RootComponent() {
  const { queryClient } = rootRoute.useRouteContext();
  useGlobalErrorHandler();
  useNotifications();

  return (
    <QueryClientProvider client={queryClient}>
      <ErrorBoundary>
        <Outlet />
        <ToastContainer />
      </ErrorBoundary>
    </QueryClientProvider>
  );
}

export const rootRoute = createRootRouteWithContext<{ queryClient: QueryClient }>()({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: "Pebble — fine your impulses, invest the difference" },
      { name: "description", content: "Pebble is the wallet that scores every spend for impulsiveness and auto-invests the penalty." },
      { name: "author", content: "Pebble" },
      { property: "og:title", content: "Pebble — fine your impulses, invest the difference" },
      { property: "og:description", content: "Score every spend. Penalize the impulse. Auto-invest the difference." },
      { property: "og:type", content: "website" },
      { name: "twitter:card", content: "summary_large_image" },
    ],
    links: [
      {
        rel: "stylesheet",
        href: appCss,
      },
    ],
  }),
  shellComponent: RootShell,
  component: RootComponent,
  notFoundComponent: NotFoundComponent,
  errorComponent: ErrorComponent,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.landing,
  component: Landing,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.login,
  component: Login,
});

const signupRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.signup,
  component: Signup,
});

function withAuth(Component: React.ComponentType) {
  return function AuthWrapped() {
    return (
      <ProtectedRoute>
        <Component />
      </ProtectedRoute>
    );
  };
}

const onboardingRiskRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.onboardingRisk,
  component: RiskProfile,
});

const onboardingWalletRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.onboardingWallet,
  component: WalletSetup,
});

const dashboardRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.dashboard,
  component: withAuth(Dashboard),
});

const logTransactionRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.logTransaction,
  component: withAuth(LogTransaction),
});

const portfolioRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.portfolio,
  component: withAuth(Portfolio),
});

const insightsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.insights,
  component: withAuth(Insights),
});

const historyRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.history,
  component: withAuth(TransactionHistory),
});

const settingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: ROUTES.settings,
  component: withAuth(Settings),
});

const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  signupRoute,
  onboardingRiskRoute,
  onboardingWalletRoute,
  dashboardRoute,
  logTransactionRoute,
  portfolioRoute,
  insightsRoute,
  historyRoute,
  settingsRoute,
]);

export const getRouter = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: (failureCount, error: any) => {
          // Don't retry auth errors or rate limit errors
          if (error?.status === 401 || error?.status === 403 || error?.status === 429) {
            return false;
          }
          return failureCount < 3;
        },
      },
    },
  });

  const router = createRouter({
    routeTree,
    context: { queryClient },
    scrollRestoration: true,
    defaultPreloadStaleTime: 0,
  });

  return router;
};

// Export router instance directly as well
export const router = getRouter();

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
