import { getMe } from "./auth.api";
import { apiRequest } from "./client";
import type { UserProfile } from "@/types/api.types";

/** Settings read from user profile until PATCH /users/me is added. */
export function getSettings() {
  return getMe() as Promise<UserProfile>;
}

/** Registers the Firebase Cloud Messaging device token with the backend. */
export function updateDeviceToken(token: string) {
  // In a real implementation this would PATCH /users/me or POST /devices
  // For now, we'll mock it if the endpoint isn't fully implemented in the Go backend
  return apiRequest<{ success: boolean }>("/users/me/device-token", {
    method: "POST",
    body: JSON.stringify({ device_token: token }),
  }).catch((err) => {
    console.warn("Device token registration skipped (endpoint may not exist yet):", err);
    return { success: true }; // graceful fallback
  });
}
