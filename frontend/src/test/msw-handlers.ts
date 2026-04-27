import { http, HttpResponse } from "msw";

const API_BASE = "http://localhost:8080";

// 各テストで `server.use(...)` で上書きするためのデフォルト handler。
// 認証・データ系の "happy path" のみ。エラー系はテスト側で書く。
export const handlers = [
  // Companies
  http.get(`${API_BASE}/api/v1/companies`, () =>
    HttpResponse.json({ companies: [] }),
  ),
  http.post(`${API_BASE}/api/v1/companies`, async ({ request }) => {
    const body = (await request.json()) as { name: string; memo?: string };
    return HttpResponse.json(
      {
        id: "c-mock",
        name: body.name,
        memo: body.memo ?? "",
        createdAt: "2026-04-26T00:00:00Z",
        updatedAt: "2026-04-26T00:00:00Z",
      },
      { status: 201 },
    );
  }),

  // Entries
  http.get(`${API_BASE}/api/v1/entries`, () =>
    HttpResponse.json({ entries: [] }),
  ),
  http.post(`${API_BASE}/api/v1/entries`, async ({ request }) => {
    const body = (await request.json()) as Record<string, string>;
    return HttpResponse.json(
      {
        id: "e-mock",
        companyId: body.companyId,
        route: body.route,
        source: body.source,
        status: "in_progress",
        stageKind: "application",
        stageLabel: "応募",
        memo: body.memo ?? "",
        createdAt: "2026-04-26T00:00:00Z",
        updatedAt: "2026-04-26T00:00:00Z",
      },
      { status: 201 },
    );
  }),
];
