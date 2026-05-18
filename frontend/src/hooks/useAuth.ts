import { useEffect } from "react";
import { useAuthStore } from "@/stores/authStore";

export function useAuth() {
  const store = useAuthStore();

  useEffect(() => {
    store.hydrate();
    if (getToken() && !store.user) {
      void store.loadProfile();
    }
  }, []);

  return store;
}

function getToken() {
  try {
    return localStorage.getItem("pebble_token");
  } catch {
    return null;
  }
}
