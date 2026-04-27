import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { server } from "@/test/msw-server";
import {
  listEntries,
  getEntry,
  createEntry,
  updateEntry,
  deleteEntry,
  ApiError,
} from "./entries";

const API = "http://localhost:8080";

describe("entries API", () => {
  describe("listEntries", () => {
    it("ユーザーのエントリー一覧を返す", async () => {
      server.use(
        http.get(`${API}/api/v1/entries`, () =>
          HttpResponse.json({
            entries: [
              {
                id: "e1",
                companyId: "c1",
                route: "本選考",
                source: "リクナビ",
                status: "in_progress",
                stageKind: "interview",
                stageLabel: "一次面接",
                memo: "",
                createdAt: "2026-04-20T00:00:00Z",
                updatedAt: "2026-04-26T00:00:00Z",
              },
            ],
          }),
        ),
      );

      const result = await listEntries();

      expect(result).toHaveLength(1);
      expect(result[0].id).toBe("e1");
      expect(result[0].stageLabel).toBe("一次面接");
    });

    it("filter パラメータを query string で送る", async () => {
      let receivedUrl = "";
      server.use(
        http.get(`${API}/api/v1/entries`, ({ request }) => {
          receivedUrl = request.url;
          return HttpResponse.json({ entries: [] });
        }),
      );

      await listEntries({ status: "in_progress", stageKind: "interview" });

      expect(receivedUrl).toContain("status=in_progress");
      expect(receivedUrl).toContain("stageKind=interview");
    });

    it("401 のとき ApiError を unauthorized=true で投げる", async () => {
      server.use(
        http.get(`${API}/api/v1/entries`, () =>
          HttpResponse.json({ message: "unauthenticated" }, { status: 401 }),
        ),
      );

      await expect(listEntries()).rejects.toMatchObject({
        status: 401,
        unauthorized: true,
      });
    });
  });

  describe("getEntry", () => {
    it("ID から1件取得する", async () => {
      server.use(
        http.get(`${API}/api/v1/entries/e1`, () =>
          HttpResponse.json({
            id: "e1",
            companyId: "c1",
            route: "本選考",
            source: "リクナビ",
            status: "in_progress",
            stageKind: "interview",
            stageLabel: "一次面接",
            memo: "",
            createdAt: "2026-04-20T00:00:00Z",
            updatedAt: "2026-04-26T00:00:00Z",
          }),
        ),
      );

      const result = await getEntry("e1");
      expect(result.id).toBe("e1");
    });

    it("404 のとき notFound=true で投げる", async () => {
      server.use(
        http.get(`${API}/api/v1/entries/missing`, () =>
          HttpResponse.json({ message: "not found" }, { status: 404 }),
        ),
      );

      await expect(getEntry("missing")).rejects.toMatchObject({
        status: 404,
        notFound: true,
      });
    });
  });

  describe("createEntry", () => {
    it("companyId / route / source を POST し、作成された Entry を返す", async () => {
      let received: Record<string, unknown> | null = null;
      server.use(
        http.post(`${API}/api/v1/entries`, async ({ request }) => {
          received = (await request.json()) as Record<string, unknown>;
          return HttpResponse.json(
            {
              id: "new-e",
              companyId: received.companyId,
              route: received.route,
              source: received.source,
              status: "in_progress",
              stageKind: "application",
              stageLabel: "応募",
              memo: received.memo ?? "",
              createdAt: "2026-04-26T00:00:00Z",
              updatedAt: "2026-04-26T00:00:00Z",
            },
            { status: 201 },
          );
        }),
      );

      const result = await createEntry({
        companyId: "c1",
        route: "本選考",
        source: "リクナビ",
        memo: "メモ",
      });

      expect(received).toEqual({
        companyId: "c1",
        route: "本選考",
        source: "リクナビ",
        memo: "メモ",
      });
      expect(result.id).toBe("new-e");
    });
  });

  describe("updateEntry", () => {
    it("PATCH で部分更新する", async () => {
      let receivedBody: Record<string, unknown> | null = null;
      server.use(
        http.patch(`${API}/api/v1/entries/e1`, async ({ request }) => {
          receivedBody = (await request.json()) as Record<string, unknown>;
          return HttpResponse.json({
            id: "e1",
            companyId: "c1",
            route: "本選考",
            source: "マイナビ",
            status: "in_progress",
            stageKind: "interview",
            stageLabel: "一次面接",
            memo: "",
            createdAt: "2026-04-20T00:00:00Z",
            updatedAt: "2026-04-26T01:00:00Z",
          });
        }),
      );

      const result = await updateEntry("e1", { source: "マイナビ" });

      expect(receivedBody).toEqual({ source: "マイナビ" });
      expect(result.source).toBe("マイナビ");
    });
  });

  describe("deleteEntry", () => {
    it("DELETE 後に void を返す", async () => {
      let called = false;
      server.use(
        http.delete(`${API}/api/v1/entries/e1`, () => {
          called = true;
          return new HttpResponse(null, { status: 204 });
        }),
      );

      await deleteEntry("e1");
      expect(called).toBe(true);
    });
  });

  describe("ApiError", () => {
    it("Error を継承し name === 'ApiError'", () => {
      const e = new ApiError(500, "boom");
      expect(e).toBeInstanceOf(Error);
      expect(e.name).toBe("ApiError");
      expect(e.status).toBe(500);
      expect(e.unauthorized).toBe(false);
      expect(e.notFound).toBe(false);
    });
  });
});
