import { getMe } from "./auth.api";
import type { UserProfile } from "@/types/api.types";

/** Settings read from user profile until PATCH /users/me is added. */
export function getSettings() {
  return getMe() as Promise<UserProfile>;
}
