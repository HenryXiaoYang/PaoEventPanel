import api from "./client";
import type { ActivitySettings } from "@/types";

export async function getActivitySettings(): Promise<ActivitySettings> {
  const { data } = await api.get<ActivitySettings>("/activity");
  return data;
}

export async function updateActivitySettings(updates: Partial<ActivitySettings>): Promise<ActivitySettings> {
  const { data } = await api.put<ActivitySettings>("/activity", updates);
  return data;
}

export async function uploadLogo(file: File): Promise<{ path: string }> {
  const formData = new FormData();
  formData.append("file", file);
  const { data } = await api.post<{ path: string }>("/upload/logo", formData);
  return data;
}

export async function uploadBackground(file: File): Promise<{ path: string }> {
  const formData = new FormData();
  formData.append("file", file);
  const { data } = await api.post<{ path: string }>("/upload/background", formData);
  return data;
}
