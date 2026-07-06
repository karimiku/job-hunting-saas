"use client";

// Inbox 一覧表示。SSR で渡された clips を初期値として、検索・並び替え・ソース絞り込みを
// クライアント側で行う (API・取得ロジックは変えない)。相対時刻は初回レンダー時刻を
// useState の lazy initializer で1回だけ固定し、以後の再レンダーで揺れないようにする。

import { useState } from "react";
import Link from "next/link";
import { ExternalLink, Inbox, Sparkles } from "lucide-react";
import { InboxClipConvert } from "./InboxClipConvert";
import { InboxClipDelete } from "./InboxClipDelete";
import { InboxToolbar, type SortOrder } from "./InboxToolbar";
import type { InboxClipResponse } from "@/lib/api/inboxClips";
import type { CompanyResponse } from "@/lib/api/companies";

// InboxClipConvert.normalizeCompanyName と同じ考え方 (全角/空白/株式会社等を除去して比較)。
function normalizeCompanyName(value: string): string {
  return value.replace(/[株式会社（）()・\s　]/g, "").toLowerCase();
}

export function filterClips(clips: InboxClipResponse[], query: string): InboxClipResponse[] {
  const q = query.trim().toLowerCase();
  if (!q) return clips;
  return clips.filter((c) =>
    [c.title, c.url, c.guess].some((value) => (value ?? "").toLowerCase().includes(q)),
  );
}

export function filterClipsBySource(
  clips: InboxClipResponse[],
  source: string | null,
): InboxClipResponse[] {
  if (!source) return clips;
  return clips.filter((c) => c.source === source);
}

export function sortClips(clips: InboxClipResponse[], order: SortOrder): InboxClipResponse[] {
  const sorted = [...clips];
  if (order === "company") {
    sorted.sort((a, b) =>
      normalizeCompanyName(a.guess || a.title).localeCompare(
        normalizeCompanyName(b.guess || b.title),
        "ja",
      ),
    );
    return sorted;
  }
  sorted.sort((a, b) => {
    const diff = new Date(a.capturedAt).getTime() - new Date(b.capturedAt).getTime();
    return order === "old" ? diff : -diff;
  });
  return sorted;
}

// 正規化した会社候補ごとにクリップをグルーピングし、2件以上あるものだけ返す。
export function findDuplicateGroups(
  clips: InboxClipResponse[],
): Map<string, InboxClipResponse[]> {
  const groups = new Map<string, InboxClipResponse[]>();
  for (const clip of clips) {
    const guess = clip.guess?.trim();
    if (!guess) continue;
    const key = normalizeCompanyName(guess);
    if (!key) continue;
    groups.set(key, [...(groups.get(key) ?? []), clip]);
  }
  for (const [key, group] of groups) {
    if (group.length < 2) groups.delete(key);
  }
  return groups;
}

function uniqueSources(clips: InboxClipResponse[]): string[] {
  const seen = new Set<string>();
  const result: string[] = [];
  for (const clip of clips) {
    if (!clip.source || seen.has(clip.source)) continue;
    seen.add(clip.source);
    result.push(clip.source);
  }
  return result;
}

export function InboxList({
  clips,
  companies,
}: {
  clips: InboxClipResponse[];
  companies: CompanyResponse[];
}) {
  const [renderedAt] = useState(() => Date.now());
  const [query, setQuery] = useState("");
  const [order, setOrder] = useState<SortOrder>("new");
  const [source, setSource] = useState<string | null>(null);

  if (clips.length === 0) return <EmptyState />;

  const duplicateGroups = findDuplicateGroups(clips);
  const sources = uniqueSources(clips);
  const visible = sortClips(filterClips(filterClipsBySource(clips, source), query), order);

  return (
    <div>
      {clips.length > 1 && (
        <InboxToolbar
          query={query}
          onQueryChange={setQuery}
          order={order}
          onOrderChange={setOrder}
          sources={sources}
          selectedSource={source}
          onSourceChange={setSource}
        />
      )}

      {visible.length === 0 ? (
        <div className="rounded-xl border border-dashed border-line bg-surface p-8 text-center text-[12px] text-ink-2">
          条件に一致するクリップがありません。
        </div>
      ) : (
        <ul className="entre-stagger flex flex-col gap-2">
          {visible.map((c) => {
            const normalizedGuess = c.guess?.trim() ? normalizeCompanyName(c.guess) : "";
            const duplicateCount = normalizedGuess
              ? (duplicateGroups.get(normalizedGuess)?.length ?? 1) - 1
              : 0;

            return (
              <li
                key={c.id}
                className="flex flex-col gap-3 rounded-xl border border-line bg-surface p-3.5 transition-colors hover:border-sage"
              >
                <div className="flex items-start gap-3">
                  <div className="grid h-9 w-9 shrink-0 place-items-center rounded-md bg-sage-wash text-sage">
                    <Inbox size={17} aria-hidden />
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="line-clamp-2 text-[12px] font-bold leading-snug">{c.title}</div>
                    <a
                      href={c.url}
                      target="_blank"
                      rel="noreferrer"
                      className="mt-1 inline-flex max-w-full items-center gap-1 font-mono text-[12px] text-ink-3 transition-colors hover:text-sage"
                    >
                      <span className="truncate">{c.url}</span>
                      <ExternalLink size={10} className="shrink-0" aria-hidden />
                    </a>
                    <div className="mt-1 flex flex-wrap items-center gap-x-2 gap-y-1 text-[12px] text-ink-2">
                      <span className="rounded-sm bg-cream-2 px-1.5 py-0.5 font-bold">{c.source}</span>
                      <span title={new Date(c.capturedAt).toLocaleString("ja-JP")}>
                        {formatRelative(c.capturedAt, renderedAt)}
                      </span>
                      {c.guess ? (
                        <span className="break-words text-sage">会社候補: {c.guess}</span>
                      ) : (
                        <span className="text-amber-700">会社名未検出</span>
                      )}
                    </div>
                    {duplicateCount > 0 && (
                      <p className="mt-1 text-[12px] font-semibold text-amber-700">
                        同じ会社の候補が他に{duplicateCount}件あります
                      </p>
                    )}
                  </div>
                </div>
                <div className="flex items-end justify-between gap-2">
                  <InboxClipDelete clip={c} />
                  <div className="min-w-0 flex-1">
                    <InboxClipConvert clip={c} companies={companies} />
                  </div>
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}

function formatRelative(iso: string, now: number): string {
  const date = new Date(iso);
  const diffMin = Math.floor((now - date.getTime()) / 60000);
  if (diffMin < 1) return "たった今";
  if (diffMin < 60) return `${diffMin}分前`;
  if (diffMin < 60 * 24) return `${Math.floor(diffMin / 60)}時間前`;
  if (diffMin < 60 * 24 * 7) return `${Math.floor(diffMin / (60 * 24))}日前`;

  const sameYear = date.getFullYear() === new Date(now).getFullYear();
  return sameYear
    ? `${date.getMonth() + 1}/${date.getDate()}`
    : `${date.getFullYear()}/${date.getMonth() + 1}/${date.getDate()}`;
}

function EmptyState() {
  return (
    <div className="flex flex-col items-center rounded-xl border border-dashed border-line bg-surface p-10 text-center">
      <div className="grid h-14 w-14 place-items-center rounded-xl bg-sage-wash text-sage">
        <Inbox size={24} aria-hidden />
      </div>
      <p className="mt-3 font-serif text-base font-extrabold">クリップは空です</p>
      <p className="mt-1 text-[12px] text-ink-2">
        保存した求人は、ここで会社名を確認して応募先にできます。
      </p>
      <div className="mt-4 flex flex-wrap justify-center gap-2">
        <Link
          href="/entry/new"
          prefetch={false}
          className="inline-flex items-center gap-1.5 rounded-lg border border-line bg-surface px-3 py-1.5 text-[12px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
        >
          <Sparkles size={13} aria-hidden />
          応募先を追加
        </Link>
      </div>
    </div>
  );
}
