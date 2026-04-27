"use client";

import { useState } from "react";
import { useEntry } from "@/hooks/useEntry";
import { useTasksByEntry } from "@/hooks/useTasksByEntry";
import { updateEntry } from "@/lib/api/entries";
import { Mascot } from "./Mascot";
import { Confetti } from "./Confetti";

const STAGE_ORDER = ["application", "document", "test", "interview", "group", "offer"] as const;
const STAGE_LABEL: Record<string, string> = {
  application: "エントリー",
  document: "書類",
  test: "テスト",
  interview: "面接",
  group: "GD",
  offer: "内定",
};
const STAGE_COLOR: Record<string, string> = {
  application: "var(--color-stage-entry)",
  document: "var(--color-stage-doc)",
  test: "var(--color-stage-es)",
  interview: "var(--color-stage-interview)",
  group: "var(--color-stage-interview-deep)",
  offer: "var(--color-stage-offer)",
};

interface Props {
  entryId: string;
}

/** Entry 詳細 — ステージ進捗バー + 「進める →」 + 内定スタンプ + Tasks 表示。 */
export function EntryDetailView({ entryId }: Props) {
  const entry = useEntry(entryId);
  const tasks = useTasksByEntry(entryId);
  const [confetti, setConfetti] = useState(0);
  const [advancing, setAdvancing] = useState(false);

  if (entry.loading) {
    return <p role="status" className="text-[12px] text-ink-3">読み込み中…</p>;
  }
  if (entry.error || !entry.data) {
    return (
      <p role="alert" className="rounded-lg bg-pink/40 p-3 text-[12px] font-semibold text-ink">
        詳細を読み込めませんでした
      </p>
    );
  }

  const e = entry.data;
  const currentIdx = STAGE_ORDER.indexOf(e.stageKind as (typeof STAGE_ORDER)[number]);
  const isOffer = e.stageKind === "offer";

  const handleAdvance = async () => {
    if (currentIdx < 0 || currentIdx >= STAGE_ORDER.length - 1) return;
    setAdvancing(true);
    const nextKind = STAGE_ORDER[currentIdx + 1];
    try {
      await updateEntry(e.id, {
        stageKind: nextKind,
        stageLabel: STAGE_LABEL[nextKind],
      });
      entry.refetch();
      if (nextKind === "offer") {
        setConfetti((n) => n + 1);
      }
    } finally {
      setAdvancing(false);
    }
  };

  return (
    <div className="relative">
      {/* Header */}
      <div className="mb-4 flex items-center gap-3">
        <div className="grid h-14 w-14 place-items-center rounded-xl bg-sage-wash font-serif text-2xl font-extrabold text-sage">
          {e.source.slice(0, 1)}
        </div>
        <div className="flex-1">
          <h1 className="font-serif text-lg font-extrabold tracking-tight">{e.source}</h1>
          <p className="mt-0.5 text-[10px] text-ink-3">{e.route}</p>
        </div>
        {isOffer && (
          <div
            className="rounded-lg border-[2.5px] border-mint bg-mint/10 px-2.5 py-1.5 font-serif text-sm font-black text-mint"
            style={{ animation: "entre-stamp 0.6s cubic-bezier(0.2, 0.8, 0.4, 1) both" }}
          >
            内定！
          </div>
        )}
      </div>

      {/* Stage progression bar */}
      <section className="mb-3 rounded-xl border border-line bg-surface p-3">
        <div className="mb-2 flex items-center justify-between">
          <p className="text-[10px] font-bold text-ink-2">選考ステータス</p>
          <button
            type="button"
            onClick={handleAdvance}
            disabled={advancing || isOffer}
            className="rounded border border-sage px-1.5 py-0.5 text-[9px] font-bold text-sage transition-opacity disabled:opacity-50"
          >
            進める →
          </button>
        </div>
        <div className="flex gap-1">
          {STAGE_ORDER.map((kind, i) => {
            const reached = i <= currentIdx;
            return (
              <div
                key={kind}
                className="grid flex-1 place-items-center rounded-md text-[8px] font-bold transition-colors"
                style={{
                  background: reached ? STAGE_COLOR[kind] : "var(--color-line)",
                  color: reached ? "#fff" : "var(--color-ink-3)",
                  height: 26,
                }}
              >
                {STAGE_LABEL[kind]}
              </div>
            );
          })}
        </div>
        <p className="mt-2 text-[10px] font-bold text-sage">
          📍 現在: <span data-testid="current-stage">{e.stageLabel}</span>
        </p>
      </section>

      {/* Memo */}
      {e.memo && (
        <section className="mb-3 rounded-xl border border-line bg-cream-2 p-3">
          <p className="mb-1 font-hand text-sm text-sage">📝 メモ</p>
          <p className="text-[11px] leading-relaxed text-ink-2">{e.memo}</p>
        </section>
      )}

      {/* Tasks */}
      <section className="mb-3 rounded-xl border border-line bg-surface p-3">
        <p className="mb-2 text-[12px] font-bold">📌 タスク</p>
        {tasks.loading && <p className="text-[11px] text-ink-3">読み込み中…</p>}
        {tasks.data?.length === 0 && (
          <p className="text-[11px] text-ink-3">まだタスクがありません</p>
        )}
        {tasks.data && tasks.data.length > 0 && (
          <ul className="flex flex-col gap-1.5">
            {tasks.data.map((t) => (
              <li
                key={t.id}
                className={`flex items-center gap-2 text-[11px] ${
                  t.status === "done" ? "line-through text-ink-3" : ""
                }`}
              >
                <span
                  className={`h-3.5 w-3.5 rounded-full border-[1.5px] ${
                    t.status === "done" ? "border-sage bg-sage" : "border-line"
                  }`}
                />
                <span className="flex-1">{t.title}</span>
                {t.dueDate && <span className="font-mono text-[9px] text-ink-3">{t.dueDate}</span>}
              </li>
            ))}
          </ul>
        )}
      </section>

      {/* Mascot encouragement */}
      <div className="flex items-center gap-3 rounded-xl border-[1.5px] border-line bg-gradient-to-br from-cream-2 to-sage-wash p-4">
        <Mascot size={48} mood={isOffer ? "cheering" : "thinking"} />
        <p className="font-hand text-sm text-sage">
          {isOffer ? "おめでとう！本当にお疲れさま 🎉" : "次のステップ、応援してます。"}
        </p>
      </div>

      <Confetti trigger={confetti} count={28} />
    </div>
  );
}
