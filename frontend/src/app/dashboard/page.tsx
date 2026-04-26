"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { signOut } from "@/lib/auth";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";
import { CountUp } from "@/components/entre/CountUp";
import { Reveal } from "@/components/entre/Reveal";
import { STAGES, STAGE_COUNTS } from "@/lib/sample-data";

const STATS = [
  { v: 24, l: "エントリー数", c: "text-sage", sub: "+3 今週" },
  { v: 12, l: "今週の予定", c: "text-pink-deep", sub: "今日 3件" },
  { v: 18, l: "選考中", c: "text-amber", sub: "面接 5件" },
  { v: 2, l: "内定", c: "text-mint", sub: "オファー待ち" },
];

const QUESTS = [
  { t: "14:00", e: "○○商事 一次面接(Web)", s: "明日締切 ES確認", co: "bg-pink", done: false },
  { t: "17:00", e: "OBの田中さんへDM返信", s: "メッセージ準備済", co: "bg-sage", done: true },
  { t: "23:59", e: "△△株式会社 ESを提出", s: "最終チェック", co: "bg-amber", done: false },
  { t: "本日中", e: "◇◇テック SPI受験", s: "45分", co: "bg-sky", done: false },
];

export default function DashboardPage() {
  const router = useRouter();
  const state = useUser();

  useEffect(() => {
    if (state.status === "guest") {
      router.replace("/login");
    }
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  const user = state.user;
  const firstName = user.name.split(/[\s　]/)[0] || user.name;

  return (
    <AppShell userName={user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[1100px] px-5 py-6 md:px-8 md:py-8">
        {/* Greeting */}
        <header className="mb-5 flex flex-col gap-3 md:mb-6 md:flex-row md:items-baseline md:justify-between">
          <div className="animate-[entre-fade-in_0.6s_both]">
            <p
              className="font-hand text-[22px] text-sage md:text-2xl"
              style={{ transform: "rotate(-1.5deg)", display: "inline-block" }}
            >
              Welcome back,
            </p>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight md:text-[28px]">
              {firstName}
              <span className="ml-1 text-sm font-medium text-ink-2 md:text-base">
                さん、今日もお疲れさまです 🌱
              </span>
            </h1>
          </div>

          <div className="flex items-center gap-2">
            <div className="flex items-center gap-1.5 rounded-full bg-gradient-to-br from-cream-2 to-sage-wash px-3.5 py-1.5 text-[11px] font-bold">
              <span style={{ animation: "entre-wiggle 2s infinite" }}>🔥</span>
              連続 7日
            </div>
            <SignOutButton />
          </div>
        </header>

        {/* Stats */}
        <section className="entre-stagger mb-5 grid grid-cols-2 gap-3 md:mb-6 md:grid-cols-4 md:gap-3">
          {STATS.map((s, i) => (
            <div
              key={s.l}
              className="rounded-xl border border-line bg-surface px-4 py-3.5"
            >
              <div className="text-[11px] font-semibold text-ink-3">{s.l}</div>
              <div className={`mt-1.5 font-serif text-3xl font-extrabold leading-tight ${s.c}`}>
                <CountUp end={s.v} duration={900 + i * 100} />.
              </div>
              <div className="mt-1 text-[10px] text-ink-2">{s.sub}</div>
            </div>
          ))}
        </section>

        {/* Quest + Status */}
        <div className="grid grid-cols-1 gap-4 md:grid-cols-[1.4fr_1fr]">
          {/* Today's quest */}
          <Reveal delay={150}>
            <div className="rounded-xl border border-line bg-surface p-5">
              <div className="mb-3 flex items-baseline justify-between">
                <h2 className="text-[13px] font-extrabold">📌 今日のクエスト</h2>
                <Link href="/task" className="text-[10px] font-bold text-sage">
                  すべて見る →
                </Link>
              </div>
              <div className="mb-3.5 h-1.5 overflow-hidden rounded-sm bg-line">
                <div
                  className="h-full rounded-sm bg-gradient-to-r from-sage-mid to-sage transition-all duration-1000"
                  style={{ width: "25%" }}
                />
              </div>

              <ul>
                {QUESTS.map((r, i) => (
                  <li
                    key={r.e}
                    className={`flex items-center gap-3 py-2.5 ${
                      i ? "border-t border-dashed border-line" : ""
                    } ${r.done ? "opacity-50" : ""}`}
                  >
                    <span
                      className={`grid h-[18px] w-[18px] place-items-center rounded-full text-[10px] text-white ${
                        r.done ? "border-[1.5px] border-sage bg-sage" : "border-[1.5px] border-line bg-transparent"
                      }`}
                    >
                      {r.done ? "✓" : ""}
                    </span>
                    <div className="min-w-0 flex-1">
                      <div className={`text-xs font-semibold ${r.done ? "line-through" : ""}`}>{r.e}</div>
                      <div className="mt-0.5 text-[10px] text-ink-3">{r.s}</div>
                    </div>
                    <span
                      className={`shrink-0 rounded-md px-2 py-0.5 font-mono text-[10px] font-bold text-white ${r.co}`}
                    >
                      {r.t}
                    </span>
                  </li>
                ))}
              </ul>
            </div>
          </Reveal>

          {/* Status pie */}
          <Reveal delay={250}>
            <div className="rounded-xl border border-line bg-surface p-5">
              <h2 className="mb-3 text-[13px] font-extrabold">選考ステータス</h2>
              <div className="flex items-center gap-3.5">
                <StatusPie />
                <ul className="flex flex-1 flex-col gap-1.5 text-[11px]">
                  {STAGES.map((s) => (
                    <li key={s.key} className="flex items-center gap-2">
                      <span
                        className="block h-2.5 w-2.5 rounded-sm"
                        style={{ background: s.color }}
                      />
                      <span className="flex-1">{s.label}</span>
                      <span className="font-mono font-bold">{STAGE_COUNTS[s.key]}</span>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </Reveal>
        </div>

        {/* Mascot encouragement */}
        <Reveal delay={350}>
          <div className="mt-4 flex flex-col items-start gap-4 rounded-xl border-[1.5px] border-line bg-gradient-to-br from-cream-2 to-sage-wash p-5 md:flex-row md:items-center md:p-6">
            <div style={{ animation: "entre-float 3s infinite" }}>
              <Mascot size={64} mood="cheering" />
            </div>
            <div className="flex-1">
              <p className="font-hand text-[18px] text-sage">あと少しですね！</p>
              <p className="mt-0.5 font-serif text-base font-extrabold">
                面接5社、内定まであと一歩。
              </p>
              <p className="mt-1 text-[11px] text-ink-2">
                今日のクエスト、お疲れさまです。明日の○○商事の一次面接、応援しています！
              </p>
            </div>
            <Link
              href="/roadmap"
              className="rounded-lg bg-sage px-3.5 py-2 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5"
            >
              ロードマップ →
            </Link>
          </div>
        </Reveal>
      </div>
    </AppShell>
  );
}

/** 円グラフ — ステータス別件数の可視化。 */
function StatusPie() {
  // 8 / 4 / 5 / 5 / 2 = 24 total
  const segments = [
    { stroke: "var(--color-stage-entry)", dasharray: "40 113", offset: 0 },
    { stroke: "var(--color-stage-doc)", dasharray: "20 113", offset: -40 },
    { stroke: "var(--color-stage-es)", dasharray: "25 113", offset: -60 },
    { stroke: "var(--color-stage-interview)", dasharray: "25 113", offset: -85 },
    { stroke: "var(--color-stage-offer)", dasharray: "10 113", offset: -110 },
  ];
  return (
    <svg width="100" height="100" viewBox="0 0 50 50" aria-label="ステータス別の選考件数">
      <circle cx="25" cy="25" r="18" fill="none" stroke="var(--color-line)" strokeWidth="6" />
      {segments.map((s, i) => (
        <circle
          key={i}
          cx="25"
          cy="25"
          r="18"
          fill="none"
          stroke={s.stroke}
          strokeWidth="6"
          strokeDasharray={s.dasharray}
          strokeDashoffset={s.offset}
          transform="rotate(-90 25 25)"
        />
      ))}
    </svg>
  );
}

function SignOutButton() {
  const router = useRouter();
  return (
    <button
      type="button"
      onClick={async () => {
        await signOut();
        router.push("/login");
      }}
      className="rounded-md border border-line bg-surface px-3 py-1.5 text-[11px] font-semibold text-ink-2 transition-colors hover:bg-line-2"
    >
      ログアウト
    </button>
  );
}
