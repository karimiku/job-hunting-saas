import { beforeEach, describe, expect, it, vi } from "vitest";

const { getSupabaseBrowserAccessToken } = vi.hoisted(() => ({
  getSupabaseBrowserAccessToken: vi.fn(),
}));

vi.mock("../supabase/client", () => ({ getSupabaseBrowserAccessToken }));

import { apiFetch } from "./client";

describe("apiFetch", () => {
  beforeEach(() => {
    getSupabaseBrowserAccessToken.mockReset();
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        Response.json({
          ok: true,
        }),
      ),
    );
  });

  it("Supabase access token を Authorization header に付ける", async () => {
    getSupabaseBrowserAccessToken.mockResolvedValue("browser-token");

    await apiFetch("/api/v1/entries");

    const [, init] = vi.mocked(fetch).mock.calls[0];
    const headers = new Headers(init?.headers);
    expect(headers.get("Authorization")).toBe("Bearer browser-token");
    expect(headers.get("Content-Type")).toBe("application/json");
  });

  it("呼び出し側が Authorization を指定した場合は上書きしない", async () => {
    getSupabaseBrowserAccessToken.mockResolvedValue("browser-token");

    await apiFetch("/api/v1/entries", {
      headers: {
        Authorization: "Bearer explicit-token",
      },
    });

    const [, init] = vi.mocked(fetch).mock.calls[0];
    const headers = new Headers(init?.headers);
    expect(headers.get("Authorization")).toBe("Bearer explicit-token");
  });
});
