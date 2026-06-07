import { describe, expect, it } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { DashboardStats, summarizeEntries } from "./DashboardStats";
import type { EntryResponse } from "@/lib/api/entries";

const e = (overrides: Partial<EntryResponse> = {}): EntryResponse => ({
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

describe("summarizeEntries", () => {
  it("件数 / 選考中 / 面接中 / 内定 を集計する", () => {
    const result = summarizeEntries([
      e({ status: "in_progress", stageKind: "interview" }),
      e({ status: "in_progress", stageKind: "interview" }),
      e({ status: "in_progress", stageKind: "document" }),
      e({ status: "offered", stageKind: "offer" }),
      e({ status: "rejected", stageKind: "document" }),
    ]);
    expect(result.total).toBe(5);
    expect(result.inProgress).toBe(3);
    expect(result.interviewing).toBe(2);
    expect(result.offered).toBe(1);
  });
});

describe("DashboardStats", () => {
  it("CountUp が完了すると集計値を表示する", async () => {
    render(
      <DashboardStats
        entries={[
          e({ status: "in_progress", stageKind: "interview" }),
          e({ status: "in_progress", stageKind: "interview" }),
          e({ status: "in_progress", stageKind: "document" }),
          e({ status: "offered", stageKind: "offer" }),
          e({ status: "rejected", stageKind: "document" }),
        ]}
      />,
    );
    // CountUp はイージング付きで非同期に値を更新するので個別に waitFor で待つ
    await waitFor(() => expect(screen.getByTestId("stat-total")).toHaveTextContent("5"));
    await waitFor(() => expect(screen.getByTestId("stat-in-progress")).toHaveTextContent("3"));
    await waitFor(() => expect(screen.getByTestId("stat-offered")).toHaveTextContent("1"));
    await waitFor(() => expect(screen.getByTestId("stat-interviewing")).toHaveTextContent("2"));
  });
});
