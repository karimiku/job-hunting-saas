"use client";

import Link from "next/link";
import { useEntries } from "@/hooks/useEntries";

const STAGE_BG: Record<string, string> = {
  application: "var(--color-stage-entry)",
  document: "var(--color-stage-doc)",
  test: "var(--color-stage-es)",
  interview: "var(--color-stage-interview)",
  group: "var(--color-stage-interview)",
  offer: "var(--color-stage-offer)",
};

/** エントリー一覧を API から取得して表示するビュー。 */
export function EntryListView() {
  const { data, loading, error } = useEntries();

  if (loading) {
    return (
      <p role="status" className="text-[12px] text-ink-3">
        読み込み中…
      </p>
    );
  }

  if (error) {
    return (
      <p role="alert" className="rounded-lg bg-pink/40 p-3 text-[12px] font-semibold text-ink">
        読み込みに失敗しました（{error.message}）
      </p>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-line bg-surface p-6 text-center text-[12px] text-ink-2">
        まだエントリーがありません。＋ボタンから1件追加しましょう。
      </div>
    );
  }

  return (
    <ul className="entre-stagger flex flex-col gap-2">
      {data.map((e) => (
        <li
          key={e.id}
          className="flex cursor-pointer items-center gap-2.5 rounded-xl border border-line bg-surface p-3 transition-all hover:translate-x-0.5 hover:border-sage"
        >
          <Link href={`/entry/${e.id}`} className="flex flex-1 items-center gap-2.5">
            <div className="grid h-9 w-9 place-items-center rounded-[10px] bg-sage-wash font-serif text-lg font-extrabold text-sage">
              {e.companyId.slice(0, 1).toUpperCase()}
            </div>
            <div className="min-w-0 flex-1">
              <div className="text-[12px] font-bold">{e.source}</div>
              <div className="mt-0.5 flex items-center gap-1.5 text-[10px] text-ink-3">
                <span
                  className="rounded-sm px-1.5 py-0.5 text-[8px] font-bold text-white"
                  style={{ background: STAGE_BG[e.stageKind] ?? "var(--color-ink-3)" }}
                >
                  {e.stageLabel}
                </span>
                <span>{e.route}</span>
              </div>
              {e.memo && <div className="mt-1 text-[10px] text-ink-2">{e.memo}</div>}
            </div>
          </Link>
          <span className="text-ink-3" aria-hidden>›</span>
        </li>
      ))}
    </ul>
  );
}
