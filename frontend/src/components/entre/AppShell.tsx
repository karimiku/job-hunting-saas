"use client";

import { ReactNode } from "react";
import { Sidebar, MobileTabBar } from "./Sidebar";

interface AppShellProps {
  children: ReactNode;
  userName?: string;
  userSubtitle?: string;
}

/** Entré アプリ全体のシェル。デスクトップ＝サイドバー、モバイル＝ボトムタブ。 */
export function AppShell({ children, userName, userSubtitle }: AppShellProps) {
  return (
    <div className="flex min-h-screen bg-cream font-sans text-ink">
      <Sidebar userName={userName} userSubtitle={userSubtitle} />
      <main className="min-w-0 flex-1 overflow-x-hidden pb-20 md:pb-0">{children}</main>
      <MobileTabBar />
    </div>
  );
}
