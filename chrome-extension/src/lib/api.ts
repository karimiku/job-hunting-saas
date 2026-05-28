// Entré バックエンド API クライアント。
// Chrome 拡張からは Session Cookie を共有できるよう host_permissions と credentials: include を使う。

const API_BASE =
  (import.meta.env.VITE_API_BASE_URL as string | undefined) ??
  "http://localhost:8080";

export interface InboxClipResponse {
  id: string;
  url: string;
  title: string;
  source: string;
  guess: string;
  capturedAt: string;
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(init.headers ?? {}),
    },
  });
  if (!res.ok) {
    throw new Error(`API ${path} failed: ${res.status}`);
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
