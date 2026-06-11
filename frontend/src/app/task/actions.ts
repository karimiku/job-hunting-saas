"use server";

// Server Action — タスクのステータスを切り替える。
// PATCH /api/v1/tasks/{taskId} を server cookie 付きで叩き、成功後に /task を revalidate する。
// revalidatePath は action レスポンスに更新済み RSC ツリーを含めるため、Client 側の
// router.refresh() は不要 (二重フルレンダーになる)。

import { revalidatePath } from "next/cache";
import { serverFetch } from "@/lib/api/server";
import type { TaskResponse } from "@/lib/api/tasks";

export interface SetTaskStatusResult {
  ok: boolean;
  status?: TaskResponse["status"];
  error?: string;
}

export interface DeleteTaskResult {
  ok: boolean;
  error?: string;
}

export interface CreateTaskFormState {
  ok?: boolean;
  error?: string;
  values?: {
    entryId: string;
    title: string;
    type: TaskResponse["type"];
    dueDate: string;
    memo: string;
  };
}

export async function deleteTaskAction(taskId: string): Promise<DeleteTaskResult> {
  try {
    await serverFetch<void>(`/api/v1/tasks/${taskId}`, {
      method: "DELETE",
    });
    revalidatePath("/task");
    revalidatePath("/dashboard");
    return { ok: true };
  } catch {
    return { ok: false, error: "タスクの削除に失敗しました" };
  }
}

function readField(form: FormData, name: string, fallback = ""): string {
  const v = form.get(name);
  return typeof v === "string" ? v : fallback;
}

export async function setTaskStatusAction(
  taskId: string,
  status: TaskResponse["status"],
): Promise<SetTaskStatusResult> {
  try {
    const updated = await serverFetch<TaskResponse>(`/api/v1/tasks/${taskId}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    });
    revalidatePath("/task");
    return { ok: true, status: updated.status };
  } catch {
    return { ok: false, error: "タスクの更新に失敗しました" };
  }
}

export async function createTaskFromTaskPageAction(
  _prev: CreateTaskFormState,
  formData: FormData,
): Promise<CreateTaskFormState> {
  const entryId = readField(formData, "entryId").trim();
  const title = readField(formData, "title").trim();
  const typeRaw = readField(formData, "type", "deadline");
  const dueDateRaw = readField(formData, "dueDate").trim();
  const memo = readField(formData, "memo").trim();
  const type: TaskResponse["type"] =
    typeRaw === "schedule" ? "schedule" : "deadline";
  const values = { entryId, title, type, dueDate: dueDateRaw, memo };

  if (!entryId) return { error: "Entry を選択してください", values };
  if (!title) return { error: "タスク名は必須です", values };

  const body: {
    title: string;
    type: TaskResponse["type"];
    dueDate?: string;
    memo?: string;
  } = { title, type };
  if (dueDateRaw) body.dueDate = `${dueDateRaw}T00:00:00.000Z`;
  if (memo) body.memo = memo;

  try {
    await serverFetch<TaskResponse>(`/api/v1/entries/${entryId}/tasks`, {
      method: "POST",
      body: JSON.stringify(body),
    });
  } catch {
    return { error: "タスクの作成に失敗しました", values };
  }

  revalidatePath("/task");
  revalidatePath("/dashboard");
  revalidatePath(`/entry/${entryId}`);
  return {
    ok: true,
    values: { entryId, title: "", type, dueDate: "", memo: "" },
  };
}
