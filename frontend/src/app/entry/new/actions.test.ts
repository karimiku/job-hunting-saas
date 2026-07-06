import { beforeEach, describe, expect, it, vi } from "vitest";

const { redirect, revalidatePath, serverFetch } = vi.hoisted(() => ({
  redirect: vi.fn((path: string) => {
    const err = new Error("NEXT_REDIRECT") as Error & { digest?: string };
    err.digest = `NEXT_REDIRECT;replace;${path}`;
    throw err;
  }),
  revalidatePath: vi.fn(),
  serverFetch: vi.fn(),
}));

vi.mock("@/lib/api/server", () => ({ serverFetch }));
vi.mock("next/cache", () => ({ revalidatePath }));
vi.mock("next/navigation", () => ({ redirect }));

import { createNewEntryAction, type NewEntryFormState } from "./actions";

const INITIAL: NewEntryFormState = {};

function form(fields: Record<string, string>): FormData {
  const fd = new FormData();
  for (const [key, value] of Object.entries(fields)) fd.set(key, value);
  return fd;
}

function postedJson(path: string): Record<string, unknown> {
  const call = serverFetch.mock.calls.find(([calledPath]) => calledPath === path);
  const init = call?.[1] as RequestInit | undefined;
  return JSON.parse(String(init?.body));
}

describe("createNewEntryAction", () => {
  beforeEach(() => {
    redirect.mockClear();
    revalidatePath.mockClear();
    serverFetch.mockReset();
  });

  it("候補外の応募経路も trim してそのまま作成 API に渡す", async () => {
    serverFetch.mockImplementation(async (path: string, init?: RequestInit) => {
      if (path === "/api/v1/entries/with-company" && init?.method === "POST") {
        return {
          id: "entry1",
          companyId: "company1",
          route: "説明会経由",
          source: "リクナビ",
          status: "in_progress",
          stageKind: "application",
          stageLabel: "応募",
          memo: "",
          createdAt: "",
          updatedAt: "",
        };
      }
      if (path === "/api/v1/entries/entry1/selection-flow" && init?.method === "PUT") {
        return {
          id: "flow1",
          entryId: "entry1",
          source: "template",
          currentStagePosition: 1,
          stages: [],
          createdAt: "",
          updatedAt: "",
        };
      }
      return undefined;
    });

    await expect(
      createNewEntryAction(
        INITIAL,
        form({
          companyName: "テスト商事",
          route: "  説明会経由  ",
          source: "リクナビ",
          memo: "",
          flowMode: "template",
          customFlowText: "",
        }),
      ),
    ).rejects.toMatchObject({ digest: expect.stringContaining("/entry/entry1") });

    expect(postedJson("/api/v1/entries/with-company")).toMatchObject({
      companyName: "テスト商事",
      route: "説明会経由",
      source: "リクナビ",
    });
    expect(revalidatePath).toHaveBeenCalledWith("/entry");
    expect(redirect).toHaveBeenCalledWith("/entry/entry1");
  });
});
