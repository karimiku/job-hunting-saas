"use client";

import { useEntries } from "@/hooks/useEntries";

const ORDER = ["application", "document", "test", "interview", "group", "offer"] as const;
const LABEL: Record<string, string> = {
  application: "エントリー",
  document: "書類選考",
  test: "テスト",
  interview: "面接",
  group: "GD",
  offer: "内定",
};
const COLOR: Record<string, string> = {
  application: "var(--color-stage-entry)",
  document: "var(--color-stage-doc)",
  test: "var(--color-stage-es)",
  interview: "var(--color-stage-interview)",
  group: "var(--color-stage-interview-deep)",
  offer: "var(--color-stage-offer)",
};

/** ステージ別件数の凡例 + ミニ円グラフ。 */
export function StatusBreakdown() {
  const { data } = useEntries();
  const counts = new Map<string, number>();
  if (data) {
    for (const e of data) {
      counts.set(e.stageKind, (counts.get(e.stageKind) ?? 0) + 1);
    }
  }
  const total = data?.length ?? 0;

  // 円グラフ用 — 各ステージの dasharray (周長 113 を total で按分)
  let offset = 0;
  const segments = ORDER.map((kind) => {
    const count = counts.get(kind) ?? 0;
    const length = total > 0 ? (count / total) * 113 : 0;
    const seg = { kind, length, offset };
    offset -= length;
    return seg;
  });

  return (
    <div className="flex items-center gap-3.5">
      <svg width="100" height="100" viewBox="0 0 50 50" aria-label="ステータス別件数">
        <circle cx="25" cy="25" r="18" fill="none" stroke="var(--color-line)" strokeWidth="6" />
        {segments.map((s) =>
          s.length > 0 ? (
            <circle
              key={s.kind}
              cx="25"
              cy="25"
              r="18"
              fill="none"
              stroke={COLOR[s.kind]}
              strokeWidth="6"
              strokeDasharray={`${s.length} 113`}
              strokeDashoffset={s.offset}
              transform="rotate(-90 25 25)"
            />
          ) : null,
        )}
      </svg>
      <ul className="flex flex-1 flex-col gap-1.5 text-[11px]">
        {ORDER.map((kind) => (
          <li key={kind} className="flex items-center gap-2">
            <span className="block h-2.5 w-2.5 rounded-sm" style={{ background: COLOR[kind] }} />
            <span className="flex-1">{LABEL[kind]}</span>
            <span data-testid={`count-${kind}`} className="font-mono font-bold">
              {counts.get(kind) ?? 0}
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}
