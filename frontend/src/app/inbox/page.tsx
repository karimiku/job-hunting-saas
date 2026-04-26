"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";
import { useInboxClips } from "@/hooks/useInboxClips";

export default function InboxPage() {
  const router = useRouter();
  const state = useUser();
  const { data, loading, error } = useInboxClips();

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  return (
    <AppShell userName={state.user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[800px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-baseline justify-between">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">Inbox</h1>
            <p className="mt-0.5 text-[11px] text-ink-3">
              Chrome拡張から保存したクリップ
            </p>
          </div>
          <Mascot size={32} mood="wink" />
        </header>

        {loading && (
          <p role="status" className="text-[12px] text-ink-3">
            読み込み中…
          </p>
        )}

        {error && (
          <p role="alert" className="rounded-lg bg-pink/40 p-3 text-[12px] font-semibold text-ink">
            読み込みに失敗しました（{error.message}）
          </p>
        )}

        {!loading && !error && (data?.length ?? 0) === 0 && <EmptyState />}

        {!loading && !error && (data?.length ?? 0) > 0 && (
          <ul className="flex flex-col gap-2">
            {data!.map((c) => (
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
                    <RelativeTime iso={c.capturedAt} />
                    {c.guess && <span className="text-sage">→ {c.guess}</span>}
                  </div>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </AppShell>
  );
}

function RelativeTime({ iso }: { iso: string }) {
  const date = new Date(iso);
  const diffMs = Date.now() - date.getTime();
  const diffMin = Math.floor(diffMs / 60000);
  let label: string;
  if (diffMin < 1) label = "たった今";
  else if (diffMin < 60) label = `${diffMin}分前`;
  else if (diffMin < 60 * 24) label = `${Math.floor(diffMin / 60)}時間前`;
  else label = `${Math.floor(diffMin / (60 * 24))}日前`;
  return <span>{label}</span>;
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
