import api from "./client";
import type { HouseRanking, RankedStudent, Stats } from "@/types";

export async function getHouseRankings(): Promise<HouseRanking[]> {
  const { data } = await api.get<HouseRanking[]>("/rankings/houses");
  return data;
}

export async function getStudentRankings(): Promise<RankedStudent[]> {
  const { data } = await api.get<RankedStudent[]>("/rankings/students");
  return data;
}

export async function getStats(): Promise<Stats> {
  const { data } = await api.get<Stats>("/rankings/stats");
  return data;
}
