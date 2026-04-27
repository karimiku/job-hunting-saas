// Server Component。entries は props で受け取る (SSR 時に集計)。
// CountUp / Reveal が animation 用の Client なので、それらは内部で混在表示される。

import type { EntryResponse } from "@/lib/api/entries";
import { Mascot } from "./Mascot";
import { Reveal } from "./Reveal";
import { CountUp } from "./CountUp";

const MILESTONES = [
  { kind: "application", label: "エントリー", color: "var(--color-stage-entry)", emoji: "📝" },
  { kind: "document", label: "書類選考", color: "var(--color-stage-doc)", emoji: "📄" },
  { kind: "test", label: "ES/テスト", color: "var(--color-stage-es)", emoji: "✏️" },
  { kind: "interview", label: "面接", color: "var(--color-stage-interview)", emoji: "🤝" },
  { kind: "offer", label: "内定", color: "var(--color-stage-offer)", emoji: "🎉" },
] as const;

const STAGE_INDEX: Record<string, number> = {
  application: 0,
  document: 1,
  test: 2,
  interview: 3,
  group: 3,
  offer: 4,
};

interface MilestoneStat {
  kind: string;
  label: string;
  color: string;
  emoji: string;
  count: number;
  done: boolean;
  current: boolean;
}

/** ロードマップ — Entry の最大ステージから現在地を判定し、各マイルストーンの社数を集計。 */
export function RoadmapView({ entries }: { entries: EntryResponse[] }) {
  const counts = new Map<string, number>();
  let maxIdx = -1;
  for (const e of entries) {
    counts.set(e.stageKind, (counts.get(e.stageKind) ?? 0) + 1);
    const idx = STAGE_INDEX[e.stageKind];
    if (idx !== undefined && idx > maxIdx) maxIdx = idx;
  }

  const milestones: MilestoneStat[] = MILESTONES.map((m, i) => ({
    ...m,
    count: counts.get(m.kind) ?? (m.kind === "test" ? counts.get("test") ?? 0 : 0),
    done: i < maxIdx,
    current: i === maxIdx,
  }));

  // 進捗 % は max ステージ / (全ステージ-1)
  const progressPct = maxIdx >= 0 ? Math.min(100, Math.round((maxIdx / 4) * 100)) : 0;

  return (
    <div>
      <Reveal>
        <div className="mb-6 rounded-xl border border-line bg-surface p-5">
          <div className="mb-2 flex items-baseline justify-between">
            <p className="text-[12px] font-bold">今のあなたの進捗</p>
            <p className="font-serif text-2xl font-extrabold text-sage">
              <CountUp end={progressPct} duration={1100} />
              <span className="text-xs">%</span>
            </p>
          </div>
          <div className="h-2.5 overflow-hidden rounded-md bg-line">
            <div
              className="h-full rounded-md bg-gradient-to-r from-sage-mid to-sage transition-[width] duration-1000"
              style={{ width: `${progressPct}%` }}
            />
          </div>
          <p className="mt-2 font-hand text-[14px] text-sage">
            {progressPct >= 100 ? "おめでとう、内定 🎉" : "あと少しですね、ファイト！"}
          </p>
        </div>
      </Reveal>

      {/* Desktop: horizontal milestone path */}
      <div className="relative hidden py-10 md:block">
        <div className="relative flex justify-between px-5">
          {milestones.map((m, i) => (
            <Reveal key={m.kind} delay={i * 120}>
              <MilestoneNode m={m} />
            </Reveal>
          ))}
        </div>
      </div>

      {/* Mobile: vertical */}
      <ul className="flex flex-col gap-3 md:hidden">
        {milestones.map((m, i) => (
          <Reveal key={m.kind} delay={i * 80}>
            <li
              className="flex items-center gap-3 rounded-xl border-[1.5px] bg-surface p-3.5"
              style={{
                borderColor: m.current ? m.color : "var(--color-line)",
                boxShadow: m.current ? `0 6px 14px -4px ${m.color}66` : "none",
              }}
            >
              <span className="text-2xl">{m.emoji}</span>
              <div className="flex-1">
                <p className="font-serif text-[15px] font-extrabold">{m.label}</p>
                <p
                  data-testid={`milestone-count-${m.kind}`}
                  className="text-[10px] text-ink-2"
                >
                  {m.count}社
                </p>
              </div>
              {m.current && (
                <span
                  className="rounded-full px-2 py-0.5 text-[9px] font-bold text-white"
                  style={{ background: m.color }}
                >
                  NOW
                </span>
              )}
              {m.done && (
                <span className="rounded-full bg-sage-soft px-2 py-0.5 text-[9px] font-bold text-sage">
                  クリア
                </span>
              )}
            </li>
          </Reveal>
        ))}
      </ul>

      {/* Mascot encouragement */}
      <Reveal delay={300}>
        <div className="mt-6 flex flex-col items-center rounded-xl border-[1.5px] border-line bg-gradient-to-br from-cream-2 to-sage-wash p-5 text-center">
          <Mascot size={64} mood={progressPct >= 100 ? "cheering" : "thinking"} />
          <p className="mt-2 font-hand text-[20px] text-sage">
            {progressPct >= 100 ? "やりましたね！" : "あと一歩！"}
          </p>
          <p className="mt-1 font-serif text-sm font-extrabold">{progressPct}% 完了</p>
        </div>
      </Reveal>
    </div>
  );
}

function MilestoneNode({ m }: { m: MilestoneStat }) {
  return (
    <div className="flex flex-1 flex-col items-center">
      <div
        className="relative z-[1] grid h-[60px] w-[60px] place-items-center rounded-full border-[3.5px] text-2xl font-extrabold"
        style={{
          background: m.done ? m.color : m.current ? "#fff" : "var(--color-line)",
          borderColor: m.color,
          color: m.done ? "#fff" : m.color,
          boxShadow: m.current ? `0 8px 20px -4px ${m.color}66` : "none",
          animation: m.current ? "entre-pulse-ring 1.8s infinite" : undefined,
        }}
      >
        {m.emoji}
      </div>
      <p className="mt-2.5 font-serif text-base font-extrabold">{m.label}</p>
      {/* mobile 側に同じ testid があるため、desktop ノードは数値表示のみ（assertion 対象外） */}
      <p className="mt-1 text-[11px] text-ink-2">{m.count}社</p>
      {m.current && (
        <span
          className="mt-1.5 rounded-full px-2.5 py-0.5 text-[10px] font-bold text-white"
          style={{ background: m.color, animation: "entre-wiggle 2s infinite" }}
        >
          NOW
        </span>
      )}
      {m.done && (
        <span className="mt-1.5 rounded-full bg-sage-soft px-2.5 py-0.5 text-[10px] font-bold text-sage">
          クリア ✓
        </span>
      )}
    </div>
  );
}
