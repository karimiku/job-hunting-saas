import { describe, it, expect, vi, beforeEach } from "vitest";

// server.ts は next/headers (server-only) を読むので serverFetch を mock してパスだけ検証する。
const serverFetch = vi.fn();
vi.mock("./server", () => ({
  serverFetch: (path: string, init?: RequestInit) => serverFetch(path, init),
}));

import { listAllTasksServer } from "./server-resources";

beforeEach(() => {
  serverFetch.mockReset();
});

describe("listAllTasksServer", () => {
  it("渡された entry ごとに tasks を引き、会社名付きで集約する", async () => {
    serverFetch.mockImplementation(async (path: string) => {
      if (path === "/api/v1/entries/e1/tasks") {
        return {
          tasks: [
            { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
          ],
        };
      }
      if (path === "/api/v1/entries/e2/tasks") {
        return {
          tasks: [
            { id: "t2", entryId: "e2", title: "SPI受験", type: "schedule", status: "done", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
          ],
        };
      }
      throw new Error(`unexpected path: ${path}`);
    });

    const tasks = await listAllTasksServer([
      { id: "e1", companyName: "○○商事" },
      { id: "e2", companyName: "△△株式会社" },
    ]);

    expect(tasks).toHaveLength(2);
    expect(tasks.find((t) => t.id === "t1")?.companyName).toBe("○○商事");
    expect(tasks.find((t) => t.id === "t2")?.companyName).toBe("△△株式会社");
  });

  it("個別 entry の tasks 取得に失敗しても取れたぶんだけ返す", async () => {
    serverFetch.mockImplementation(async (path: string) => {
      if (path === "/api/v1/entries/e1/tasks") {
        return {
          tasks: [
            { id: "t1", entryId: "e1", title: "面接準備", type: "schedule", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
          ],
        };
      }
      // e2 のタスク取得は失敗させる
      throw new Error("boom");
    });

    const tasks = await listAllTasksServer([
      { id: "e1", companyName: "A社" },
      { id: "e2", companyName: "B社" },
    ]);
    expect(tasks).toHaveLength(1);
    expect(tasks[0].id).toBe("t1");
  });

  it("entry が無ければ tasks API を呼ばず空配列を返す", async () => {
    const tasks = await listAllTasksServer([]);
    expect(tasks).toEqual([]);
    expect(serverFetch).not.toHaveBeenCalled();
  });
});
