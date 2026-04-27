// Server Component / Server Action 用の fetch ラッパー。
// クライアントの session cookie を Next.js Server から backend に転送する。
//
// なぜ client.ts と分けるか:
//   - "use client" 配下のコードに next/headers (server-only) を import すると build 失敗する。
//   - Server Component で credentials:"include" を使ってもブラウザの cookie は当たらない。
//     Next.js Server から backend へは fetch の Cookie ヘッダで明示的に渡す必要がある。

import { cookies } from "next/headers";
import { ApiError } from "./client";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

/**
 * Server Component から backend を叩く。
 * リクエストの Cookie ヘッダを backend に転送し、`credentials:"include"` 相当を再現する。
 *
 * `cache: "no-store"` を default にして、ユーザー固有データが要求間で漏れないようにする。
 * 公開キャッシュが欲しい場合は呼び出し側で `next: { revalidate }` 等を渡す。
 */
export async function serverFetch<T>(
  path: string,
  init: RequestInit = {},
): Promise<T> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();

  const res = await fetch(`${API_BASE}${path}`, {
    cache: "no-store",
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(cookieHeader ? { cookie: cookieHeader } : {}),
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
