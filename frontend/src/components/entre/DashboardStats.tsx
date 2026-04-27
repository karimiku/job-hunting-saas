// Server Component。エントリー集計を SSR で計算し、子の CountUp (Client) に数字を渡す。

import type { EntryResponse } from "@/lib/api/entries";
import { CountUp } from "./CountUp";

interface Stats {
  total: number;
  inProgress: number;
  interviewing: number;
  offered: number;
}

export function summarizeEntries(entries: EntryResponse[]): Stats {
  return {
    total: entries.length,
    inProgress: entries.filter((e) => e.status === "in_progress").length,
    interviewing: entries.filter(
      (e) => e.stageKind === "interview" || e.stageKind === "group",
    ).length,
    offered: entries.filter(
      (e) => e.status === "offered" || e.status === "accepted",
    ).length,
  };
}

/** Entry を集計して 4 つの stat タイルを表示する。 */
export function DashboardStats({ entries }: { entries: EntryResponse[] }) {
  const s = summarizeEntries(entries);
  const stats = [
    { v: s.total, l: "エントリー数", c: "text-sage", testId: "stat-total" },
    { v: s.interviewing, l: "面接中", c: "text-pink-deep", testId: "stat-interviewing" },
    { v: s.inProgress, l: "選考中", c: "text-amber", testId: "stat-in-progress" },
    { v: s.offered, l: "内定", c: "text-mint", testId: "stat-offered" },
  ];

  return (
    <div className="entre-stagger grid grid-cols-2 gap-3 md:grid-cols-4">
      {stats.map((s, i) => (
        <div key={s.l} className="rounded-xl border border-line bg-surface px-4 py-3.5">
          <div className="text-[11px] font-semibold text-ink-3">{s.l}</div>
          <div
            data-testid={s.testId}
            className={`mt-1.5 font-serif text-3xl font-extrabold leading-tight ${s.c}`}
          >
            <CountUp end={s.v} duration={900 + i * 100} />.
          </div>
        </div>
      ))}
    </div>
  );
}
