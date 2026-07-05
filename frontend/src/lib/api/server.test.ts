import { beforeEach, describe, expect, it, vi } from "vitest";

const { cookies, headers, getSupabaseServerAccessToken } = vi.hoisted(() => {
  process.env.BACKEND_API_BASE_URL = "http://localhost:8080";
  return {
    cookies: vi.fn(),
    headers: vi.fn(),
    getSupabaseServerAccessToken: vi.fn(),
  };
});

vi.mock("next/headers", () => ({ cookies, headers }));
vi.mock("../supabase/server", () => ({ getSupabaseServerAccessToken }));

import { serverFetch } from "./server";

describe("serverFetch", () => {
  beforeEach(() => {
    process.env.BACKEND_API_BASE_URL = "http://localhost:8080";
    delete process.env.NEXT_PUBLIC_API_BASE_URL;
    delete process.env.VERCEL;
    delete process.env.VERCEL_URL;
    delete process.env.BACKEND_API_ALLOWED_HOSTS;
    cookies.mockReset();
    headers.mockReset();
    getSupabaseServerAccessToken.mockReset();
    cookies.mockResolvedValue({ toString: () => "session=legacy" });
    headers.mockResolvedValue({
      get: (name: string) => (name === "origin" ? "http://localhost:3000" : null),
    });
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
    getSupabaseServerAccessToken.mockResolvedValue("server-token");

    await serverFetch("/api/v1/entries");

    expect(vi.mocked(fetch).mock.calls[0]?.[0]).toBe("http://localhost:8080/api/v1/entries");
    const [, init] = vi.mocked(fetch).mock.calls[0];
    const headers = new Headers(init?.headers);
    expect(headers.get("Authorization")).toBe("Bearer server-token");
    expect(headers.get("cookie")).toBe("session=legacy");
    expect(headers.get("origin")).toBe("http://localhost:3000");
  });

  it("呼び出し側が Authorization を指定した場合は上書きしない", async () => {
    getSupabaseServerAccessToken.mockResolvedValue("server-token");

    await serverFetch("/api/v1/entries", {
      headers: {
        Authorization: "Bearer explicit-token",
      },
    });

    const [, init] = vi.mocked(fetch).mock.calls[0];
    const headers = new Headers(init?.headers);
    expect(headers.get("Authorization")).toBe("Bearer explicit-token");
  });
});

describe("serverFetch with Vercel service path prefix", () => {
  beforeEach(async () => {
    vi.resetModules();
    process.env.BACKEND_API_BASE_URL = "https://entre.kamiriku.com/backend";
    delete process.env.NEXT_PUBLIC_API_BASE_URL;
    delete process.env.VERCEL;
    delete process.env.VERCEL_URL;
    process.env.BACKEND_API_ALLOWED_HOSTS = "entre.kamiriku.com";
    cookies.mockReset();
    headers.mockReset();
    getSupabaseServerAccessToken.mockReset();
    cookies.mockResolvedValue({ toString: () => "" });
    headers.mockResolvedValue({ get: () => null });
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        Response.json({
          ok: true,
        }),
      ),
    );
  });

  it("BACKEND_API_BASE_URL の path prefix を保持する", async () => {
    const { serverFetch } = await import("./server");

    await serverFetch("/api/v1/entries");

    expect(vi.mocked(fetch).mock.calls[0]?.[0]).toBe(
      "https://entre.kamiriku.com/backend/api/v1/entries",
    );
  });
});

describe("serverFetch on Vercel preview without explicit backend URL", () => {
  beforeEach(async () => {
    vi.resetModules();
    delete process.env.BACKEND_API_BASE_URL;
    delete process.env.NEXT_PUBLIC_API_BASE_URL;
    process.env.VERCEL = "1";
    process.env.VERCEL_URL = "job-hunting-saas-git-supabase-kamirikus-projects.vercel.app";
    process.env.BACKEND_API_ALLOWED_HOSTS = "*.vercel.app";
    cookies.mockReset();
    headers.mockReset();
    getSupabaseServerAccessToken.mockReset();
    cookies.mockResolvedValue({ toString: () => "" });
    headers.mockResolvedValue({
      get: (name: string) => (name === "x-forwarded-proto" ? "https" : null),
    });
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        Response.json({
          ok: true,
        }),
      ),
    );
  });

  it("VERCEL_URL から同一deploymentの /backend を組み立てる", async () => {
    const { serverFetch } = await import("./server");

    await serverFetch("/api/v1/entries");

    expect(vi.mocked(fetch).mock.calls[0]?.[0]).toBe(
      "https://job-hunting-saas-git-supabase-kamirikus-projects.vercel.app/backend/api/v1/entries",
    );
  });
});

describe("serverFetch ignores attacker-controlled host headers (V4)", () => {
  beforeEach(async () => {
    vi.resetModules();
    delete process.env.BACKEND_API_BASE_URL;
    delete process.env.NEXT_PUBLIC_API_BASE_URL;
    process.env.VERCEL = "1";
    // VERCEL_URL is intentionally left unset to simulate the scenario where the
    // only remaining signal would be attacker-controlled request headers.
    delete process.env.VERCEL_URL;
    delete process.env.BACKEND_API_ALLOWED_HOSTS;
    cookies.mockReset();
    headers.mockReset();
    getSupabaseServerAccessToken.mockReset();
    cookies.mockResolvedValue({ toString: () => "session=legacy" });
    headers.mockResolvedValue({
      get: (name: string) => {
        if (name === "x-forwarded-host") return "evil.attacker.example";
        if (name === "host") return "evil.attacker.example";
        if (name === "x-forwarded-proto") return "https";
        return null;
      },
    });
    getSupabaseServerAccessToken.mockResolvedValue("server-token");
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        Response.json({
          ok: true,
        }),
      ),
    );
  });

  it("x-forwarded-host / host ヘッダを無視し、Bearer/cookie を攻撃者ホストへ送らない", async () => {
    const { serverFetch } = await import("./server");

    await serverFetch("/api/v1/entries");

    const [url, init] = vi.mocked(fetch).mock.calls[0];
    expect(url).toBe("http://localhost:8080/api/v1/entries");
    expect(String(url)).not.toContain("evil.attacker.example");
    const sentHeaders = new Headers(init?.headers);
    expect(sentHeaders.get("Authorization")).toBe("Bearer server-token");
    expect(sentHeaders.get("cookie")).toBe("session=legacy");
  });
});
