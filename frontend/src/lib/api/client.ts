// Client Component 専用の fetch ラッパー。Session Cookie を含めるため credentials: include を必ず付ける。
// HTTP エラーは ApiError に統一して投げる（呼び出し側で .unauthorized / .notFound で分岐できる）。
// Server Component からは ./server.ts の serverFetch を使うこと (cookie 転送のため)。

import { ApiError } from "./client-types";
export { ApiError } from "./client-types";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export async function apiFetch<T>(
  path: string,
  init: RequestInit = {},
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(init.headers ?? {}),
    },
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
