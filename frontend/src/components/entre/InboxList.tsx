// Inbox 一覧表示。Server Component で props だけ受け取る純粋表示コンポーネント。
// 相対時刻もサーバ側で計算するので SSR/CSR 不一致が起きない (Date.now() を Client で呼ばない)。

import { cache } from "react";
import { Mascot } from "./Mascot";
import type { InboxClipResponse } from "@/lib/api/inboxClips";

// React.cache で 1 リクエスト内 memoize。Date.now() 自体は impure だが、cache() で
// 包むことで「同一リクエストでは同じ値」を保証でき、components-and-hooks-must-be-pure 規則も満たす。
const getRenderedAt = cache(() => Date.now());

export function InboxList({ clips }: { clips: InboxClipResponse[] }) {
  if (clips.length === 0) return <EmptyState />;

  const renderedAt = getRenderedAt();

  return (
    <ul className="flex flex-col gap-2">
      {clips.map((c) => (
        <li
          key={c.id}
          className="flex cursor-pointer items-start gap-3 rounded-xl border border-line bg-surface p-3.5 transition-colors hover:border-sage"
        >
          <div className="grid h-9 w-9 shrink-0 place-items-center rounded-md bg-sage-wash text-base">
            ✉
          </div>
          <div className="min-w-0 flex-1">
            <div className="truncate text-[12px] font-bold">{c.title}</div>
            <div className="mt-0.5 truncate font-mono text-[10px] text-ink-3">
              {c.url}
            </div>
            <div className="mt-1 flex items-center gap-2 text-[10px] text-ink-2">
              <span className="rounded-sm bg-cream-2 px-1.5 py-0.5 font-bold">{c.source}</span>
              <span>{formatRelative(c.capturedAt, renderedAt)}</span>
              {c.guess && <span className="text-sage">→ {c.guess}</span>}
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
      <div style={{ animation: "entre-float 3s infinite" }}>
        <Mascot size={80} mood="sleeping" />
      </div>
      <p className="mt-3 font-serif text-base font-extrabold">クリップは空です</p>
      <p className="mt-1 text-[11px] text-ink-2">
        Chrome拡張で気になる求人ページを保存すると、ここに溜まります。
      </p>
    </div>
  );
}
