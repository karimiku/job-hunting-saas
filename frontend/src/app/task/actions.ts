"use server";

// Server Action — タスクのステータスを切り替える。
// PATCH /api/v1/tasks/{taskId} を server cookie 付きで叩き、成功後に /task を revalidate する
// (Client Component 側は router.refresh() で SSR を再評価し、最新の tasks を受け取る)。

import { revalidatePath } from "next/cache";
import { serverFetch } from "@/lib/api/server";
import type { TaskResponse } from "@/lib/api/tasks";

export interface SetTaskStatusResult {
  ok: boolean;
  status?: TaskResponse["status"];
  error?: string;
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
