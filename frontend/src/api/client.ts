import type { ApiError } from "@/types/api.types";

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1";

const TOKEN_KEY = "pebble_token";
const REFRESH_KEY = "pebble_refresh_token";
const USER_ID_KEY = "pebble_user_id";

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(TOKEN_KEY);
}

export function setAuth(token: string, userId: string, refreshToken?: string) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_ID_KEY, userId);
  if (refreshToken) localStorage.setItem(REFRESH_KEY, refreshToken);
}

export function clearAuth() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(REFRESH_KEY);
  localStorage.removeItem(USER_ID_KEY);
}

export function getUserId(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(USER_ID_KEY);
}

export class ApiClientError extends Error {
  status: number;
  code: string;
  constructor(status: number, message: string, code: string) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

// Global error listeners for cross-cutting concerns (toast notifications, logging)
type ErrorListener = (error: ApiClientError) => void;
const errorListeners: ErrorListener[] = [];

/** Register a global error handler (e.g., toast notifications). Returns unsubscribe fn. */
export function onApiError(listener: ErrorListener): () => void {
  errorListeners.push(listener);
  return () => {
    const idx = errorListeners.indexOf(listener);
    if (idx >= 0) errorListeners.splice(idx, 1);
  };
}

function notifyErrorListeners(error: ApiClientError) {
  for (const listener of errorListeners) {
    try { listener(error); } catch { /* swallow listener errors */ }
  }
}

/** Attempt to refresh the access token using the stored refresh token. */
async function tryRefreshToken(): Promise<boolean> {
  const refreshToken = localStorage.getItem(REFRESH_KEY);
  if (!refreshToken) return false;

  try {
    const res = await fetch(`${BASE_URL}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!res.ok) return false;

    const data = await res.json();
    if (data.token) {
      localStorage.setItem(TOKEN_KEY, data.token);
      if (data.refresh_token) localStorage.setItem(REFRESH_KEY, data.refresh_token);
      return true;
    }
    return false;
  } catch {
    return false;
  }
}

// Prevent concurrent refresh attempts
let refreshPromise: Promise<boolean> | null = null;

export async function apiRequest<T>(
  path: string,
  options: RequestInit = {},
  auth = true,
): Promise<T> {
  const headers = new Headers(options.headers);
  if (!headers.has("Content-Type") && !(options.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }
  if (auth) {
    const token = getToken();
    if (token) headers.set("Authorization", `Bearer ${token}`);
  }

  const res = await fetch(`${BASE_URL}${path}`, { ...options, headers });

  // Handle 401: attempt token refresh once, then retry the original request
  if (res.status === 401 && auth) {
    if (!refreshPromise) {
      refreshPromise = tryRefreshToken().finally(() => { refreshPromise = null; });
    }
    const refreshed = await refreshPromise;
    if (refreshed) {
      // Retry with new token
      const retryHeaders = new Headers(options.headers);
      if (!retryHeaders.has("Content-Type") && !(options.body instanceof FormData)) {
        retryHeaders.set("Content-Type", "application/json");
      }
      const newToken = getToken();
      if (newToken) retryHeaders.set("Authorization", `Bearer ${newToken}`);

      const retryRes = await fetch(`${BASE_URL}${path}`, { ...options, headers: retryHeaders });
      if (retryRes.ok) {
        if (retryRes.status === 204) return undefined as T;
        const text = await retryRes.text();
        if (!text) return undefined as T;
        return JSON.parse(text) as T;
      }
    }

    // Refresh failed — clear auth and redirect to login
    clearAuth();
    const authError = new ApiClientError(401, "Session expired. Please log in again.", "SESSION_EXPIRED");
    notifyErrorListeners(authError);
    throw authError;
  }

  if (!res.ok) {
    let message = res.statusText;
    let code = "REQUEST_FAILED";
    try {
      const err = (await res.json()) as ApiError;
      message = err.error ?? message;
      code = err.code ?? code;
    } catch {
      /* ignore */
    }
    const error = new ApiClientError(res.status, message, code);
    notifyErrorListeners(error);
    throw error;
  }

  if (res.status === 204) return undefined as T;
  const text = await res.text();
  if (!text) return undefined as T;
  return JSON.parse(text) as T;
}

export { BASE_URL };

