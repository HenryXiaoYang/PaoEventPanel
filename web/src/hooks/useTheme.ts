import { useEffect } from "react";
import { usePolling } from "./usePolling";
import { getActivitySettings } from "@/api/activity";
import { useDarkMode } from "./useDarkMode";
import type { ActivitySettings } from "@/types";

const SETTINGS_POLL_INTERVAL = 30000;

const defaultSettings: ActivitySettings = {
  id: 1,
  activity_name: "Campus Event",
  logo_path: "",
  bg_mode: "gradient",
  bg_color: "#f8fafc",
  bg_gradient_start: "#ffffff",
  bg_gradient_end: "#f3f4f6",
  bg_gradient_angle: 135,
  bg_image_path: "",
  bg_blur: false,
  dark_bg_mode: "gradient",
  dark_bg_color: "#0f1117",
  dark_bg_gradient_start: "#1a1c25",
  dark_bg_gradient_end: "#0f1117",
  dark_bg_gradient_angle: 135,
  dark_bg_image_path: "",
  dark_bg_blur: false,
  primary_color: "#7CB99A",
  secondary_color: "#f3f4f6",
  accent_color: "#9ABBE0",
  muted_color: "#A0AEC0",
  dark_primary_color: "#8ECFAB",
  dark_secondary_color: "#252833",
  dark_accent_color: "#A8C8E8",
  dark_muted_color: "#5c6173",
  card_opacity: 90,
  card_blur: 12,
  dark_card_opacity: 85,
  dark_card_blur: 16,
};

export function useTheme() {
  const { data: settings, refresh } = usePolling(getActivitySettings, SETTINGS_POLL_INTERVAL);
  const { isDark } = useDarkMode();
  const current = settings ?? defaultSettings;

  useEffect(() => {
    const root = document.documentElement;
    if (isDark) {
      root.style.setProperty("--theme-primary", current.dark_primary_color || current.primary_color);
      root.style.setProperty("--theme-secondary", current.dark_secondary_color || current.secondary_color);
      root.style.setProperty("--theme-accent", current.dark_accent_color || current.accent_color);
      root.style.setProperty("--theme-muted", current.dark_muted_color || current.muted_color);
      root.style.setProperty("--card-opacity", String((current.dark_card_opacity ?? 85) / 100));
      root.style.setProperty("--card-blur", `${current.dark_card_blur ?? 16}px`);
    } else {
      root.style.setProperty("--theme-primary", current.primary_color);
      root.style.setProperty("--theme-secondary", current.secondary_color);
      root.style.setProperty("--theme-accent", current.accent_color);
      root.style.setProperty("--theme-muted", current.muted_color);
      root.style.setProperty("--card-opacity", String((current.card_opacity ?? 90) / 100));
      root.style.setProperty("--card-blur", `${current.card_blur ?? 12}px`);
    }
  }, [
    isDark,
    current.primary_color, current.secondary_color, current.accent_color, current.muted_color,
    current.dark_primary_color, current.dark_secondary_color, current.dark_accent_color, current.dark_muted_color,
    current.card_opacity, current.card_blur, current.dark_card_opacity, current.dark_card_blur,
  ]);

  return { settings: current, refreshTheme: refresh };
}
