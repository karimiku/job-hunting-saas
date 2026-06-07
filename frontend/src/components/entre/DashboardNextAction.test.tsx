import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import {
  DashboardNextAction,
  getDashboardNextAction,
} from "./DashboardNextAction";

describe("getDashboardNextAction", () => {
  it("保存クリップがあれば最優先でEntry化へ案内する", () => {
    expect(
      getDashboardNextAction({ inboxCount: 2, entryCount: 4, openTaskCount: 3 }),
    ).toMatchObject({
      activeStep: "inbox",
      href: "/inbox",
      cta: "保存箱を開く",
    });
  });

  it("保存クリップもEntryも無ければEntry追加へ案内する", () => {
    expect(
      getDashboardNextAction({ inboxCount: 0, entryCount: 0, openTaskCount: 0 }),
    ).toMatchObject({
      activeStep: "entry",
      href: "/entry/new",
      cta: "Entryを追加",
    });
  });

  it("Entryはあるが未完了タスクが無ければタスク追加へ案内する", () => {
    expect(
      getDashboardNextAction({ inboxCount: 0, entryCount: 3, openTaskCount: 0 }),
    ).toMatchObject({
      activeStep: "task",
      href: "/task",
      cta: "タスクを追加",
    });
  });

  it("未完了タスクがあればタスク確認へ案内する", () => {
    expect(
      getDashboardNextAction({ inboxCount: 0, entryCount: 3, openTaskCount: 5 }),
    ).toMatchObject({
      activeStep: "task",
      href: "/task",
      cta: "タスクを見る",
    });
  });
});

describe("DashboardNextAction", () => {
  it("現在の次アクションだけを表示する", () => {
    render(
      <DashboardNextAction inboxCount={1} entryCount={0} openTaskCount={0} />,
    );

    expect(screen.getByText("次にやること")).toBeInTheDocument();
    expect(screen.getByText("保存クリップ 1件をEntryにする")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /保存箱を開く/ })).toHaveAttribute(
      "href",
      "/inbox",
    );
    expect(screen.queryByText("1. 保存")).toBeNull();
  });
});
