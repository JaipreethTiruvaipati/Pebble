import { useQuery } from "@tanstack/react-query";
import * as settingsApi from "@/api/settings.api";

export function useSettings() {
  return useQuery({
    queryKey: ["settings", "profile"],
    queryFn: settingsApi.getSettings,
  });
}
