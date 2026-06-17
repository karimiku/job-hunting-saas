import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import {
  KanbanBoard,
  kanbanStageUpdateInput,
  normalizeKanbanStageKind,
} from "./KanbanBoard";
import type { EntryResponse } from "@/lib/api/entries";

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
      />,
    );
    expect(screen.getByTestId("column-count-application")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-interview")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-group")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-other")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-offer")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-document")).toHaveTextContent("0");
  });

  it("カードはキーボード操作可能な role=button として描画される", () => {
    render(<KanbanBoard initialEntries={[e("application", "リクナビ")]} />);
    const buttons = screen.getAllByRole("button");
    expect(buttons.length).toBeGreaterThanOrEqual(1);
  });

  it("group ステージのエントリーは group 列に表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[e("interview", "A"), e("group", "B")]}
      />,
    );
    expect(screen.getByTestId("column-count-interview")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-group")).toHaveTextContent("1");
  });

  it("未知の stageKind は other 列にフォールバックする", () => {
    render(<KanbanBoard initialEntries={[e("coding_test", "A")]} />);

    expect(screen.getByTestId("column-count-other")).toHaveTextContent("1");
  });

  it("カードに会社名を主表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[{ ...e("application", "リクナビ"), companyName: "テスト商事" }]}
      />,
    );
    expect(screen.getAllByText("テスト商事").length).toBeGreaterThan(0);
  });

  it("カードに Entry の stageLabel を表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[
          e("interview", "マイナビ", {
            companyName: "テスト商事",
            stageLabel: "一次面接",
          }),
        ]}
      />,
    );

    expect(screen.getAllByText("一次面接").length).toBeGreaterThan(0);
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
