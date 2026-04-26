"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { ENTRIES, STAGES } from "@/lib/sample-data";

export default function EntryListPage() {
  const router = useRouter();
  const state = useUser();
  const [filter, setFilter] = useState<string>("all");

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  const filtered =
    filter === "all" ? ENTRIES : ENTRIES.filter((e) => e.stage === filter);

  return (
    <AppShell userName={state.user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[900px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-baseline justify-between">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">
              Entry
              <span className="ml-2 text-[11px] font-medium text-ink-3">
                {ENTRIES.length}社
              </span>
            </h1>
          </div>
          <button
            type="button"
            className="rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5"
          >
            ＋ 追加
          </button>
        </header>

        {/* Filter chips */}
        <div className="mb-4 flex gap-1.5 overflow-x-auto pb-1">
          <button
            type="button"
            onClick={() => setFilter("all")}
            className={`whitespace-nowrap rounded-full border px-3 py-1.5 text-[10px] font-bold transition-colors ${
              filter === "all"
                ? "border-sage bg-sage text-white"
                : "border-line bg-surface text-ink-2"
            }`}
          >
            全て
          </button>
          {STAGES.map((s) => (
            <button
              key={s.key}
              type="button"
              onClick={() => setFilter(s.key)}
              className={`whitespace-nowrap rounded-full border px-3 py-1.5 text-[10px] font-bold transition-colors ${
                filter === s.key
                  ? "border-sage bg-sage text-white"
                  : "border-line bg-surface text-ink-2"
              }`}
            >
              {s.label}
            </button>
          ))}
        </div>

        {/* Entries list */}
        <ul className="entre-stagger flex flex-col gap-2">
          {filtered.map((e) => (
            <li
              key={e.id}
              className="flex cursor-pointer items-center gap-2.5 rounded-xl border border-line bg-surface p-3 transition-all hover:translate-x-0.5 hover:border-sage"
            >
              <div className="grid h-9 w-9 place-items-center rounded-[10px] bg-sage-wash font-serif text-lg font-extrabold text-sage">
                {e.logo}
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-1.5">
                  <span className="truncate text-[12px] font-bold">{e.co}</span>
                  {e.fav && (
                    <span
                      className="text-[9px] text-amber"
                      style={{ animation: "entre-glow 2s infinite" }}
                      aria-label="favorite"
                    >
                      ★
                    </span>
                  )}
                </div>
                <div className="mt-0.5 flex items-center gap-1.5 text-[9px] text-ink-3">
                  <span
                    className="rounded-sm px-1.5 py-0.5 text-[8px] font-bold text-white"
                    style={{ background: e.color }}
                  >
                    {e.stageLabel}
                  </span>
                  <span className="font-mono">{e.due}</span>
                </div>
                <div className="mt-1 text-[10px] text-ink-2">{e.task}</div>
              </div>
              <span className="text-ink-3" aria-hidden>›</span>
            </li>
          ))}
        </ul>
      </div>
    </AppShell>
  );
}
