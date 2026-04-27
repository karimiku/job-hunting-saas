import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { server } from "@/test/msw-server";
import { DashboardStats } from "./DashboardStats";

const API = "http://localhost:8080";
const e = (overrides: Record<string, unknown>) => ({
  id: String(Math.random()),
  companyId: "c1",
  route: "本選考",
  source: "リクナビ",
  status: "in_progress",
  stageKind: "interview",
  stageLabel: "面接",
  memo: "",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

describe("DashboardStats", () => {
  it("API のエントリーから 件数 / 選考中 / 内定数を集計表示する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({
          entries: [
            e({ status: "in_progress", stageKind: "interview" }),
            e({ status: "in_progress", stageKind: "interview" }),
            e({ status: "in_progress", stageKind: "document" }),
            e({ status: "offered", stageKind: "offer" }),
            e({ status: "rejected", stageKind: "document" }),
          ],
        }),
      ),
    );

    render(<DashboardStats />);
    // CountUp はイージング付きで非同期に値を更新するので個別に waitFor で待つ
    await waitFor(() => expect(screen.getByTestId("stat-total")).toHaveTextContent("5"));
    await waitFor(() => expect(screen.getByTestId("stat-in-progress")).toHaveTextContent("3"));
    await waitFor(() => expect(screen.getByTestId("stat-offered")).toHaveTextContent("1"));
    await waitFor(() => expect(screen.getByTestId("stat-interviewing")).toHaveTextContent("2"));
  });

  it("読み込み中は status 役割で読み込み中を出す", () => {
    server.use(http.get(`${API}/api/v1/entries`, () => HttpResponse.json({ entries: [] })));
    render(<DashboardStats />);
    expect(screen.getByRole("status")).toBeInTheDocument();
  });
});
