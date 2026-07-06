"use client";

// SSR で渡された entries を初期値として、検索・ステージ絞り込み・状態絞り込み・並び替えを
// クライアント側で行う (API・取得ロジックは変えない)。

import { useState } from "react";
import Link from "next/link";
import { ArrowRight, ExternalLink, Inbox, Plus } from "lucide-react";
import {
  companyDisplayName,
  entrySourceUrl,
  type EntryResponse,
} from "@/lib/api/entries";
import {
  ENTRY_STATUS_LABEL,
  isKanbanStageKind,
  KANBAN_STAGE_LABEL,
  KANBAN_STAGE_ORDER,
  kanbanStageIndexOf,
  STAGE_BG,
  STAGE_ORDER,
  stageIndexOf,
  type KanbanStageKind,
} from "@/lib/entry-stage";
import { EntryListToolbar } from "./EntryListToolbar";

function normalizedStageKind(entry: Pick<EntryResponse, "stageKind">): KanbanStageKind {
  return isKanbanStageKind(entry.stageKind) ? entry.stageKind : "other";
}

export function filterEntries(entries: EntryResponse[], query: string): EntryResponse[] {
  const q = query.trim().toLowerCase();
  if (!q) return entries;
  return entries.filter((e) =>
    [companyDisplayName(e), e.route, e.source, e.memo].some((value) =>
      (value ?? "").toLowerCase().includes(q),
    ),
  );
}

export function filterEntriesByStage(
  entries: EntryResponse[],
  stage: KanbanStageKind | null,
): EntryResponse[] {
  if (!stage) return entries;
  return entries.filter((e) => normalizedStageKind(e) === stage);
}

export type EntryStatusGroup = "all" | "in_progress" | "offer" | "closed";

// offered/accepted は「内定」、rejected/withdrawn は「終了」としてまとめる。
// 未知のステータスは終了扱い（選考中・内定のどちらでもない = 進行が止まっている）にする。
const STATUS_GROUP_OF: Record<string, EntryStatusGroup> = {
  in_progress: "in_progress",
  offered: "offer",
  accepted: "offer",
  rejected: "closed",
  withdrawn: "closed",
};

export function filterEntriesByStatusGroup(
  entries: EntryResponse[],
  group: EntryStatusGroup,
): EntryResponse[] {
  if (group === "all") return entries;
  return entries.filter((e) => (STATUS_GROUP_OF[e.status] ?? "closed") === group);
}

export type EntrySortOrder = "updated" | "company" | "stage";

export function sortEntries(entries: EntryResponse[], order: EntrySortOrder): EntryResponse[] {
  const sorted = [...entries];
  if (order === "company") {
    sorted.sort((a, b) => companyDisplayName(a).localeCompare(companyDisplayName(b), "ja"));
    return sorted;
  }
  if (order === "stage") {
    sorted.sort(
      (a, b) =>
        kanbanStageIndexOf(normalizedStageKind(a)) - kanbanStageIndexOf(normalizedStageKind(b)),
    );
    return sorted;
  }
  sorted.sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime());
  return sorted;
}

function presentStages(entries: EntryResponse[]): { value: KanbanStageKind; label: string }[] {
  const present = new Set(entries.map((e) => normalizedStageKind(e)));
  return KANBAN_STAGE_ORDER.filter((kind) => present.has(kind)).map((value) => ({
    value,
    label: KANBAN_STAGE_LABEL[value],
  }));
}

export function EntryListView({ entries }: { entries: EntryResponse[] }) {
  const [query, setQuery] = useState("");
  const [stage, setStage] = useState<KanbanStageKind | null>(null);
  const [statusGroup, setStatusGroup] = useState<EntryStatusGroup>("all");
  const [order, setOrder] = useState<EntrySortOrder>("updated");

  if (entries.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-line bg-surface p-8 text-center">
        <p className="font-serif text-base font-extrabold">まだ応募先がありません</p>
        <p className="mx-auto mt-1 max-w-[420px] text-[12px] leading-relaxed text-ink-2">
          保存した求人は保存箱から応募先にできます。直接追加もできます。
        </p>
        <div className="mt-4 flex flex-wrap justify-center gap-2">
          <Link
            href="/inbox"
            prefetch={false}
            className="inline-flex items-center gap-1.5 rounded-lg border border-sage bg-sage-wash px-3 py-1.5 text-[12px] font-bold text-sage transition-colors hover:bg-sage hover:text-white"
          >
            <Inbox size={13} aria-hidden />
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
    );
  }

  const visible = sortEntries(
    filterEntriesByStatusGroup(filterEntriesByStage(filterEntries(entries, query), stage), statusGroup),
    order,
  );

  return (
    <div>
      {entries.length > 1 && (
        <EntryListToolbar
          query={query}
          onQueryChange={setQuery}
          order={order}
          onOrderChange={setOrder}
          stages={presentStages(entries)}
          selectedStage={stage}
          onStageChange={setStage}
          statusGroup={statusGroup}
          onStatusGroupChange={setStatusGroup}
        />
      )}

      {visible.length === 0 ? (
        <div className="rounded-xl border border-dashed border-line bg-surface p-8 text-center text-[12px] text-ink-2">
          条件に一致する応募先がありません。
        </div>
      ) : (
        <ul className="entre-stagger flex flex-col gap-2">
          {visible.map((e) => {
            const sourceUrl = entrySourceUrl(e);
            return (
              <li
                key={e.id}
                className="flex items-center gap-2.5 rounded-xl border border-line bg-surface p-3 transition-all hover:translate-x-0.5 hover:border-sage"
              >
                <Link
                  href={`/entry/${e.id}`}
                  prefetch={false}
                  className="flex min-w-0 flex-1 items-start gap-2.5"
                >
                  <div className="min-w-0 flex-1">
                    <div className="flex min-w-0 items-center gap-1.5">
                      <div className="truncate text-[12px] font-bold">{companyDisplayName(e)}</div>
                      {e.status !== "in_progress" && ENTRY_STATUS_LABEL[e.status] && (
                        <span className="shrink-0 rounded-full bg-cream px-1.5 py-0.5 text-[11px] font-black text-ink-3">
                          {ENTRY_STATUS_LABEL[e.status]}
                        </span>
                      )}
                    </div>
                    <div className="mt-0.5 flex min-w-0 items-center gap-1.5 text-[12px] text-ink-3">
                      <span
                        className="rounded-sm px-1.5 py-0.5 text-[11px] font-bold text-white"
                        style={{ background: STAGE_BG[e.stageKind] ?? "var(--color-ink-3)" }}
                      >
                        {e.stageLabel}
                      </span>
                      <span>{e.route}</span>
                      <span aria-hidden>·</span>
                      <span className="truncate">{e.source}</span>
                    </div>
                    <div
                      className="mt-2 grid gap-0.5"
                      style={{ gridTemplateColumns: `repeat(${STAGE_ORDER.length}, minmax(0, 1fr))` }}
                      aria-label={`選考ステージ ${e.stageLabel}`}
                    >
                      {STAGE_ORDER.map((kind, index) => {
                        const reached = index <= stageIndexOf(e.stageKind);
                        return (
                          <span
                            key={kind}
                            className="h-1.5 rounded-full"
                            style={{
                              background: reached ? STAGE_BG[kind] : "var(--color-line-2)",
                            }}
                          />
                        );
                      })}
                    </div>
                    <p className="mt-1 text-[12px] font-bold text-ink-3">
                      {e.stageLabel} ・ {STAGE_ORDER.length}ステップ中{stageIndexOf(e.stageKind) + 1}
                    </p>
                    {e.memo && <div className="mt-1 text-[12px] text-ink-2">{e.memo}</div>}
                  </div>
                  <span className="hidden shrink-0 items-center gap-1 rounded-md bg-cream px-2 py-1 text-[12px] font-bold text-ink-3 md:inline-flex">
                    詳細を見る
                    <ArrowRight size={12} aria-hidden />
                  </span>
                </Link>
                {sourceUrl && (
                  <a
                    href={sourceUrl}
                    target="_blank"
                    rel="noreferrer"
                    aria-label={`${companyDisplayName(e)} の応募元ページを開く`}
                    className="grid h-8 w-8 shrink-0 place-items-center rounded-md border border-line text-ink-3 transition-colors hover:border-sage hover:text-sage"
                  >
                    <ExternalLink size={13} aria-hidden />
                  </a>
                )}
                <span className="text-ink-3" aria-hidden>›</span>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
