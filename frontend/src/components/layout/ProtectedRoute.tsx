import { useEffect } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useAuthStore } from "@/stores/authStore";
import { ROUTES } from "@/routes";
import type { ReactNode } from "react";

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const navigate = useNavigate();
  const { isAuthenticated, isLoading, loadProfile, user } = useAuthStore();

  useEffect(() => {
    useAuthStore.getState().hydrate();
    if (!localStorage.getItem("pebble_token")) {
      navigate({ to: ROUTES.login });
      return;
    }
    if (!user) void loadProfile();
  }, [loadProfile, navigate, user]);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate({ to: ROUTES.login });
    }
  }, [isAuthenticated, isLoading, navigate]);

  if (isLoading || !isAuthenticated) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background text-muted-foreground">
        Loading…
      </div>
    );
  }

  return <>{children}</>;
}
