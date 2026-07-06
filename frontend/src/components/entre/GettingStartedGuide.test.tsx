import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import {
  buildGettingStartedSteps,
  GettingStartedGuide,
} from "./GettingStartedGuide";

describe("buildGettingStartedSteps", () => {
  it("何も無ければ全ステップ未完了", () => {
    const steps = buildGettingStartedSteps({ hasEntries: false, hasTasks: false });
    expect(steps.map((s) => s.done)).toEqual([false, false, false]);
    expect(steps.map((s) => s.index)).toEqual([1, 2, 3]);
  });

  it("応募先があれば①②が完了、タスクがなければ③は未完了", () => {
    const steps = buildGettingStartedSteps({ hasEntries: true, hasTasks: false });
    expect(steps.map((s) => s.done)).toEqual([true, true, false]);
  });

  it("応募先とタスクが両方あれば全ステップ完了", () => {
    const steps = buildGettingStartedSteps({ hasEntries: true, hasTasks: true });
    expect(steps.map((s) => s.done)).toEqual([true, true, true]);
  });
});

describe("GettingStartedGuide", () => {
  it("未完了ステップにリンクを表示する", () => {
    render(
      <GettingStartedGuide hasEntries={false} hasTasks={false} hasDoneTasks={false} />,
    );
    expect(screen.getByText("はじめかた")).toBeInTheDocument();
    expect(screen.getByText("0/3 完了")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /登録する/ })).toHaveAttribute(
      "href",
      "/entry/new",
    );
    expect(screen.getByRole("link", { name: /ボードを見る/ })).toHaveAttribute(
      "href",
      "/kanban",
    );
    expect(screen.getByRole("link", { name: /タスクを追加する/ })).toHaveAttribute(
      "href",
      "/task",
    );
  });

  it("完了したステップはリンクの代わりにチェック表示にする", () => {
    render(
      <GettingStartedGuide hasEntries={true} hasTasks={false} hasDoneTasks={false} />,
    );
    expect(screen.getByText("2/3 完了")).toBeInTheDocument();
    expect(screen.queryByRole("link", { name: /登録する/ })).toBeNull();
    expect(screen.queryByRole("link", { name: /ボードを見る/ })).toBeNull();
    expect(screen.getByRole("link", { name: /タスクを追加する/ })).toBeInTheDocument();
  });

  it("hasDoneTasks が true なら実行中の一言を添える", () => {
    render(
      <GettingStartedGuide hasEntries={true} hasTasks={true} hasDoneTasks={true} />,
    );
    expect(screen.getByText("3/3 完了・タスク実行中")).toBeInTheDocument();
  });

  it("compact 時はステップの説明文を省く", () => {
    render(
      <GettingStartedGuide
        hasEntries={false}
        hasTasks={false}
        hasDoneTasks={false}
        compact
      />,
    );
    expect(screen.queryByText("会社名だけでOK・30秒で終わります。")).toBeNull();
  });
});
