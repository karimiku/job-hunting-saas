import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const { exchangeCodeForSession, createSupabaseServerClient } = vi.hoisted(() => ({
  exchangeCodeForSession: vi.fn(),
  createSupabaseServerClient: vi.fn(),
}));

vi.mock("@/lib/supabase/server", () => ({ createSupabaseServerClient }));

import { GET } from "./route";

describe("GET /auth/callback (V5: open redirect)", () => {
  beforeEach(() => {
    exchangeCodeForSession.mockReset();
    createSupabaseServerClient.mockReset();
    createSupabaseServerClient.mockResolvedValue({
      auth: { exchangeCodeForSession },
    });
    delete process.env.NEXT_PUBLIC_SITE_URL;
  });

  afterEach(() => {
    delete process.env.NEXT_PUBLIC_SITE_URL;
  });

  it("x-forwarded-host ヘッダを無視し、request.url 由来の origin へリダイレクトする", async () => {
    exchangeCodeForSession.mockResolvedValue({ error: null });

    const request = new Request("https://entre.kamiriku.com/auth/callback?code=abc", {
      headers: {
        "x-forwarded-host": "evil.attacker.example",
        host: "evil.attacker.example",
      },
    });

    const res = await GET(request);

    expect(res.status).toBeGreaterThanOrEqual(300);
    expect(res.status).toBeLessThan(400);
    const location = res.headers.get("location");
    expect(location).toBe("https://entre.kamiriku.com/dashboard");
    expect(location).not.toContain("evil.attacker.example");
  });

  it("NEXT_PUBLIC_SITE_URL が設定されていれば、それを redirect origin として使う", async () => {
    process.env.NEXT_PUBLIC_SITE_URL = "https://entre.kamiriku.com";
    exchangeCodeForSession.mockResolvedValue({ error: null });

    const request = new Request("http://internal-service:8080/auth/callback?code=abc", {
      headers: {
        "x-forwarded-host": "evil.attacker.example",
      },
    });

    const res = await GET(request);

    const location = res.headers.get("location");
    expect(location).toBe("https://entre.kamiriku.com/dashboard");
  });

  it("認証エラー時も攻撃者ホストへリダイレクトしない", async () => {
    exchangeCodeForSession.mockResolvedValue({ error: new Error("invalid_code") });

    const request = new Request("https://entre.kamiriku.com/auth/callback?code=abc", {
      headers: {
        "x-forwarded-host": "evil.attacker.example",
      },
    });

    const res = await GET(request);

    const location = res.headers.get("location");
    expect(location).toBe("https://entre.kamiriku.com/login?error=auth_callback");
  });
});
