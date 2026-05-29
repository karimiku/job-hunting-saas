import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import {
  DashboardNextAction,
  getDashboardNextAction,
} from "./DashboardNextAction";

describe("getDashboardNextAction", () => {
  it("Inbox の保存クリップがあれば最優先で Entry 化へ案内する", () => {
    expect(
      getDashboardNextAction({ inboxCount: 2, entryCount: 4, openTaskCount: 3 }),
    ).toMatchObject({
      activeStep: "inbox",
      href: "/inbox",
      cta: "Inboxで確認",
    });
  });

  it("保存クリップも Entry も無ければ Entry 追加へ案内する", () => {
    expect(
      getDashboardNextAction({ inboxCount: 0, entryCount: 0, openTaskCount: 0 }),
    ).toMatchObject({
      activeStep: "entry",
      href: "/entry/new",
      cta: "Entryを追加",
    });
  });

  it("Entry はあるが未完了 Task が無ければ Task 追加へ案内する", () => {
    expect(
      getDashboardNextAction({ inboxCount: 0, entryCount: 3, openTaskCount: 0 }),
    ).toMatchObject({
      activeStep: "task",
      href: "/task",
      cta: "Taskを追加",
    });
  });

  it("未完了 Task があれば Task 確認へ案内する", () => {
    expect(
      getDashboardNextAction({ inboxCount: 0, entryCount: 3, openTaskCount: 5 }),
    ).toMatchObject({
      activeStep: "task",
      href: "/task",
      cta: "Taskを見る",
    });
  });
});

describe("DashboardNextAction", () => {
  it("応募管理の流れと現在の次アクションを表示する", () => {
    render(
      <DashboardNextAction inboxCount={1} entryCount={0} openTaskCount={0} />,
    );

    expect(screen.getByText("次にやること")).toBeInTheDocument();
    expect(screen.getByText("保存クリップ 1件を応募先にする")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Inboxで確認/ })).toHaveAttribute(
      "href",
      "/inbox",
    );
    expect(screen.getByText("1. 保存")).toBeInTheDocument();
    expect(screen.getByText("2. 応募先化")).toBeInTheDocument();
    expect(screen.getByText("3. 締切管理")).toBeInTheDocument();
  });
});
