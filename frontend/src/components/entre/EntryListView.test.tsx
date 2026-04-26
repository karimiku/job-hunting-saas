import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { server } from "@/test/msw-server";
import { EntryListView } from "./EntryListView";

const API = "http://localhost:8080";

describe("EntryListView", () => {
  it("初期表示は読み込み中", () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({ entries: [] }),
      ),
    );
    render(<EntryListView />);
    expect(screen.getByRole("status")).toHaveTextContent("読み込み中");
  });

  it("API のエントリーを表示する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({
          entries: [
            { id: "e1", companyId: "c1", route: "本選考", source: "リクナビ", status: "in_progress", stageKind: "interview", stageLabel: "一次面接", memo: "", createdAt: "x", updatedAt: "x" },
            { id: "e2", companyId: "c2", route: "本選考", source: "マイナビ", status: "in_progress", stageKind: "document", stageLabel: "ES提出", memo: "", createdAt: "x", updatedAt: "x" },
          ],
        }),
      ),
    );

    render(<EntryListView />);

    await waitFor(() => {
      expect(screen.getByText("一次面接")).toBeInTheDocument();
      expect(screen.getByText("ES提出")).toBeInTheDocument();
    });
  });

  it("API 失敗時はエラー表示する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({ message: "boom" }, { status: 500 }),
      ),
    );

    render(<EntryListView />);

    await waitFor(() => {
      expect(screen.getByRole("alert")).toHaveTextContent("読み込みに失敗");
    });
  });

  it("エントリー 0件のとき空状態を表示する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({ entries: [] }),
      ),
    );

    render(<EntryListView />);

    await waitFor(() => {
      expect(screen.getByText(/まだエントリーがありません/)).toBeInTheDocument();
    });
  });
});
