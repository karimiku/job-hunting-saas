import { apiFetch } from "./client";

export interface TaskResponse {
  id: string;
  entryId: string;
  title: string;
  type: "deadline" | "schedule";
  status: "todo" | "done";
  dueDate: string | null;
  memo: string;
  createdAt: string;
  updatedAt: string;
}

export async function listTasksByEntry(entryId: string): Promise<TaskResponse[]> {
  const res = await apiFetch<{ tasks: TaskResponse[] }>(`/api/v1/entries/${entryId}/tasks`);
  return res.tasks;
}

export interface CreateTaskInput {
  title: string;
  type: "deadline" | "schedule";
  dueDate?: string;
  memo?: string;
}

export async function createTask(entryId: string, input: CreateTaskInput): Promise<TaskResponse> {
  return apiFetch<TaskResponse>(`/api/v1/entries/${entryId}/tasks`, {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export interface UpdateTaskInput {
  title?: string;
  status?: "todo" | "done";
  dueDate?: string | null;
  memo?: string;
}

export async function updateTask(taskId: string, input: UpdateTaskInput): Promise<TaskResponse> {
  return apiFetch<TaskResponse>(`/api/v1/tasks/${taskId}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
}

export async function deleteTask(taskId: string): Promise<void> {
  await apiFetch<void>(`/api/v1/tasks/${taskId}`, { method: "DELETE" });
}
