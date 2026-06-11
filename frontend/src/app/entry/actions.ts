"use server";

// Entry 詳細 / Kanban から使う Server Action。
// クライアントから backend を直接叩く (cross-origin + CORS preflight が毎回乗る) のをやめ、
// 同一 origin の Server Action 経由にする。revalidatePath がレスポンスに更新済み RSC ツリーを
// 含めるため、Client 側の router.refresh() も不要になる。

import { revalidatePath } from "next/cache";
import { serverFetch } from "@/lib/api/server";
import type { EntryResponse, UpdateEntryInput } from "@/lib/api/entries";
import type { CreateTaskInput, TaskResponse } from "@/lib/api/tasks";

export interface UpdateEntryResult {
  ok: boolean;
  error?: string;
}

export interface CreateTaskResult {
  ok: boolean;
  task?: TaskResponse;
  error?: string;
}

// entry のステージ/ステータスを PATCH し、entry を表示する各画面を revalidate する。
export async function updateEntryAction(
  entryId: string,
  input: UpdateEntryInput,
): Promise<UpdateEntryResult> {
  try {
    await serverFetch<EntryResponse>(`/api/v1/entries/${entryId}`, {
      method: "PATCH",
      body: JSON.stringify(input),
    });
  } catch {
    return { ok: false, error: "選考ステータスの更新に失敗しました" };
  }
  revalidatePath(`/entry/${entryId}`);
  revalidatePath("/entry");
  revalidatePath("/kanban");
  revalidatePath("/dashboard");
  return { ok: true };
}

// entry 配下にタスクを作成する。作成したタスクを返し、楽観表示の確定に使えるようにする。
export async function createTaskForEntryAction(
  entryId: string,
  input: CreateTaskInput,
): Promise<CreateTaskResult> {
  let task: TaskResponse;
  try {
    task = await serverFetch<TaskResponse>(`/api/v1/entries/${entryId}/tasks`, {
      method: "POST",
      body: JSON.stringify(input),
    });
  } catch {
    return { ok: false, error: "タスクの追加に失敗しました" };
  }
  revalidatePath(`/entry/${entryId}`);
  revalidatePath("/task");
  revalidatePath("/dashboard");
  return { ok: true, task };
}
