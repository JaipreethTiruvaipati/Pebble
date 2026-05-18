import { apiRequest, setAuth } from "./client";
import type { AuthResponse, UserProfile } from "@/types/api.types";

export async function login(email: string, password: string, referralCode?: string) {
  const res = await apiRequest<AuthResponse>(
    "/auth/login",
    {
      method: "POST",
      body: JSON.stringify({ email, password, referral_code: referralCode ?? "" }),
    },
    false,
  );
  setAuth(res.token, res.user_id);
  return res;
}

export async function verifyOtp(phone: string, otp: string) {
  const res = await apiRequest<AuthResponse>(
    "/auth/verify-otp",
    { method: "POST", body: JSON.stringify({ phone, otp }) },
    false,
  );
  setAuth(res.token, res.user_id);
  return res;
}

export async function getMe() {
  return apiRequest<UserProfile>("/users/me");
}

export { clearAuth as logout } from "./client";
