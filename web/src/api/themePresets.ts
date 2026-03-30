import api from "./client";
import type { ThemePreset, ActivitySettings } from "@/types";

export async function getThemePresets(): Promise<ThemePreset[]> {
  const { data } = await api.get<ThemePreset[]>("/theme-presets");
  return data;
}

export async function createThemePreset(name: string): Promise<ThemePreset> {
  const { data } = await api.post<ThemePreset>("/theme-presets", { name });
  return data;
}

export async function deleteThemePreset(id: number): Promise<void> {
  await api.delete(`/theme-presets/${id}`);
}

export async function applyThemePreset(id: number): Promise<ActivitySettings> {
  const { data } = await api.post<ActivitySettings>(`/theme-presets/${id}/apply`);
  return data;
}

export async function applyBuiltinPreset(name: string): Promise<ActivitySettings> {
  const { data } = await api.post<ActivitySettings>("/theme-presets/builtin/apply", { name });
  return data;
}
