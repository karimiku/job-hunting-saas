"use client";

import Link from "next/link";
import { useState, type CSSProperties, type ReactNode } from "react";
import { DeerMascot, MiniMascot } from "./mascot";
import { Reveal } from "./motion";

// ═════════════════════════════════════════════════════════════════════════
// Entré LP — cream + sage edition (mockup v2)
// ═════════════════════════════════════════════════════════════════════════

export function LandingPage() {
  return (
    <div
      className="lp-scope"
      style={{ position: "relative", minHeight: "100%", background: "var(--lp-cream)" }}
    >
      <LPNav />
      <HeroSection />
      <ProblemSection />
      <MobileFirstSection />
      <ToolsSection />
      <UsageSection />
      <FAQSection />
      <FinalCTA />
      <LPFooter />
    </div>
  );
}

// ── NAV ──────────────────────────────────────────────────────────────────

function LPNav() {
  return (
    <div
      style={{
        position: "sticky",
        top: 0,
        zIndex: 50,
        background: "rgba(246, 241, 231, 0.82)",
        backdropFilter: "blur(12px)",
        WebkitBackdropFilter: "blur(12px)",
        borderBottom: "1px solid var(--lp-line-2)",
      }}
    >
      <div
        style={{
          maxWidth: 1200,
          margin: "0 auto",
          padding: "18px 32px",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
          <MiniMascot size={32} />
          <div
            style={{
              fontFamily: "var(--lp-font-serif)",
              fontSize: 26,
              fontWeight: 700,
              letterSpacing: "-0.02em",
              color: "var(--lp-ink)",
            }}
          >
            Entr<span style={{ color: "var(--lp-sage)" }}>é</span>
          </div>
          <span
            style={{
              fontSize: 10,
              fontWeight: 700,
              letterSpacing: "0.12em",
              padding: "3px 8px",
              borderRadius: 4,
              background: "var(--lp-sage-soft)",
              color: "var(--lp-sage-2)",
              marginLeft: 4,
            }}
          >
            BETA
          </span>
        </div>
        <nav
          className="lp-hide-sm"
          style={{ display: "flex", gap: 36, fontSize: 14, color: "var(--lp-ink-2)", fontWeight: 500 }}
        >
          <a href="#features">機能</a>
          <a href="#pricing">料金</a>
          <a href="#usage">使い方</a>
          <a href="#faq">よくある質問</a>
        </nav>
        <div style={{ display: "flex", gap: 10 }}>
          <Link
            href="/login"
            className="lp-btn lp-btn-ghost lp-hide-sm"
            style={{ height: 40, padding: "0 18px", fontSize: 14, display: "inline-flex", alignItems: "center" }}
          >
            ログイン
          </Link>
          <Link
            href="/login"
            className="lp-btn lp-btn-sage"
            style={{ height: 40, padding: "0 18px", fontSize: 14, display: "inline-flex", alignItems: "center" }}
          >
            無料で始める
          </Link>
        </div>
      </div>
    </div>
  );
}

// ── HERO ─────────────────────────────────────────────────────────────────

function HeroSection() {
  return (
    <section style={{ position: "relative", padding: "80px 32px 80px" }}>
      <ScatteredSparkles />
      <div
        style={{
          maxWidth: 1260,
          margin: "0 auto",
          display: "grid",
          gridTemplateColumns: "0.92fr 1.08fr",
          gap: 40,
          alignItems: "center",
        }}
        className="lp-hero-grid"
      >
        {/* Left: copy + CTAs */}
        <div style={{ position: "relative" }}>
          <Reveal>
            <div
              className="lp-hand"
              style={{
                fontSize: 24,
                color: "var(--lp-sage-2)",
                marginBottom: 20,
                display: "inline-flex",
                alignItems: "center",
                gap: 10,
              }}
            >
              <Swoosh />
              就活の、いちばんの味方に。
              <Sparkle />
            </div>
          </Reveal>
          <Reveal delay={80}>
            <h1
              style={{
                fontFamily: "var(--lp-font-serif)",
                fontSize: "clamp(44px, 5.8vw, 68px)",
                lineHeight: 1.15,
                letterSpacing: "-0.035em",
                fontWeight: 700,
                margin: 0,
                color: "var(--lp-ink)",
              }}
            >
              散らかった就活、<br />
              ぜんぶ1枚に。
            </h1>
          </Reveal>
          <Reveal delay={160}>
            <p
              style={{
                marginTop: 28,
                fontSize: 15,
                lineHeight: 2,
                color: "var(--lp-ink-2)",
                maxWidth: 440,
              }}
            >
              エントリーや締切、選考ステータス、メモ、URLまでひとまとめ。
              マイナビ、リクナビ、ONE CAREER、企業サイト…あちこちの情報を、さっと整理できます。
            </p>
          </Reveal>

          <Reveal delay={240}>
            <div style={{ marginTop: 32, display: "flex", flexWrap: "wrap", gap: 12 }}>
              <Link href="/login" className="lp-btn lp-btn-sage">
                無料ではじめる
              </Link>
            </div>
          </Reveal>

          <Reveal delay={320}>
            <div
              style={{
                marginTop: 26,
                display: "flex",
                flexWrap: "wrap",
                gap: 22,
                color: "var(--lp-ink-2)",
                fontSize: 13,
              }}
            >
              <TrustChip icon={<YenIcon />} label="完全無料" />
              <TrustChip icon={<UserIcon />} label="登録カンタン" />
              <TrustChip icon={<ShieldIcon />} label="データは安全に保存" />
            </div>
          </Reveal>
        </div>

        {/* Right: phone + desktop stage */}
        <Reveal delay={120} dir="right">
          <div
            style={{
              position: "relative",
              height: 640,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <HeroBlob />

            {/* Desktop behind-right */}
            <div
              style={{
                position: "absolute",
                right: "0%",
                top: "42%",
                zIndex: 1,
                transform: "translateY(-50%)",
              }}
            >
              <DesktopMock />
            </div>

            {/* Phone in front */}
            <div
              style={{
                position: "absolute",
                left: "4%",
                top: "50%",
                transform: "translateY(-50%)",
                zIndex: 3,
              }}
            >
              <PhoneDashboard />
            </div>

            {/* Mascot top-right with handwriting */}
            <div
              style={{
                position: "absolute",
                right: "0%",
                top: "-6%",
                zIndex: 4,
                display: "flex",
                alignItems: "flex-end",
                gap: 6,
              }}
              className="lp-hide-sm"
            >
              <div
                className="lp-hand"
                style={{
                  fontSize: 17,
                  color: "var(--lp-sage-2)",
                  textAlign: "right",
                  lineHeight: 1.4,
                  marginBottom: 16,
                }}
              >
                PCでも見やすく、<br />しっかり使える。
              </div>
              <div className="lp-float-y">
                <DeerMascot size={96} mood="wave" />
              </div>
            </div>

            {/* Curvy arrow connecting handwriting to desktop */}
            <svg
              aria-hidden
              className="lp-hide-sm"
              style={{ position: "absolute", right: "18%", top: "28%", zIndex: 2 }}
              width="64"
              height="42"
              viewBox="0 0 64 42"
            >
              <path
                d="M2 4 Q 30 -4 54 20 L54 36 M48 28 L54 36 L60 28"
                stroke="var(--lp-sage-2)"
                strokeWidth="1.6"
                fill="none"
                strokeLinecap="round"
                strokeLinejoin="round"
                opacity="0.7"
              />
            </svg>
          </div>
        </Reveal>
      </div>

      <style>{`
        @media (max-width: 980px) {
          .lp-hero-grid { grid-template-columns: 1fr !important; gap: 40px !important; }
          .lp-hero-grid > div:nth-child(2) { height: 520px !important; }
        }
      `}</style>
    </section>
  );
}

function ScatteredSparkles() {
  return (
    <div aria-hidden className="lp-sparkles">
      <Sparkle size={18} style={{ position: "absolute", left: "4%", top: "10%" }} />
      <Sparkle size={12} style={{ position: "absolute", left: "44%", top: "6%" }} />
      <Sparkle size={14} style={{ position: "absolute", left: "2%", top: "58%" }} />
      <Sparkle size={10} style={{ position: "absolute", left: "26%", bottom: "8%" }} />
      <Sparkle size={16} style={{ position: "absolute", right: "6%", top: "72%" }} />
      <Sparkle size={10} style={{ position: "absolute", right: "42%", bottom: "14%" }} />
    </div>
  );
}

function TrustChip({ icon, label }: { icon: ReactNode; label: string }) {
  return (
    <div style={{ display: "inline-flex", alignItems: "center", gap: 8 }}>
      <span
        style={{
          width: 26,
          height: 26,
          borderRadius: 999,
          background: "var(--lp-sage-soft)",
          color: "var(--lp-sage-2)",
          display: "inline-flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        {icon}
      </span>
      <span style={{ fontWeight: 600, color: "var(--lp-ink-2)" }}>{label}</span>
    </div>
  );
}

function Decorations() {
  return (
    <div aria-hidden style={{ position: "absolute", inset: 0, pointerEvents: "none" }}>
      <Sparkle size={18} style={{ position: "absolute", left: "6%", top: "14%" }} />
      <Sparkle size={14} style={{ position: "absolute", left: "42%", top: "8%" }} />
      <Sparkle size={16} style={{ position: "absolute", right: "8%", top: "48%" }} />
      <Sparkle size={12} style={{ position: "absolute", left: "2%", bottom: "8%" }} />
    </div>
  );
}

function HeroBlob() {
  return (
    <div
      aria-hidden
      style={{
        position: "absolute",
        inset: 0,
        pointerEvents: "none",
      }}
    >
      <div
        style={{
          position: "absolute",
          left: "18%",
          top: "14%",
          width: 360,
          height: 280,
          borderRadius: "60% 40% 55% 45% / 50% 50% 50% 50%",
          background: "var(--lp-sage-soft)",
          opacity: 0.55,
          filter: "blur(3px)",
        }}
      />
      <div
        style={{
          position: "absolute",
          right: "4%",
          bottom: "8%",
          width: 220,
          height: 140,
          borderRadius: "60% 40% 55% 45% / 50% 50% 50% 50%",
          background: "var(--lp-sage-tint)",
          opacity: 0.7,
        }}
      />
    </div>
  );
}

// ── PHONE MOCKUPS ───────────────────────────────────────────────────────

function PhoneFrame({ children }: { children: ReactNode }) {
  return (
    <div className="lp-phone">
      <div className="lp-phone-notch-status">
        <span>9:41</span>
        <span style={{ display: "inline-flex", alignItems: "center", gap: 5 }}>
          <SignalBars />
          <WifiIcon />
          <BatteryIcon />
        </span>
      </div>
      <div className="lp-phone-screen">{children}</div>
    </div>
  );
}

function SignalBars() {
  return (
    <svg width="14" height="10" viewBox="0 0 14 10" aria-hidden>
      <rect x="0" y="6" width="2" height="4" rx="0.5" fill="currentColor" />
      <rect x="4" y="4" width="2" height="6" rx="0.5" fill="currentColor" />
      <rect x="8" y="2" width="2" height="8" rx="0.5" fill="currentColor" />
      <rect x="12" y="0" width="2" height="10" rx="0.5" fill="currentColor" />
    </svg>
  );
}

function WifiIcon() {
  return (
    <svg width="13" height="10" viewBox="0 0 13 10" aria-hidden>
      <path d="M1 4 Q6.5 -1 12 4" stroke="currentColor" strokeWidth="1.2" fill="none" strokeLinecap="round" />
      <path d="M3 6.5 Q6.5 3 10 6.5" stroke="currentColor" strokeWidth="1.2" fill="none" strokeLinecap="round" />
      <circle cx="6.5" cy="9" r="0.9" fill="currentColor" />
    </svg>
  );
}

function BatteryIcon() {
  return (
    <svg width="22" height="10" viewBox="0 0 22 10" aria-hidden>
      <rect x="0.5" y="0.5" width="18" height="9" rx="2" stroke="currentColor" strokeWidth="1" fill="none" />
      <rect x="2" y="2" width="14" height="6" rx="1" fill="currentColor" />
      <rect x="19.5" y="3" width="2" height="4" rx="0.5" fill="currentColor" />
    </svg>
  );
}

function PhoneTopBar({ title, rightIcon }: { title: string; rightIcon?: ReactNode }) {
  return (
    <>
      <div style={{ height: 46 }} />
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          padding: "4px 16px 12px",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: 6 }}>
          <MiniMascot size={18} />
          <span style={{ fontFamily: "var(--lp-font-serif)", fontWeight: 700, fontSize: 15 }}>
            Entr<span style={{ color: "var(--lp-sage)" }}>é</span>
          </span>
          <span
            style={{
              fontSize: 8,
              fontWeight: 700,
              background: "var(--lp-sage-soft)",
              color: "var(--lp-sage-2)",
              padding: "2px 5px",
              borderRadius: 3,
            }}
          >
            BETA
          </span>
        </div>
        <div style={{ color: "var(--lp-ink-3)" }}>{rightIcon || <BellIcon />}</div>
      </div>
      <div style={{ padding: "0 16px 8px", fontSize: 18, fontWeight: 700 }}>{title}</div>
    </>
  );
}

function PhoneTabBar({ active }: { active: string }) {
  const tabs = [
    { id: "dashboard", label: "ダッシュボード", icon: <HomeIcon /> },
    { id: "entry", label: "エントリー", icon: <ListIcon /> },
    { id: "task", label: "タスク", icon: <CheckBoxIcon /> },
    { id: "inbox", label: "Inbox", icon: <InboxIcon /> },
    { id: "more", label: "その他", icon: <MoreIcon /> },
  ];
  return (
    <div
      style={{
        position: "absolute",
        bottom: 0,
        left: 0,
        right: 0,
        padding: "8px 8px 18px",
        background: "#fff",
        borderTop: "1px solid var(--lp-line)",
        display: "flex",
        justifyContent: "space-around",
      }}
    >
      {tabs.map((t) => (
        <div
          key={t.id}
          style={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            gap: 2,
            color: active === t.id ? "var(--lp-sage)" : "var(--lp-ink-3)",
            fontSize: 9,
            fontWeight: 600,
          }}
        >
          {t.icon}
          <span>{t.label}</span>
        </div>
      ))}
    </div>
  );
}

function PhoneDashboard() {
  return (
    <PhoneFrame>
      <PhoneTopBar title="ダッシュボード" />
      <div style={{ padding: "0 14px 90px", display: "grid", gap: 10 }}>
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8 }}>
          <StatMini label="エントリー数" value="24" unit="社" compact />
          <StatMini label="選考中" value="12" unit="社" compact />
          <StatMini label="内定" value="2" unit="社" compact />
          <StatMini label="今月の締切" value="5" unit="件" compact />
        </div>

        <Card>
          <div style={{ fontSize: 11, color: "var(--lp-ink-3)", fontWeight: 600, marginBottom: 8 }}>選考ステータス</div>
          <div style={{ display: "flex", alignItems: "center", gap: 14 }}>
            <Donut />
            <div style={{ display: "grid", gap: 3, fontSize: 10 }}>
              <LegendRow dot="var(--lp-s-entry-ink)" label="エントリー" count={12} />
              <LegendRow dot="var(--lp-s-doc-ink)" label="書類選考" count={6} />
              <LegendRow dot="var(--lp-s-test-ink)" label="面接" count={4} />
              <LegendRow dot="var(--lp-s-offer-ink)" label="内定" count={2} />
              <LegendRow dot="var(--lp-s-reject-ink)" label="辞退" count={1} />
            </div>
          </div>
        </Card>

        <Card>
          <div style={{ fontSize: 11, color: "var(--lp-ink-3)", fontWeight: 600, marginBottom: 6 }}>直近の締切</div>
          <DeadlineRow date="5/20 (火) 23:59" company="株式会社○○○○" tag="エントリーシート" />
          <DeadlineRow date="5/25 (日) 10:00" company="株式会社△△△△" tag="一次面接（Web）" />
          <DeadlineRow date="5/27 (火) 12:00" company="株式会社□□□□" tag="インターン説明会" />
          <div style={{ textAlign: "right", fontSize: 10, color: "var(--lp-sage-2)", fontWeight: 700, marginTop: 4 }}>
            すべて見る →
          </div>
        </Card>
      </div>
      <div
        style={{
          position: "absolute",
          right: 14,
          bottom: 72,
          width: 40,
          height: 40,
          borderRadius: "50%",
          background: "var(--lp-sage)",
          color: "#fff",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          fontSize: 20,
          fontWeight: 700,
          boxShadow: "0 8px 20px -4px rgba(83, 119, 96, 0.45)",
        }}
      >
        +
      </div>
      <PhoneTabBar active="dashboard" />
    </PhoneFrame>
  );
}

function Card({ children, style }: { children: ReactNode; style?: CSSProperties }) {
  return (
    <div
      style={{
        background: "#fff",
        border: "1px solid var(--lp-line)",
        borderRadius: 12,
        padding: 12,
        ...style,
      }}
    >
      {children}
    </div>
  );
}

function Donut() {
  const bg =
    "conic-gradient(var(--lp-s-entry-ink) 0 48%, var(--lp-s-doc-ink) 48% 72%, var(--lp-s-test-ink) 72% 88%, var(--lp-s-offer-ink) 88% 96%, var(--lp-s-reject-ink) 96% 100%)";
  return (
    <div style={{ position: "relative", width: 72, height: 72, flexShrink: 0 }}>
      <div style={{ width: 72, height: 72, borderRadius: "50%", background: bg }} />
      <div
        style={{
          position: "absolute",
          inset: 14,
          borderRadius: "50%",
          background: "#fff",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          fontSize: 14,
          fontWeight: 700,
        }}
      >
        25
      </div>
    </div>
  );
}

function LegendRow({ dot, label, count }: { dot: string; label: string; count: number }) {
  return (
    <div style={{ display: "flex", alignItems: "center", gap: 6 }}>
      <span className="lp-dot" style={{ background: dot, width: 6, height: 6 }} />
      <span style={{ color: "var(--lp-ink-2)" }}>{label}</span>
      <span style={{ marginLeft: "auto", fontWeight: 700 }}>{count}</span>
    </div>
  );
}

function DeadlineRow({ date, company, tag }: { date: string; company: string; tag: string }) {
  return (
    <div style={{ padding: "6px 0", borderBottom: "1px dashed var(--lp-line)", fontSize: 10 }}>
      <div style={{ color: "var(--lp-ink-2)", fontWeight: 700 }}>{date}</div>
      <div style={{ display: "flex", justifyContent: "space-between", marginTop: 2 }}>
        <span>{company}</span>
        <span style={{ color: "var(--lp-ink-3)" }}>{tag}</span>
      </div>
    </div>
  );
}

function StatMini({
  label,
  value,
  unit,
  compact = false,
}: {
  label: string;
  value: string;
  unit: string;
  compact?: boolean;
}) {
  return (
    <div
      style={{
        border: "1px solid var(--lp-line)",
        borderRadius: compact ? 10 : 10,
        padding: compact ? "8px 10px" : "10px 12px",
        background: "#fff",
      }}
    >
      <div style={{ fontSize: compact ? 9 : 10, color: "var(--lp-ink-3)", fontWeight: 600 }}>
        {label}
      </div>
      <div style={{ display: "flex", alignItems: "baseline", gap: 4, marginTop: 2 }}>
        <span style={{ fontSize: compact ? 18 : 22, fontWeight: 800, fontFamily: "var(--lp-font-serif)" }}>
          {value}
        </span>
        <span style={{ fontSize: compact ? 10 : 11, color: "var(--lp-ink-3)" }}>{unit}</span>
      </div>
    </div>
  );
}

// ── DESKTOP MOCKUP ───────────────────────────────────────────────────────

function DesktopMock({ size = "md" }: { size?: "md" | "sm" }) {
  const scale = size === "sm" ? 0.82 : 1;
  return (
    <div className="lp-desktop" style={{ transform: `scale(${scale})`, transformOrigin: "center" }}>
      <div className="lp-desktop-screen">
        <div className="lp-desktop-screen-inner" style={{ display: "grid", gridTemplateColumns: "64px 1fr", fontSize: 10 }}>
          <div
            style={{
              background: "var(--lp-cream-2)",
              borderRight: "1px solid var(--lp-line)",
              padding: "10px 6px",
              display: "grid",
              gap: 6,
              fontSize: 8,
              color: "var(--lp-ink-2)",
            }}
          >
            <div style={{ display: "flex", alignItems: "center", gap: 4 }}>
              <MiniMascot size={12} />
              <span style={{ fontFamily: "var(--lp-font-serif)", fontWeight: 700, fontSize: 10 }}>
                Entr<span style={{ color: "var(--lp-sage)" }}>é</span>
              </span>
            </div>
            <DesktopNav label="ダッシュボード" active icon={<HomeIcon />} />
            <DesktopNav label="エントリー" icon={<ListIcon />} />
            <DesktopNav label="タスク" icon={<CheckBoxIcon />} />
            <DesktopNav label="Inbox" icon={<InboxIcon />} />
            <DesktopNav label="ブックマーク" icon={<BookmarkIcon />} />
          </div>
          <div style={{ padding: "10px 12px" }}>
            <div style={{ fontSize: 11, fontWeight: 700, marginBottom: 8 }}>ダッシュボード</div>
            <div style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 6, marginBottom: 8 }}>
              <DeskStat label="エントリー数" value="24" unit="社" />
              <DeskStat label="選考中" value="12" unit="社" />
              <DeskStat label="内定" value="2" unit="社" />
              <DeskStat label="今月の締切" value="5" unit="件" />
            </div>
            <div style={{ fontSize: 9, fontWeight: 700, marginBottom: 4, color: "var(--lp-ink-2)" }}>選考ステータス</div>
            <div style={{ display: "grid", gridTemplateColumns: "repeat(5, 1fr)", gap: 4 }}>
              <DeskStatus bg="var(--lp-s-entry)" ink="var(--lp-s-entry-ink)" label="E" value={12} />
              <DeskStatus bg="var(--lp-s-doc)" ink="var(--lp-s-doc-ink)" label="書" value={6} />
              <DeskStatus bg="var(--lp-s-test)" ink="var(--lp-s-test-ink)" label="面" value={4} />
              <DeskStatus bg="var(--lp-s-offer)" ink="var(--lp-s-offer-ink)" label="内" value={2} />
              <DeskStatus bg="var(--lp-s-reject)" ink="var(--lp-s-reject-ink)" label="辞" value={1} />
            </div>
            <div style={{ textAlign: "right", fontSize: 8, color: "var(--lp-sage-2)", marginTop: 4, fontWeight: 700 }}>
              すべて見る →
            </div>
          </div>
        </div>
      </div>
      <div className="lp-desktop-neck" />
      <div className="lp-desktop-base" />
    </div>
  );
}

function DesktopNav({ label, icon, active = false }: { label: string; icon: ReactNode; active?: boolean }) {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        gap: 4,
        padding: "4px 5px",
        borderRadius: 4,
        background: active ? "var(--lp-sage-soft)" : "transparent",
        color: active ? "var(--lp-sage-2)" : "var(--lp-ink-2)",
        fontWeight: active ? 700 : 500,
        fontSize: 8,
      }}
    >
      <span style={{ transform: "scale(0.7)" }}>{icon}</span>
      <span>{label}</span>
    </div>
  );
}

function DeskStat({ label, value, unit }: { label: string; value: string; unit: string }) {
  return (
    <div style={{ border: "1px solid var(--lp-line)", borderRadius: 6, padding: "5px 6px" }}>
      <div style={{ fontSize: 7, color: "var(--lp-ink-3)", fontWeight: 600 }}>{label}</div>
      <div style={{ display: "flex", alignItems: "baseline", gap: 2, marginTop: 1 }}>
        <span style={{ fontSize: 13, fontWeight: 800, fontFamily: "var(--lp-font-serif)" }}>{value}</span>
        <span style={{ fontSize: 7, color: "var(--lp-ink-3)" }}>{unit}</span>
      </div>
    </div>
  );
}

function DeskStatus({
  bg,
  ink,
  label,
  value,
}: {
  bg: string;
  ink: string;
  label: string;
  value: number;
}) {
  return (
    <div style={{ background: bg, borderRadius: 5, padding: "4px", textAlign: "center" }}>
      <div style={{ fontSize: 7, fontWeight: 700, color: ink }}>{label}</div>
      <div style={{ fontSize: 12, fontWeight: 800, color: ink, fontFamily: "var(--lp-font-serif)" }}>
        {value}
      </div>
    </div>
  );
}

// ── PROBLEM ──────────────────────────────────────────────────────────────

function ProblemSection() {
  const problems = [
    { title: "情報がバラバラで、\n管理が大変", icon: <DocsIcon /> },
    { title: "締切や予定を\nうっかり忘れてしまう", icon: <CalendarIcon /> },
    { title: "メモやURLを\n探すのに時間がかかる", icon: <SearchIcon /> },
    { title: "PCとスマホで\n情報が分散している", icon: <DevicesIcon /> },
  ];
  return (
    <section style={{ padding: "50px 32px 80px", position: "relative" }}>
      <div style={{ maxWidth: 1200, margin: "0 auto" }}>
        <Reveal>
          <h2
            style={{
              textAlign: "center",
              fontFamily: "var(--lp-font-serif)",
              fontSize: "clamp(24px, 3vw, 30px)",
              fontWeight: 700,
              letterSpacing: "-0.02em",
              marginBottom: 40,
              color: "var(--lp-ink)",
            }}
          >
            こんなお悩み、ありませんか？
          </h2>
        </Reveal>
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(4, 1fr) 190px",
            gap: 16,
            alignItems: "stretch",
          }}
          className="lp-problem-grid"
        >
          {problems.map((p, i) => (
            <Reveal key={i} delay={i * 60}>
              <div
                className="lp-card"
                style={{
                  padding: "26px 18px",
                  height: "100%",
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                  gap: 16,
                  textAlign: "center",
                }}
              >
                <div
                  style={{
                    width: 64,
                    height: 64,
                    borderRadius: 16,
                    background: "var(--lp-sage-tint)",
                    display: "inline-flex",
                    alignItems: "center",
                    justifyContent: "center",
                    color: "var(--lp-sage-2)",
                  }}
                >
                  {p.icon}
                </div>
                <div
                  style={{
                    fontSize: 14,
                    fontWeight: 700,
                    color: "var(--lp-ink)",
                    lineHeight: 1.65,
                    whiteSpace: "pre-line",
                  }}
                >
                  {p.title}
                </div>
              </div>
            </Reveal>
          ))}
          <Reveal delay={300}>
            <div
              style={{
                height: "100%",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                position: "relative",
              }}
            >
              <div className="lp-float-y">
                <DeerMascot size={160} mood="worried" />
              </div>
            </div>
          </Reveal>
        </div>
      </div>

      <style>{`
        @media (max-width: 1024px) {
          .lp-problem-grid { grid-template-columns: 1fr 1fr !important; }
          .lp-problem-grid > div:last-child { grid-column: 1 / -1; }
        }
        @media (max-width: 520px) {
          .lp-problem-grid { grid-template-columns: 1fr !important; }
        }
      `}</style>
    </section>
  );
}

// ── MOBILE FIRST ─────────────────────────────────────────────────────────

function MobileFirstSection() {
  return (
    <section id="features" style={{ padding: "80px 32px 90px" }}>
      <div
        style={{
          maxWidth: 1200,
          margin: "0 auto",
          display: "grid",
          gridTemplateColumns: "0.95fr 1.05fr",
          gap: 50,
          alignItems: "center",
        }}
        className="lp-mf-grid"
      >
        <div>
          <Reveal>
            <div
              style={{
                fontSize: 12,
                fontWeight: 800,
                letterSpacing: "0.16em",
                color: "var(--lp-sage-2)",
                marginBottom: 14,
              }}
            >
              MOBILE FIRST
            </div>
          </Reveal>
          <Reveal delay={80}>
            <h2
              style={{
                fontFamily: "var(--lp-font-serif)",
                fontSize: "clamp(30px, 4vw, 40px)",
                lineHeight: 1.28,
                letterSpacing: "-0.03em",
                fontWeight: 700,
                margin: 0,
              }}
            >
              スマホがメイン。<br />
              PCでもしっかりサポート。
            </h2>
          </Reveal>
          <Reveal delay={160}>
            <p style={{ marginTop: 22, fontSize: 15, lineHeight: 1.95, color: "var(--lp-ink-2)" }}>
              いつでも、どこでも、サッと記録。
              <br />
              すきま時間を味方にして、就活をもっとスマートに。
            </p>
          </Reveal>
          <div style={{ marginTop: 26, display: "grid", gap: 14 }}>
            {[
              "スマホでサッと記録・確認",
              "通知で締切や予定をお知らせ",
              "データはクラウドで自動同期",
            ].map((t, i) => (
              <Reveal key={t} delay={200 + i * 60}>
                <div style={{ display: "flex", alignItems: "center", gap: 10, fontSize: 14, color: "var(--lp-ink-2)" }}>
                  <CheckBadge /> {t}
                </div>
              </Reveal>
            ))}
          </div>
        </div>

        <Reveal delay={120} dir="right">
          <div
            style={{
              position: "relative",
              height: 560,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <div
              aria-hidden
              style={{
                position: "absolute",
                right: "8%",
                top: "18%",
                width: 300,
                height: 200,
                borderRadius: "60% 40% 55% 45% / 50% 50% 50% 50%",
                background: "var(--lp-sage-tint)",
                opacity: 0.65,
                filter: "blur(4px)",
              }}
            />
            <div style={{ position: "absolute", left: "0%", top: "50%", transform: "translateY(-50%)", zIndex: 3 }}>
              <PhoneDashboard />
            </div>
            <div
              style={{
                position: "absolute",
                right: "-2%",
                top: "58%",
                transform: "translateY(-50%)",
                zIndex: 1,
              }}
            >
              <DesktopMock />
            </div>
            <div
              style={{
                position: "absolute",
                right: "8%",
                top: "2%",
                zIndex: 4,
              }}
              className="lp-hide-sm"
            >
              <SyncChip />
            </div>
          </div>
        </Reveal>
      </div>

      <style>{`
        @media (max-width: 980px) {
          .lp-mf-grid { grid-template-columns: 1fr !important; gap: 40px !important; }
        }
      `}</style>
    </section>
  );
}

function SyncChip() {
  return (
    <div
      style={{
        background: "var(--lp-sage-soft)",
        color: "var(--lp-sage-2)",
        padding: "18px 22px",
        borderRadius: 999,
        fontSize: 14,
        fontWeight: 700,
        boxShadow: "0 10px 24px -12px rgba(83, 119, 96, 0.3)",
        textAlign: "center",
        lineHeight: 1.5,
      }}
    >
      自動同期で<br />いつでも最新
    </div>
  );
}

function PhoneDashboardCompact() {
  // reuse PhoneDashboard but the outer frame (lp-phone is already set)
  return <PhoneDashboard />;
}

// ── 4 TOOLS ──────────────────────────────────────────────────────────────

function ToolsSection() {
  const tools: Array<{
    title: string;
    sub: string;
    body: string;
    icon: ReactNode;
    bg: string;
    note?: string;
  }> = [
    {
      title: "Entry",
      sub: "エントリー管理",
      body: "企業や求人をまとめて管理。選考ステータスも一目でわかる。",
      icon: <EntryIconLg />,
      bg: "var(--lp-s-entry)",
    },
    {
      title: "Task",
      sub: "タスク管理",
      body: "締切や面接予定を登録して、リマインドで忘れない。",
      icon: <TaskIconLg />,
      bg: "var(--lp-s-doc)",
    },
    {
      title: "Inbox",
      sub: "連絡・やること管理",
      body: "企業からの連絡や自分のやることを整理。",
      icon: <InboxIconLg />,
      bg: "var(--lp-s-reject)",
    },
    {
      title: "Chrome拡張",
      sub: "ワンクリック登録",
      body: "ブラウザからカンタンに求人情報を保存。",
      icon: <PuzzleIconLg />,
      bg: "var(--lp-s-interview)",
      note: "PCで試す",
    },
  ];
  return (
    <section style={{ padding: "20px 32px 80px" }}>
      <div style={{ maxWidth: 1200, margin: "0 auto" }}>
        <Reveal>
          <div
            style={{
              textAlign: "center",
              fontFamily: "var(--lp-font-serif)",
              fontSize: "clamp(24px, 3vw, 30px)",
              fontWeight: 700,
              letterSpacing: "-0.02em",
              marginBottom: 36,
              color: "var(--lp-ink)",
            }}
          >
            就活を支える、4つのツール
          </div>
        </Reveal>
        <div
          style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 16 }}
          className="lp-tools-grid"
        >
          {tools.map((t, i) => (
            <Reveal key={t.title} delay={i * 60}>
              <div
                className="lp-card"
                style={{
                  padding: 22,
                  height: "100%",
                  display: "flex",
                  flexDirection: "column",
                  gap: 12,
                  position: "relative",
                }}
              >
                {t.note && (
                  <span
                    style={{
                      position: "absolute",
                      top: 14,
                      right: 14,
                      fontSize: 10,
                      fontWeight: 700,
                      color: "var(--lp-sage-2)",
                      background: "var(--lp-sage-soft)",
                      padding: "3px 8px",
                      borderRadius: 999,
                      display: "inline-flex",
                      alignItems: "center",
                      gap: 4,
                    }}
                  >
                    <DesktopIcon /> {t.note}
                  </span>
                )}
                <div
                  style={{
                    width: 48,
                    height: 48,
                    borderRadius: 12,
                    background: t.bg,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                  }}
                >
                  {t.icon}
                </div>
                <div>
                  <div style={{ fontFamily: "var(--lp-font-serif)", fontSize: 22, fontWeight: 700 }}>
                    {t.title}
                  </div>
                  <div style={{ fontSize: 12, color: "var(--lp-ink-3)", fontWeight: 600, marginTop: 2 }}>
                    {t.sub}
                  </div>
                </div>
                <div style={{ fontSize: 13, color: "var(--lp-ink-2)", lineHeight: 1.7 }}>{t.body}</div>
              </div>
            </Reveal>
          ))}
        </div>
      </div>

      <style>{`
        @media (max-width: 880px) { .lp-tools-grid { grid-template-columns: 1fr 1fr !important; } }
        @media (max-width: 520px) { .lp-tools-grid { grid-template-columns: 1fr !important; } }
      `}</style>
    </section>
  );
}

// ── USAGE FLOW ───────────────────────────────────────────────────────────

function UsageSection() {
  const steps = [
    { n: 1, title: "アカウントを作成", body: "Googleで登録してすぐに使えます。", icon: <GoogleG size={20} color /> },
    { n: 2, title: "情報をまとめる", body: "エントリー・締切・メモ・URLをまとめて記録。", icon: <StackIcon /> },
    { n: 3, title: "ステータスを管理", body: "選考の進捗をひと目でチェック。", icon: <ChecklistIcon /> },
    { n: 4, title: "通知を受け取る", body: "締切や予定を通知でお知らせ。うっかりを防止。", icon: <BellFlagIcon /> },
  ];
  return (
    <section id="usage" style={{ padding: "20px 32px 100px" }}>
      <div style={{ maxWidth: 1200, margin: "0 auto" }}>
        <Reveal>
          <div
            style={{
              textAlign: "center",
              fontFamily: "var(--lp-font-serif)",
              fontSize: "clamp(24px, 3vw, 30px)",
              fontWeight: 700,
              letterSpacing: "-0.02em",
              marginBottom: 36,
              color: "var(--lp-ink)",
            }}
          >
            使い方は、とってもシンプル。
          </div>
        </Reveal>
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(4, 1fr) 0.9fr",
            gap: 14,
            alignItems: "stretch",
          }}
          className="lp-usage-grid"
        >
          {steps.map((s, i) => (
            <Reveal key={s.n} delay={i * 60}>
              <div className="lp-card" style={{ padding: 20, height: "100%", position: "relative" }}>
                <div
                  style={{
                    width: 26,
                    height: 26,
                    borderRadius: "50%",
                    background: "var(--lp-sage)",
                    color: "#fff",
                    fontSize: 13,
                    fontWeight: 700,
                    display: "inline-flex",
                    alignItems: "center",
                    justifyContent: "center",
                    marginBottom: 10,
                  }}
                >
                  {s.n}
                </div>
                <div style={{ fontSize: 14, fontWeight: 700, marginBottom: 10 }}>{s.title}</div>
                <div
                  style={{
                    width: 52,
                    height: 52,
                    borderRadius: 12,
                    background: "var(--lp-sage-tint)",
                    display: "inline-flex",
                    alignItems: "center",
                    justifyContent: "center",
                    color: "var(--lp-sage-2)",
                    marginBottom: 10,
                  }}
                >
                  {s.icon}
                </div>
                <div style={{ fontSize: 12, color: "var(--lp-ink-2)", lineHeight: 1.7 }}>{s.body}</div>
                {i < 3 && (
                  <div
                    aria-hidden
                    className="lp-hide-sm"
                    style={{
                      position: "absolute",
                      right: -18,
                      top: "50%",
                      transform: "translateY(-50%)",
                      color: "var(--lp-ink-4)",
                      fontSize: 16,
                    }}
                  >
                    →
                  </div>
                )}
              </div>
            </Reveal>
          ))}
          <Reveal delay={300}>
            <div
              style={{
                padding: 20,
                height: "100%",
                background: "var(--lp-sage-tint)",
                border: "1px solid rgba(107, 144, 121, 0.18)",
                borderRadius: 20,
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
                textAlign: "center",
                gap: 8,
              }}
            >
              <div
                style={{
                  fontFamily: "var(--lp-font-serif)",
                  fontSize: 18,
                  fontWeight: 700,
                  color: "var(--lp-ink)",
                  lineHeight: 1.5,
                }}
              >
                あなたの就活を、<br />やさしく自動化。
              </div>
              <div className="lp-float-y" style={{ marginTop: 4 }}>
                <DeerMascot size={108} mood="sparkle" />
              </div>
            </div>
          </Reveal>
        </div>
      </div>

      <style>{`
        @media (max-width: 1024px) { .lp-usage-grid { grid-template-columns: 1fr 1fr !important; } }
        @media (max-width: 520px) { .lp-usage-grid { grid-template-columns: 1fr !important; } }
      `}</style>
    </section>
  );
}

// ── FAQ ──────────────────────────────────────────────────────────────────

function FAQSection() {
  const [open, setOpen] = useState<number>(-1);
  const faqs = [
    {
      q: "本当に無料で使えますか？",
      a: "はい、MVP期間中は完全無料でお使いいただけます。クレジットカード登録も不要です。",
    },
    {
      q: "データは安全に保管されますか？",
      a: "通信は暗号化され、データは安全な環境に保管されます。退会時には全データが削除されます。",
    },
    {
      q: "対応している端末は何ですか？",
      a: "PC・スマホ（iOS/Android）のブラウザからご利用いただけます。Chrome拡張のみPC専用です。",
    },
  ];
  return (
    <section id="faq" style={{ padding: "60px 32px 80px" }}>
      <div style={{ maxWidth: 1200, margin: "0 auto" }}>
        <Reveal>
          <h2
            style={{
              textAlign: "center",
              fontFamily: "var(--lp-font-serif)",
              fontSize: "clamp(26px, 3.2vw, 32px)",
              fontWeight: 700,
              letterSpacing: "-0.02em",
              marginBottom: 36,
              color: "var(--lp-ink)",
            }}
          >
            よくある質問
          </h2>
        </Reveal>
        <div
          style={{ display: "grid", gridTemplateColumns: "1.4fr 1fr", gap: 24, alignItems: "start" }}
          className="lp-faq-grid"
        >
          <div style={{ display: "grid", gap: 10 }}>
            {faqs.map((f, i) => (
              <Reveal key={i} delay={i * 60}>
                <div className="lp-faq-item">
                  <button
                    type="button"
                    onClick={() => setOpen(open === i ? -1 : i)}
                    style={{
                      width: "100%",
                      padding: "16px 20px",
                      background: "transparent",
                      border: "none",
                      cursor: "pointer",
                      fontFamily: "inherit",
                      display: "flex",
                      alignItems: "center",
                      gap: 14,
                      textAlign: "left",
                    }}
                  >
                    <span
                      style={{
                        width: 24,
                        height: 24,
                        borderRadius: 999,
                        background: "var(--lp-sage-soft)",
                        color: "var(--lp-sage-2)",
                        display: "inline-flex",
                        alignItems: "center",
                        justifyContent: "center",
                        flexShrink: 0,
                      }}
                    >
                      <QIcon />
                    </span>
                    <span style={{ fontSize: 14, fontWeight: 700, flex: 1, color: "var(--lp-ink)" }}>
                      {f.q}
                    </span>
                    <span
                      style={{
                        color: "var(--lp-ink-3)",
                        transform: open === i ? "rotate(180deg)" : "none",
                        transition: "transform 200ms",
                      }}
                    >
                      <ChevronDown />
                    </span>
                  </button>
                  {open === i && (
                    <div
                      style={{
                        padding: "0 20px 18px 58px",
                        fontSize: 13,
                        lineHeight: 1.85,
                        color: "var(--lp-ink-2)",
                      }}
                    >
                      {f.a}
                    </div>
                  )}
                </div>
              </Reveal>
            ))}
            <Reveal delay={220}>
              <a
                href="#"
                style={{
                  marginTop: 6,
                  fontSize: 13,
                  fontWeight: 700,
                  color: "var(--lp-sage-2)",
                  display: "inline-block",
                }}
              >
                すべての質問を見る →
              </a>
            </Reveal>
          </div>

          <Reveal delay={180} dir="right">
            <div
              className="lp-card-soft"
              style={{
                padding: "28px 26px",
                display: "flex",
                alignItems: "center",
                gap: 18,
                background: "var(--lp-sage-tint)",
                border: "1px solid rgba(107, 144, 121, 0.15)",
                borderRadius: 20,
                minHeight: 200,
              }}
            >
              <div
                style={{
                  flex: 1,
                  fontSize: 15,
                  fontWeight: 700,
                  lineHeight: 1.7,
                  color: "var(--lp-ink)",
                  fontFamily: "var(--lp-font-serif)",
                }}
              >
                ご不明点は<br />サポートまで<br />お気軽にどうぞ！
              </div>
              <div className="lp-float-y">
                <DeerMascot size={130} mood="happy" />
              </div>
            </div>
          </Reveal>
        </div>
      </div>

      <style>{`
        @media (max-width: 980px) { .lp-faq-grid { grid-template-columns: 1fr !important; } }
      `}</style>
    </section>
  );
}

// ── FINAL CTA ────────────────────────────────────────────────────────────

function FinalCTA() {
  return (
    <section id="pricing" style={{ padding: "40px 32px 110px" }}>
      <div
        style={{
          maxWidth: 1120,
          margin: "0 auto",
          padding: "52px 56px",
          borderRadius: 28,
          background: "var(--lp-sage-tint)",
          border: "1px solid rgba(107, 144, 121, 0.2)",
          display: "grid",
          gridTemplateColumns: "1.1fr 0.9fr 180px",
          gap: 30,
          alignItems: "center",
          position: "relative",
          overflow: "visible",
          boxShadow: "0 30px 60px -40px rgba(83, 119, 96, 0.35)",
        }}
        className="lp-cta-card"
      >
        <div>
          <Reveal>
            <h2
              style={{
                fontFamily: "var(--lp-font-serif)",
                fontSize: "clamp(30px, 4vw, 42px)",
                lineHeight: 1.3,
                letterSpacing: "-0.03em",
                fontWeight: 700,
                margin: 0,
                color: "var(--lp-ink)",
              }}
            >
              さあ、Entr<span style={{ color: "var(--lp-sage)" }}>é</span>で<br />
              就活をもっとシンプルに。
            </h2>
          </Reveal>
          <Reveal delay={160}>
            <div
              style={{
                marginTop: 22,
                display: "flex",
                flexWrap: "wrap",
                gap: 18,
                color: "var(--lp-ink-2)",
                fontSize: 13,
              }}
            >
              <TrustChip icon={<YenIcon />} label="完全無料" />
              <TrustChip icon={<UserIcon />} label="登録カンタン" />
              <TrustChip icon={<ShieldIcon />} label="データは安全に保存" />
            </div>
          </Reveal>
        </div>

        <Reveal delay={120} dir="right">
          <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
            <Link href="/login" className="lp-btn lp-btn-sage" style={{ width: "100%" }}>
              無料ではじめる
            </Link>
          </div>
        </Reveal>

        <Reveal delay={200} dir="right">
          <div
            style={{
              display: "flex",
              justifyContent: "center",
              alignItems: "flex-end",
              position: "relative",
            }}
            className="lp-hide-sm"
          >
            <div className="lp-float-y">
              <DeerMascot size={160} mood="sparkle" />
            </div>
          </div>
        </Reveal>
      </div>

      <style>{`
        @media (max-width: 980px) {
          .lp-cta-card { grid-template-columns: 1fr !important; padding: 36px 26px !important; text-align: center; }
        }
      `}</style>
    </section>
  );
}

// ── FOOTER ───────────────────────────────────────────────────────────────

function LPFooter() {
  const cols = [
    { t: "プロダクト", items: ["機能一覧", "Chrome拡張", "価格（無料）", "ロードマップ"] },
    { t: "サポート", items: ["FAQ", "お問い合わせ", "データエクスポート", "退会方法"] },
    { t: "会社", items: ["運営について", "プライバシー", "利用規約"] },
  ];
  return (
    <footer
      style={{
        padding: "60px 32px 40px",
        background: "var(--lp-ink)",
        color: "#fff",
        position: "relative",
        zIndex: 3,
      }}
    >
      <div
        style={{
          maxWidth: 1240,
          margin: "0 auto",
          display: "grid",
          gridTemplateColumns: "1.4fr 1fr 1fr 1fr",
          gap: 40,
        }}
        className="lp-footer-grid"
      >
        <div>
          <div style={{ display: "flex", alignItems: "center", gap: 10, marginBottom: 14 }}>
            <MiniMascot size={36} />
            <div
              style={{
                fontWeight: 800,
                fontSize: 22,
                fontFamily: "var(--lp-font-serif)",
                letterSpacing: "-0.02em",
              }}
            >
              Entr<span style={{ color: "var(--lp-sage)" }}>é</span>
            </div>
          </div>
          <div
            style={{
              fontSize: 13,
              color: "rgba(255,255,255,0.7)",
              lineHeight: 1.7,
              maxWidth: 320,
            }}
          >
            複数の就活サイトに分散する応募・選考・タスクを、Entry単位で横断管理するSaaS。
          </div>
        </div>
        {cols.map((c) => (
          <div key={c.t}>
            <div style={{ fontSize: 12, fontWeight: 800, marginBottom: 12 }}>{c.t}</div>
            <div style={{ display: "grid", gap: 8, fontSize: 13, color: "rgba(255,255,255,0.75)" }}>
              {c.items.map((i) => (
                <a key={i} href="#">
                  {i}
                </a>
              ))}
            </div>
          </div>
        ))}
      </div>
      <div
        style={{
          maxWidth: 1240,
          margin: "48px auto 0",
          paddingTop: 24,
          borderTop: "1px solid rgba(255,255,255,0.12)",
          display: "flex",
          justifyContent: "space-between",
          fontSize: 12,
          color: "rgba(255,255,255,0.55)",
        }}
      >
        <div>© 2026 Entré</div>
        <div className="lp-mono">v0.3.0-beta</div>
      </div>

      <style>{`
        @media (max-width: 880px) { .lp-footer-grid { grid-template-columns: 1fr 1fr !important; } }
      `}</style>
    </footer>
  );
}

// ═════════════════════════════════════════════════════════════════════════
// Inline icons & decorations
// ═════════════════════════════════════════════════════════════════════════

function Swoosh({ flip = false, size = 22 }: { flip?: boolean; size?: number }) {
  return (
    <svg
      width={size * 1.8}
      height={size * 0.6}
      viewBox="0 0 36 12"
      style={{ transform: flip ? "scaleX(-1)" : undefined }}
      aria-hidden
    >
      <path d="M2 6 Q10 1 18 6 T34 6" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinecap="round" />
    </svg>
  );
}

function Sparkle({ size = 16, color, style }: { size?: number; color?: string; style?: CSSProperties }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" aria-hidden style={style}>
      <path
        d="M12 2 L14 10 L22 12 L14 14 L12 22 L10 14 L2 12 L10 10 Z"
        fill={color || "var(--lp-sage)"}
        opacity="0.65"
      />
    </svg>
  );
}

function GoogleG({ size = 18, color = false }: { size?: number; color?: boolean }) {
  if (color) {
    return (
      <svg width={size} height={size} viewBox="0 0 24 24" aria-hidden>
        <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
        <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
        <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" />
        <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
      </svg>
    );
  }
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" aria-hidden>
      <path fill="#fff" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
      <path fill="#fff" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" opacity="0.9" />
      <path fill="#fff" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" opacity="0.8" />
      <path fill="#fff" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" opacity="0.7" />
    </svg>
  );
}

function MailIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" aria-hidden>
      <rect x="3" y="5" width="18" height="14" rx="2" stroke="currentColor" strokeWidth="1.8" fill="none" />
      <path d="M3 7 L12 13 L21 7" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

function YenIcon() {
  return (
    <svg width="13" height="13" viewBox="0 0 24 24" aria-hidden>
      <path d="M7 4 L12 12 L17 4 M4 12 H20 M4 16 H20 M12 12 V20" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" fill="none" />
    </svg>
  );
}

function UserIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" aria-hidden>
      <circle cx="12" cy="8" r="4" stroke="currentColor" strokeWidth="1.8" fill="none" />
      <path d="M4 20 Q4 14 12 14 Q20 14 20 20" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinecap="round" />
    </svg>
  );
}

function ShieldIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" aria-hidden>
      <path d="M12 3 L20 6 V12 Q20 18 12 21 Q4 18 4 12 V6 Z" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinejoin="round" />
      <path d="M9 12 L11 14 L15 10" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" fill="none" />
    </svg>
  );
}

function BellIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" aria-hidden>
      <path
        d="M12 3 a6 6 0 0 0 -6 6 v4 l-2 3 h16 l-2 -3 v-4 a6 6 0 0 0 -6 -6 M10 20 a2 2 0 0 0 4 0"
        stroke="currentColor"
        strokeWidth="1.8"
        fill="none"
        strokeLinecap="round"
      />
    </svg>
  );
}

function HomeIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" aria-hidden>
      <path d="M4 11 L12 4 L20 11 V20 H14 V14 H10 V20 H4 Z" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinejoin="round" />
    </svg>
  );
}

function ListIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" aria-hidden>
      <path d="M4 7 H6 M4 12 H6 M4 17 H6 M9 7 H20 M9 12 H20 M9 17 H20" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  );
}

function CheckBoxIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" aria-hidden>
      <rect x="4" y="4" width="16" height="16" rx="3" stroke="currentColor" strokeWidth="1.8" fill="none" />
      <path d="M8 12 L11 15 L16 9" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinecap="round" />
    </svg>
  );
}

function InboxIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" aria-hidden>
      <path
        d="M4 6 H20 L20 15 H16 L14 17 H10 L8 15 H4 Z M4 15 L7 6 H17 L20 15"
        stroke="currentColor"
        strokeWidth="1.8"
        fill="none"
        strokeLinejoin="round"
      />
    </svg>
  );
}

function MoreIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" aria-hidden>
      <circle cx="6" cy="12" r="1.6" fill="currentColor" />
      <circle cx="12" cy="12" r="1.6" fill="currentColor" />
      <circle cx="18" cy="12" r="1.6" fill="currentColor" />
    </svg>
  );
}

function BookmarkIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" aria-hidden>
      <path d="M7 4 H17 V20 L12 16 L7 20 Z" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinejoin="round" />
    </svg>
  );
}

function CheckBadge() {
  return (
    <span
      style={{
        width: 22,
        height: 22,
        borderRadius: 999,
        background: "var(--lp-sage)",
        color: "#fff",
        display: "inline-flex",
        alignItems: "center",
        justifyContent: "center",
        fontSize: 12,
        fontWeight: 700,
        flexShrink: 0,
      }}
    >
      ✓
    </span>
  );
}

function EntryIconLg() {
  return (
    <svg width="28" height="28" viewBox="0 0 24 24" aria-hidden>
      <rect x="4" y="4" width="16" height="16" rx="2.5" stroke="var(--lp-s-entry-ink)" strokeWidth="1.8" fill="none" />
      <path d="M8 9 H16 M8 13 H16 M8 17 H12" stroke="var(--lp-s-entry-ink)" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  );
}

function TaskIconLg() {
  return (
    <svg width="28" height="28" viewBox="0 0 24 24" aria-hidden>
      <rect x="4" y="4" width="16" height="16" rx="2.5" stroke="var(--lp-s-doc-ink)" strokeWidth="1.8" fill="none" />
      <path d="M8 9 L10 11 L14 7 M8 15 L10 17 L14 13" stroke="var(--lp-s-doc-ink)" strokeWidth="1.8" fill="none" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

function InboxIconLg() {
  return (
    <svg width="28" height="28" viewBox="0 0 24 24" aria-hidden>
      <path
        d="M4 7 L20 7 L20 17 L4 17 Z M4 7 L12 13 L20 7"
        stroke="var(--lp-s-reject-ink)"
        strokeWidth="1.8"
        fill="none"
        strokeLinejoin="round"
      />
    </svg>
  );
}

function PuzzleIconLg() {
  return (
    <svg width="28" height="28" viewBox="0 0 24 24" aria-hidden>
      <path
        d="M9 4 H15 A2 2 0 0 1 17 6 V8 A2 2 0 0 0 19 10 H20 V16 A2 2 0 0 1 18 18 H16 A2 2 0 0 0 14 20 H9 A2 2 0 0 1 7 18 V16 A2 2 0 0 0 5 14 H4 V9 A2 2 0 0 1 6 7 H8 A2 2 0 0 0 10 5 V4 Z"
        stroke="var(--lp-s-interview-ink)"
        strokeWidth="1.6"
        fill="none"
        strokeLinejoin="round"
      />
    </svg>
  );
}

function StackIcon() {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" aria-hidden>
      <path d="M12 4 L20 8 L12 12 L4 8 Z M4 12 L12 16 L20 12 M4 16 L12 20 L20 16" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinejoin="round" />
    </svg>
  );
}

function ChecklistIcon() {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" aria-hidden>
      <path d="M4 6 L6 8 L9 5 M4 12 L6 14 L9 11 M4 18 L6 20 L9 17 M12 7 H20 M12 13 H20 M12 19 H18" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

function BellFlagIcon() {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" aria-hidden>
      <path d="M5 3 V21 M5 4 L17 4 L15 8 L17 12 L5 12" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinejoin="round" strokeLinecap="round" />
    </svg>
  );
}

function DesktopIcon() {
  return (
    <svg width="10" height="10" viewBox="0 0 24 24" aria-hidden>
      <rect x="3" y="4" width="18" height="12" rx="2" stroke="currentColor" strokeWidth="2" fill="none" />
      <path d="M8 20 H16 M12 16 V20" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
    </svg>
  );
}

function QIcon() {
  return (
    <svg width="12" height="12" viewBox="0 0 24 24" aria-hidden>
      <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="2" fill="none" />
      <path d="M9 10 Q9 7 12 7 Q15 7 15 10 Q15 12 13 13 Q12 13 12 14 V15" stroke="currentColor" strokeWidth="1.8" fill="none" strokeLinecap="round" />
      <circle cx="12" cy="17.5" r="0.9" fill="currentColor" />
    </svg>
  );
}

function ChevronDown() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" aria-hidden>
      <path d="M6 9 L12 15 L18 9" stroke="currentColor" strokeWidth="2" fill="none" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

// ── Problem icons ──────────────────────────────────────────────────────

function DocsIcon() {
  return (
    <svg width="26" height="26" viewBox="0 0 24 24" aria-hidden>
      <path d="M7 3 H13 L17 7 V17 H7 Z" stroke="currentColor" strokeWidth="1.6" fill="none" strokeLinejoin="round" />
      <path d="M5 7 V21 H15" stroke="currentColor" strokeWidth="1.6" fill="none" strokeLinejoin="round" />
      <path d="M13 3 V7 H17" stroke="currentColor" strokeWidth="1.6" fill="none" />
      <path d="M9 10 H14 M9 13 H14" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
    </svg>
  );
}

function CalendarIcon() {
  return (
    <svg width="26" height="26" viewBox="0 0 24 24" aria-hidden>
      <rect x="3" y="5" width="18" height="15" rx="2" stroke="currentColor" strokeWidth="1.6" fill="none" />
      <path d="M8 3 V7 M16 3 V7 M3 10 H21" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
      <rect x="7" y="13" width="3" height="3" rx="0.6" fill="currentColor" opacity="0.6" />
    </svg>
  );
}

function SearchIcon() {
  return (
    <svg width="26" height="26" viewBox="0 0 24 24" aria-hidden>
      <circle cx="10" cy="10" r="6" stroke="currentColor" strokeWidth="1.6" fill="none" />
      <path d="M14.5 14.5 L20 20" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  );
}

function DevicesIcon() {
  return (
    <svg width="26" height="26" viewBox="0 0 24 24" aria-hidden>
      <rect x="2" y="5" width="14" height="10" rx="1.4" stroke="currentColor" strokeWidth="1.6" fill="none" />
      <rect x="15" y="9" width="7" height="12" rx="1.4" stroke="currentColor" strokeWidth="1.6" fill="none" />
      <path d="M5 18 H12 M8.5 15 V18" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
    </svg>
  );
}
