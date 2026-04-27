import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { server } from "@/test/msw-server";
import { EntryDetailView } from "./EntryDetailView";

const API = "http://localhost:8080";
const sample = (overrides = {}) => ({
  id: "e1",
  companyId: "c1",
  route: "本選考",
  source: "リクナビ",
  status: "in_progress",
  stageKind: "interview",
  stageLabel: "一次面接",
  memo: "テストメモ",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

describe("EntryDetailView", () => {
  it("API から取得した詳細を表示する", async () => {
    server.use(http.get(`${API}/api/v1/entries/e1`, () => HttpResponse.json(sample())));
    server.use(http.get(`${API}/api/v1/entries/e1/tasks`, () => HttpResponse.json({ tasks: [] })));

    render(<EntryDetailView entryId="e1" />);
    await waitFor(() => expect(screen.getByText("一次面接")).toBeInTheDocument());
    expect(screen.getByText("テストメモ")).toBeInTheDocument();
  });

  it("「進める →」クリックで PATCH が走り stageKind が更新される", async () => {
    let patchedBody: Record<string, unknown> | null = null;
    server.use(
      http.get(`${API}/api/v1/entries/e1`, () => HttpResponse.json(sample())),
      http.get(`${API}/api/v1/entries/e1/tasks`, () => HttpResponse.json({ tasks: [] })),
      http.patch(`${API}/api/v1/entries/e1`, async ({ request }) => {
        patchedBody = (await request.json()) as Record<string, unknown>;
        return HttpResponse.json(sample({ stageKind: "offer", stageLabel: "内定" }));
      }),
    );

    render(<EntryDetailView entryId="e1" />);
    await waitFor(() => expect(screen.getByText("一次面接")).toBeInTheDocument());

    const advance = screen.getByRole("button", { name: /進める/ });
    await userEvent.click(advance);

    await waitFor(() => expect(patchedBody).not.toBeNull());
    expect(patchedBody).toMatchObject({ stageKind: expect.any(String) });
  });

  it("内定到達時はスタンプを表示する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries/e1`, () => HttpResponse.json(sample({ stageKind: "offer", stageLabel: "内定" }))),
      http.get(`${API}/api/v1/entries/e1/tasks`, () => HttpResponse.json({ tasks: [] })),
    );
    render(<EntryDetailView entryId="e1" />);
    await waitFor(() => expect(screen.getByText("内定！")).toBeInTheDocument());
  });

  it("読み込み失敗時は alert を表示する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries/e1`, () => HttpResponse.json({ message: "boom" }, { status: 500 })),
      http.get(`${API}/api/v1/entries/e1/tasks`, () => HttpResponse.json({ tasks: [] })),
    );
    render(<EntryDetailView entryId="e1" />);
    await waitFor(() => expect(screen.getByRole("alert")).toBeInTheDocument());
  });
});
