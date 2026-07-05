import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import {
  KanbanBoard,
  kanbanStageUpdateInput,
  nextTaskFor,
  normalizeKanbanStageKind,
} from "./KanbanBoard";
import type { EntryResponse } from "@/lib/api/entries";
import type { TaskWithEntry } from "@/lib/api/server-resources";

const e = (
  kind: string,
  source: string,
  overrides: Partial<EntryResponse> = {},
): EntryResponse => ({
  id: `e-${source}`,
  companyId: `c-${source}`,
  route: "本選考",
  source,
  status: "in_progress",
  stageKind: kind,
  stageLabel: kind,
  memo: "",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

const t = (overrides: Partial<TaskWithEntry> = {}): TaskWithEntry => ({
  id: "t-1",
  entryId: "e-リクナビ",
  title: "ES提出",
  type: "deadline",
  status: "todo",
  dueDate: "2026-07-10",
  memo: "",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

describe("KanbanBoard", () => {
  it("共通ステージ定義の列に initialEntries を振り分けて表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[
          e("application", "リクナビ"),
          e("interview", "マイナビ"),
          e("group", "ONE CAREER"),
          e("other", "外資就活", { stageLabel: "独自選考" }),
          e("offer", "OfferBox"),
        ]}
        tasks={[]}
      />,
    );
    expect(screen.getByTestId("column-count-application")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-interview")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-group")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-other")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-offer")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-document")).toHaveTextContent("0");
  });

  it("0件の列は励ましトーンのコピーを表示する", () => {
    render(<KanbanBoard initialEntries={[e("application", "リクナビ")]} tasks={[]} />);
    expect(screen.getAllByText("まだこの段階の応募先はありません").length).toBeGreaterThan(0);
  });

  it("カードはキーボード操作可能な role=button として描画される", () => {
    render(<KanbanBoard initialEntries={[e("application", "リクナビ")]} tasks={[]} />);
    const buttons = screen.getAllByRole("button");
    expect(buttons.length).toBeGreaterThanOrEqual(1);
  });

  it("group ステージのエントリーは group 列に表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[e("interview", "A"), e("group", "B")]}
        tasks={[]}
      />,
    );
    expect(screen.getByTestId("column-count-interview")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-group")).toHaveTextContent("1");
  });

  it("未知の stageKind は other 列にフォールバックする", () => {
    render(<KanbanBoard initialEntries={[e("coding_test", "A")]} tasks={[]} />);

    expect(screen.getByTestId("column-count-other")).toHaveTextContent("1");
  });

  it("カードに会社名を主表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[{ ...e("application", "リクナビ"), companyName: "テスト商事" }]}
        tasks={[]}
      />,
    );
    expect(screen.getAllByText("テスト商事").length).toBeGreaterThan(0);
  });

  it("カードは列見出しと重複する stageLabel バッジを表示しない", () => {
    render(
      <KanbanBoard
        initialEntries={[
          e("interview", "マイナビ", {
            companyName: "テスト商事",
            stageLabel: "一次面接",
          }),
        ]}
        tasks={[]}
      />,
    );

    expect(screen.queryAllByText("一次面接").length).toBe(0);
  });

  it("未完了タスクがなければ「予定なし」を表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[e("application", "リクナビ")]}
        tasks={[]}
      />,
    );
    expect(screen.getAllByText("予定なし").length).toBeGreaterThan(0);
  });

  it("未完了タスクがあれば次の締切バッジを表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[e("application", "リクナビ")]}
        tasks={[t({ title: "ES提出", dueDate: "2026-07-10", status: "todo" })]}
      />,
    );
    expect(screen.getAllByText("締切 7/10（ES提出）").length).toBeGreaterThan(0);
  });

  it("応募元URLがあるカードは元ページへのリンクを表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[
          {
            ...e("application", "リクナビ"),
            companyName: "テスト商事",
            sourceUrl: "https://job.rikunabi.com/2027/company/r123/",
          },
        ]}
        tasks={[]}
      />,
    );
    expect(screen.getByRole("link", { name: "応募元" })).toHaveAttribute(
      "href",
      "https://job.rikunabi.com/2027/company/r123/",
    );
  });
});

describe("kanbanStageUpdateInput", () => {
  it("ドラッグ先 stageKind に応じて stageLabel と status を作る", () => {
    expect(kanbanStageUpdateInput("group")).toEqual({
      stageKind: "group",
      stageLabel: "GD",
      status: "in_progress",
    });
    expect(kanbanStageUpdateInput("offer")).toEqual({
      stageKind: "offer",
      stageLabel: "内定",
      status: "offered",
    });
    expect(kanbanStageUpdateInput("other")).toEqual({
      stageKind: "other",
      stageLabel: "その他",
      status: "in_progress",
    });
  });
});

describe("normalizeKanbanStageKind", () => {
  it("未知の stageKind を other として扱う", () => {
    expect(normalizeKanbanStageKind("coding_test")).toBe("other");
    expect(normalizeKanbanStageKind("document")).toBe("document");
  });
});

describe("nextTaskFor", () => {
  const now = new Date("2026-07-06T00:00:00Z");

  it("未完了タスクがなければ null を返す", () => {
    expect(nextTaskFor("e-リクナビ", [], now)).toBeNull();
  });

  it("未完了タスクが全て done なら null を返す", () => {
    const tasks = [t({ status: "done", dueDate: "2026-07-10" })];
    expect(nextTaskFor("e-リクナビ", tasks, now)).toBeNull();
  });

  it("dueDate がないタスクは対象外", () => {
    const tasks = [t({ status: "todo", dueDate: null })];
    expect(nextTaskFor("e-リクナビ", tasks, now)).toBeNull();
  });

  it("他 entry のタスクは対象外", () => {
    const tasks = [t({ entryId: "e-マイナビ", dueDate: "2026-07-10" })];
    expect(nextTaskFor("e-リクナビ", tasks, now)).toBeNull();
  });

  it("複数の未完了タスクから最も近い dueDate のものを返す", () => {
    const tasks = [
      t({ id: "t-far", title: "面接対策", dueDate: "2026-08-01" }),
      t({ id: "t-near", title: "ES提出", dueDate: "2026-07-08" }),
      t({ id: "t-mid", title: "説明会", dueDate: "2026-07-20" }),
    ];
    expect(nextTaskFor("e-リクナビ", tasks, now)).toEqual({
      title: "ES提出",
      dueDate: "2026-07-08",
    });
  });
});
