import { useEffect, useState } from "react";
import { messaging, getToken, onMessage } from "@/lib/firebaseClient";
import { updateDeviceToken } from "@/api/settings.api";
import { toast } from "@/components/ui/Toast";
import { useAuthStore } from "@/stores/authStore";

/**
 * Hook to manage Firebase Cloud Messaging push notifications.
 * Requests permission, registers the device token, and listens for foreground messages.
 */
export function useNotifications() {
  const { isAuthenticated } = useAuthStore();
  const [permission, setPermission] = useState<NotificationPermission | "unsupported">(
    typeof window !== "undefined" && "Notification" in window ? window.Notification.permission : "unsupported"
  );

  useEffect(() => {
    if (!isAuthenticated || !messaging || permission === "unsupported") return;

    const registerToken = async () => {
      try {
        const status = await Notification.requestPermission();
        setPermission(status);

        if (status === "granted") {
          // You would typically get this from your Firebase console -> Project Settings -> Cloud Messaging
          const vapidKey = import.meta.env.VITE_FIREBASE_VAPID_KEY;
          const currentToken = await getToken(messaging, { vapidKey });

          if (currentToken) {
            await updateDeviceToken(currentToken);
            console.info("FCM device token registered.");
          } else {
            console.warn("No registration token available. Request permission to generate one.");
          }
        }
      } catch (error) {
        console.error("An error occurred while retrieving token: ", error);
      }
    };

    if (permission === "default") {
      void registerToken();
    } else if (permission === "granted") {
      // Re-register if already granted to ensure backend has the latest token
      void registerToken();
    }

    // Handle foreground messages
    const unsubscribe = onMessage(messaging, (payload) => {
      console.log("Message received in foreground:", payload);
      
      const title = payload.notification?.title || "New Notification";
      const body = payload.notification?.body || "";
      
      // We can use different toast types based on the data payload if needed
      toast(`${title}: ${body}`, "info");
    });

    return () => {
      unsubscribe();
    };
  }, [isAuthenticated, permission]);

  return {
    permission,
    requestPermission: async () => {
      if (permission === "unsupported") return;
      const status = await Notification.requestPermission();
      setPermission(status);
    }
  };
}
