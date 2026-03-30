import api from "./client";
import type { User } from "@/types";

export async function getUsers(): Promise<User[]> {
  const { data } = await api.get<User[]>("/users");
  return data;
}

export async function createUser(username: string, password: string, role: string): Promise<User> {
  const { data } = await api.post<User>("/users", { username, password, role });
  return data;
}

export async function deleteUser(id: number): Promise<void> {
  await api.delete(`/users/${id}`);
}
