"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { Reveal } from "@/components/entre/Reveal";
import { STAGES, KANBAN_CARDS } from "@/lib/sample-data";

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
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">
              選考カンバン
            </h1>
            <p className="mt-0.5 text-[11px] text-ink-3">ステータスごとに整理</p>
          </div>
          <button
            type="button"
            className="rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5"
          >
            ＋ エントリー
          </button>
        </header>

        {/* Kanban — 5 columns */}
        <div className="grid gap-2.5 md:grid-cols-5 grid-cols-[repeat(5,minmax(220px,1fr))] overflow-x-auto pb-2">
          {STAGES.map((s, colIdx) => (
            <Reveal key={s.key} delay={colIdx * 80}>
              <div className="flex h-full flex-col gap-2 rounded-xl border border-line bg-surface p-2.5">
                <div className="flex items-center gap-2 border-b border-dashed border-line px-1 pb-2">
                  <span
                    className="block h-2 w-2 rounded-full"
                    style={{ background: s.color }}
                  />
                  <span className="text-[11px] font-extrabold">{s.label}</span>
                  <span className="ml-auto font-mono text-[10px] text-ink-3">
                    {KANBAN_CARDS[s.key].length}
                  </span>
                </div>

                <ul className="flex flex-col gap-1.5">
                  {KANBAN_CARDS[s.key].map((c, i) => (
                    <li
                      key={`${s.key}-${i}`}
                      className="cursor-grab rounded-[10px] border border-line bg-cream p-2.5 transition-all hover:-translate-y-0.5 hover:shadow-[0_6px_14px_-4px_rgba(0,0,0,0.15)]"
                    >
                      <div className="mb-1.5 flex items-center gap-2">
                        <div className="grid h-6 w-6 place-items-center rounded-md bg-sage-wash font-serif text-xs font-extrabold text-sage">
                          {c.l}
                        </div>
                        <div className="truncate text-[10px] font-bold">{c.co}</div>
                      </div>
                      <div className="flex justify-between text-[9px] text-ink-3">
                        <span className="font-mono">{c.d}</span>
                        <span aria-hidden>📝</span>
                      </div>
                    </li>
                  ))}
                </ul>
              </div>
            </Reveal>
          ))}
        </div>
      </div>
    </AppShell>
  );
}
