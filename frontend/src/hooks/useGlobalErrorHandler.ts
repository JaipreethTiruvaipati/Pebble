import { useEffect } from "react";
import { useNavigate } from "@tanstack/react-router";
import { onApiError, type ApiClientError } from "@/api/client";
import { toast } from "@/components/ui/Toast";
import { ROUTES } from "@/routes";

/**
 * Global API error interceptor — mount once at app root.
 * Shows toast notifications for API errors and redirects on session expiry.
 */
export function useGlobalErrorHandler() {
  const navigate = useNavigate();

  useEffect(() => {
    const unsubscribe = onApiError((error: ApiClientError) => {
      switch (error.status) {
        case 401:
          toast("Session expired — please log in again.", "error");
          navigate({ to: ROUTES.login });
          break;
        case 403:
          toast("You don't have permission to do that.", "error");
          break;
        case 404:
          // Silently ignore 404s for polling endpoints
          break;
        case 422:
          toast(error.message || "Invalid input. Please check your data.", "error");
          break;
        case 429:
          toast("Too many requests — please slow down.", "error");
          break;
        case 500:
        case 502:
        case 503:
          toast("Server error — please try again later.", "error");
          break;
        default:
          if (error.status >= 400) {
            toast(error.message || "Something went wrong.", "error");
          }
      }
    });

    return unsubscribe;
  }, [navigate]);
}
