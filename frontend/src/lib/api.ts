// バックエンド (Go) への fetch ラッパー。
// Session Cookie を含めるため credentials: "include" を必ず付ける。
//
// 同一 origin の /backend 経由で backend へ届く (next.config.ts の rewrite が proxy する)。
// cross-origin で直接叩くと CORS preflight が毎回1往復乗るため相対パスを default にする。
// テスト (jsdom は相対 URL の fetch 不可) は NEXT_PUBLIC_CLIENT_API_BASE で絶対 URL に上書きする。

const API_BASE = process.env.NEXT_PUBLIC_CLIENT_API_BASE ?? "/backend";

export async function apiFetch(path: string, init: RequestInit = {}): Promise<Response> {
  return fetch(`${API_BASE}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(init.headers ?? {}),
    },
  });
}
