"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";
import { Reveal } from "@/components/entre/Reveal";
import { CountUp } from "@/components/entre/CountUp";
import { ENTRIES, MILESTONES } from "@/lib/sample-data";

export default function RoadmapPage() {
  const router = useRouter();
  const state = useUser();

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  const interviewing = ENTRIES.filter((e) => e.stageIdx === 3);

  return (
    <AppShell userName={state.user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[1100px] px-5 py-7 md:px-9 md:py-9">
        <header className="mb-6 md:mb-7">
          <p
            className="font-hand text-[22px] text-sage md:text-[26px]"
            style={{ transform: "rotate(-1.5deg)", display: "inline-block" }}
          >
            your journey,
          </p>
          <h1 className="font-serif text-3xl font-extrabold tracking-tight md:text-[34px]">
            就活ロードマップ
          </h1>
          <p className="mt-1 text-[12px] text-ink-2">
            内定までの道のりを、一緒に歩いていきましょう。
          </p>
        </header>

        {/* Overall progress (mobile-friendly summary) */}
        <Reveal>
          <div className="mb-6 rounded-xl border border-line bg-surface p-5">
            <div className="mb-2 flex items-baseline justify-between">
              <p className="text-[12px] font-bold">今のあなたの進捗</p>
              <p className="font-serif text-2xl font-extrabold text-sage">
                <CountUp end={62} duration={1100} />
                <span className="text-xs">%</span>
              </p>
            </div>
            <div className="h-2.5 overflow-hidden rounded-md bg-line">
              <div
                className="h-full rounded-md bg-gradient-to-r from-sage-mid to-sage transition-[width] duration-1000"
                style={{ width: "62%" }}
              />
            </div>
            <p className="mt-2 font-hand text-[14px] text-sage">あと少しですね、ファイト！</p>
          </div>
        </Reveal>

        {/* Horizontal milestone path (desktop) / vertical (mobile) */}
        <div className="mb-6">
          {/* Desktop: horizontal */}
          <div className="relative hidden py-10 md:block">
            <svg
              className="pointer-events-none absolute left-0 right-0 top-12 h-20 w-full"
              viewBox="0 0 1000 80"
              preserveAspectRatio="none"
              aria-hidden
            >
              <path
                d="M50 40 Q 200 0 350 40 T 650 40 T 950 40"
                stroke="var(--color-line)"
                strokeWidth="3"
                fill="none"
                strokeDasharray="6 6"
              />
              <path
                d="M50 40 Q 200 0 350 40 T 650 40"
                stroke="var(--color-sage)"
                strokeWidth="3"
                fill="none"
              />
            </svg>
            <div className="relative flex justify-between px-5">
              {MILESTONES.map((m, i) => (
                <Reveal key={m.l} delay={i * 120}>
                  <MilestoneNode m={m} />
                </Reveal>
              ))}
            </div>
          </div>

          {/* Mobile: vertical zigzag */}
          <div className="relative pl-8 md:hidden">
            <svg
              className="absolute left-3.5 top-2 h-[calc(100%-16px)] w-6"
              viewBox="0 0 24 600"
              preserveAspectRatio="none"
              aria-hidden
            >
              <path
                d="M12 0 Q 24 75 12 150 Q 0 225 12 300 Q 24 375 12 450 Q 0 525 12 600"
                stroke="var(--color-line)"
                strokeWidth="2"
                fill="none"
                strokeDasharray="4 4"
              />
            </svg>
            <ul className="flex flex-col gap-4">
              {MILESTONES.map((m, i) => (
                <Reveal key={m.l} delay={i * 80}>
                  <li className="relative">
                    <span
                      className="absolute -left-[30px] top-1.5 grid h-6 w-6 place-items-center rounded-full border-[2.5px] text-[11px] font-extrabold"
                      style={{
                        background: m.done ? m.c : m.current ? "#fff" : "var(--color-line)",
                        borderColor: m.c,
                        color: m.done ? "#fff" : (m.c as string),
                        animation: m.current ? "entre-pulse-ring 1.8s infinite" : undefined,
                      }}
                    >
                      {m.done ? "✓" : i + 1}
                    </span>
                    <div
                      className="rounded-xl border-[1.5px] bg-surface p-3.5"
                      style={{
                        borderColor: m.current ? m.c : "var(--color-line)",
                        boxShadow: m.current ? `0 6px 14px -4px ${m.c}66` : "none",
                      }}
                    >
                      <div className="mb-1.5 flex items-center gap-2">
                        <span className="text-lg">{m.emoji}</span>
                        <span className="font-serif text-[15px] font-extrabold">{m.l}</span>
                        {m.current && (
                          <span
                            className="ml-auto rounded-full px-2 py-0.5 text-[9px] font-bold text-white"
                            style={{
                              background: m.c,
                              animation: "entre-wiggle 2s infinite",
                            }}
                          >
                            NOW
                          </span>
                        )}
                        {m.done && (
                          <span className="ml-auto rounded-full bg-sage-soft px-2 py-0.5 text-[9px] font-bold text-sage">
                            クリア
                          </span>
                        )}
                      </div>
                      <p className="text-[10px] text-ink-2">
                        {m.done && `${m.n}社クリアしました 🎉`}
                        {m.current && `${m.n}社 進行中 — 一次面接が間もなく始まります`}
                        {!m.done && !m.current && `あと一歩。${m.n}社の内定があなたを待っています`}
                      </p>
                    </div>
                  </li>
                </Reveal>
              ))}
            </ul>
          </div>
        </div>

        {/* Detail cards */}
        <div className="grid gap-4 md:grid-cols-[2fr_1fr]">
          <Reveal>
            <div className="rounded-xl border border-line bg-surface p-5">
              <div className="mb-3 flex items-baseline justify-between">
                <h2 className="text-[13px] font-extrabold">
                  📍 今のフェーズ — 面接（{interviewing.length}社）
                </h2>
                <Link href="/kanban" className="text-[11px] font-bold text-sage">
                  カンバンで見る →
                </Link>
              </div>
              <ul>
                {interviewing.map((e, i) => (
                  <li
                    key={e.id}
                    className={`flex items-center gap-3 py-2.5 ${
                      i ? "border-t border-dashed border-line" : ""
                    }`}
                  >
                    <div className="grid h-9 w-9 place-items-center rounded-lg bg-sage-wash font-serif text-lg font-extrabold text-sage">
                      {e.logo}
                    </div>
                    <div className="min-w-0 flex-1">
                      <div className="text-xs font-bold">{e.co}</div>
                      <div className="mt-0.5 text-[10px] text-ink-3">
                        {e.task} · {e.due}
                      </div>
                    </div>
                    <span
                      className="self-center rounded-md px-2.5 py-1 text-[10px] font-bold text-white"
                      style={{ background: e.color }}
                    >
                      {e.stageLabel}
                    </span>
                  </li>
                ))}
              </ul>
            </div>
          </Reveal>

          <Reveal delay={150}>
            <div className="flex flex-col items-center rounded-xl border-[1.5px] border-line bg-gradient-to-br from-cream-2 to-sage-wash p-5 text-center">
              <div style={{ animation: "entre-float 3s infinite" }}>
                <Mascot size={70} mood="cheering" />
              </div>
              <p className="mt-2 font-hand text-[20px] text-sage">あと一歩！</p>
              <p className="mt-1 font-serif text-sm font-extrabold">62% 完了</p>
              <p className="mt-1 text-[11px] leading-relaxed text-ink-2">
                面接5社突破できれば、
                <br />
                第一志望もぐっと近づきます。
              </p>
            </div>
          </Reveal>
        </div>
      </div>
    </AppShell>
  );
}

function MilestoneNode({ m }: { m: (typeof MILESTONES)[number] }) {
  return (
    <div className="flex flex-1 flex-col items-center">
      <div
        className="relative z-[1] grid h-[60px] w-[60px] place-items-center rounded-full border-[3.5px] text-2xl font-extrabold"
        style={{
          background: m.done ? m.c : m.current ? "#fff" : "var(--color-line)",
          borderColor: m.c,
          color: m.done ? "#fff" : (m.c as string),
          boxShadow: m.current ? `0 8px 20px -4px ${m.c}66` : "none",
          animation: m.current ? "entre-pulse-ring 1.8s infinite" : undefined,
        }}
      >
        {m.emoji}
      </div>
      <p className="mt-2.5 font-serif text-base font-extrabold">{m.l}</p>
      <p className="mt-1 text-[11px] text-ink-2">{m.n}社</p>
      {m.current && (
        <span
          className="mt-1.5 rounded-full px-2.5 py-0.5 text-[10px] font-bold text-white"
          style={{ background: m.c, animation: "entre-wiggle 2s infinite" }}
        >
          NOW
        </span>
      )}
      {m.done && (
        <span className="mt-1.5 rounded-full bg-sage-soft px-2.5 py-0.5 text-[10px] font-bold text-sage">
          クリア ✓
        </span>
      )}
    </div>
  );
}
