"use client";

// 応募先一覧の検索・ステージ絞り込み・状態絞り込み・並び替えUI。状態は親 (EntryListView) が持ち、
// ここは値と onChange を受け取るだけの制御コンポーネント。

import { Search } from "lucide-react";
import type { KanbanStageKind } from "@/lib/entry-stage";
import type { EntrySortOrder, EntryStatusGroup } from "./EntryListView";

const SORT_OPTIONS: { value: EntrySortOrder; label: string }[] = [
  { value: "updated", label: "更新が新しい順" },
  { value: "company", label: "会社名順" },
  { value: "stage", label: "ステージ順" },
];

const STATUS_GROUP_OPTIONS: { value: EntryStatusGroup; label: string }[] = [
  { value: "all", label: "すべて" },
  { value: "in_progress", label: "選考中" },
  { value: "offer", label: "内定" },
  { value: "closed", label: "終了" },
];

export function EntryListToolbar({
  query,
  onQueryChange,
  order,
  onOrderChange,
  stages,
  selectedStage,
  onStageChange,
  statusGroup,
  onStatusGroupChange,
}: {
  query: string;
  onQueryChange: (value: string) => void;
  order: EntrySortOrder;
  onOrderChange: (value: EntrySortOrder) => void;
  stages: { value: KanbanStageKind; label: string }[];
  selectedStage: KanbanStageKind | null;
  onStageChange: (value: KanbanStageKind | null) => void;
  statusGroup: EntryStatusGroup;
  onStatusGroupChange: (value: EntryStatusGroup) => void;
}) {
  return (
    <div className="mb-3 flex flex-col gap-2">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
        <label className="relative flex-1">
          <span className="sr-only">応募先を検索</span>
          <Search
            size={13}
            className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-ink-3"
            aria-hidden
          />
          <input
            type="search"
            value={query}
            onChange={(e) => onQueryChange(e.target.value)}
            placeholder="会社名・応募経路・媒体・メモで検索"
            className="h-9 w-full rounded-md border border-line bg-surface pl-8 pr-2.5 text-[12px] font-semibold outline-none transition-colors focus:border-sage focus:ring-2 focus:ring-sage/20"
          />
        </label>
        <label className="shrink-0">
          <span className="sr-only">並び替え</span>
          <select
            value={order}
            onChange={(e) => onOrderChange(e.target.value as EntrySortOrder)}
            className="h-9 rounded-md border border-line bg-surface px-2.5 text-[12px] font-bold text-ink-2 outline-none transition-colors focus:border-sage focus:ring-2 focus:ring-sage/20"
          >
            {SORT_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
      </div>

      <div className="-mx-5 overflow-x-auto px-5 md:mx-0 md:px-0">
        <div className="flex min-w-max gap-1.5 pb-1">
          {STATUS_GROUP_OPTIONS.map((option) => {
            const selected = statusGroup === option.value;
            return (
              <button
                key={option.value}
                type="button"
                onClick={() => onStatusGroupChange(option.value)}
                aria-pressed={selected}
                className={`h-8 rounded-full border px-3 text-[12px] font-black transition-colors ${
                  selected
                    ? "border-sage bg-sage text-white"
                    : "border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
                }`}
              >
                {option.label}
              </button>
            );
          })}
        </div>
      </div>

      {stages.length > 1 && (
        <div className="-mx-5 overflow-x-auto px-5 md:mx-0 md:px-0">
          <div className="flex min-w-max gap-1.5 pb-1">
            <button
              type="button"
              onClick={() => onStageChange(null)}
              aria-pressed={selectedStage === null}
              className={`h-8 rounded-full border px-3 text-[12px] font-black transition-colors ${
                selectedStage === null
                  ? "border-sage bg-sage text-white"
                  : "border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
              }`}
            >
              すべてのステージ
            </button>
            {stages.map((stage) => {
              const selected = selectedStage === stage.value;
              return (
                <button
                  key={stage.value}
                  type="button"
                  onClick={() => onStageChange(stage.value)}
                  aria-pressed={selected}
                  className={`h-8 rounded-full border px-3 text-[12px] font-black transition-colors ${
                    selected
                      ? "border-sage bg-sage text-white"
                      : "border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
                  }`}
                >
                  {stage.label}
                </button>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
