"use client";

import Link from "next/link";
import { useEntries } from "@/hooks/useEntries";
import type { EntryResponse } from "@/lib/api/entries";
import { Reveal } from "./Reveal";

const COLUMNS = [
  { kind: "application", label: "エントリー", color: "var(--color-stage-entry)" },
  { kind: "document", label: "書類選考", color: "var(--color-stage-doc)" },
  { kind: "test", label: "テスト/ES", color: "var(--color-stage-es)" },
  { kind: "interview", label: "面接", color: "var(--color-stage-interview)" },
  { kind: "offer", label: "内定", color: "var(--color-stage-offer)" },
] as const;

/** Entry を stageKind ごとに振り分けるカンバン。 */
export function KanbanBoard() {
  const { data, loading, error } = useEntries();

  if (loading) {
    return <p role="status" className="text-[12px] text-ink-3">読み込み中…</p>;
  }
  if (error) {
    return (
      <p role="alert" className="rounded-lg bg-pink/40 p-3 text-[12px] font-semibold text-ink">
        読み込みに失敗しました
      </p>
    );
  }
  const entries = data ?? [];
  const byKind = new Map<string, EntryResponse[]>();
  for (const e of entries) {
    if (!byKind.has(e.stageKind)) byKind.set(e.stageKind, []);
    byKind.get(e.stageKind)!.push(e);
  }
  // group 等 未列挙のものは interview 列に寄せる
  const groupCol = byKind.get("group") ?? [];
  byKind.set("interview", [...(byKind.get("interview") ?? []), ...groupCol]);

  return (
    <div className="grid gap-2.5 md:grid-cols-5 grid-cols-[repeat(5,minmax(220px,1fr))] overflow-x-auto pb-2">
      {COLUMNS.map((col, i) => {
        const cards = byKind.get(col.kind) ?? [];
        return (
          <Reveal key={col.kind} delay={i * 80}>
            <div className="flex h-full flex-col gap-2 rounded-xl border border-line bg-surface p-2.5">
              <div className="flex items-center gap-2 border-b border-dashed border-line px-1 pb-2">
                <span className="block h-2 w-2 rounded-full" style={{ background: col.color }} />
                <span className="text-[11px] font-extrabold">{col.label}</span>
                <span
                  data-testid={`column-count-${col.kind}`}
                  className="ml-auto font-mono text-[10px] text-ink-3"
                >
                  {cards.length}
                </span>
              </div>
              <ul className="flex flex-col gap-1.5">
                {cards.map((c) => (
                  <li key={c.id}>
                    <Link
                      href={`/entry/${c.id}`}
                      className="block cursor-grab rounded-[10px] border border-line bg-cream p-2.5 transition-all hover:-translate-y-0.5 hover:shadow-[0_6px_14px_-4px_rgba(0,0,0,0.15)]"
                    >
                      <div className="mb-1.5 flex items-center gap-2">
                        <div className="grid h-6 w-6 place-items-center rounded-md bg-sage-wash font-serif text-xs font-extrabold text-sage">
                          {c.source.slice(0, 1)}
                        </div>
                        <div className="truncate text-[10px] font-bold">{c.source}</div>
                      </div>
                      <div className="flex justify-between text-[9px] text-ink-3">
                        <span>{c.route}</span>
                        <span aria-hidden>›</span>
                      </div>
                    </Link>
                  </li>
                ))}
                {cards.length === 0 && (
                  <li className="rounded-md border border-dashed border-line p-2 text-center text-[9px] text-ink-3">
                    まだなし
                  </li>
                )}
              </ul>
            </div>
          </Reveal>
        );
      })}
    </div>
  );
}
