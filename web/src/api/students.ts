import api from "./client";
import type { Student } from "@/types";

export async function getStudents(houseId?: number): Promise<Student[]> {
  const params = houseId ? { house_id: houseId } : {};
  const { data } = await api.get<Student[]>("/students", { params });
  return data;
}

export async function createStudent(name: string, houseId: number, lapCount?: number): Promise<Student> {
  const { data } = await api.post<Student>("/students", { name, house_id: houseId, lap_count: lapCount ?? 0 });
  return data;
}

export async function updateStudent(id: number, updates: { name?: string; house_id?: number }): Promise<Student> {
  const { data } = await api.put<Student>(`/students/${id}`, updates);
  return data;
}

export async function deleteStudent(id: number): Promise<void> {
  await api.delete(`/students/${id}`);
}

export async function addLaps(id: number, delta: number): Promise<Student> {
  const { data } = await api.post<Student>(`/students/${id}/laps`, { delta });
  return data;
}

export async function setLaps(id: number, lapCount: number): Promise<Student> {
  const { data } = await api.put<Student>(`/students/${id}/laps`, { lap_count: lapCount });
  return data;
}
