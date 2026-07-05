// Client Component 専用の fetch ラッパー。
// Supabase session があれば access token を Authorization: Bearer として backend に渡す。
// HTTP エラーは ApiError に統一して投げる（呼び出し側で .unauthorized / .notFound で分岐できる）。
// Server Component からは ./server.ts の serverFetch を使うこと。

import { ApiError } from "./client-types";
import { getSupabaseBrowserAccessToken } from "../supabase/client";
export { ApiError } from "./client-types";

// 同一 origin の /backend 経由で backend へ届く (next.config.ts の rewrite が proxy する)。
// cross-origin で直接叩くと CORS preflight が毎回1往復乗るため相対パスを default にする。
// テスト (jsdom は相対 URL の fetch 不可) は NEXT_PUBLIC_CLIENT_API_BASE で絶対 URL に上書きする。
const API_BASE = process.env.NEXT_PUBLIC_CLIENT_API_BASE ?? "/backend";

export async function apiFetch<T>(
  path: string,
  init: RequestInit = {},
): Promise<T> {
  const headers = new Headers(init.headers);
  if (!headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (!headers.has("Authorization")) {
    const accessToken = await getSupabaseBrowserAccessToken();
    if (accessToken) {
      headers.set("Authorization", `Bearer ${accessToken}`);
    }
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    credentials: "include",
    headers,
  });
  if (!res.ok) {
    let message = `HTTP ${res.status}`;
    try {
      const body = (await res.json()) as { message?: string };
      if (body?.message) message = body.message;
    } catch {
      // ignore JSON parse errors — fall back to default message
    }
    throw new ApiError(res.status, message);
  }
  if (res.status === 204) {
    return undefined as T;
  }
  return res.json() as Promise<T>;
}
