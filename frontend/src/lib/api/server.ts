// Server Component / Server Action 用の fetch ラッパー。
// クライアントの session cookie を Next.js Server から backend に転送する。
//
// なぜ client.ts と分けるか:
//   - "use client" 配下のコードに next/headers (server-only) を import すると build 失敗する。
//   - Server Component で credentials:"include" を使ってもブラウザの cookie は当たらない。
//     Next.js Server から backend へは fetch の Cookie ヘッダで明示的に渡す必要がある。

import { cookies, headers } from "next/headers";
import { ApiError } from "./client";

const DEFAULT_BACKEND_ALLOWED_HOSTS = ["localhost", "127.0.0.1", "api.entre.kamiriku.com"];
const API_BASE = serverBackendOrigin();

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
  const incomingHeaders = await headers();
  const cookieHeader = cookieStore.toString();
  const outgoingHeaders = new Headers(init.headers);
  if (!outgoingHeaders.has("Content-Type")) {
    outgoingHeaders.set("Content-Type", "application/json");
  }
  if (cookieHeader) {
    outgoingHeaders.set("cookie", cookieHeader);
  }
  for (const header of ["origin", "referer"]) {
    const value = incomingHeaders.get(header);
    if (value && !outgoingHeaders.has(header)) {
      outgoingHeaders.set(header, value);
    }
  }

  const res = await fetch(`${API_BASE}${path}`, {
    cache: "no-store",
    ...init,
    headers: outgoingHeaders,
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

function serverBackendOrigin(): string {
  const raw = process.env.BACKEND_API_BASE_URL ?? "http://localhost:8080";
  const parsed = new URL(raw);
  if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
    throw new Error("BACKEND_API_BASE_URL must use http or https");
  }
  if (parsed.username || parsed.password || parsed.pathname !== "/" || parsed.search || parsed.hash) {
    throw new Error("BACKEND_API_BASE_URL must be an origin without credentials, path, query, or hash");
  }
  assertBackendHostAllowed(parsed.hostname);
  return parsed.origin;
}

function assertBackendHostAllowed(hostname: string) {
  const allowedHosts = (
    process.env.BACKEND_API_ALLOWED_HOSTS ?? DEFAULT_BACKEND_ALLOWED_HOSTS.join(",")
  )
    .split(",")
    .map((host) => host.trim().toLowerCase())
    .filter((host) => host.length > 0);

  const normalized = hostname.toLowerCase();
  if (allowedHosts.some((allowedHost) => hostMatches(normalized, allowedHost))) return;
  throw new Error(`BACKEND_API_BASE_URL host ${hostname} is not in BACKEND_API_ALLOWED_HOSTS`);
}

function hostMatches(hostname: string, allowedHost: string): boolean {
  if (allowedHost.startsWith("*.")) {
    const suffix = allowedHost.slice(1);
    return hostname.endsWith(suffix) && hostname.length > suffix.length;
  }
  return hostname === allowedHost;
}
