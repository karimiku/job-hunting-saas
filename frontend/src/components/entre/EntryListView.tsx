// Server Component で render される純粋表示コンポーネント。
// データは props で渡される (page.tsx 側で SSR 取得済み)。

import Link from "next/link";
import { ArrowRight, ExternalLink, Inbox, Plus } from "lucide-react";
import {
  companyDisplayName,
  entrySourceUrl,
  type EntryResponse,
} from "@/lib/api/entries";
import {
  ENTRY_STATUS_LABEL,
  STAGE_BG,
  STAGE_ORDER,
  stageIndexOf,
} from "@/lib/entry-stage";

export function EntryListView({ entries }: { entries: EntryResponse[] }) {
  if (entries.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-line bg-surface p-8 text-center">
        <p className="font-serif text-base font-extrabold">まだ Entry がありません</p>
        <p className="mx-auto mt-1 max-w-[420px] text-[11px] leading-relaxed text-ink-2">
          保存した求人は保存箱から Entry にできます。直接追加することもできます。
        </p>
        <div className="mt-4 flex flex-wrap justify-center gap-2">
          <Link
            href="/inbox"
            prefetch={false}
            className="inline-flex items-center gap-1.5 rounded-lg border border-sage bg-sage-wash px-3 py-1.5 text-[11px] font-bold text-sage transition-colors hover:bg-sage hover:text-white"
          >
            <Inbox size={13} aria-hidden />
            保存箱を見る
          </Link>
          <Link
            href="/entry/new"
            prefetch={false}
            className="inline-flex items-center gap-1.5 rounded-lg border border-line bg-surface px-3 py-1.5 text-[11px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
          >
            <Plus size={13} aria-hidden />
            Entryを追加
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div>
      <ul className="entre-stagger flex flex-col gap-2">
        {entries.map((e) => {
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
                    <span className="shrink-0 rounded-full bg-cream px-1.5 py-0.5 text-[8px] font-black text-ink-3">
                      {ENTRY_STATUS_LABEL[e.status] ?? e.status}
                    </span>
                  </div>
                  <div className="mt-0.5 flex min-w-0 items-center gap-1.5 text-[10px] text-ink-3">
                    <span
                      className="rounded-sm px-1.5 py-0.5 text-[8px] font-bold text-white"
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
                  {e.memo && <div className="mt-1 text-[10px] text-ink-2">{e.memo}</div>}
                </div>
                <span className="hidden shrink-0 items-center gap-1 rounded-md bg-cream px-2 py-1 text-[10px] font-bold text-ink-3 md:inline-flex">
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
    </div>
  );
}
