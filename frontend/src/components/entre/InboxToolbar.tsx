"use client";

// 保存箱の検索・並び替え・ソース絞り込みUI。状態は親 (InboxList) が持ち、ここは
// 値と onChange を受け取るだけの制御コンポーネント。

import { Search } from "lucide-react";

export type SortOrder = "new" | "old" | "company";

const SORT_OPTIONS: { value: SortOrder; label: string }[] = [
  { value: "new", label: "新しい順" },
  { value: "old", label: "古い順" },
  { value: "company", label: "会社名順" },
];

export function InboxToolbar({
  query,
  onQueryChange,
  order,
  onOrderChange,
  sources,
  selectedSource,
  onSourceChange,
}: {
  query: string;
  onQueryChange: (value: string) => void;
  order: SortOrder;
  onOrderChange: (value: SortOrder) => void;
  sources: string[];
  selectedSource: string | null;
  onSourceChange: (value: string | null) => void;
}) {
  return (
    <div className="mb-3 flex flex-col gap-2">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
        <label className="relative flex-1">
          <span className="sr-only">クリップを検索</span>
          <Search
            size={13}
            className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-ink-3"
            aria-hidden
          />
          <input
            type="search"
            value={query}
            onChange={(e) => onQueryChange(e.target.value)}
            placeholder="タイトル・URL・会社名で検索"
            className="h-9 w-full rounded-md border border-line bg-surface pl-8 pr-2.5 text-[12px] font-semibold outline-none transition-colors focus:border-sage focus:ring-2 focus:ring-sage/20"
          />
        </label>
        <label className="shrink-0">
          <span className="sr-only">並び替え</span>
          <select
            value={order}
            onChange={(e) => onOrderChange(e.target.value as SortOrder)}
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

      {sources.length > 1 && (
        <div className="-mx-5 overflow-x-auto px-5 md:mx-0 md:px-0">
          <div className="flex min-w-max gap-1.5 pb-1">
            <button
              type="button"
              onClick={() => onSourceChange(null)}
              aria-pressed={selectedSource === null}
              className={`h-8 rounded-full border px-3 text-[12px] font-black transition-colors ${
                selectedSource === null
                  ? "border-sage bg-sage text-white"
                  : "border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
              }`}
            >
              すべて
            </button>
            {sources.map((sourceOption) => {
              const selected = selectedSource === sourceOption;
              return (
                <button
                  key={sourceOption}
                  type="button"
                  onClick={() => onSourceChange(sourceOption)}
                  aria-pressed={selected}
                  className={`h-8 rounded-full border px-3 text-[12px] font-black transition-colors ${
                    selected
                      ? "border-sage bg-sage text-white"
                      : "border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
                  }`}
                >
                  {sourceOption}
                </button>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
