import { useEffect, useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useAuthStore } from "@/stores/authStore";
import { ROUTES } from "@/routes";
import type { ReactNode } from "react";

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const navigate = useNavigate();
  const { isAuthenticated, isLoading, loadProfile, user } = useAuthStore();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    useAuthStore.getState().hydrate();
    if (!localStorage.getItem("pebble_token")) {
      navigate({ to: ROUTES.login });
      return;
    }
    if (!user) void loadProfile();
  }, [loadProfile, navigate, user]);

  useEffect(() => {
    if (mounted && !isLoading && !isAuthenticated) {
      navigate({ to: ROUTES.login });
    }
  }, [isAuthenticated, isLoading, navigate, mounted]);

  if (!mounted || isLoading || !isAuthenticated) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background text-muted-foreground">
        Loading…
      </div>
    );
  }

  return <>{children}</>;
}
