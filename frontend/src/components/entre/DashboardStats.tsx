"use client";

import { useEntries } from "@/hooks/useEntries";
import { CountUp } from "./CountUp";

/** Entry を集計して 4 つの stat タイルを表示する。 */
export function DashboardStats() {
  const { data, loading } = useEntries();

  if (loading || !data) {
    return (
      <p role="status" className="text-[12px] text-ink-3">
        読み込み中…
      </p>
    );
  }

  const total = data.length;
  const inProgress = data.filter((e) => e.status === "in_progress").length;
  const interviewing = data.filter(
    (e) => e.stageKind === "interview" || e.stageKind === "group",
  ).length;
  const offered = data.filter(
    (e) => e.status === "offered" || e.status === "accepted",
  ).length;

  const stats = [
    { v: total, l: "エントリー数", c: "text-sage", testId: "stat-total" },
    { v: interviewing, l: "面接中", c: "text-pink-deep", testId: "stat-interviewing" },
    { v: inProgress, l: "選考中", c: "text-amber", testId: "stat-in-progress" },
    { v: offered, l: "内定", c: "text-mint", testId: "stat-offered" },
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

const STAGE_LABEL_BY_KIND: Record<string, string> = {
  application: "エントリー",
  document: "書類",
  test: "テスト",
  interview: "面接",
  group: "GD",
  offer: "内定",
};

/** ダッシュボード補助: ステータス円グラフ用の集計を返す。 */
export function useStageBreakdown(): { kind: string; label: string; count: number }[] {
  const { data } = useEntries();
  if (!data) return [];
  const counts = new Map<string, number>();
  for (const e of data) {
    counts.set(e.stageKind, (counts.get(e.stageKind) ?? 0) + 1);
  }
  return Array.from(counts.entries()).map(([kind, count]) => ({
    kind,
    label: STAGE_LABEL_BY_KIND[kind] ?? kind,
    count,
  }));
}
