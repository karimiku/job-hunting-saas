// Entré バックエンド API クライアント。
// Chrome 拡張からは Session Cookie を共有できるよう host_permissions と credentials: include を使う。

const API_BASE =
  (import.meta.env.VITE_API_BASE_URL as string | undefined) ??
  "http://localhost:8080";

// Web アプリ (ログイン導線) の URL。401 のときにログインページを開くのに使う。
export const WEB_BASE =
  (import.meta.env.VITE_WEB_BASE_URL as string | undefined) ??
  "http://localhost:3000";

export interface InboxClipResponse {
  id: string;
  url: string;
  title: string;
  source: string;
  guess: string;
  capturedAt: string;
}

// エラーの種類。popup 側でユーザー向けメッセージ/回復導線に変換する。
//   unauthorized: 401 — Web 未ログイン
//   forbidden   : 403 — 権限/拡張許可/再ログインが必要
//   server      : 5xx — サーバー側エラー（再試行）
//   network     : fetch 自体が失敗（サーバー停止・CORS 拒否・オフライン等。CORS と通信断は区別不能）
//   client      : その他 4xx
export type ApiErrorKind =
  | "unauthorized"
  | "forbidden"
  | "server"
  | "network"
  | "client";

export class ApiRequestError extends Error {
  readonly kind: ApiErrorKind;
  readonly status?: number;
  constructor(kind: ApiErrorKind, status?: number, message?: string) {
    super(message ?? `API request failed: ${kind}${status ? ` (${status})` : ""}`);
    this.name = "ApiRequestError";
    this.kind = kind;
    this.status = status;
  }
}

export function errorKindFromStatus(status: number): ApiErrorKind {
  if (status === 401) return "unauthorized";
  if (status === 403) return "forbidden";
  if (status >= 500) return "server";
  return "client";
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  let res: Response;
  try {
    res = await fetch(`${API_BASE}${path}`, {
      ...init,
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...(init.headers ?? {}),
      },
    });
  } catch {
    // fetch が reject する＝レスポンスに到達できていない（サーバー停止 / CORS 拒否 / オフライン）。
    throw new ApiRequestError("network");
  }

  if (!res.ok) {
    throw new ApiRequestError(errorKindFromStatus(res.status), res.status);
  }
  return res.json() as Promise<T>;
}

export interface CreateInboxClipInput {
  url: string;
  title: string;
  source: string;
  guess?: string;
}

export async function createInboxClip(input: CreateInboxClipInput): Promise<InboxClipResponse> {
  return request<InboxClipResponse>("/api/v1/inbox/clips", {
    method: "POST",
    body: JSON.stringify(input),
  });
}
