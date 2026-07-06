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

export interface RescheduleTaskResult {
  ok: boolean;
  dueDate?: TaskResponse["dueDate"];
  error?: string;
}

export interface UpdateTaskActionInput {
  title?: string;
  type?: TaskResponse["type"];
  dueDate?: string | null;
  memo?: string;
}

export interface UpdateTaskResult {
  ok: boolean;
  task?: TaskResponse;
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

export async function deleteTaskAction(
  taskId: string,
  entryId?: string,
): Promise<DeleteTaskResult> {
  try {
    await serverFetch<void>(`/api/v1/tasks/${taskId}`, {
      method: "DELETE",
    });
    revalidatePath("/task");
    revalidatePath(`/task/${taskId}`);
    revalidatePath("/dashboard");
    if (entryId) revalidatePath(`/entry/${entryId}`);
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
  entryId?: string,
): Promise<SetTaskStatusResult> {
  try {
    const updated = await serverFetch<TaskResponse>(`/api/v1/tasks/${taskId}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    });
    revalidatePath("/task");
    revalidatePath(`/task/${taskId}`);
    if (entryId) revalidatePath(`/entry/${entryId}`);
    return { ok: true, status: updated.status };
  } catch {
    return { ok: false, error: "タスクの更新に失敗しました" };
  }
}

const DATE_ONLY = /^\d{4}-\d{2}-\d{2}$/;

// dueDate は "YYYY-MM-DD" を受け取り、他アクションと同じ 00:00:00.000Z 起点の ISO に変換して PATCH する。
export async function rescheduleTaskAction(
  taskId: string,
  dueDate: string,
  entryId?: string,
): Promise<RescheduleTaskResult> {
  if (!DATE_ONLY.test(dueDate)) {
    return { ok: false, error: "期日の形式が不正です" };
  }
  try {
    const updated = await serverFetch<TaskResponse>(`/api/v1/tasks/${taskId}`, {
      method: "PATCH",
      body: JSON.stringify({ dueDate: `${dueDate}T00:00:00.000Z` }),
    });
    revalidatePath("/task");
    revalidatePath(`/task/${taskId}`);
    if (entryId) revalidatePath(`/entry/${entryId}`);
    return { ok: true, dueDate: updated.dueDate };
  } catch {
    return { ok: false, error: "タスクの延期に失敗しました" };
  }
}

// 編集フォームからのまとめて更新。dueDate は "YYYY-MM-DD" または null (クリア) を受け取る。
export async function updateTaskAction(
  taskId: string,
  input: UpdateTaskActionInput,
  entryId?: string,
): Promise<UpdateTaskResult> {
  const body: {
    title?: string;
    type?: TaskResponse["type"];
    dueDate?: string | null;
    memo?: string;
  } = {};
  if (input.title !== undefined) body.title = input.title;
  if (input.type !== undefined) body.type = input.type;
  if (input.memo !== undefined) body.memo = input.memo;
  if (input.dueDate !== undefined) {
    if (input.dueDate && !DATE_ONLY.test(input.dueDate)) {
      return { ok: false, error: "期日の形式が不正です" };
    }
    body.dueDate = input.dueDate ? `${input.dueDate}T00:00:00.000Z` : null;
  }

  try {
    const updated = await serverFetch<TaskResponse>(`/api/v1/tasks/${taskId}`, {
      method: "PATCH",
      body: JSON.stringify(body),
    });
    revalidatePath("/task");
    revalidatePath(`/task/${taskId}`);
    if (entryId) revalidatePath(`/entry/${entryId}`);
    return { ok: true, task: updated };
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
