import { describe, it, expect, vi, beforeEach } from "vitest";
import { ApiError } from "@/lib/api/client-types";

// serverFetch を差し替えて呼び出し履歴を検証する。
// （実体は next/headers の cookies を使うため、ここでは読み込ませない）
// vi.mock は巻き上げられるため、mock 本体は vi.hoisted で先に生成する。
const { serverFetch } = vi.hoisted(() => ({ serverFetch: vi.fn() }));
vi.mock("@/lib/api/server", () => ({ serverFetch }));
// revalidatePath はリクエストコンテキスト外だと throw するので no-op に。
vi.mock("next/cache", () => ({ revalidatePath: vi.fn() }));

import {
  convertInboxClipAction,
  deleteInboxClipAction,
  type ConvertClipFormState,
} from "./actions";

const INITIAL: ConvertClipFormState = {};

function form(fields: Record<string, string>): FormData {
  const fd = new FormData();
  for (const [k, val] of Object.entries(fields)) fd.set(k, val);
  return fd;
}

function callSignatures(): string[] {
  return serverFetch.mock.calls
    .filter(([p]) => typeof p === "string")
    .map(([p, i]) => `${(i as RequestInit | undefined)?.method ?? "GET"} ${p}`);
}

async function callAndCapture(
  fd: FormData,
): Promise<{ thrown?: unknown; result?: ConvertClipFormState }> {
  try {
    return { result: await convertInboxClipAction(INITIAL, fd) };
  } catch (thrown) {
    return { thrown };
  }
}

const companyResp = { id: "co1", name: "テスト商事", memo: "", createdAt: "", updatedAt: "" };
const entryResp = {
  id: "en1",
  companyId: "co1",
  route: "本選考",
  source: "マイナビ",
  status: "in_progress",
  stageKind: "application",
  stageLabel: "応募",
  memo: "",
  createdAt: "",
  updatedAt: "",
};
const flowResp = {
  id: "flow1",
  entryId: "en1",
  source: "template",
  currentStagePosition: 1,
  stages: [],
  createdAt: "",
  updatedAt: "",
};

describe("convertInboxClipAction", () => {
  beforeEach(() => serverFetch.mockReset());

  it("新規会社の場合は atomic endpoint で entry を作成し、その後 clip を削除して作成 Entry にリダイレクトする", async () => {
    serverFetch.mockImplementation(async (path?: string, init?: RequestInit) => {
      if (path === "/api/v1/entries/with-company" && init?.method === "POST") return entryResp;
      if (path === "/api/v1/entries/en1/selection-flow" && init?.method === "PUT") return flowResp;
      return undefined;
    });

    const fd = form({
      clipId: "clip1",
      companyName: "テスト商事",
      route: "本選考",
      source: "マイナビ",
      memo: "memo",
    });

    // 成功時は redirect() が throw する。digest に遷移先が載る。
    const { thrown } = await callAndCapture(fd);
    expect((thrown as { digest?: string } | undefined)?.digest).toContain("/entry/en1");

    const calls = callSignatures();
    expect(calls).toContain("POST /api/v1/entries/with-company");
    expect(calls).toContain("PUT /api/v1/entries/en1/selection-flow");
    expect(calls).not.toContain("POST /api/v1/companies");
    expect(calls).toContain("DELETE /api/v1/inbox/clips/clip1");
  });

  it("新規会社の atomic 作成に失敗したら clip は削除しない", async () => {
    serverFetch.mockImplementation(async (path?: string, init?: RequestInit) => {
      if (path === "/api/v1/entries/with-company" && init?.method === "POST") throw new ApiError(500, "boom");
      return undefined;
    });

    const fd = form({
      clipId: "clip1",
      companyName: "テスト商事",
      route: "本選考",
      source: "マイナビ",
      memo: "memo",
    });

    const { result } = await callAndCapture(fd);
    expect(result?.error).toBeTruthy();

    const calls = callSignatures();
    expect(calls).toContain("POST /api/v1/entries/with-company");
    expect(calls).not.toContain("POST /api/v1/companies");
    expect(calls).not.toContain("DELETE /api/v1/companies/co1");
    expect(calls.some((c) => c.startsWith("DELETE /api/v1/inbox/clips/"))).toBe(false);
  });

  it("既存会社IDがある場合は company を新規作成せず entry を作成する", async () => {
    serverFetch.mockImplementation(async (path?: string, init?: RequestInit) => {
      if (path === "/api/v1/companies/co1" && !init?.method) return companyResp;
      if (path === "/api/v1/entries" && init?.method === "POST") return entryResp;
      if (path === "/api/v1/entries/en1/selection-flow" && init?.method === "PUT") return flowResp;
      return undefined;
    });

    const fd = form({
      clipId: "clip1",
      existingCompanyId: "co1",
      companyName: "テスト商事",
      route: "本選考",
      source: "マイナビ",
      memo: "memo",
    });

    const { thrown } = await callAndCapture(fd);
    expect((thrown as { digest?: string } | undefined)?.digest).toContain("/entry/en1");

    const calls = callSignatures();
    expect(calls).toContain("GET /api/v1/companies/co1");
    expect(calls).toContain("POST /api/v1/entries");
    expect(calls).toContain("PUT /api/v1/entries/en1/selection-flow");
    expect(calls).not.toContain("POST /api/v1/companies");
  });

  it("会社名が空ならエラーを返し、API を呼ばない", async () => {
    const fd = form({ clipId: "clip1", companyName: "   ", route: "本選考", source: "マイナビ", memo: "" });
    const { result } = await callAndCapture(fd);
    expect(result?.error).toContain("会社名");
    expect(serverFetch).not.toHaveBeenCalled();
  });
});

describe("deleteInboxClipAction", () => {
  beforeEach(() => serverFetch.mockReset());

  it("clip を DELETE して空 state を返す", async () => {
    serverFetch.mockResolvedValue(undefined);
    const result = await deleteInboxClipAction({}, form({ clipId: "clip1" }));
    expect(result.error).toBeUndefined();
    expect(callSignatures()).toContain("DELETE /api/v1/inbox/clips/clip1");
  });

  it("clipId が空なら API を呼ばずエラーを返す", async () => {
    const result = await deleteInboxClipAction({}, form({ clipId: "  " }));
    expect(result.error).toBeTruthy();
    expect(serverFetch).not.toHaveBeenCalled();
  });

  it("DELETE 失敗時はエラーメッセージを返す", async () => {
    serverFetch.mockImplementationOnce(async () => {
      throw new ApiError(500, "boom");
    });
    serverFetch.mockResolvedValue(undefined);
    const result = await deleteInboxClipAction({}, form({ clipId: "clip1" }));
    expect(result.error).toBeTruthy();
    // vitest 4.1.5 が spy の reject 結果を unhandled 誤検知するため resolve で締める。
    await serverFetch("/__settle__");
  });
});
