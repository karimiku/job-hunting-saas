"use client";

import { ReactNode } from "react";
import { Sidebar, MobileTabBar, type NavCounts } from "./Sidebar";

interface AppShellProps {
  children: ReactNode;
  userName?: string;
  userSubtitle?: string;
  /** サイドバーのバッジ用カウント。Server Component で集計して渡す。 */
  navCounts?: Partial<NavCounts>;
}

/** Entré アプリ全体のシェル。デスクトップ＝サイドバー、モバイル＝ボトムタブ。 */
export function AppShell({ children, userName, userSubtitle, navCounts }: AppShellProps) {
  return (
    <div className="flex min-h-screen bg-cream font-sans text-ink">
      <Sidebar userName={userName} userSubtitle={userSubtitle} navCounts={navCounts} />
      <main className="min-w-0 flex-1 overflow-x-hidden pb-28 md:pb-0">{children}</main>
      <MobileTabBar />
    </div>
  );
}
