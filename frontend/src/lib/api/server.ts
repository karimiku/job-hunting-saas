// Server Component / Server Action 用の fetch ラッパー。
// Supabase session があれば access token を Authorization: Bearer として backend に渡す。
//
// なぜ client.ts と分けるか:
//   - "use client" 配下のコードに next/headers (server-only) を import すると build 失敗する。
//   - Server Component では Supabase SSR client が request cookie から session を読む必要がある。

import { cookies, headers } from "next/headers";
import { ApiError } from "./client-types";
import { getSupabaseServerAccessToken } from "../supabase/server";

const DEFAULT_BACKEND_ALLOWED_HOSTS = [
  "localhost",
  "127.0.0.1",
  "entre.kamiriku.com",
  "*.vercel.app",
  "job-hunting-saas.vercel.app",
  "job-hunting-saas-kamirikus-projects.vercel.app",
  "api.entre.kamiriku.com",
  "entre-backend-gfsd4pzoxq-an.a.run.app",
];

/**
 * Server Component から backend を叩く。
 * Supabase access token を Authorization header に付ける。
 * dev auth / legacy Firebase rollback 用に Cookie ヘッダ転送も残す。
 *
 * `cache: "no-store"` を default にして、ユーザー固有データが要求間で漏れないようにする。
 * 公開キャッシュが欲しい場合は呼び出し側で `next: { revalidate }` 等を渡す。
 */
export async function serverFetch<T>(
  path: string,
  init: RequestInit = {},
): Promise<T> {
  const startedAt = Date.now();
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
  if (!outgoingHeaders.has("Authorization")) {
    const accessToken = await getSupabaseServerAccessToken();
    if (accessToken) {
      outgoingHeaders.set("Authorization", `Bearer ${accessToken}`);
    }
  }
  for (const header of ["origin", "referer"]) {
    const value = incomingHeaders.get(header);
    if (value && !outgoingHeaders.has(header)) {
      outgoingHeaders.set(header, value);
    }
  }
  const apiBase = serverBackendBaseURL();

  const res = await fetch(`${apiBase}${path}`, {
    cache: "no-store",
    ...init,
    headers: outgoingHeaders,
  });

  logServerFetchTiming(path, init.method ?? "GET", res, Date.now() - startedAt);

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

function logServerFetchTiming(path: string, method: string, res: Response, durationMs: number) {
  if (process.env.ENTRE_SERVER_FETCH_LOG !== "1") return;

  const serverTiming = res.headers.get("server-timing");
  const timingSuffix = serverTiming ? ` server-timing="${serverTiming}"` : "";
  console.info(
    `[serverFetch] ${method.toUpperCase()} ${path} -> ${res.status} ${durationMs}ms${timingSuffix}`,
  );
}

function serverBackendBaseURL(): string {
  const explicitURL = firstEnvValue(["BACKEND_API_BASE_URL", "NEXT_PUBLIC_API_BASE_URL"]);
  const raw = explicitURL ?? vercelServiceBackendBaseURL() ?? "http://localhost:8080";
  const parsed = new URL(raw);
  if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
    throw new Error("BACKEND_API_BASE_URL must use http or https");
  }
  if (parsed.username || parsed.password || parsed.search || parsed.hash) {
    throw new Error("BACKEND_API_BASE_URL must not include credentials, query, or hash");
  }
  assertBackendHostAllowed(parsed.hostname);
  return `${parsed.origin}${parsed.pathname.replace(/\/$/, "")}`;
}

// 資格情報 (cookie / Bearer) を転送する宛先ホストは、攻撃者が制御できる
// リクエストヘッダ (x-forwarded-host / host) からは絶対に決定しない。
// VERCEL_URL はプラットフォームが実行環境に設定する値でクライアント入力ではないため
// 信頼できる。未設定なら Vercel service 経路は使わず、呼び出し元で localhost にフォールバックする。
function vercelServiceBackendBaseURL(): string | undefined {
  if (process.env.VERCEL !== "1") return undefined;

  const deploymentHost = process.env.VERCEL_URL?.trim();
  if (!deploymentHost) return undefined;

  return `https://${deploymentHost}/backend`;
}

function firstEnvValue(envNames: string[]): string | undefined {
  for (const envName of envNames) {
    const value = process.env[envName]?.trim();
    if (value) return value;
  }
  return undefined;
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
