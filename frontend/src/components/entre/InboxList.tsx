// Inbox 一覧表示。Server Component で props だけ受け取る純粋表示コンポーネント。
// 相対時刻もサーバ側で計算するので SSR/CSR 不一致が起きない (Date.now() を Client で呼ばない)。

import { cache } from "react";
import Link from "next/link";
import { ExternalLink, Inbox, Sparkles } from "lucide-react";
import { InboxClipConvert } from "./InboxClipConvert";
import { InboxClipDelete } from "./InboxClipDelete";
import type { InboxClipResponse } from "@/lib/api/inboxClips";
import type { CompanyResponse } from "@/lib/api/companies";

// React.cache で 1 リクエスト内 memoize。Date.now() 自体は impure だが、cache() で
// 包むことで「同一リクエストでは同じ値」を保証でき、components-and-hooks-must-be-pure 規則も満たす。
const getRenderedAt = cache(() => Date.now());

export function InboxList({
  clips,
  companies,
}: {
  clips: InboxClipResponse[];
  companies: CompanyResponse[];
}) {
  if (clips.length === 0) return <EmptyState />;

  const renderedAt = getRenderedAt();

  return (
    <ul className="entre-stagger flex flex-col gap-2">
      {clips.map((c) => (
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
              <div className="mt-1 flex items-center gap-2 text-[12px] text-ink-2">
                <span className="rounded-sm bg-cream-2 px-1.5 py-0.5 font-bold">{c.source}</span>
                <span>{formatRelative(c.capturedAt, renderedAt)}</span>
                {c.guess ? (
                  <span className="truncate text-sage">会社候補: {c.guess}</span>
                ) : (
                  <span className="text-amber-700">会社名未検出</span>
                )}
              </div>
            </div>
          </div>
          <div className="flex items-end justify-between gap-2">
            <InboxClipDelete clip={c} />
            <div className="min-w-0 flex-1">
              <InboxClipConvert clip={c} companies={companies} />
            </div>
          </div>
        </li>
      ))}
    </ul>
  );
}

function formatRelative(iso: string, now: number): string {
  const date = new Date(iso);
  const diffMin = Math.floor((now - date.getTime()) / 60000);
  if (diffMin < 1) return "たった今";
  if (diffMin < 60) return `${diffMin}分前`;
  if (diffMin < 60 * 24) return `${Math.floor(diffMin / 60)}時間前`;
  return `${Math.floor(diffMin / (60 * 24))}日前`;
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
