"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";

interface Clip {
  id: string;
  url: string;
  title: string;
  source: string;
  captured: string;
  guess: string;
}

const CLIPS: Clip[] = [
  { id: "c1", url: "recruit.example.com/jobs/3921", title: "【26卒】エンジニア職 本選考 | パネル製作所", source: "リクナビ", captured: "15分前", guess: "パネル製作所" },
  { id: "c2", url: "mypage.ats-i-web.com/candidate/...", title: "マイページ | 選考スケジュール", source: "i-web", captured: "1時間前", guess: "オリーブ商事（候補）" },
  { id: "c3", url: "one-career.example/report/interview-report", title: "面接レポート — ブリック出版 1次", source: "ONE CAREER", captured: "今朝", guess: "ブリック出版" },
];

export default function InboxPage() {
  const router = useRouter();
  const state = useUser();

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
          <span className="text-[20px]">
            <Mascot size={32} mood="wink" />
          </span>
        </header>

        {CLIPS.length === 0 ? (
          <EmptyState />
        ) : (
          <ul className="flex flex-col gap-2">
            {CLIPS.map((c) => (
              <li
                key={c.id}
                className="flex cursor-pointer items-start gap-3 rounded-xl border border-line bg-surface p-3.5 transition-colors hover:border-sage"
              >
                <div className="grid h-9 w-9 shrink-0 place-items-center rounded-md bg-sage-wash text-base">
                  ✉
                </div>
                <div className="min-w-0 flex-1">
                  <div className="truncate text-[12px] font-bold">{c.title}</div>
                  <div className="mt-0.5 truncate font-mono text-[10px] text-ink-3">{c.url}</div>
                  <div className="mt-1 flex items-center gap-2 text-[10px] text-ink-2">
                    <span className="rounded-sm bg-cream-2 px-1.5 py-0.5 font-bold">
                      {c.source}
                    </span>
                    <span>{c.captured}</span>
                    <span className="text-sage">→ {c.guess}</span>
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
