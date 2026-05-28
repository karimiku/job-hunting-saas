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
  it("entries を引き、各 entry の tasks を会社名付きで集約する", async () => {
    serverFetch.mockImplementation(async (path: string) => {
      if (path === "/api/v1/entries") {
        return {
          entries: [
            { id: "e1", companyId: "c1", source: "リクナビ" },
            { id: "e2", companyId: "c2", source: "マイナビ" },
          ],
        };
      }
      if (path === "/api/v1/companies") {
        return {
          companies: [
            { id: "c1", name: "○○商事" },
            { id: "c2", name: "△△株式会社" },
          ],
        };
      }
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

    const tasks = await listAllTasksServer();

    expect(tasks).toHaveLength(2);
    expect(tasks.find((t) => t.id === "t1")?.companyName).toBe("○○商事");
    expect(tasks.find((t) => t.id === "t2")?.companyName).toBe("△△株式会社");
  });

  it("個別 entry の tasks 取得に失敗しても取れたぶんだけ返す", async () => {
    serverFetch.mockImplementation(async (path: string) => {
      if (path === "/api/v1/entries") {
        return { entries: [{ id: "e1", companyId: "c1", source: "x" }, { id: "e2", companyId: "c2", source: "y" }] };
      }
      if (path === "/api/v1/companies") {
        return { companies: [{ id: "c1", name: "A社" }, { id: "c2", name: "B社" }] };
      }
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

    const tasks = await listAllTasksServer();
    expect(tasks).toHaveLength(1);
    expect(tasks[0].id).toBe("t1");
  });

  it("タスクが1件も無ければ空配列を返す", async () => {
    serverFetch.mockImplementation(async (path: string) => {
      if (path === "/api/v1/entries") return { entries: [] };
      if (path === "/api/v1/companies") return { companies: [] };
      throw new Error(`unexpected path: ${path}`);
    });

    const tasks = await listAllTasksServer();
    expect(tasks).toEqual([]);
  });
});
