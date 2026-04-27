"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import {
  DndContext,
  DragEndEvent,
  PointerSensor,
  useDraggable,
  useDroppable,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { CSS } from "@dnd-kit/utilities";
import { useEntries } from "@/hooks/useEntries";
import { updateEntry, type EntryResponse } from "@/lib/api/entries";
import { Reveal } from "./Reveal";

const COLUMNS = [
  { kind: "application", label: "エントリー", color: "var(--color-stage-entry)" },
  { kind: "document", label: "書類選考", color: "var(--color-stage-doc)" },
  { kind: "test", label: "テスト/ES", color: "var(--color-stage-es)" },
  { kind: "interview", label: "面接", color: "var(--color-stage-interview)" },
  { kind: "offer", label: "内定", color: "var(--color-stage-offer)" },
] as const;

const STAGE_LABEL: Record<string, string> = {
  application: "エントリー",
  document: "書類選考",
  test: "ES提出",
  interview: "面接",
  offer: "内定",
};

/** Entry を stageKind ごとに振り分けるカンバン。
 *  カードはドラッグで列間移動でき、ドロップ時に PATCH /entries/{id} で永続化する。
 *  楽観的更新 (overrides) で UI を即時反映、API 失敗時はロールバック。 */
export function KanbanBoard() {
  const { data, loading, error, refetch } = useEntries();
  // entryId → 楽観的に上書きされた stageKind。API 成功で消し、失敗で消してロールバック。
  const [overrides, setOverrides] = useState<Record<string, string>>({});

  // 5px 以上動かさないとドラッグ開始しない（クリックでカード詳細に飛べるように）
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
  );

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

  // overrides を適用した entries
  const entries = (data ?? []).map((e) => ({
    ...e,
    stageKind: overrides[e.id] ?? e.stageKind,
  }));

  // group は interview 列に寄せる
  const byKind = new Map<string, EntryResponse[]>();
  for (const e of entries) {
    const key = e.stageKind === "group" ? "interview" : e.stageKind;
    if (!byKind.has(key)) byKind.set(key, []);
    byKind.get(key)!.push(e);
  }

  async function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (!over) return;

    const entryId = String(active.id);
    const targetKind = String(over.id);
    const current = entries.find((e) => e.id === entryId);
    if (!current || current.stageKind === targetKind) return;

    // 楽観的反映
    setOverrides((prev) => ({ ...prev, [entryId]: targetKind }));

    try {
      await updateEntry(entryId, {
        stageKind: targetKind,
        stageLabel: STAGE_LABEL[targetKind] ?? targetKind,
      });
      // 成功 — refetch で server state に同期 + override をクリア
      refetch();
      setOverrides((prev) => {
        const next = { ...prev };
        delete next[entryId];
        return next;
      });
    } catch {
      // 失敗 — override 削除でロールバック
      setOverrides((prev) => {
        const next = { ...prev };
        delete next[entryId];
        return next;
      });
    }
  }

  return (
    <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
      <div className="grid gap-2.5 md:grid-cols-5 grid-cols-[repeat(5,minmax(220px,1fr))] overflow-x-auto pb-2">
        {COLUMNS.map((col, i) => (
          <Reveal key={col.kind} delay={i * 80}>
            <KanbanColumn col={col} cards={byKind.get(col.kind) ?? []} />
          </Reveal>
        ))}
      </div>
    </DndContext>
  );
}

function KanbanColumn({
  col,
  cards,
}: {
  col: (typeof COLUMNS)[number];
  cards: EntryResponse[];
}) {
  const { setNodeRef, isOver } = useDroppable({ id: col.kind });
  return (
    <div
      ref={setNodeRef}
      className={`flex h-full flex-col gap-2 rounded-xl border bg-surface p-2.5 transition-colors ${
        isOver ? "border-sage bg-sage-wash" : "border-line"
      }`}
    >
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
      <ul className="flex flex-col gap-1.5 min-h-[40px]">
        {cards.map((c) => (
          <KanbanCard key={c.id} entry={c} />
        ))}
        {cards.length === 0 && (
          <li className="rounded-md border border-dashed border-line p-2 text-center text-[9px] text-ink-3">
            ここに置く
          </li>
        )}
      </ul>
    </div>
  );
}

function KanbanCard({ entry }: { entry: EntryResponse }) {
  const router = useRouter();
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: entry.id,
  });

  const style = {
    transform: CSS.Translate.toString(transform),
    opacity: isDragging ? 0.4 : 1,
    zIndex: isDragging ? 50 : "auto",
    cursor: isDragging ? "grabbing" : "grab",
  } as const;

  // クリックは詳細遷移、ドラッグはカード移動。
  // PointerSensor の activationConstraint で 5px 未満は click として扱われる。
  // dnd-kit が drag 中の click を抑制してくれるので、単純に navigate するだけで良い。
  const handleClick = () => {
    if (isDragging) return;
    router.push(`/entry/${entry.id}`);
  };

  return (
    <li
      ref={setNodeRef}
      style={style}
      onClick={handleClick}
      onKeyDown={(e) => {
        if (e.key === "Enter") handleClick();
      }}
      {...attributes}
      {...listeners}
      className="rounded-[10px] border border-line bg-cream p-2.5 transition-shadow hover:shadow-[0_6px_14px_-4px_rgba(0,0,0,0.15)] focus:outline-none focus:ring-2 focus:ring-sage"
    >
      <div className="mb-1.5 flex items-center gap-2">
        <div className="grid h-6 w-6 place-items-center rounded-md bg-sage-wash font-serif text-xs font-extrabold text-sage">
          {entry.source.slice(0, 1)}
        </div>
        <div className="truncate text-[10px] font-bold">{entry.source}</div>
      </div>
      <div className="flex justify-between text-[9px] text-ink-3">
        <span>{entry.route}</span>
        <span aria-hidden>⇆</span>
      </div>
    </li>
  );
}
