import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import {
  buildDashboardEntries,
  DashboardEntries,
} from "./DashboardEntries";
import type { EntryResponse } from "@/lib/api/entries";
import type { TaskWithEntry } from "@/lib/api/server-resources";

const entry = (overrides: Partial<EntryResponse> = {}): EntryResponse => ({
  id: "e1",
  companyId: "c1",
  companyName: "○○商事",
  route: "本選考",
  source: "リクナビ",
  status: "in_progress",
  stageKind: "application",
  stageLabel: "エントリー",
  memo: "",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

const task = (overrides: Partial<TaskWithEntry> = {}): TaskWithEntry => ({
  id: "t1",
  entryId: "e1",
  title: "ES提出",
  type: "deadline",
  status: "todo",
  dueDate: "2026-06-10T00:00:00Z",
  memo: "",
  companyName: "○○商事",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

describe("buildDashboardEntries", () => {
  it("未完了タスクの近いEntryを優先し、タスクがなければステージが進んだEntryを優先する", () => {
    const result = buildDashboardEntries(
      [
        entry({ id: "late", companyName: "遅い会社", stageKind: "offer", stageLabel: "内定" }),
        entry({ id: "soon", companyName: "近い会社", stageKind: "document", stageLabel: "書類" }),
        entry({ id: "none", companyName: "タスクなし", stageKind: "interview", stageLabel: "面接" }),
      ],
      [
        task({ id: "t-late", entryId: "late", dueDate: "2026-07-01T00:00:00Z" }),
        task({ id: "t-soon", entryId: "soon", dueDate: "2026-06-01T00:00:00Z" }),
      ],
    );

    expect(result.map((item) => item.id)).toEqual(["soon", "late", "none"]);
    expect(result[0]).toMatchObject({
      company: "近い会社",
      openTaskCount: 1,
      nearestDue: "6/1",
    });
  });
});

describe("DashboardEntries", () => {
  it("ホームにEntryと未完了タスク数を表示する", () => {
    render(
      <DashboardEntries
        entries={[entry({ companyName: "テスト商事", stageKind: "interview", stageLabel: "面接" })]}
        tasks={[task({ entryId: "e1" })]}
      />,
    );

    expect(screen.getByText("進行中のEntry")).toBeInTheDocument();
    expect(screen.getByText("テスト商事")).toBeInTheDocument();
    expect(screen.getByText("面接")).toBeInTheDocument();
    expect(screen.getByText("未完了 1")).toBeInTheDocument();
  });
});
