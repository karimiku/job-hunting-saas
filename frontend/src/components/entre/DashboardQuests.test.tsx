import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import {
  DashboardQuests,
  buildQuests,
  questProgress,
} from "./DashboardQuests";
import type { TaskWithEntry } from "@/lib/api/server-resources";

const t = (overrides: Partial<TaskWithEntry> = {}): TaskWithEntry => ({
  id: String(Math.random()),
  entryId: "e1",
  title: "ES提出",
  type: "deadline",
  status: "todo",
  dueDate: null,
  memo: "",
  createdAt: "x",
  updatedAt: "x",
  companyName: "○○商事",
  ...overrides,
});

describe("buildQuests", () => {
  it("未完了を期限昇順、完了を末尾に並べる", () => {
    const now = new Date("2026-05-29T00:00:00Z");
    const quests = buildQuests(
      [
        t({ id: "done", status: "done", dueDate: "2026-05-20" }),
        t({ id: "late", dueDate: "2026-06-10" }),
        t({ id: "soon", dueDate: "2026-05-30" }),
        t({ id: "none", dueDate: null }),
      ],
      now,
    );
    expect(quests.map((q) => q.id)).toEqual(["soon", "late", "none", "done"]);
  });

  it("会社名 + タイトルをラベルにする", () => {
    const [q] = buildQuests([t({ companyName: "△△株式会社", title: "SPI受験" })]);
    expect(q.label).toBe("△△株式会社 SPI受験");
  });

  it("最大 5 件に絞る", () => {
    const quests = buildQuests(Array.from({ length: 8 }, () => t()));
    expect(quests).toHaveLength(5);
  });

  it("期限の近さでバッジ色を決める", () => {
    const now = new Date("2026-05-29T00:00:00Z");
    expect(buildQuests([t({ dueDate: "2026-05-28" })], now)[0].color).toBe("bg-pink");
    expect(buildQuests([t({ dueDate: "2026-05-31" })], now)[0].color).toBe("bg-amber");
    expect(buildQuests([t({ dueDate: "2026-06-15" })], now)[0].color).toBe("bg-sky");
    expect(buildQuests([t({ dueDate: null })], now)[0].color).toBe("bg-sage");
  });

  it("期限切れは「M/D ・n日超過」を due ラベルにし bg-pink のまま強調する", () => {
    const now = new Date("2026-05-29T00:00:00Z");
    const [q] = buildQuests([t({ dueDate: "2026-05-20" })], now);
    expect(q.due).toBe("5/20 ・9日超過");
    expect(q.color).toBe("bg-pink");
  });

  it("今日が期限なら超過表記にしない", () => {
    const now = new Date("2026-05-29T00:00:00Z");
    const [q] = buildQuests([t({ dueDate: "2026-05-29" })], now);
    expect(q.due).toBe("5/29");
  });
});

describe("questProgress", () => {
  it("完了率を整数%で返す", () => {
    expect(questProgress([])).toBe(0);
    expect(
      questProgress([t({ status: "done" }), t({ status: "todo" })]),
    ).toBe(50);
  });
});

describe("DashboardQuests", () => {
  it("見出しは「直近のタスク」", () => {
    render(<DashboardQuests tasks={[]} />);
    expect(screen.getByText("直近のタスク")).toBeInTheDocument();
  });

  it("実タスクをクエストとして描画する", () => {
    render(
      <DashboardQuests
        tasks={[t({ title: "一次面接", companyName: "○○商事", dueDate: null })]}
      />,
    );
    expect(screen.getByText("○○商事 一次面接")).toBeInTheDocument();
    expect(screen.queryByTestId("quest-empty")).toBeNull();
  });

  it("タスクが無ければ応募先の登録を促す空状態を表示する", () => {
    render(<DashboardQuests tasks={[]} />);
    expect(screen.getByTestId("quest-empty")).toBeInTheDocument();
    expect(screen.getByText("応募先ごとに締切や予定を追加すると、近い順に表示されます。")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "応募先を確認" })).toHaveAttribute("href", "/entry");
  });
});
