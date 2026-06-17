import { describe, it, expect, vi, beforeEach } from "vitest";

// server.ts は next/headers (server-only) を読むので serverFetch を mock してパスだけ検証する。
const serverFetch = vi.fn();
vi.mock("./server", () => ({
  serverFetch: (path: string, init?: RequestInit) => serverFetch(path, init),
}));

import {
  attachCompanyNamesToTasks,
  buildNavCounts,
  getAppPageDataServer,
  getNavCountsServer,
  getTaskPageDataServer,
  getTaskServer,
  listAllTasksServer,
  listEntriesWithCompanyNamesServer,
  listTasksServer,
} from "./server-resources";
import { ApiError } from "./client-types";

beforeEach(() => {
  serverFetch.mockReset();
});

describe("listAllTasksServer", () => {
  it("全タスクAPIを1回だけ呼び、会社名付きで返す", async () => {
    serverFetch.mockImplementation(async (path: string) => {
      if (path === "/api/v1/tasks") {
        return {
          tasks: [
            { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
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
    expect(serverFetch).toHaveBeenCalledTimes(1);
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/tasks", undefined);
  });

  it("全タスクAPI取得に失敗したら呼び出し側にエラーを返す", async () => {
    serverFetch.mockRejectedValue(new Error("boom"));

    await expect(listAllTasksServer([
      { id: "e1", companyName: "A社" },
      { id: "e2", companyName: "B社" },
    ])).rejects.toThrow("boom");
  });

  it("entry が無ければ tasks API を呼ばず空配列を返す", async () => {
    const tasks = await listAllTasksServer([]);
    expect(tasks).toEqual([]);
    expect(serverFetch).not.toHaveBeenCalled();
  });
});

describe("listTasksServer", () => {
  it("全タスクAPIを entries に依存せず呼び出す", async () => {
    serverFetch.mockResolvedValue({
      tasks: [
        { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
      ],
    });

    const tasks = await listTasksServer();

    expect(tasks).toHaveLength(1);
    expect(serverFetch).toHaveBeenCalledTimes(1);
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/tasks", undefined);
  });
});

describe("getTaskServer", () => {
  it("タスク詳細APIを呼ぶ", async () => {
    serverFetch.mockResolvedValue({
      id: "t1",
      entryId: "e1",
      title: "ES提出",
      type: "deadline",
      status: "todo",
      dueDate: null,
      memo: "",
      createdAt: "x",
      updatedAt: "x",
    });

    await expect(getTaskServer("t1")).resolves.toEqual(
      expect.objectContaining({ id: "t1" }),
    );
    expect(serverFetch).toHaveBeenCalledTimes(1);
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/tasks/t1", undefined);
  });
});

describe("getTaskPageDataServer", () => {
  it("/task 用の集約APIを1回だけ呼び、tasks に会社名を付けて返す", async () => {
    serverFetch.mockResolvedValue({
      user: { id: "u1", email: "student@example.com", name: "Student" },
      entries: [
        {
          id: "e1",
          companyId: "c1",
          companyName: "○○商事",
          route: "",
          source: "",
          status: "open",
          stageKind: "pre_entry",
          stageLabel: "",
          memo: "",
          createdAt: "x",
          updatedAt: "x",
        },
      ],
      tasks: [
        { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
      ],
    });

    const pageData = await getTaskPageDataServer();

    expect(pageData?.user.email).toBe("student@example.com");
    expect(pageData?.entries).toHaveLength(1);
    expect(pageData?.tasks[0]?.companyName).toBe("○○商事");
    expect(serverFetch).toHaveBeenCalledTimes(1);
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/page-data/task", undefined);
  });

  it("401 は未ログインとして null を返す", async () => {
    serverFetch.mockRejectedValue(new ApiError(401, "Unauthorized"));

    await expect(getTaskPageDataServer()).resolves.toBeNull();
  });
});

describe("getAppPageDataServer", () => {
  it("主要画面用の集約APIを1回だけ呼び、navCounts と会社名付き tasks を返す", async () => {
    serverFetch.mockResolvedValue({
      user: { id: "u1", email: "student@example.com", name: "Student" },
      entries: [
        {
          id: "e1",
          companyId: "c1",
          companyName: "○○商事",
          route: "",
          source: "",
          status: "open",
          stageKind: "pre_entry",
          stageLabel: "",
          memo: "",
          createdAt: "x",
          updatedAt: "x",
        },
      ],
      tasks: [
        { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
        { id: "t2", entryId: "e1", title: "完了済み", type: "schedule", status: "done", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
      ],
      clips: [{ id: "clip1", url: "https://example.com", title: "求人", source: "web", guess: "", capturedAt: "x" }],
      companies: [{ id: "c1", name: "○○商事", memo: "", createdAt: "x", updatedAt: "x" }],
    });

    const pageData = await getAppPageDataServer();

    expect(pageData?.user.email).toBe("student@example.com");
    expect(pageData?.tasks[0]?.companyName).toBe("○○商事");
    expect(pageData?.companies).toHaveLength(1);
    expect(pageData?.navCounts).toEqual({ entry: 1, task: 1, inbox: 1 });
    expect(serverFetch).toHaveBeenCalledTimes(1);
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/page-data/app", undefined);
  });

  it("401 は未ログインとして null を返す", async () => {
    serverFetch.mockRejectedValue(new ApiError(401, "Unauthorized"));

    await expect(getAppPageDataServer()).resolves.toBeNull();
  });
});

describe("listEntriesWithCompanyNamesServer", () => {
  it("backend が会社名を同梱する entries API だけを呼ぶ", async () => {
    serverFetch.mockResolvedValue({
      entries: [
        {
          id: "e1",
          companyId: "c1",
          companyName: "○○商事",
          route: "",
          source: "",
          status: "open",
          stageKind: "pre_entry",
          stageLabel: "",
          memo: "",
          createdAt: "x",
          updatedAt: "x",
        },
      ],
    });

    await expect(listEntriesWithCompanyNamesServer()).resolves.toEqual([
      expect.objectContaining({ id: "e1", companyName: "○○商事" }),
    ]);
    expect(serverFetch).toHaveBeenCalledTimes(1);
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/entries", undefined);
  });
});

describe("attachCompanyNamesToTasks", () => {
  it("取得済み tasks と entries を entryId で join する", () => {
    const tasks = attachCompanyNamesToTasks(
      [
        { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
        { id: "t2", entryId: "missing", title: "SPI受験", type: "schedule", status: "done", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
      ],
      [{ id: "e1", companyName: "○○商事" }],
    );

    expect(tasks.find((t) => t.id === "t1")?.companyName).toBe("○○商事");
    expect(tasks.find((t) => t.id === "t2")?.companyName).toBeUndefined();
  });
});

describe("getNavCountsServer", () => {
  it("主要画面用の集約APIを1回だけ呼んで件数を返す", async () => {
    serverFetch.mockResolvedValue({
      user: { id: "u1", email: "student@example.com", name: "Student" },
      entries: [
        {
          id: "e1",
          companyId: "c1",
          route: "",
          source: "",
          status: "open",
          stageKind: "pre_entry",
          stageLabel: "",
          memo: "",
          createdAt: "x",
          updatedAt: "x",
        },
        {
          id: "e2",
          companyId: "c2",
          route: "",
          source: "",
          status: "open",
          stageKind: "pre_entry",
          stageLabel: "",
          memo: "",
          createdAt: "x",
          updatedAt: "x",
        },
      ],
      tasks: [
        { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
        { id: "t2", entryId: "e2", title: "SPI受験", type: "schedule", status: "done", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
      ],
      clips: [{ id: "clip1", url: "https://example.com", title: "求人", source: "web", guess: "", capturedAt: "x" }],
      companies: [],
    });

    await expect(getNavCountsServer()).resolves.toEqual({
      entry: 2,
      task: 1,
      inbox: 1,
    });
    expect(serverFetch).toHaveBeenCalledTimes(1);
    expect(serverFetch).toHaveBeenCalledWith("/api/v1/page-data/app", undefined);
  });
});

describe("buildNavCounts", () => {
  it("取得済み entries, tasks, clips からサイドバー件数を作る", () => {
    const navCounts = buildNavCounts(
      [
        {
          id: "e1",
          companyId: "c1",
          route: "",
          source: "",
          status: "open",
          stageKind: "pre_entry",
          stageLabel: "",
          memo: "",
          createdAt: "x",
          updatedAt: "x",
        },
        {
          id: "e2",
          companyId: "c2",
          route: "",
          source: "",
          status: "open",
          stageKind: "pre_entry",
          stageLabel: "",
          memo: "",
          createdAt: "x",
          updatedAt: "x",
        },
      ],
      [
        { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
        { id: "t2", entryId: "e2", title: "SPI受験", type: "schedule", status: "done", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
      ],
      [
        {
          id: "clip1",
          url: "https://example.com",
          title: "clip",
          source: "web",
          guess: "",
          capturedAt: "x",
        },
      ],
    );

    expect(navCounts).toEqual({
      entry: 2,
      task: 1,
      inbox: 1,
    });
  });
});
