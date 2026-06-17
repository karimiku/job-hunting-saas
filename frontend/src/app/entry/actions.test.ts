import { beforeEach, describe, expect, it, vi } from "vitest";

const { redirect, revalidatePath, serverFetch } = vi.hoisted(() => ({
  redirect: vi.fn(),
  revalidatePath: vi.fn(),
  serverFetch: vi.fn(),
}));

vi.mock("@/lib/api/server", () => ({ serverFetch }));
vi.mock("next/cache", () => ({ revalidatePath }));
vi.mock("next/navigation", () => ({ redirect }));

import { deleteEntryAction, updateEntryAction } from "./actions";

describe("entry actions", () => {
  beforeEach(() => {
    redirect.mockReset();
    revalidatePath.mockReset();
    serverFetch.mockReset();
  });

  it("updateEntryAction は entry を PATCH して関連画面を revalidate する", async () => {
    serverFetch.mockResolvedValue({
      id: "e1",
      status: "in_progress",
      stageKind: "interview",
      stageLabel: "一次面接",
    });

    const result = await updateEntryAction("e1", { stageKind: "interview" });

    expect(result).toEqual({ ok: true });
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/entries/e1", {
      method: "PATCH",
      body: JSON.stringify({ stageKind: "interview" }),
    });
    expect(revalidatePath).toHaveBeenCalledWith("/entry/e1");
    expect(revalidatePath).toHaveBeenCalledWith("/entry");
    expect(revalidatePath).toHaveBeenCalledWith("/kanban");
    expect(revalidatePath).toHaveBeenCalledWith("/dashboard");
  });

  it("deleteEntryAction は entry を DELETE して一覧に戻す", async () => {
    serverFetch.mockResolvedValue(undefined);

    await deleteEntryAction("e1");

    expect(serverFetch).toHaveBeenCalledWith("/api/v1/entries/e1", {
      method: "DELETE",
    });
    expect(revalidatePath).toHaveBeenCalledWith("/entry");
    expect(revalidatePath).toHaveBeenCalledWith("/task");
    expect(revalidatePath).toHaveBeenCalledWith("/kanban");
    expect(revalidatePath).toHaveBeenCalledWith("/dashboard");
    expect(redirect).toHaveBeenCalledWith("/entry");
  });

  it("deleteEntryAction は DELETE 失敗時にエラーを返す", async () => {
    serverFetch.mockRejectedValue(new Error("boom"));

    const result = await deleteEntryAction("e1");

    expect(result).toEqual({ ok: false, error: "Entryの削除に失敗しました" });
    expect(redirect).not.toHaveBeenCalled();
  });
});
