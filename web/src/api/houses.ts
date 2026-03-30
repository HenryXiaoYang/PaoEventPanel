import api from "./client";
import type { House } from "@/types";

export async function getHouses(): Promise<House[]> {
  const { data } = await api.get<House[]>("/houses");
  return data;
}

export async function updateHouseColor(id: number, color: string): Promise<House> {
  const { data } = await api.put<House>(`/houses/${id}/color`, { color });
  return data;
}
