"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useRef, useState } from "react";
import { AlertCircle, ExternalLink, Plus } from "lucide-react";
import {
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  PointerSensor,
  useDraggable,
  useDroppable,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  companyDisplayName,
  entrySourceUrl,
  type EntryResponse,
} from "@/lib/api/entries";
import { updateEntryAction } from "@/app/entry/actions";
import {
  KANBAN_STAGE_COLOR,
  KANBAN_STAGE_LABEL,
  KANBAN_STAGE_ORDER,
  isKanbanStageKind,
  statusForKanbanStage,
  type KanbanStageKind,
} from "@/lib/entry-stage";

const COLUMNS = KANBAN_STAGE_ORDER.map((kind) => ({
  kind,
  label: KANBAN_STAGE_LABEL[kind],
  color: KANBAN_STAGE_COLOR[kind],
}));

interface Props {
  initialEntries: EntryResponse[];
}

type EntryOverride = Pick<EntryResponse, "stageKind" | "stageLabel" | "status">;

/** Entry を stageKind ごとに振り分けるカンバン。
 *  initialEntries は SSR で取得した snapshot。ドラッグで列間移動 → Server Action 経由で
 *  PATCH /entries/{id} を永続化、楽観的更新 (overrides) で UI を即時反映する。
 *  Action 内の revalidatePath("/kanban") がレスポンスに更新済み RSC ツリーを含めるため
 *  router.refresh() は不要 (overrides は成功時に clear)。
 *  API 失敗時は overrides を消してロールバック。 */
export function KanbanBoard({ initialEntries }: Props) {
  // entryId → 楽観的に上書きされたステージ情報。API 成功で消し、失敗で消してロールバック。
  const [overrides, setOverrides] = useState<Record<string, EntryOverride>>({});
  const [activeId, setActiveId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  // 同一カード連続移動の race 対策: entryId ごとに最新リクエスト ID を採番。
  // 自分が最新でないリクエストは override clear / refresh をスキップする。
  const requestIdsRef = useRef<Map<string, number>>(new Map());

  // 5px 以上動かさないとドラッグ開始しない（クリックでカード詳細に飛べるように）
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
  );

  // overrides を適用した entries (initialEntries が SSR で来ている前提)
  const entries = initialEntries.map((e) => ({
    ...e,
    ...overrides[e.id],
  }));

  const byKind = new Map<KanbanStageKind, EntryResponse[]>();
  for (const e of entries) {
    const key = normalizeKanbanStageKind(e.stageKind);
    if (!byKind.has(key)) byKind.set(key, []);
    byKind.get(key)!.push(e);
  }

  const activeEntry = activeId ? entries.find((e) => e.id === activeId) ?? null : null;

  function handleDragStart(event: DragStartEvent) {
    setActiveId(String(event.active.id));
  }

  async function handleDragEnd(event: DragEndEvent) {
    setActiveId(null);
    const { active, over } = event;
    if (!over) return;

    const entryId = String(active.id);
    const targetKind = normalizeKanbanStageKind(String(over.id));
    const current = entries.find((e) => e.id === entryId);
    if (!current || normalizeKanbanStageKind(current.stageKind) === targetKind) return;
    const next = kanbanStageUpdateInput(targetKind);

    // この entry に対する最新リクエスト番号を採番。
    // 後続のドラッグで番号が更新されたら、自分は古いリクエストとして clear/refetch をスキップする。
    const requestId = (requestIdsRef.current.get(entryId) ?? 0) + 1;
    requestIdsRef.current.set(entryId, requestId);

    // 楽観的反映
    setError(null);
    setOverrides((prev) => ({ ...prev, [entryId]: next }));

    const result = await updateEntryAction(entryId, next);
    // 自分が最新のリクエストのときだけ override を clear / ロールバックする。
    // そうでない場合は後続のドラッグが進行中なので、その完了に任せる。
    if (requestIdsRef.current.get(entryId) !== requestId) return;
    if (!result.ok) {
      setError("選考フェーズの更新に失敗しました。通信状態を確認して、もう一度動かしてください。");
    }
    // 成功時は action レスポンスの revalidate 済みツリーが initialEntries を更新しているので
    // override を消しても表示は維持される。
    setOverrides((prev) => {
      const next = { ...prev };
      delete next[entryId];
      return next;
    });
  }

  return (
    <DndContext
      id="entre-kanban"
      sensors={sensors}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      {initialEntries.length === 0 && (
        <div className="mb-3 rounded-xl border border-dashed border-line bg-surface p-6 text-center">
          <p className="font-serif text-base font-extrabold">ボードに表示する応募先がありません</p>
          <p className="mx-auto mt-1 max-w-[460px] text-[12px] leading-relaxed text-ink-2">
            保存箱のクリップを応募先にするか、手動で応募先を追加すると選考フェーズごとに並びます。
          </p>
          <div className="mt-4 flex flex-wrap justify-center gap-2">
            <Link
              href="/inbox"
              prefetch={false}
              className="rounded-lg border border-sage bg-sage-wash px-3 py-1.5 text-[12px] font-bold text-sage transition-colors hover:bg-sage hover:text-white"
            >
              保存箱を見る
            </Link>
            <Link
              href="/entry/new"
              prefetch={false}
              className="inline-flex items-center gap-1.5 rounded-lg border border-line bg-surface px-3 py-1.5 text-[12px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
            >
              <Plus size={13} aria-hidden />
              応募先を追加
            </Link>
          </div>
        </div>
      )}

      {error && (
        <p
          role="alert"
          className="mb-3 inline-flex items-center gap-1.5 rounded-lg bg-pink/40 px-3 py-2 text-[12px] font-semibold text-ink"
        >
          <AlertCircle size={13} aria-hidden />
          {error}
        </p>
      )}

      <div data-testid="kanban-mobile-list" className="flex flex-col gap-2.5 md:hidden">
        {COLUMNS.map((col) => (
          <MobileKanbanSection key={col.kind} col={col} cards={byKind.get(col.kind) ?? []} />
        ))}
      </div>

      <div
        data-testid="kanban-desktop-board"
        className="hidden gap-2.5 overflow-x-auto pb-2 md:grid"
        style={{ gridTemplateColumns: `repeat(${COLUMNS.length}, minmax(10rem, 1fr))` }}
      >
        {COLUMNS.map((col) => (
          <KanbanColumn key={col.kind} col={col} cards={byKind.get(col.kind) ?? []} activeId={activeId} />
        ))}
      </div>

      {/* DragOverlay は親の transform / overflow に影響されない portal 的な場所で描画される。 */}
      <DragOverlay dropAnimation={null}>
        {activeEntry ? <KanbanCardPreview entry={activeEntry} /> : null}
      </DragOverlay>
    </DndContext>
  );
}

function MobileKanbanSection({
  col,
  cards,
}: {
  col: (typeof COLUMNS)[number];
  cards: EntryResponse[];
}) {
  return (
    <section className="rounded-xl border border-line bg-surface p-3">
      <div className="mb-2 flex items-center gap-2">
        <span className="block h-2 w-2 rounded-full" style={{ background: col.color }} />
        <h2 className="text-[12px] font-extrabold">{col.label}</h2>
        <span className="ml-auto rounded-md bg-cream px-2 py-0.5 font-mono text-[12px] font-bold text-ink-3">
          {cards.length}
        </span>
      </div>
      {cards.length === 0 ? (
        <p className="rounded-md border border-dashed border-line bg-cream px-3 py-3 text-center text-[12px] text-ink-3">
          まだこの段階の応募先はありません
        </p>
      ) : (
        <ul className="flex flex-col gap-1.5">
          {cards.map((entry) => (
            <li key={entry.id}>
              <Link
                href={`/entry/${entry.id}`}
                prefetch={false}
                className="block rounded-lg border border-line bg-cream p-2.5 transition-colors hover:border-sage"
              >
                <CardContent entry={entry} showSourceLink={false} />
              </Link>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}

function KanbanColumn({
  col,
  cards,
  activeId,
}: {
  col: (typeof COLUMNS)[number];
  cards: EntryResponse[];
  activeId: string | null;
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
        <span className="text-[12px] font-extrabold">{col.label}</span>
        <span
          data-testid={`column-count-${col.kind}`}
          className="ml-auto font-mono text-[12px] text-ink-3"
        >
          {cards.length}
        </span>
      </div>
      <ul className="flex flex-col gap-1.5 min-h-[60px]">
        {cards.map((c) => (
          <KanbanCard key={c.id} entry={c} dragging={c.id === activeId} />
        ))}
        {cards.length === 0 && (
          <li className="rounded-md border border-dashed border-line p-2 text-center text-[12px] text-ink-3">
            まだこの段階の応募先はありません
          </li>
        )}
      </ul>
    </div>
  );
}

function KanbanCard({ entry, dragging }: { entry: EntryResponse; dragging: boolean }) {
  const router = useRouter();
  const { attributes, listeners, setNodeRef } = useDraggable({ id: entry.id });

  // クリックは詳細遷移、ドラッグはカード移動。
  // PointerSensor の activationConstraint で 5px 未満は click として扱われる。
  const handleClick = () => {
    if (dragging) return;
    router.push(`/entry/${entry.id}`);
  };

  return (
    <li
      ref={setNodeRef}
      onClick={handleClick}
      onKeyDown={(e) => {
        if (e.key === "Enter") handleClick();
      }}
      {...attributes}
      {...listeners}
      style={{
        // ドラッグ中はカード位置の元の場所に半透明スロットを残す（DragOverlay でゴーストが追従）
        opacity: dragging ? 0.3 : 1,
        cursor: "grab",
        touchAction: "none",
      }}
      className="rounded-[10px] border border-line bg-cream p-2.5 transition-shadow hover:shadow-[0_6px_14px_-4px_rgba(0,0,0,0.15)] focus:outline-none focus:ring-2 focus:ring-sage"
    >
      <CardContent entry={entry} />
    </li>
  );
}

/** DragOverlay 内に描画される、ポインタ追従のゴーストカード。
 *  親 transform に影響されないので、どこにドラッグしてもズレない。 */
function KanbanCardPreview({ entry }: { entry: EntryResponse }) {
  return (
    <div
      style={{ cursor: "grabbing" }}
      className="rounded-[10px] border-[1.5px] border-sage bg-cream p-2.5 shadow-[0_12px_24px_-8px_rgba(79,110,88,0.4)]"
    >
      <CardContent entry={entry} />
    </div>
  );
}

function CardContent({
  entry,
  showSourceLink = true,
}: {
  entry: EntryResponse;
  showSourceLink?: boolean;
}) {
  const sourceUrl = entrySourceUrl(entry);
  const stageKind = normalizeKanbanStageKind(entry.stageKind);
  const stageLabel = entry.stageLabel || KANBAN_STAGE_LABEL[stageKind];
  return (
    <>
      <div className="mb-1.5 truncate text-[12px] font-bold">{companyDisplayName(entry)}</div>
      <div className="mb-1.5 flex min-w-0">
        <span
          className="max-w-full truncate rounded border border-line bg-surface px-1.5 py-0.5 text-[12px] font-bold text-ink-3"
          style={{ borderColor: KANBAN_STAGE_COLOR[stageKind] }}
        >
          {stageLabel}
        </span>
      </div>
      <div className="flex justify-between gap-2 text-[12px] text-ink-3">
        <span className="truncate">
          {entry.route} · {entry.source}
        </span>
        <span aria-hidden>⇆</span>
      </div>
      {showSourceLink && sourceUrl && (
        <a
          href={sourceUrl}
          target="_blank"
          rel="noreferrer"
          onClick={(event) => event.stopPropagation()}
          onPointerDown={(event) => event.stopPropagation()}
          className="mt-1.5 inline-flex max-w-full items-center gap-1 rounded border border-line bg-surface px-1.5 py-0.5 text-[12px] font-bold text-ink-3 transition-colors hover:border-sage hover:text-sage"
        >
          <span className="truncate">応募元</span>
          <ExternalLink size={10} className="shrink-0" aria-hidden />
        </a>
      )}
    </>
  );
}

export function normalizeKanbanStageKind(value: string): KanbanStageKind {
  return isKanbanStageKind(value) ? value : "other";
}

export function kanbanStageUpdateInput(stageKind: KanbanStageKind): EntryOverride {
  return {
    stageKind,
    stageLabel: KANBAN_STAGE_LABEL[stageKind],
    status: statusForKanbanStage(stageKind),
  };
}
