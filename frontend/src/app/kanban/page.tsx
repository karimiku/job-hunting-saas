"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { KanbanBoard } from "@/components/entre/KanbanBoard";

export default function KanbanPage() {
  const router = useRouter();
  const state = useUser();

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  return (
    <AppShell userName={state.user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[1400px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-5 flex items-baseline justify-between">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">選考カンバン</h1>
            <p className="mt-0.5 text-[11px] text-ink-3">ステータスごとに整理</p>
          </div>
          <button
            type="button"
            className="rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5"
          >
            ＋ エントリー
          </button>
        </header>

        <KanbanBoard />
      </div>
    </AppShell>
  );
}
