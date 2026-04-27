import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { KanbanBoard } from "./KanbanBoard";
import type { EntryResponse } from "@/lib/api/entries";

const e = (kind: string, source: string): EntryResponse => ({
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
});

describe("KanbanBoard", () => {
  it("ステージごとの 5 列に initialEntries を振り分けて表示する", () => {
    render(
      <KanbanBoard
        initialEntries={[
          e("application", "リクナビ"),
          e("interview", "マイナビ"),
          e("interview", "ONE CAREER"),
          e("offer", "OfferBox"),
        ]}
      />,
    );
    expect(screen.getByTestId("column-count-application")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-interview")).toHaveTextContent("2");
    expect(screen.getByTestId("column-count-offer")).toHaveTextContent("1");
    expect(screen.getByTestId("column-count-document")).toHaveTextContent("0");
  });

  it("カードはキーボード操作可能な role=button として描画される", () => {
    render(<KanbanBoard initialEntries={[e("application", "リクナビ")]} />);
    const buttons = screen.getAllByRole("button");
    expect(buttons.length).toBeGreaterThanOrEqual(1);
  });

  it("group ステージのエントリーは interview 列に集約される", () => {
    render(
      <KanbanBoard
        initialEntries={[e("interview", "A"), e("group", "B")]}
      />,
    );
    expect(screen.getByTestId("column-count-interview")).toHaveTextContent("2");
  });
});
