import { create } from "zustand";

interface DarkModeState {
  isDark: boolean;
  toggle: () => void;
}

export const useDarkMode = create<DarkModeState>((set, get) => {
  const stored = localStorage.getItem("theme");
  const isDark = stored === "dark";
  if (isDark) {
    document.documentElement.classList.add("dark");
  }

  return {
    isDark,
    toggle: () => {
      const next = !get().isDark;

      const applyTheme = () => {
        if (next) {
          document.documentElement.classList.add("dark");
          localStorage.setItem("theme", "dark");
        } else {
          document.documentElement.classList.remove("dark");
          localStorage.setItem("theme", "light");
        }
        set({ isDark: next });
      };

      const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
      const hasViewTransition = "startViewTransition" in document;

      if (!hasViewTransition || reduceMotion) {
        applyTheme();
        return;
      }

      (document as any).startViewTransition(applyTheme);
    },
  };
});
