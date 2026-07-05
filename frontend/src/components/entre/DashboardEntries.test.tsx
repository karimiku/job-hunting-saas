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
  it("ホームに応募先と未完了タスク数を表示する", () => {
    render(
      <DashboardEntries
        entries={[entry({ companyName: "テスト商事", stageKind: "interview", stageLabel: "面接" })]}
        tasks={[task({ entryId: "e1" })]}
      />,
    );

    expect(screen.getByText("進行中の応募先")).toBeInTheDocument();
    expect(screen.getByText("テスト商事")).toBeInTheDocument();
    expect(screen.getByText("面接")).toBeInTheDocument();
    expect(screen.getByText("未完了 1件")).toBeInTheDocument();
  });

  it("タスク期日にラベルを付けて表示する", () => {
    render(
      <DashboardEntries
        entries={[entry({ companyName: "テスト商事", stageKind: "interview", stageLabel: "面接" })]}
        tasks={[task({ entryId: "e1", dueDate: "2026-07-03T00:00:00Z" })]}
      />,
    );

    expect(screen.getByText("締切 7/3")).toBeInTheDocument();
  });

  it("未完了タスクがなければ期日ラベルなしで「期日なし」を表示する", () => {
    render(
      <DashboardEntries
        entries={[entry({ companyName: "タスクなし", stageKind: "interview", stageLabel: "面接" })]}
        tasks={[]}
      />,
    );

    expect(screen.getByText("期日なし")).toBeInTheDocument();
  });

  it("進捗バーの近くにステージ名と「Nステップ中M」を表示する", () => {
    render(
      <DashboardEntries
        entries={[entry({ companyName: "テスト商事", stageKind: "interview", stageLabel: "面接" })]}
        tasks={[task({ entryId: "e1" })]}
      />,
    );

    expect(screen.getByText("面接")).toBeInTheDocument();
    expect(screen.getByText("6ステップ中4")).toBeInTheDocument();
  });

  it("選考中は自明なため status を表示せず、確定した結果のみ表示する", () => {
    render(
      <DashboardEntries
        entries={[
          entry({ companyName: "選考中の会社", stageKind: "interview", stageLabel: "面接" }),
        ]}
        tasks={[]}
      />,
    );

    expect(screen.queryByText("選考中")).not.toBeInTheDocument();
  });

  it("確定した結果（内定獲得等）は status ラベルを表示する", () => {
    render(
      <DashboardEntries
        entries={[
          entry({
            companyName: "内定の会社",
            stageKind: "offer",
            stageLabel: "内定",
            status: "offered",
          }),
        ]}
        tasks={[]}
      />,
    );

    expect(screen.getByText("内定獲得")).toBeInTheDocument();
  });

  it("未知の status 値は素で出さない", () => {
    render(
      <DashboardEntries
        entries={[
          entry({
            companyName: "不明ステータスの会社",
            stageKind: "interview",
            stageLabel: "面接",
            status: "active",
          }),
        ]}
        tasks={[]}
      />,
    );

    expect(screen.queryByText("active")).not.toBeInTheDocument();
  });

  it("応募先が無ければ登録を促す空状態を表示する", () => {
    render(<DashboardEntries entries={[]} tasks={[]} />);

    expect(screen.getByText("応募先はまだ登録されていません")).toBeInTheDocument();
    const link = screen.getByRole("link", { name: "応募先を追加" });
    expect(link).toHaveAttribute("href", "/entry/new");
  });
});
