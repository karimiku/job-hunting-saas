import { beforeEach, describe, expect, it, vi } from "vitest";

const { serverFetch, revalidatePath } = vi.hoisted(() => ({
  serverFetch: vi.fn(),
  revalidatePath: vi.fn(),
}));

vi.mock("@/lib/api/server", () => ({ serverFetch }));
vi.mock("next/cache", () => ({ revalidatePath }));

import {
  createTaskFromTaskPageAction,
  deleteTaskAction,
  setTaskStatusAction,
} from "./actions";

function form(fields: Record<string, string>): FormData {
  const fd = new FormData();
  for (const [key, value] of Object.entries(fields)) fd.set(key, value);
  return fd;
}

describe("task actions", () => {
  beforeEach(() => {
    serverFetch.mockReset();
    revalidatePath.mockReset();
  });

  it("createTaskFromTaskPageAction は entry 配下にタスクを作成して関連画面を revalidate する", async () => {
    serverFetch.mockResolvedValue({
      id: "t1",
      entryId: "e1",
      title: "一次面接",
      type: "schedule",
      status: "todo",
      dueDate: "2026-06-01T00:00:00Z",
      memo: "",
      createdAt: "x",
      updatedAt: "x",
    });

    const result = await createTaskFromTaskPageAction(
      {},
      form({
        entryId: "e1",
        title: "一次面接",
        type: "schedule",
        dueDate: "2026-06-01",
        memo: "オンライン",
      }),
    );

    expect(result.ok).toBe(true);
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/entries/e1/tasks", {
      method: "POST",
      body: JSON.stringify({
        title: "一次面接",
        type: "schedule",
        dueDate: "2026-06-01T00:00:00.000Z",
        memo: "オンライン",
      }),
    });
    expect(revalidatePath).toHaveBeenCalledWith("/task");
    expect(revalidatePath).toHaveBeenCalledWith("/dashboard");
    expect(revalidatePath).toHaveBeenCalledWith("/entry/e1");
  });

  it("createTaskFromTaskPageAction は title が空なら API を呼ばない", async () => {
    const result = await createTaskFromTaskPageAction(
      {},
      form({ entryId: "e1", title: "   ", type: "deadline" }),
    );

    expect(result.error).toContain("タスク名");
    expect(serverFetch).not.toHaveBeenCalled();
  });

  it("setTaskStatusAction は status を PATCH する", async () => {
    serverFetch.mockResolvedValue({ status: "done" });

    const result = await setTaskStatusAction("t1", "done", "e1");

    expect(result).toEqual({ ok: true, status: "done" });
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/tasks/t1", {
      method: "PATCH",
      body: JSON.stringify({ status: "done" }),
    });
    expect(revalidatePath).toHaveBeenCalledWith("/task");
    expect(revalidatePath).toHaveBeenCalledWith("/task/t1");
    expect(revalidatePath).toHaveBeenCalledWith("/entry/e1");
  });

  it("deleteTaskAction は task を DELETE して関連画面を revalidate する", async () => {
    serverFetch.mockResolvedValue(undefined);

    const result = await deleteTaskAction("t1", "e1");

    expect(result).toEqual({ ok: true });
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/tasks/t1", {
      method: "DELETE",
    });
    expect(revalidatePath).toHaveBeenCalledWith("/task");
    expect(revalidatePath).toHaveBeenCalledWith("/task/t1");
    expect(revalidatePath).toHaveBeenCalledWith("/dashboard");
    expect(revalidatePath).toHaveBeenCalledWith("/entry/e1");
  });
});
