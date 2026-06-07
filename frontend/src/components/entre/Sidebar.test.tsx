import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { Sidebar } from "./Sidebar";

describe("Sidebar", () => {
  it("navCounts を渡すと実カウントをバッジ表示する", () => {
    render(<Sidebar navCounts={{ entry: 12, task: 3, inbox: 5 }} />);
    expect(screen.getByTestId("nav-count-entry")).toHaveTextContent("12");
    expect(screen.getByTestId("nav-count-task")).toHaveTextContent("3");
    expect(screen.getByTestId("nav-count-inbox")).toHaveTextContent("5");
  });

  it("カウントが 0 でもバッジを表示する", () => {
    render(<Sidebar navCounts={{ entry: 0, task: 0, inbox: 0 }} />);
    expect(screen.getByTestId("nav-count-entry")).toHaveTextContent("0");
    expect(screen.getByTestId("nav-count-task")).toHaveTextContent("0");
    expect(screen.getByTestId("nav-count-inbox")).toHaveTextContent("0");
  });

  it("navCounts 未指定ならバッジを出さない（固定値を出さない）", () => {
    render(<Sidebar />);
    expect(screen.queryByTestId("nav-count-entry")).toBeNull();
    expect(screen.queryByTestId("nav-count-task")).toBeNull();
    expect(screen.queryByTestId("nav-count-inbox")).toBeNull();
  });

  it("userName を表示する", () => {
    render(<Sidebar userName="山田 太郎" navCounts={{ entry: 1, task: 1, inbox: 1 }} />);
    expect(screen.getByText("山田 太郎")).toBeInTheDocument();
  });
});
