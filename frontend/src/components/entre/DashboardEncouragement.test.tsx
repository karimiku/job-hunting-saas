import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import {
  DashboardEncouragement,
  buildEncouragement,
} from "./DashboardEncouragement";

describe("buildEncouragement", () => {
  it("内定があれば祝福コピーを返す", () => {
    const c = buildEncouragement({ interviewing: 2, offered: 1, openTasks: 3 });
    expect(c.headline).toContain("内定 1 件");
  });

  it("面接中があれば面接コピー（内定が無い場合）", () => {
    const c = buildEncouragement({ interviewing: 5, offered: 0, openTasks: 2 });
    expect(c.headline).toContain("面接 5 社");
    expect(c.body).toContain("2 件");
  });

  it("選考のみで未完了タスクがあればタスク促しコピー", () => {
    const c = buildEncouragement({ interviewing: 0, offered: 0, openTasks: 4 });
    expect(c.headline).toContain("4 件");
  });

  it("何も無ければ汎用励ましコピー", () => {
    const c = buildEncouragement({ interviewing: 0, offered: 0, openTasks: 0 });
    expect(c.headline).toContain("完了");
  });
});

describe("DashboardEncouragement", () => {
  it("実データから見出しを描画する", () => {
    render(<DashboardEncouragement interviewing={3} offered={0} openTasks={0} />);
    expect(screen.getByTestId("encouragement-headline")).toHaveTextContent("面接 3 社");
  });
});
