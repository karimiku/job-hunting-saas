import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { server } from "@/test/msw-server";
import { listInboxClips, createInboxClip, deleteInboxClip } from "./inboxClips";

const API = "http://localhost:8080";

describe("inboxClips API", () => {
  it("listInboxClips は配列を返す", async () => {
    server.use(
      http.get(`${API}/api/v1/inbox/clips`, () =>
        HttpResponse.json({
          clips: [
            { id: "c1", url: "https://example.com/jobs/1", title: "Title", source: "マイナビ", guess: "○○商事", capturedAt: "2026-04-26T00:00:00Z" },
          ],
        }),
      ),
    );
    const result = await listInboxClips();
    expect(result).toHaveLength(1);
    expect(result[0].source).toBe("マイナビ");
  });

  it("createInboxClip は POST してレスポンスを返す", async () => {
    let body: Record<string, unknown> | null = null;
    server.use(
      http.post(`${API}/api/v1/inbox/clips`, async ({ request }) => {
        body = (await request.json()) as Record<string, unknown>;
        return HttpResponse.json(
          { id: "c-new", url: body.url, title: body.title, source: body.source, guess: body.guess ?? "", capturedAt: "2026-04-26T00:00:00Z" },
          { status: 201 },
        );
      }),
    );
    const result = await createInboxClip({
      url: "https://example.com/jobs/1",
      title: "Title",
      source: "マイナビ",
      guess: "○○商事",
    });
    expect(body).toMatchObject({ url: "https://example.com/jobs/1", source: "マイナビ" });
    expect(result.id).toBe("c-new");
  });

  it("deleteInboxClip は DELETE する", async () => {
    let called = false;
    server.use(
      http.delete(`${API}/api/v1/inbox/clips/c1`, () => {
        called = true;
        return new HttpResponse(null, { status: 204 });
      }),
    );
    await deleteInboxClip("c1");
    expect(called).toBe(true);
  });
});
