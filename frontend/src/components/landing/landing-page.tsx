"use client";

import { useEffect, useState, type CSSProperties } from "react";
import { BackgroundScene } from "./background-scene";
import { Mascot, MiniMascot, Stamp, type Mood } from "./mascot";
import { CountUp, Parallax, Reveal, StaggerText } from "./motion";

// ═════════════════════════════════════════════════════════════════════════
// Entré LP — ported from design bundle (project/lib/landing.jsx)
// All CTAs are intentionally inert until backend auth wiring lands.
// ═════════════════════════════════════════════════════════════════════════

export function LandingPage() {
  return (
    <div className="lp-scope" style={{ position: "relative", minHeight: "100%" }}>
      <BackgroundScene />
      <div style={{ position: "relative", zIndex: 2 }}>
        <LPNav />
        <HeroMain />
        <MarqueeStrip />
        <ProblemSection />
        <MascotShowcase />
        <FeatureTour />
        <LiveDashboardSection />
        <ExtensionDemo />
        <TimelineSection />
        <PersonaSection />
        <FAQSection />
        <FinalCTA />
        <LPFooter />
      </div>
    </div>
  );
}

// ── NAV ──────────────────────────────────────────────────────────────────
function LPNav() {
  const [scrolled, setScrolled] = useState(false);
  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 40);
    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, []);
  return (
    <div
      style={{
        position: "sticky",
        top: 0,
        zIndex: 50,
        background: scrolled ? "rgba(244,247,251,0.8)" : "transparent",
        backdropFilter: scrolled ? "blur(12px)" : "none",
        WebkitBackdropFilter: scrolled ? "blur(12px)" : "none",
        borderBottom: scrolled
          ? "1px solid var(--lp-line-2)"
          : "1px solid transparent",
        transition: "background 220ms, border-color 220ms",
      }}
    >
      <div
        style={{
          maxWidth: 1240,
          margin: "0 auto",
          padding: "16px 32px",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
          <MiniMascot size={36} />
          <div
            style={{
              fontWeight: 800,
              fontSize: 20,
              letterSpacing: "-0.02em",
              fontFamily: "var(--lp-font-serif)",
            }}
          >
            Entr<span style={{ color: "var(--lp-blue)" }}>é</span>
            <span
              style={{
                fontSize: 10,
                fontWeight: 600,
                color: "var(--lp-ink-3)",
                fontFamily: "var(--lp-font-mono)",
                marginLeft: 6,
                letterSpacing: "0.1em",
              }}
            >
              BETA
            </span>
          </div>
        </div>
        <nav
          style={{
            display: "flex",
            gap: 28,
            fontSize: 14,
            color: "var(--lp-ink-2)",
            fontWeight: 500,
          }}
          className="lp-nav-links"
        >
          <a href="#features">機能</a>
          <a href="#screens">画面</a>
          <a href="#extension">拡張</a>
          <a href="#persona">こんな人に</a>
          <a href="#faq">FAQ</a>
        </nav>
        <div style={{ display: "flex", gap: 10, alignItems: "center" }}>
          <button
            type="button"
            className="lp-btn lp-btn-glass"
            style={{ height: 38, fontSize: 13 }}
          >
            ログイン
          </button>
          <button
            type="button"
            className="lp-btn lp-btn-blue"
            style={{ height: 40, fontSize: 13 }}
          >
            無料で始める
          </button>
        </div>
      </div>
    </div>
  );
}

// ── HERO ─────────────────────────────────────────────────────────────────
function HeroMain() {
  const [count, setCount] = useState(0);
  useEffect(() => {
    let n = 0;
    const target = 12;
    const id = setInterval(() => {
      n++;
      setCount(n);
      if (n >= target) clearInterval(id);
    }, 90);
    return () => clearInterval(id);
  }, []);

  return (
    <section
      style={{ position: "relative", overflow: "hidden", padding: "40px 0 100px" }}
    >
      <div
        aria-hidden
        style={{ position: "absolute", inset: 0, pointerEvents: "none" }}
      >
        <svg
          style={{ position: "absolute", top: 180, left: "8%" }}
          width="60"
          height="60"
          viewBox="0 0 60 60"
          className="lp-spin-slow"
        >
          <path
            d="M30 4 L36 24 L56 30 L36 36 L30 56 L24 36 L4 30 L24 24 Z"
            fill="var(--lp-blue)"
            opacity="0.6"
          />
        </svg>
        <svg
          style={{ position: "absolute", top: 620, right: "12%" }}
          width="40"
          height="40"
          viewBox="0 0 40 40"
          className="lp-spin-slow"
        >
          <path
            d="M20 2 L24 16 L38 20 L24 24 L20 38 L16 24 L2 20 L16 16 Z"
            fill="var(--lp-sun)"
          />
        </svg>
      </div>

      <div
        style={{
          maxWidth: 1240,
          margin: "0 auto",
          padding: "0 32px",
          position: "relative",
        }}
      >
        <div className="lp-hero-grid">
          <div>
            <Reveal>
              <div
                className="lp-chip"
                style={{
                  marginBottom: 28,
                  height: 32,
                  padding: "0 14px",
                  fontSize: 12,
                  background: "var(--lp-glass-strong)",
                  color: "var(--lp-blue)",
                  borderColor: "rgba(74,108,247,0.2)",
                }}
              >
                <span
                  className="lp-dot"
                  style={{ background: "var(--lp-blue)" }}
                />
                Entré — Entryから、就活がはじまる。
              </div>
            </Reveal>
            <Reveal delay={100}>
              <h1
                style={{
                  fontFamily: "var(--lp-font-serif)",
                  fontSize: "clamp(56px, 9vw, 128px)",
                  lineHeight: 0.95,
                  margin: 0,
                  letterSpacing: "-0.045em",
                  fontWeight: 700,
                }}
              >
                <StaggerText text="散らかった" perChar={60} />
                <br />
                <StaggerText text="就活、" perChar={60} delay={300} />
                <br />
                <span style={{ position: "relative", display: "inline-block" }}>
                  <span className="lp-hl-blue lp-hl-sweep">
                    <StaggerText text="ぜんぶ 1枚" perChar={60} delay={600} />
                  </span>
                  <svg
                    style={{
                      position: "absolute",
                      left: -10,
                      bottom: -18,
                      width: "108%",
                    }}
                    viewBox="0 0 420 20"
                    preserveAspectRatio="none"
                  >
                    <path
                      className="lp-draw-stroke"
                      pathLength={1}
                      d="M4 14 Q 100 2, 210 12 T 416 10"
                      fill="none"
                      stroke="var(--lp-blue)"
                      strokeWidth="3.5"
                      strokeLinecap="round"
                    />
                  </svg>
                </span>
                に。
              </h1>
            </Reveal>
            <Reveal delay={200}>
              <p
                style={{
                  fontSize: 18,
                  lineHeight: 1.85,
                  color: "var(--lp-ink-2)",
                  marginTop: 32,
                  maxWidth: 520,
                }}
              >
                マイナビ、リクナビ、ワンキャリ、企業HP…
                <br />
                応募が増えるほど、締切もメモも散らばる。
                <br />
                <b style={{ color: "var(--lp-ink)" }}>応募ごとに1箇所</b>
                へ集める、就活の台帳。
              </p>
            </Reveal>
            <Reveal delay={300}>
              <div
                style={{
                  display: "flex",
                  gap: 14,
                  marginTop: 40,
                  alignItems: "center",
                  flexWrap: "wrap",
                }}
              >
                <button
                  type="button"
                  className="lp-btn lp-btn-blue"
                  style={{ height: 56, padding: "0 28px", fontSize: 15 }}
                >
                  Googleで無料ではじめる →
                </button>
                <button
                  type="button"
                  className="lp-btn lp-btn-glass"
                  style={{ height: 56, padding: "0 24px", fontSize: 15 }}
                >
                  ▶ 触って試す
                </button>
              </div>
            </Reveal>
            <Reveal delay={400}>
              <div
                style={{
                  display: "flex",
                  gap: 28,
                  marginTop: 40,
                  fontSize: 13,
                  color: "var(--lp-ink-3)",
                  flexWrap: "wrap",
                }}
              >
                {["完全無料", "Googleログイン", "退会＝全削除"].map((t) => (
                  <div
                    key={t}
                    style={{ display: "flex", alignItems: "center", gap: 6 }}
                  >
                    <span style={{ color: "var(--lp-blue)", fontWeight: 800 }}>
                      ✓
                    </span>{" "}
                    {t}
                  </div>
                ))}
              </div>
            </Reveal>
          </div>

          <div
            style={{
              position: "relative",
              minHeight: 560,
              display: "grid",
              placeItems: "center",
            }}
          >
            <div
              className="lp-glass-strong"
              style={{
                position: "absolute",
                top: 20,
                right: 30,
                zIndex: 3,
                borderRadius: 18,
                padding: "10px 14px",
                fontFamily: "var(--lp-font-serif)",
                fontSize: 15,
                fontWeight: 700,
                boxShadow: "0 12px 32px -12px rgba(26,31,46,0.25)",
                transform: "rotate(3deg)",
              }}
            >
              次、<span style={{ color: "var(--lp-blue)" }}>なにやるんだっけ？</span>
              <svg
                style={{ position: "absolute", bottom: -14, left: 30 }}
                width="24"
                height="18"
                viewBox="0 0 24 18"
              >
                <path d="M2 0 L22 0 L8 16 Z" fill="rgba(255,255,255,0.78)" />
              </svg>
            </div>

            <Parallax speed={-0.1}>
              <div
                className="lp-float-y"
                style={{ position: "relative", zIndex: 2 }}
              >
                <Mascot size={280} mood="happy" />
              </div>
            </Parallax>

            <div
              className="lp-float-y lp-glass-strong"
              style={{
                position: "absolute",
                top: 40,
                left: 0,
                zIndex: 1,
                borderRadius: 18,
                padding: 14,
                width: 180,
                boxShadow: "0 16px 40px -16px rgba(26,31,46,0.2)",
                transform: "rotate(-6deg)",
                animationDelay: "1s",
              }}
            >
              <div
                style={{
                  fontSize: 10,
                  fontFamily: "var(--lp-font-mono)",
                  color: "var(--lp-ink-3)",
                }}
              >
                締切まで
              </div>
              <div
                style={{
                  fontFamily: "var(--lp-font-serif)",
                  fontSize: 30,
                  fontWeight: 700,
                  lineHeight: 1,
                  color: "var(--lp-blue)",
                }}
              >
                23
                <span style={{ fontSize: 14, color: "var(--lp-ink-3)" }}>:59</span>
              </div>
              <div
                style={{
                  fontSize: 11,
                  color: "var(--lp-ink-2)",
                  marginTop: 6,
                }}
              >
                オリーブ商事 / ES
              </div>
            </div>

            <div
              className="lp-float-y"
              style={{
                position: "absolute",
                bottom: 60,
                right: 10,
                zIndex: 1,
                background: "linear-gradient(135deg, var(--lp-sun) 0%, #FFE88B 100%)",
                borderRadius: 18,
                padding: 14,
                width: 170,
                boxShadow: "0 16px 40px -12px rgba(255,217,74,0.5)",
                transform: "rotate(5deg)",
                animationDelay: "0.4s",
                border: "1px solid rgba(26,31,46,0.06)",
              }}
            >
              <div style={{ fontSize: 10, fontFamily: "var(--lp-font-mono)" }}>
                今、<b>{count}</b>件進行中
              </div>
              <div
                style={{
                  fontFamily: "var(--lp-font-serif)",
                  fontSize: 16,
                  fontWeight: 700,
                  lineHeight: 1.3,
                  marginTop: 6,
                }}
              >
                ダッシュボード
              </div>
              <div style={{ display: "flex", gap: 3, marginTop: 8 }}>
                {[1, 2, 3, 4, 5, 6].map((i) => (
                  <div
                    key={i}
                    style={{
                      flex: 1,
                      height: 6,
                      borderRadius: 3,
                      background:
                        i <= 4 ? "var(--lp-ink)" : "rgba(26,31,46,0.12)",
                    }}
                  />
                ))}
              </div>
            </div>

            <div
              className="lp-wiggle"
              style={{ position: "absolute", bottom: 0, left: 40, zIndex: 4 }}
            >
              <Stamp text="FREE" color="var(--lp-blue)" size={90} />
            </div>
          </div>
        </div>
      </div>

      <style>{`
        .lp-hero-grid { display: grid; grid-template-columns: 1.1fr 0.9fr; gap: 48px; align-items: center; }
        @media (max-width: 880px) {
          .lp-hero-grid { grid-template-columns: 1fr; gap: 40px; }
          .lp-nav-links { display: none !important; }
        }
      `}</style>
    </section>
  );
}

// ── MARQUEE ──────────────────────────────────────────────────────────────
function MarqueeStrip() {
  const companies = [
    { n: "オリーブ商事", c1: "#4A6CF7", c2: "#7B95FF", g: "O" },
    { n: "パネル製作所", c1: "#E85B5B", c2: "#FF8A8A", g: "P" },
    { n: "ブリック出版", c1: "#F0A830", c2: "#FFD94A", g: "B" },
    { n: "ソラノ航空", c1: "#62A8E8", c2: "#A8D8F5", g: "S" },
    { n: "モモイロ食品", c1: "#F09AB8", c2: "#FFC0D4", g: "M" },
    { n: "ハルカゼ証券", c1: "#62D8B6", c2: "#C8F1E2", g: "H" },
    { n: "ナナメ製薬", c1: "#9A7BE8", c2: "#C8B4F0", g: "N" },
    { n: "アカネ鉄道", c1: "#E85B5B", c2: "#FFB398", g: "A" },
    { n: "キリコ精密", c1: "#1A1F2E", c2: "#3B4358", g: "K" },
    { n: "ユメカ薬品", c1: "#FFB398", c2: "#FFD4B8", g: "Y" },
    { n: "テラス不動産", c1: "#62D8B6", c2: "#8FE5C9", g: "T" },
    { n: "ミドリ農産", c1: "#5CA85C", c2: "#8FCC8F", g: "Mi" },
    { n: "シロヤマ建設", c1: "#6B7389", c2: "#9AA2B5", g: "Sh" },
    { n: "フジノ物流", c1: "#4A6CF7", c2: "#62A8E8", g: "F" },
    { n: "ウミノ水産", c1: "#0891B2", c2: "#22D3EE", g: "U" },
    { n: "クレマチス化成", c1: "#9A7BE8", c2: "#D4BFF5", g: "C" },
  ];

  const Tile = ({ c }: { c: (typeof companies)[number] }) => (
    <div
      style={{
        display: "inline-flex",
        alignItems: "center",
        gap: 12,
        padding: "12px 20px 12px 12px",
        borderRadius: 14,
        background: "rgba(255,255,255,0.7)",
        backdropFilter: "blur(16px) saturate(160%)",
        WebkitBackdropFilter: "blur(16px) saturate(160%)",
        border: "1px solid rgba(255,255,255,0.9)",
        boxShadow: "0 4px 16px -4px rgba(26,31,46,0.08)",
        whiteSpace: "nowrap",
        flexShrink: 0,
      }}
    >
      <div
        style={{
          width: 36,
          height: 36,
          borderRadius: 10,
          background: `linear-gradient(135deg, ${c.c1}, ${c.c2})`,
          display: "grid",
          placeItems: "center",
          color: "#fff",
          fontWeight: 800,
          fontSize: 14,
          fontFamily: "var(--lp-font-serif)",
          boxShadow: `0 4px 12px -4px ${c.c1}66`,
          flexShrink: 0,
        }}
      >
        {c.g}
      </div>
      <div style={{ fontSize: 14, fontWeight: 700, color: "var(--lp-ink)" }}>
        {c.n}
      </div>
    </div>
  );

  const row1 = companies.slice(0, 8);
  const row2 = companies.slice(8);

  const Lane = ({
    items,
    reverse,
  }: {
    items: typeof companies;
    reverse?: boolean;
  }) => (
    <div className={`lp-marquee ${reverse ? "lp-marquee-rev" : ""}`}>
      <div className="lp-marquee-track">
        {[...items, ...items, ...items].map((c, i) => (
          <Tile key={i} c={c} />
        ))}
      </div>
    </div>
  );

  return (
    <section style={{ padding: "60px 0 40px", position: "relative" }}>
      <div
        style={{
          textAlign: "center",
          marginBottom: 28,
          fontSize: 11,
          color: "var(--lp-ink-3)",
          fontFamily: "var(--lp-font-mono)",
          letterSpacing: "0.2em",
          fontWeight: 600,
        }}
      >
        ⎯⎯ あなたのエントリー企業、ぜんぶここに ⎯⎯
      </div>
      <div style={{ display: "grid", gap: 14 }}>
        <Lane items={row1} />
        <Lane items={row2} reverse />
      </div>
    </section>
  );
}

// ── PROBLEM ──────────────────────────────────────────────────────────────
function ProblemSection() {
  const pains = [
    { q: "どの会社が今どのフェーズだっけ？", a: "サイト横断で「Entry」単位に統合" },
    { q: "ES締切、気づいたら過ぎてた…", a: "Taskで作業も予定も一元管理・通知" },
    { q: "マイページURL、どこに保存した？", a: "証跡Clipを応募に紐付けて保存" },
    { q: "スプシに続かない…", a: "拡張で1クリック、失敗しても必ず保存" },
  ];
  const rots = [-1.2, 0.8, -0.6, 1.4];

  return (
    <section style={{ position: "relative", padding: "140px 32px 100px" }}>
      <div style={{ maxWidth: 1240, margin: "0 auto" }}>
        <Reveal>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: 14,
              marginBottom: 20,
            }}
          >
            <span
              style={{ width: 28, height: 2, background: "var(--lp-blue)" }}
            />
            <span
              style={{
                fontSize: 12,
                fontWeight: 800,
                color: "var(--lp-blue)",
                letterSpacing: "0.12em",
              }}
            >
              WHY · なぜ散らかる？
            </span>
          </div>
        </Reveal>
        <Reveal delay={100}>
          <h2
            style={{
              fontFamily: "var(--lp-font-serif)",
              fontSize: "clamp(40px, 6vw, 72px)",
              lineHeight: 1.1,
              margin: 0,
              letterSpacing: "-0.035em",
              fontWeight: 700,
              maxWidth: 900,
            }}
          >
            ひとつ、ふたつ、
            <br />
            気がつけば<span className="lp-hl-blue">30社</span>。
          </h2>
        </Reveal>
        <Reveal delay={200}>
          <p
            style={{
              fontSize: 17,
              lineHeight: 1.85,
              color: "var(--lp-ink-2)",
              marginTop: 24,
              maxWidth: 640,
            }}
          >
            応募するサイトはひとつじゃない。企業の採用ページ、ATS（i-web/SONAR）、スカウトメール。
            「どこに何があるか」は、どんどん分からなくなる。
          </p>
        </Reveal>

        <div className="lp-pain-grid">
          {pains.map((p, i) => (
            <Reveal key={i} delay={i * 100}>
              <PainCard rot={rots[i]} index={i} pain={p} />
            </Reveal>
          ))}
        </div>
      </div>

      <style>{`
        .lp-pain-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-top: 64px; }
        @media (max-width: 1024px) { .lp-pain-grid { grid-template-columns: repeat(2, 1fr); } }
        @media (max-width: 560px) { .lp-pain-grid { grid-template-columns: 1fr; } }
      `}</style>
    </section>
  );
}

function PainCard({
  rot,
  index,
  pain,
}: {
  rot: number;
  index: number;
  pain: { q: string; a: string };
}) {
  const base: CSSProperties = {
    padding: 26,
    display: "flex",
    flexDirection: "column",
    gap: 18,
    height: "100%",
    borderRadius: 20,
    boxShadow: "0 12px 32px -16px rgba(26,31,46,0.15)",
    transform: `rotate(${rot}deg)`,
    transition:
      "transform 220ms cubic-bezier(0.34,1.56,0.64,1), box-shadow 220ms",
  };
  return (
    <div
      className="lp-glass-strong"
      style={base}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = "rotate(0) translateY(-6px)";
        e.currentTarget.style.boxShadow =
          "0 24px 48px -16px rgba(26,31,46,0.25)";
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = `rotate(${rot}deg)`;
        e.currentTarget.style.boxShadow =
          "0 12px 32px -16px rgba(26,31,46,0.15)";
      }}
    >
      <div
        style={{
          fontFamily: "var(--lp-font-mono)",
          fontSize: 11,
          color: "var(--lp-ink-3)",
        }}
      >
        0{index + 1} / 04
      </div>
      <div
        style={{
          fontFamily: "var(--lp-font-serif)",
          fontSize: 22,
          lineHeight: 1.4,
          fontWeight: 700,
        }}
      >
        「{pain.q}」
      </div>
      <div
        style={{
          marginTop: "auto",
          padding: "12px 14px",
          borderRadius: 12,
          background: "rgba(74,108,247,0.08)",
          borderLeft: "3px solid var(--lp-blue)",
          fontSize: 13,
          lineHeight: 1.6,
          color: "var(--lp-ink)",
          fontWeight: 500,
        }}
      >
        → {pain.a}
      </div>
    </div>
  );
}

// ── MASCOT SHOWCASE ──────────────────────────────────────────────────────
function MascotShowcase() {
  const moods: { m: Mood; t: string; c: string }[] = [
    { m: "happy", t: "いつもゴキゲン", c: "#FFC83D" },
    { m: "bow", t: "よろしくお願いします", c: "#F5B6C7" },
    { m: "wow", t: "締切が近いと驚く", c: "#8FD9A8" },
    { m: "wink", t: "Entryが増えると喜ぶ", c: "#7ECDF5" },
    { m: "cheer", t: "内定で大喜び", c: "#FF9A6C" },
    { m: "sleep", t: "早朝はまだ眠い", c: "#C8B8EC" },
  ];
  return (
    <section
      style={{ padding: "120px 32px", position: "relative", overflow: "hidden" }}
    >
      <Parallax
        speed={0.3}
        axis="x"
        style={{
          position: "absolute",
          inset: 0,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          pointerEvents: "none",
          zIndex: 0,
        }}
      >
        <div
          aria-hidden
          style={{
            fontFamily: "var(--lp-font-serif)",
            fontSize: 280,
            fontWeight: 700,
            color: "rgba(74,108,247,0.06)",
            letterSpacing: "-0.04em",
            whiteSpace: "nowrap",
          }}
        >
          シカくん · シカくん · シカくん
        </div>
      </Parallax>

      <div
        style={{
          maxWidth: 1240,
          margin: "0 auto",
          position: "relative",
          zIndex: 1,
        }}
      >
        <div className="lp-mascot-grid">
          <div>
            <Reveal>
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: 14,
                  marginBottom: 20,
                }}
              >
                <span
                  style={{ width: 28, height: 2, background: "var(--lp-sun)" }}
                />
                <span
                  style={{
                    fontSize: 12,
                    fontWeight: 800,
                    color: "var(--lp-ink-2)",
                    letterSpacing: "0.12em",
                  }}
                >
                  MASCOT
                </span>
              </div>
            </Reveal>
            <Reveal delay={100}>
              <h2
                style={{
                  fontFamily: "var(--lp-font-serif)",
                  fontSize: "clamp(44px, 6vw, 80px)",
                  lineHeight: 1.05,
                  margin: 0,
                  letterSpacing: "-0.035em",
                  fontWeight: 700,
                }}
              >
                はじめまして、
                <br />
                <span
                  style={{
                    background:
                      "linear-gradient(135deg, #F0B800 0%, #FFD94A 100%)",
                    WebkitBackgroundClip: "text",
                    backgroundClip: "text",
                    color: "transparent",
                    WebkitTextFillColor: "transparent",
                  }}
                >
                  シカくん
                </span>
                です。
              </h2>
            </Reveal>
            <Reveal delay={200}>
              <p
                style={{
                  fontSize: 17,
                  lineHeight: 1.85,
                  color: "var(--lp-ink-2)",
                  marginTop: 24,
                  maxWidth: 520,
                }}
              >
                応募（Entry）を集める封筒型マスコット。
                <br />
                ダッシュボードで見守ってくれて、締切前には
                <b style={{ color: "var(--lp-ink)" }}>ちゃんと慌ててくれる</b>。
                就活は、ひとりじゃなくてもいい。
              </p>
            </Reveal>
            <Reveal delay={300}>
              <div
                style={{
                  display: "grid",
                  gridTemplateColumns: "repeat(2, 1fr)",
                  gap: 14,
                  marginTop: 40,
                }}
              >
                {moods.map((m) => (
                  <div
                    key={m.m}
                    className="lp-glass-strong"
                    style={{
                      borderRadius: 16,
                      padding: 16,
                      display: "flex",
                      alignItems: "center",
                      gap: 14,
                    }}
                  >
                    <div
                      style={{
                        background: `color-mix(in srgb, ${m.c} 35%, transparent)`,
                        borderRadius: 12,
                        padding: 6,
                        boxShadow: "0 4px 12px -4px rgba(26,31,46,0.2)",
                        border: `1.5px solid ${m.c}`,
                      }}
                    >
                      <Mascot size={48} mood={m.m} />
                    </div>
                    <div style={{ fontSize: 13, fontWeight: 600 }}>{m.t}</div>
                  </div>
                ))}
              </div>
            </Reveal>
          </div>

          <MascotCarousel />
        </div>
      </div>

      <style>{`
        .lp-mascot-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 64px; align-items: center; }
        @media (max-width: 880px) { .lp-mascot-grid { grid-template-columns: 1fr; gap: 40px; } }
      `}</style>
    </section>
  );
}

function MascotCarousel() {
  const moods: Mood[] = ["happy", "bow", "wink", "wow", "cheer", "sleep"];
  const [i, setI] = useState(0);
  useEffect(() => {
    const id = setInterval(() => setI((x) => (x + 1) % moods.length), 1800);
    return () => clearInterval(id);
  }, [moods.length]);
  return (
    <div
      style={{
        position: "relative",
        display: "grid",
        placeItems: "center",
        minHeight: 440,
      }}
    >
      <div
        style={{
          width: 380,
          height: 380,
          borderRadius: "50%",
          background:
            "radial-gradient(circle, var(--lp-sun) 0%, var(--lp-peach) 70%, transparent 100%)",
          opacity: 0.25,
          position: "absolute",
          filter: "blur(40px)",
        }}
      />
      <div style={{ position: "relative" }} className="lp-float-y">
        <Mascot size={360} mood={moods[i]} />
      </div>
      <div
        className="lp-spin-slow"
        style={{
          position: "absolute",
          width: 520,
          height: 520,
          pointerEvents: "none",
        }}
      >
        <div
          style={{
            position: "absolute",
            top: 0,
            left: "50%",
            transform: "translateX(-50%)",
          }}
        >
          <Stamp text="NEW" color="var(--lp-sun)" size={64} />
        </div>
        <div style={{ position: "absolute", bottom: 10, right: 0 }}>
          <Stamp text="BETA" color="var(--lp-peach)" size={58} />
        </div>
      </div>
    </div>
  );
}

// ── FEATURE TOUR ─────────────────────────────────────────────────────────
function FeatureTour() {
  const [tab, setTab] = useState<"entry" | "task" | "inbox" | "ext">("entry");
  const tabs = [
    { id: "entry", label: "Entry", sub: "応募", caption: "応募単位で、フェーズ・メモ・履歴を集約" },
    { id: "task", label: "Task", sub: "締切/予定", caption: "作業と予定をまとめて、漏れゼロへ" },
    { id: "inbox", label: "Inbox", sub: "未整理", caption: "保存は常に成功。あとから整理すればOK" },
    { id: "ext", label: "拡張", sub: "Chrome", caption: "閲覧中のページから、1クリック保存" },
  ] as const;
  const details: Record<typeof tab, string[]> = {
    entry: [
      "会社 / 応募ルート / source（媒体）を1行で集約",
      "現在Stage（応募・書類・面接…）をワンタップ更新",
      "Stage履歴はタイムラインで後から見返せる",
    ],
    task: [
      "「締切（deadline）」と「予定（schedule）」を区別",
      "ダッシュボードで直近の予定と近い締切を別枠表示",
      "メール通知は3日前・24時間前・日次まとめ",
    ],
    inbox: [
      "Entry未選択でも、URL + タイトル + 時刻で保存成功",
      "候補企業を推定 → ワンクリックで割当",
      "30日で自動アーカイブ、Inbox件数を常時表示",
    ],
    ext: [
      "ドメイン + タイトル + 固定語で source を自動推定",
      "候補提示のみ、確定はあなたの操作",
      "抽出失敗でもInboxに必ず退避する",
    ],
  };
  const current = tabs.find((t) => t.id === tab)!;

  return (
    <section id="features" style={{ padding: "120px 32px", position: "relative" }}>
      <div style={{ maxWidth: 1240, margin: "0 auto" }}>
        <Reveal>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: 14,
              marginBottom: 20,
            }}
          >
            <span
              style={{ width: 28, height: 2, background: "var(--lp-peach)" }}
            />
            <span
              style={{
                fontSize: 12,
                fontWeight: 800,
                color: "var(--lp-peach)",
                letterSpacing: "0.12em",
              }}
            >
              FEATURES · 4つの道具
            </span>
          </div>
        </Reveal>
        <Reveal delay={100}>
          <h2
            style={{
              fontFamily: "var(--lp-font-serif)",
              fontSize: "clamp(40px, 6vw, 72px)",
              lineHeight: 1.1,
              margin: 0,
              letterSpacing: "-0.035em",
              fontWeight: 700,
            }}
          >
            就活の動かし方、
            <br />
            <span className="lp-hl-mint">4つの道具</span>で足りる。
          </h2>
        </Reveal>

        <div
          style={{
            display: "flex",
            gap: 10,
            margin: "48px 0 28px",
            flexWrap: "wrap",
          }}
        >
          {tabs.map((t) => (
            <button
              key={t.id}
              onClick={() => setTab(t.id)}
              style={{
                padding: "14px 22px",
                borderRadius: 14,
                border: "1.5px solid var(--lp-ink)",
                background:
                  tab === t.id ? "var(--lp-sun)" : "var(--lp-surface)",
                color: "var(--lp-ink)",
                fontWeight: 700,
                fontSize: 14,
                fontFamily: "inherit",
                cursor: "pointer",
                boxShadow:
                  tab === t.id ? "4px 4px 0 var(--lp-ink)" : "2px 2px 0 var(--lp-ink)",
                transform: tab === t.id ? "translate(-2px, -2px)" : "none",
                transition: "all 180ms",
                display: "flex",
                alignItems: "baseline",
                gap: 8,
              }}
            >
              <span>{t.label}</span>
              <span
                style={{
                  fontSize: 11,
                  fontWeight: 500,
                  color: "var(--lp-ink-2)",
                }}
              >
                {t.sub}
              </span>
            </button>
          ))}
        </div>

        <div
          className="lp-bold-card lp-feature-panel"
          style={{ overflow: "hidden", minHeight: 460 }}
        >
          <div
            style={{
              padding: 48,
              display: "flex",
              flexDirection: "column",
              justifyContent: "center",
              borderRight: "1px solid rgba(26,31,46,0.08)",
            }}
          >
            <div
              style={{
                fontSize: 12,
                color: "var(--lp-ink-3)",
                marginBottom: 14,
                fontFamily: "var(--lp-font-mono)",
                letterSpacing: "0.1em",
              }}
            >
              {current.label.toUpperCase()}
            </div>
            <h3
              style={{
                fontFamily: "var(--lp-font-serif)",
                fontSize: 36,
                lineHeight: 1.3,
                margin: 0,
                fontWeight: 700,
                letterSpacing: "-0.02em",
              }}
            >
              {current.caption}
            </h3>
            <ul
              style={{
                listStyle: "none",
                padding: 0,
                margin: "28px 0 0",
                display: "grid",
                gap: 14,
              }}
            >
              {details[tab].map((p) => (
                <li
                  key={p}
                  style={{
                    display: "flex",
                    gap: 12,
                    fontSize: 15,
                    lineHeight: 1.6,
                    color: "var(--lp-ink-2)",
                  }}
                >
                  <span
                    style={{
                      flexShrink: 0,
                      width: 22,
                      height: 22,
                      borderRadius: 999,
                      background: "var(--lp-sun)",
                      color: "var(--lp-ink)",
                      border: "1.5px solid var(--lp-ink)",
                      display: "grid",
                      placeItems: "center",
                      fontSize: 12,
                      fontWeight: 800,
                      marginTop: 1,
                    }}
                  >
                    ✓
                  </span>
                  {p}
                </li>
              ))}
            </ul>
          </div>
          <div
            style={{
              background: "rgba(255,255,255,0.15)",
              padding: 32,
              display: "grid",
              placeItems: "center",
            }}
          >
            <FeatureVisual tab={tab} />
          </div>
        </div>
      </div>

      <style>{`
        .lp-feature-panel { display: grid; grid-template-columns: 0.9fr 1.1fr; }
        @media (max-width: 880px) { .lp-feature-panel { grid-template-columns: 1fr; } }
      `}</style>
    </section>
  );
}

function FeatureVisual({ tab }: { tab: "entry" | "task" | "inbox" | "ext" }) {
  if (tab === "entry") return <FVEntry />;
  if (tab === "task") return <FVTask />;
  if (tab === "inbox") return <FVInbox />;
  return <FVExt />;
}

const SAMPLE_ENTRIES = [
  {
    id: "e1",
    company: "パネル製作所",
    logo: "🟧",
    route: "本選考",
    source: "リクナビ",
    stageLabel: "二次面接",
    color: "var(--lp-s-interview)",
  },
  {
    id: "e2",
    company: "オリーブ商事",
    logo: "🫒",
    route: "本選考",
    source: "マイナビ",
    stageLabel: "ES提出",
    color: "var(--lp-s-doc)",
  },
  {
    id: "e3",
    company: "ブリック出版",
    logo: "📕",
    route: "本選考",
    source: "企業HP",
    stageLabel: "SPI",
    color: "var(--lp-s-test)",
  },
  {
    id: "e4",
    company: "キナリ設計",
    logo: "🪵",
    route: "インターン",
    source: "ONE CAREER",
    stageLabel: "書類選考中",
    color: "var(--lp-s-app)",
  },
  {
    id: "e5",
    company: "モモイロ食品",
    logo: "🍑",
    route: "本選考",
    source: "OfferBox",
    stageLabel: "最終面接",
    color: "var(--lp-s-interview)",
  },
];

function FVEntry() {
  return (
    <div style={{ width: "100%", maxWidth: 420, display: "grid", gap: 12 }}>
      {SAMPLE_ENTRIES.map((e, i) => (
        <div
          key={e.id}
          style={{
            background: "var(--lp-surface)",
            borderRadius: 14,
            padding: "14px 16px",
            border: "1.5px solid var(--lp-ink)",
            display: "flex",
            alignItems: "center",
            gap: 12,
            boxShadow: "3px 3px 0 var(--lp-ink)",
            animation: `lpSlideIn 0.5s ease ${i * 0.08}s both`,
          }}
        >
          <div
            style={{
              width: 36,
              height: 36,
              borderRadius: 10,
              background: "var(--lp-surface-2)",
              border: "1px solid var(--lp-line)",
              display: "grid",
              placeItems: "center",
              fontSize: 20,
              flexShrink: 0,
            }}
          >
            {e.logo}
          </div>
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ fontSize: 14, fontWeight: 700 }}>{e.company}</div>
            <div
              style={{
                fontSize: 11,
                color: "var(--lp-ink-3)",
                marginTop: 2,
              }}
            >
              {e.route} · {e.source}
            </div>
          </div>
          <span
            style={{
              display: "inline-flex",
              alignItems: "center",
              gap: 6,
              height: 22,
              padding: "0 8px",
              borderRadius: 999,
              fontSize: 11,
              fontWeight: 600,
              background: "var(--lp-bg-alt)",
              border: "1px solid var(--lp-line)",
            }}
          >
            <span
              style={{
                width: 6,
                height: 6,
                borderRadius: 999,
                background: e.color,
              }}
            />
            {e.stageLabel}
          </span>
        </div>
      ))}
    </div>
  );
}

function FVTask() {
  const [checked, setChecked] = useState<number[]>([]);
  const tasks = [
    { c: "オリーブ商事", t: "ES提出", d: "明日 23:59", urgent: true },
    { c: "パネル製作所", t: "二次面接", d: "4/20 14:00", urgent: true },
    { c: "ブリック出版", t: "SPI受検", d: "4/24 23:59", urgent: false },
    { c: "ソラノ航空", t: "GD", d: "4/26 10:00", urgent: false },
    { c: "モモイロ食品", t: "最終面接", d: "4/28 15:30", urgent: false },
  ];
  useEffect(() => {
    let i = 0;
    const id = setInterval(() => {
      setChecked((c) => (c.length >= 3 ? [] : [...c, i]));
      i = (i + 1) % tasks.length;
    }, 1400);
    return () => clearInterval(id);
  }, [tasks.length]);
  return (
    <div className="lp-bold-card" style={{ width: "100%", maxWidth: 420, padding: 20 }}>
      <div
        style={{
          fontSize: 12,
          color: "var(--lp-ink-3)",
          fontWeight: 700,
          marginBottom: 14,
          fontFamily: "var(--lp-font-mono)",
        }}
      >
        今週のタスク
      </div>
      {tasks.map((t, i) => {
        const done = checked.includes(i);
        return (
          <div
            key={i}
            style={{
              display: "flex",
              alignItems: "center",
              gap: 12,
              padding: "12px 0",
              borderTop: i ? "1px solid var(--lp-line-2)" : "none",
              opacity: done ? 0.4 : 1,
              transition: "opacity 300ms",
            }}
          >
            <div
              style={{
                width: 20,
                height: 20,
                borderRadius: 6,
                border: "2px solid var(--lp-ink)",
                background: done ? "var(--lp-sun)" : "var(--lp-surface)",
                display: "grid",
                placeItems: "center",
                fontWeight: 800,
                fontSize: 12,
                transition: "all 220ms",
              }}
            >
              {done && "✓"}
            </div>
            <div style={{ flex: 1 }}>
              <div
                style={{
                  fontSize: 14,
                  fontWeight: 600,
                  textDecoration: done ? "line-through" : "none",
                }}
              >
                {t.t}
              </div>
              <div style={{ fontSize: 11, color: "var(--lp-ink-3)" }}>
                {t.c}
              </div>
            </div>
            <span
              style={{
                display: "inline-flex",
                alignItems: "center",
                gap: 6,
                height: 22,
                padding: "0 8px",
                borderRadius: 6,
                fontSize: 11,
                fontWeight: 600,
                fontFamily: "var(--lp-font-mono)",
                background: t.urgent
                  ? "rgba(255,179,152,0.2)"
                  : "var(--lp-surface-2)",
                color: t.urgent ? "#B8542A" : "var(--lp-ink-2)",
                border: `1px solid ${t.urgent ? "rgba(255,179,152,0.5)" : "var(--lp-line)"}`,
              }}
            >
              {t.d}
            </span>
          </div>
        );
      })}
    </div>
  );
}

function FVInbox() {
  const items = [
    {
      site: "リクナビ",
      title: "エンジニア職 本選考 | パネル製作所",
      guess: "パネル製作所",
      confident: true,
    },
    {
      site: "i-web",
      title: "マイページ | 選考スケジュール",
      guess: "オリーブ商事（候補）",
      confident: false,
    },
    {
      site: "ONE CAREER",
      title: "面接レポート - ブリック出版",
      guess: "ブリック出版",
      confident: true,
    },
  ];
  return (
    <div style={{ width: "100%", maxWidth: 420, display: "grid", gap: 10 }}>
      {items.map((c, i) => (
        <div key={i} className="lp-bold-card" style={{ padding: 14 }}>
          <div style={{ display: "flex", gap: 10, alignItems: "flex-start" }}>
            <div
              style={{
                width: 30,
                height: 30,
                borderRadius: 7,
                background: "var(--lp-bg-alt)",
                border: "1.5px solid var(--lp-ink)",
                display: "grid",
                placeItems: "center",
                fontSize: 13,
              }}
            >
              🔗
            </div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <div
                style={{
                  fontSize: 12,
                  fontWeight: 600,
                  whiteSpace: "nowrap",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                }}
              >
                {c.title}
              </div>
              <div
                style={{
                  fontSize: 10,
                  color: "var(--lp-ink-3)",
                  marginTop: 2,
                  fontFamily: "var(--lp-font-mono)",
                }}
              >
                {c.site}
              </div>
              <div style={{ marginTop: 8, display: "flex", gap: 6 }}>
                <button
                  style={{
                    height: 28,
                    padding: "0 12px",
                    borderRadius: 8,
                    border: "1.5px solid var(--lp-ink)",
                    background: c.confident ? "var(--lp-sun)" : "var(--lp-surface)",
                    color: "var(--lp-ink)",
                    fontSize: 11,
                    fontWeight: 700,
                    cursor: "pointer",
                    fontFamily: "inherit",
                  }}
                >
                  {c.confident ? "→ " : ""}
                  {c.guess}
                </button>
                <button
                  style={{
                    height: 28,
                    padding: "0 10px",
                    borderRadius: 8,
                    border: "1px solid var(--lp-line)",
                    background: "transparent",
                    fontSize: 11,
                    cursor: "pointer",
                    fontFamily: "inherit",
                  }}
                >
                  別のEntryに
                </button>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

function FVExt() {
  return (
    <div className="lp-bold-card" style={{ width: 320, overflow: "hidden" }}>
      <div
        style={{
          background: "var(--lp-ink)",
          color: "#fff",
          padding: "10px 14px",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
        }}
      >
        <div
          style={{
            fontSize: 12,
            fontWeight: 700,
            display: "flex",
            alignItems: "center",
            gap: 8,
          }}
        >
          <MiniMascot size={22} />
          <span>保存する</span>
        </div>
        <div style={{ fontSize: 10, opacity: 0.6 }}>⌘ + S</div>
      </div>
      <div style={{ padding: 14 }}>
        <div
          style={{
            fontSize: 10,
            color: "var(--lp-ink-3)",
            fontFamily: "var(--lp-font-mono)",
            marginBottom: 4,
          }}
        >
          DETECTED PAGE
        </div>
        <div style={{ fontSize: 12, fontWeight: 600, marginBottom: 10 }}>
          【25卒】エンジニア本選考｜パネル製作所
        </div>
        <div
          style={{
            fontSize: 10,
            color: "var(--lp-ink-3)",
            marginBottom: 14,
            fontFamily: "var(--lp-font-mono)",
          }}
        >
          recruit.example.com/jobs/3921
        </div>
        <div
          style={{
            fontSize: 10,
            color: "var(--lp-ink-3)",
            fontFamily: "var(--lp-font-mono)",
            marginBottom: 6,
          }}
        >
          COMPANY (GUESS)
        </div>
        <div
          style={{
            padding: "8px 10px",
            background: "var(--lp-sun)",
            borderRadius: 8,
            border: "1.5px solid var(--lp-ink)",
            fontSize: 13,
            fontWeight: 700,
            marginBottom: 10,
          }}
        >
          パネル製作所{" "}
          <span style={{ fontSize: 10, fontWeight: 600 }}>・高信頼度</span>
        </div>
        <div
          style={{
            fontSize: 10,
            color: "var(--lp-ink-3)",
            fontFamily: "var(--lp-font-mono)",
            marginBottom: 6,
          }}
        >
          SOURCE
        </div>
        <div
          style={{
            padding: "8px 10px",
            background: "var(--lp-bg-alt)",
            borderRadius: 8,
            fontSize: 12,
            marginBottom: 14,
            border: "1px solid var(--lp-line)",
          }}
        >
          リクナビ
        </div>
        <div style={{ display: "flex", gap: 6 }}>
          <button
            style={{
              flex: 1,
              height: 38,
              borderRadius: 8,
              background: "var(--lp-ink)",
              color: "#fff",
              border: "none",
              fontWeight: 700,
              fontSize: 12,
              cursor: "pointer",
              fontFamily: "inherit",
            }}
          >
            Entryに保存
          </button>
          <button
            style={{
              height: 38,
              padding: "0 12px",
              borderRadius: 8,
              border: "1.5px solid var(--lp-ink)",
              background: "var(--lp-surface)",
              fontSize: 12,
              cursor: "pointer",
              fontFamily: "inherit",
              fontWeight: 700,
            }}
          >
            Inbox
          </button>
        </div>
      </div>
    </div>
  );
}

// ── LIVE DASHBOARD SECTION ───────────────────────────────────────────────
function LiveDashboardSection() {
  const otherScreens = [
    { id: "entry", t: "Entry詳細", c: "var(--lp-s-interview)" },
    { id: "inbox", t: "Inbox", c: "var(--lp-s-offer)" },
    { id: "mobile", t: "モバイル", c: "var(--lp-s-doc)" },
    { id: "newentry", t: "新規Entry", c: "var(--lp-s-app)" },
    { id: "onboarding", t: "オンボ", c: "var(--lp-s-group)" },
  ];
  return (
    <section
      id="screens"
      style={{
        position: "relative",
        padding: "120px 32px",
        background: "transparent",
      }}
    >
      <div style={{ maxWidth: 1240, margin: "0 auto" }}>
        <div className="lp-dashboard-grid">
          <div>
            <Reveal>
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: 14,
                  marginBottom: 20,
                }}
              >
                <span
                  style={{ width: 28, height: 2, background: "var(--lp-peach)" }}
                />
                <span
                  style={{
                    fontSize: 12,
                    fontWeight: 800,
                    color: "var(--lp-peach)",
                    letterSpacing: "0.12em",
                  }}
                >
                  DASHBOARD
                </span>
              </div>
            </Reveal>
            <Reveal delay={100}>
              <h2
                style={{
                  fontFamily: "var(--lp-font-serif)",
                  fontSize: "clamp(40px, 5vw, 64px)",
                  lineHeight: 1.1,
                  margin: 0,
                  letterSpacing: "-0.035em",
                  fontWeight: 700,
                }}
              >
                開いたら、
                <br />
                もう整っている。
              </h2>
            </Reveal>
            <Reveal delay={200}>
              <p
                style={{
                  fontSize: 16,
                  lineHeight: 1.85,
                  color: "var(--lp-ink-2)",
                  marginTop: 24,
                }}
              >
                進行中のEntry、近い締切、直近の予定、未整理のInbox。
                <br />
                就活の「今日やること」がひと目で分かる、
                <b>1枚の台帳</b>。
              </p>
            </Reveal>
            <Reveal delay={300}>
              <button
                type="button"
                className="lp-btn lp-btn-sun"
                style={{
                  height: 52,
                  padding: "0 24px",
                  fontSize: 14,
                  marginTop: 32,
                  display: "inline-flex",
                }}
              >
                ダッシュボードを開く →
              </button>
            </Reveal>
          </div>

          <Reveal delay={200}>
            <BigDashboardPreview />
          </Reveal>
        </div>

        <Reveal>
          <div
            style={{
              marginTop: 96,
              paddingTop: 64,
              borderTop: "2px dashed var(--lp-line)",
            }}
          >
            <h3
              style={{
                fontFamily: "var(--lp-font-serif)",
                fontSize: 28,
                fontWeight: 700,
                margin: "0 0 32px",
              }}
            >
              他の画面も試せます
            </h3>
            <div
              style={{
                display: "grid",
                gridTemplateColumns:
                  "repeat(auto-fit, minmax(200px, 1fr))",
                gap: 14,
              }}
            >
              {otherScreens.map((s) => (
                <button
                  type="button"
                  key={s.id}
                  className="lp-bold-card"
                  style={{
                    padding: 20,
                    textAlign: "left",
                    cursor: "pointer",
                    fontFamily: "inherit",
                    transition: "transform 180ms cubic-bezier(0.34,1.56,0.64,1)",
                    display: "block",
                    width: "100%",
                  }}
                  onMouseEnter={(e) =>
                    (e.currentTarget.style.transform = "translate(-3px, -3px)")
                  }
                  onMouseLeave={(e) =>
                    (e.currentTarget.style.transform = "none")
                  }
                >
                  <div
                    style={{
                      width: 10,
                      height: 10,
                      borderRadius: 999,
                      background: s.c,
                      marginBottom: 16,
                    }}
                  />
                  <div
                    style={{
                      fontFamily: "var(--lp-font-serif)",
                      fontSize: 20,
                      fontWeight: 700,
                    }}
                  >
                    {s.t}
                  </div>
                  <div
                    style={{
                      fontSize: 12,
                      color: "var(--lp-peach)",
                      fontWeight: 700,
                      marginTop: 8,
                    }}
                  >
                    触ってみる →
                  </div>
                </button>
              ))}
            </div>
          </div>
        </Reveal>
      </div>

      <style>{`
        .lp-dashboard-grid { display: grid; grid-template-columns: 1fr 1.3fr; gap: 64px; align-items: center; }
        @media (max-width: 1024px) { .lp-dashboard-grid { grid-template-columns: 1fr; gap: 48px; } }
      `}</style>
    </section>
  );
}

function BigDashboardPreview() {
  const [hover, setHover] = useState<string | null>(null);
  const hotTasks = [
    { c: "オリーブ商事", t: "ES提出", d: "明日 23:59" },
    { c: "パネル製作所", t: "二次面接", d: "4/20 14:00" },
  ];
  return (
    <div className="lp-bold-card" style={{ padding: 24, position: "relative" }}>
      {/* window chrome */}
      <div style={{ display: "flex", gap: 6, marginBottom: 18 }}>
        {["#E6A494", "#E5D08A", "#A8C49A"].map((c) => (
          <div
            key={c}
            style={{
              width: 12,
              height: 12,
              borderRadius: 6,
              background: c,
              border: "1.5px solid var(--lp-ink)",
            }}
          />
        ))}
        <div
          style={{
            flex: 1,
            textAlign: "center",
            fontSize: 11,
            fontFamily: "var(--lp-font-mono)",
            color: "var(--lp-ink-3)",
          }}
        >
          dashboard
        </div>
      </div>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          marginBottom: 16,
        }}
      >
        <div
          style={{
            fontFamily: "var(--lp-font-serif)",
            fontSize: 22,
            fontWeight: 700,
          }}
        >
          こんにちは 👋
        </div>
        <div style={{ display: "flex", gap: 6 }}>
          <div
            className="lp-sticker"
            style={{ transform: "rotate(-6deg)", background: "var(--lp-sun)" }}
          >
            12件 進行中
          </div>
          <div
            className="lp-sticker"
            style={{
              transform: "rotate(3deg)",
              background: "var(--lp-peach)",
              color: "#fff",
            }}
          >
            Inbox 3
          </div>
        </div>
      </div>

      {/* Hot tasks */}
      <div
        style={{
          padding: 16,
          background: "rgba(255, 217, 74, 0.12)",
          borderRadius: 12,
          border: "1.5px solid var(--lp-ink)",
          marginBottom: 12,
        }}
      >
        <div
          style={{
            fontSize: 11,
            fontWeight: 700,
            marginBottom: 10,
            fontFamily: "var(--lp-font-mono)",
          }}
        >
          🔥 近い締切
        </div>
        {hotTasks.map((r, i) => (
          <div
            key={i}
            style={{
              display: "flex",
              alignItems: "center",
              gap: 10,
              padding: "8px 0",
              borderTop: i ? "1px solid var(--lp-line-2)" : "none",
            }}
          >
            <div
              style={{
                width: 5,
                height: 22,
                borderRadius: 3,
                background: "var(--lp-peach)",
              }}
            />
            <div style={{ fontSize: 13, fontWeight: 700, flex: 1 }}>{r.c}</div>
            <div style={{ fontSize: 12, color: "var(--lp-ink-2)" }}>{r.t}</div>
            <div
              style={{
                fontSize: 11,
                fontFamily: "var(--lp-font-mono)",
                color: "var(--lp-peach)",
                fontWeight: 700,
              }}
            >
              {r.d}
            </div>
          </div>
        ))}
      </div>

      {/* Entries grid */}
      <div style={{ display: "grid", gap: 8 }}>
        {SAMPLE_ENTRIES.slice(0, 5).map((e) => (
          <div
            key={e.id}
            onMouseEnter={() => setHover(e.id)}
            onMouseLeave={() => setHover(null)}
            style={{
              display: "flex",
              alignItems: "center",
              gap: 12,
              padding: "12px 14px",
              background:
                hover === e.id ? "rgba(255, 217, 74, 0.12)" : "var(--lp-surface)",
              borderRadius: 10,
              border: "1.5px solid var(--lp-ink)",
              transition: "all 180ms",
              cursor: "pointer",
              transform: hover === e.id ? "translateX(4px)" : "none",
            }}
          >
            <div
              style={{
                width: 32,
                height: 32,
                borderRadius: 8,
                background: "var(--lp-surface-2)",
                border: "1px solid var(--lp-line)",
                display: "grid",
                placeItems: "center",
                fontSize: 18,
                flexShrink: 0,
              }}
            >
              {e.logo}
            </div>
            <div style={{ fontSize: 13, fontWeight: 700, flex: 1 }}>
              {e.company}
            </div>
            <div style={{ fontSize: 11, color: "var(--lp-ink-3)" }}>
              {e.route}
            </div>
            <span
              style={{
                display: "inline-flex",
                alignItems: "center",
                gap: 6,
                height: 22,
                padding: "0 8px",
                borderRadius: 999,
                fontSize: 11,
                fontWeight: 600,
                background: "var(--lp-bg-alt)",
                border: "1px solid var(--lp-line)",
              }}
            >
              <span
                style={{
                  width: 6,
                  height: 6,
                  borderRadius: 999,
                  background: e.color,
                }}
              />
              {e.stageLabel}
            </span>
          </div>
        ))}
      </div>

      {/* Mascot floating */}
      <div
        className="lp-bob"
        style={{ position: "absolute", bottom: -30, right: -20, zIndex: 2 }}
      >
        <Mascot size={100} mood="happy" />
      </div>
    </div>
  );
}

// ── EXTENSION DEMO ───────────────────────────────────────────────────────
function ExtensionDemo() {
  const steps = [
    { k: "①", t: "就活サイトで気になる企業を開く" },
    { k: "②", t: "拡張アイコンをクリック → 候補が自動提示" },
    { k: "③", t: "Entryに割当 or Inboxへ、1クリック保存" },
  ];
  return (
    <section id="extension" style={{ padding: "120px 32px" }}>
      <div style={{ maxWidth: 1240, margin: "0 auto" }}>
        <div className="lp-extension-grid">
          <Reveal>
            <div className="lp-bold-card" style={{ overflow: "hidden" }}>
              {/* Fake browser chrome */}
              <div
                style={{
                  background: "var(--lp-bg-alt)",
                  padding: 12,
                  borderBottom: "1.5px solid var(--lp-ink)",
                  display: "flex",
                  gap: 10,
                  alignItems: "center",
                }}
              >
                <div style={{ display: "flex", gap: 6 }}>
                  {["#E6A494", "#E5D08A", "#A8C49A"].map((c) => (
                    <div
                      key={c}
                      style={{
                        width: 10,
                        height: 10,
                        borderRadius: 5,
                        background: c,
                        border: "1.5px solid var(--lp-ink)",
                      }}
                    />
                  ))}
                </div>
                <div
                  style={{
                    flex: 1,
                    background: "var(--lp-surface)",
                    borderRadius: 999,
                    padding: "4px 12px",
                    fontSize: 11,
                    color: "var(--lp-ink-3)",
                    fontFamily: "var(--lp-font-mono)",
                    border: "1.5px solid var(--lp-ink)",
                    whiteSpace: "nowrap",
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                  }}
                >
                  🔒 recruit.example.com/jobs/3921
                </div>
                <div
                  className="lp-wiggle"
                  style={{
                    width: 30,
                    height: 30,
                    borderRadius: 8,
                    background: "var(--lp-sun)",
                    display: "grid",
                    placeItems: "center",
                    border: "1.5px solid var(--lp-ink)",
                    flexShrink: 0,
                  }}
                >
                  <MiniMascot size={22} />
                </div>
              </div>
              <div
                style={{
                  padding: 32,
                  background: "#fff",
                  position: "relative",
                  minHeight: 380,
                }}
              >
                <div
                  style={{
                    fontSize: 10,
                    color: "var(--lp-ink-3)",
                    fontFamily: "var(--lp-font-mono)",
                  }}
                >
                  リクナビ — 求人詳細
                </div>
                <h3
                  style={{
                    fontSize: 22,
                    fontWeight: 700,
                    margin: "8px 0 4px",
                  }}
                >
                  【25卒】エンジニア職 本選考
                </h3>
                <div style={{ fontSize: 13, color: "var(--lp-ink-2)" }}>
                  パネル製作所株式会社
                </div>
                <div style={{ marginTop: 16, display: "grid", gap: 6 }}>
                  <div
                    style={{
                      height: 8,
                      background: "var(--lp-bg-alt)",
                      borderRadius: 4,
                      width: "80%",
                    }}
                  />
                  <div
                    style={{
                      height: 8,
                      background: "var(--lp-bg-alt)",
                      borderRadius: 4,
                      width: "65%",
                    }}
                  />
                  <div
                    style={{
                      height: 8,
                      background: "var(--lp-bg-alt)",
                      borderRadius: 4,
                      width: "72%",
                    }}
                  />
                  <div
                    style={{
                      height: 8,
                      background: "var(--lp-bg-alt)",
                      borderRadius: 4,
                      width: "58%",
                    }}
                  />
                </div>
                {/* Popup */}
                <div
                  className="lp-float-y"
                  style={{ position: "absolute", top: 20, right: 20, width: 280 }}
                >
                  <FVExt />
                </div>
              </div>
            </div>
          </Reveal>

          <div>
            <Reveal>
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: 14,
                  marginBottom: 20,
                }}
              >
                <span
                  style={{ width: 28, height: 2, background: "var(--lp-peach)" }}
                />
                <span
                  style={{
                    fontSize: 12,
                    fontWeight: 800,
                    color: "var(--lp-peach)",
                    letterSpacing: "0.12em",
                  }}
                >
                  CHROME EXTENSION
                </span>
              </div>
            </Reveal>
            <Reveal delay={100}>
              <h2
                style={{
                  fontFamily: "var(--lp-font-serif)",
                  fontSize: "clamp(36px, 5vw, 56px)",
                  lineHeight: 1.15,
                  margin: 0,
                  letterSpacing: "-0.03em",
                  fontWeight: 700,
                }}
              >
                見ているページを、
                <br />
                <span className="lp-hl">そのまま応募に</span>。
              </h2>
            </Reveal>
            <Reveal delay={200}>
              <p
                style={{
                  fontSize: 16,
                  lineHeight: 1.85,
                  color: "var(--lp-ink-2)",
                  marginTop: 24,
                  maxWidth: 480,
                }}
              >
                URL・タイトル・ドメインから、企業と媒体を自動推定。ATSドメイン（i-web/SONAR）にも対応。抽出に失敗しても、必ずInboxに退避します。
              </p>
            </Reveal>
            <div style={{ display: "grid", gap: 14, marginTop: 32 }}>
              {steps.map((s, i) => (
                <Reveal key={s.k} delay={300 + i * 80}>
                  <div
                    style={{
                      display: "flex",
                      gap: 16,
                      alignItems: "center",
                    }}
                  >
                    <div
                      style={{
                        width: 44,
                        height: 44,
                        borderRadius: 12,
                        background: "var(--lp-sun)",
                        border: "1.5px solid var(--lp-ink)",
                        color: "var(--lp-ink)",
                        display: "grid",
                        placeItems: "center",
                        fontFamily: "var(--lp-font-serif)",
                        fontWeight: 700,
                        fontSize: 18,
                        boxShadow: "3px 3px 0 var(--lp-ink)",
                        flexShrink: 0,
                      }}
                    >
                      {s.k}
                    </div>
                    <div style={{ fontSize: 16, fontWeight: 600 }}>{s.t}</div>
                  </div>
                </Reveal>
              ))}
            </div>
          </div>
        </div>
      </div>

      <style>{`
        .lp-extension-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 64px; align-items: center; }
        @media (max-width: 1024px) { .lp-extension-grid { grid-template-columns: 1fr; gap: 48px; } }
      `}</style>
    </section>
  );
}

// ── TIMELINE ─────────────────────────────────────────────────────────────
function TimelineSection() {
  const steps = [
    { t: "応募", d: "気になる企業を1クリックで保存", c: "var(--lp-s-app)" },
    { t: "書類", d: "ES締切はTaskで管理、3日前に通知", c: "var(--lp-s-doc)" },
    { t: "テスト", d: "SPI/適性検査の受検期限も逃さない", c: "var(--lp-s-test)" },
    { t: "面接", d: "現在Stageをワンタップで更新", c: "var(--lp-s-interview)" },
    { t: "内定", d: "シカくんが一緒に喜んでくれる", c: "var(--lp-s-offer)" },
  ];
  return (
    <section
      style={{ padding: "120px 32px", background: "transparent", position: "relative" }}
    >
      <div style={{ maxWidth: 1240, margin: "0 auto" }}>
        <Reveal>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: 14,
              marginBottom: 20,
            }}
          >
            <span
              style={{ width: 28, height: 2, background: "var(--lp-peach)" }}
            />
            <span
              style={{
                fontSize: 12,
                fontWeight: 800,
                color: "var(--lp-peach)",
                letterSpacing: "0.12em",
              }}
            >
              FLOW
            </span>
          </div>
        </Reveal>
        <Reveal delay={100}>
          <h2
            style={{
              fontFamily: "var(--lp-font-serif)",
              fontSize: "clamp(36px, 5vw, 60px)",
              lineHeight: 1.1,
              margin: 0,
              letterSpacing: "-0.035em",
              fontWeight: 700,
            }}
          >
            応募から内定まで、
            <br />
            ぜんぶ見える。
          </h2>
        </Reveal>

        <div style={{ marginTop: 72, position: "relative" }}>
          <div
            style={{
              position: "absolute",
              top: 44,
              left: "10%",
              right: "10%",
              height: 3,
              background:
                "repeating-linear-gradient(90deg, var(--lp-ink) 0 8px, transparent 8px 16px)",
            }}
          />
          <div className="lp-timeline-grid">
            {steps.map((s, i) => (
              <Reveal key={i} delay={i * 120}>
                <div
                  style={{
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    textAlign: "center",
                    gap: 14,
                  }}
                >
                  <div
                    style={{
                      width: 88,
                      height: 88,
                      borderRadius: "50%",
                      background: s.c,
                      border: "2px solid var(--lp-ink)",
                      display: "grid",
                      placeItems: "center",
                      fontFamily: "var(--lp-font-serif)",
                      fontSize: 28,
                      fontWeight: 700,
                      boxShadow: "4px 4px 0 var(--lp-ink)",
                      position: "relative",
                    }}
                  >
                    {i + 1}
                    {i === 4 && (
                      <div
                        style={{ position: "absolute", top: -18, right: -18 }}
                      >
                        <Stamp text="YES" color="var(--lp-sun)" size={50} />
                      </div>
                    )}
                  </div>
                  <div
                    style={{
                      fontFamily: "var(--lp-font-serif)",
                      fontSize: 22,
                      fontWeight: 700,
                    }}
                  >
                    {s.t}
                  </div>
                  <div
                    style={{
                      fontSize: 13,
                      color: "var(--lp-ink-2)",
                      lineHeight: 1.6,
                    }}
                  >
                    {s.d}
                  </div>
                </div>
              </Reveal>
            ))}
          </div>
        </div>
      </div>

      <style>{`
        .lp-timeline-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 24px; position: relative; }
        @media (max-width: 880px) { .lp-timeline-grid { grid-template-columns: repeat(2, 1fr); gap: 32px; } }
      `}</style>
    </section>
  );
}

// ── PERSONA ──────────────────────────────────────────────────────────────
function PersonaSection() {
  const personas: {
    name: string;
    tag: string;
    quote: string;
    feature: string;
    mood: Mood;
    accent: string;
  }[] = [
    {
      name: "Aさん / 理系院1",
      tag: "インターン10社",
      quote:
        "サマーインターンで8社応募。研究室と両立で、もうカレンダーが破綻してた。",
      feature: "Taskの通知で締切を守れる",
      mood: "wow",
      accent: "var(--lp-s-interview)",
    },
    {
      name: "Bさん / 文系3年",
      tag: "Notion挫折経験者",
      quote:
        "Notionで台帳作ろうとして、2週間で挫折。手入力がつらすぎる。",
      feature: "Chrome拡張で1クリック保存",
      mood: "wink",
      accent: "var(--lp-sun)",
    },
    {
      name: "Cさん / 就活中盤",
      tag: "本選考30社",
      quote:
        "本選考30社、フェーズも会社名もごちゃごちゃ。お祈り済みと進行中を分けたい。",
      feature: "Entry単位で進行中/アーカイブ",
      mood: "happy",
      accent: "var(--lp-s-offer)",
    },
  ];
  const rots = [-1, 0.8, -0.4];
  return (
    <section id="persona" style={{ padding: "120px 32px" }}>
      <div style={{ maxWidth: 1240, margin: "0 auto" }}>
        <Reveal>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: 14,
              marginBottom: 20,
            }}
          >
            <span
              style={{ width: 28, height: 2, background: "var(--lp-peach)" }}
            />
            <span
              style={{
                fontSize: 12,
                fontWeight: 800,
                color: "var(--lp-peach)",
                letterSpacing: "0.12em",
              }}
            >
              PERSONA
            </span>
          </div>
        </Reveal>
        <Reveal delay={100}>
          <h2
            style={{
              fontFamily: "var(--lp-font-serif)",
              fontSize: "clamp(40px, 5vw, 64px)",
              lineHeight: 1.1,
              margin: 0,
              letterSpacing: "-0.035em",
              fontWeight: 700,
            }}
          >
            こんな就活生に、
            <br />
            特に効く。
          </h2>
        </Reveal>
        <div className="lp-persona-grid">
          {personas.map((p, i) => (
            <Reveal key={p.name} delay={i * 120}>
              <PersonaCard p={p} rot={rots[i]} />
            </Reveal>
          ))}
        </div>
      </div>

      <style>{`
        .lp-persona-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; margin-top: 64px; }
        @media (max-width: 1024px) { .lp-persona-grid { grid-template-columns: 1fr; } }
      `}</style>
    </section>
  );
}

function PersonaCard({
  p,
  rot,
}: {
  p: {
    name: string;
    tag: string;
    quote: string;
    feature: string;
    mood: Mood;
    accent: string;
  };
  rot: number;
}) {
  return (
    <div
      className="lp-bold-card"
      style={{
        padding: 28,
        display: "flex",
        flexDirection: "column",
        gap: 18,
        minHeight: 340,
        transform: `rotate(${rot}deg)`,
        transition: "transform 220ms cubic-bezier(0.34,1.56,0.64,1)",
      }}
      onMouseEnter={(e) =>
        (e.currentTarget.style.transform = "rotate(0) translateY(-4px)")
      }
      onMouseLeave={(e) =>
        (e.currentTarget.style.transform = `rotate(${rot}deg)`)
      }
    >
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
        }}
      >
        <div
          style={{
            background: p.accent,
            borderRadius: 10,
            padding: 4,
            border: "1.5px solid var(--lp-ink)",
          }}
        >
          <Mascot size={52} mood={p.mood} accent={p.accent} />
        </div>
        <div className="lp-sticker" style={{ background: "var(--lp-sun-soft)" }}>
          {p.tag}
        </div>
      </div>
      <div
        style={{
          fontFamily: "var(--lp-font-serif)",
          fontSize: 20,
          fontWeight: 700,
        }}
      >
        {p.name}
      </div>
      <div
        style={{
          fontSize: 15,
          lineHeight: 1.7,
          color: "var(--lp-ink-2)",
          fontStyle: "italic",
        }}
      >
        「{p.quote}」
      </div>
      <div
        style={{
          marginTop: "auto",
          padding: "12px 14px",
          borderRadius: 10,
          background: "var(--lp-sun-soft)",
          fontSize: 13,
          fontWeight: 700,
          borderLeft: "3px solid var(--lp-peach)",
        }}
      >
        ✓ {p.feature}
      </div>
    </div>
  );
}

// ── FAQ ──────────────────────────────────────────────────────────────────
function FAQSection() {
  const [open, setOpen] = useState(0);
  const faqs = [
    {
      q: "本当に無料ですか？",
      a: "はい、MVP期間中は完全無料でお使いいただけます。将来的に有料プランを設ける可能性がありますが、基本機能は無料で提供し続ける予定です。",
    },
    {
      q: "Chrome拡張がないと使えませんか？",
      a: "いいえ。すべての機能は手動入力で完結できます。Chrome拡張は「入力の加速装置」であり、必須ではありません。",
    },
    {
      q: "データのエクスポートはできますか？",
      a: "CSV形式でCompany / Entry / Task / Clipのすべてをエクスポートできます。退会時は全データが削除されます。",
    },
    {
      q: "スマホだけでも使えますか？",
      a: "レスポンシブ対応しているため、スマホからダッシュボード閲覧・タスク処理・手動Clip追加がすべて可能です。",
    },
    {
      q: "i-web/SONARなどATSのページも保存できますか？",
      a: "はい。ドメイン + タイトル + 固定語で複数シグナル推定しています。推定に失敗してもInboxに退避されます。",
    },
  ];
  return (
    <section
      id="faq"
      style={{ padding: "120px 32px", background: "transparent" }}
    >
      <div style={{ maxWidth: 880, margin: "0 auto" }}>
        <Reveal>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: 14,
              marginBottom: 20,
            }}
          >
            <span
              style={{ width: 28, height: 2, background: "var(--lp-peach)" }}
            />
            <span
              style={{
                fontSize: 12,
                fontWeight: 800,
                color: "var(--lp-peach)",
                letterSpacing: "0.12em",
              }}
            >
              FAQ
            </span>
          </div>
        </Reveal>
        <Reveal delay={100}>
          <h2
            style={{
              fontFamily: "var(--lp-font-serif)",
              fontSize: "clamp(36px, 5vw, 56px)",
              lineHeight: 1.1,
              margin: 0,
              letterSpacing: "-0.035em",
              fontWeight: 700,
            }}
          >
            よくある質問。
          </h2>
        </Reveal>
        <div style={{ marginTop: 48, display: "grid", gap: 12 }}>
          {faqs.map((f, i) => (
            <Reveal key={i} delay={i * 60}>
              <div className="lp-bold-card" style={{ overflow: "hidden" }}>
                <button
                  onClick={() => setOpen(open === i ? -1 : i)}
                  style={{
                    width: "100%",
                    padding: "20px 24px",
                    textAlign: "left",
                    background: open === i ? "var(--lp-sun)" : "transparent",
                    border: "none",
                    cursor: "pointer",
                    fontFamily: "inherit",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    gap: 14,
                    transition: "background 180ms",
                  }}
                >
                  <div
                    style={{ display: "flex", alignItems: "center", gap: 12 }}
                  >
                    <span
                      style={{
                        fontFamily: "var(--lp-font-serif)",
                        fontSize: 14,
                        fontWeight: 700,
                        color: "var(--lp-peach)",
                      }}
                    >
                      Q{i + 1}
                    </span>
                    <span style={{ fontWeight: 700, fontSize: 16 }}>{f.q}</span>
                  </div>
                  <div
                    style={{
                      fontSize: 22,
                      transform: open === i ? "rotate(45deg)" : "none",
                      transition: "transform 220ms",
                    }}
                  >
                    +
                  </div>
                </button>
                {open === i && (
                  <div
                    style={{
                      padding: "16px 24px 22px 56px",
                      fontSize: 14,
                      lineHeight: 1.8,
                      color: "var(--lp-ink-2)",
                      borderTop: "1px solid var(--lp-line)",
                    }}
                  >
                    {f.a}
                  </div>
                )}
              </div>
            </Reveal>
          ))}
        </div>
      </div>
    </section>
  );
}

// ── FINAL CTA ────────────────────────────────────────────────────────────
function FinalCTA() {
  return (
    <section style={{ padding: "120px 32px 80px" }}>
      <div style={{ maxWidth: 1240, margin: "0 auto" }}>
        <Reveal>
          <div
            style={{
              background:
                "linear-gradient(135deg, #E8F0FF 0%, #F0E8FF 45%, #FFE8DC 100%)",
              color: "var(--lp-ink)",
              borderRadius: 32,
              padding: "96px 64px",
              position: "relative",
              overflow: "hidden",
              border: "1px solid rgba(255,255,255,0.8)",
              boxShadow: "0 40px 80px -40px rgba(74,108,247,0.3)",
            }}
          >
            <Parallax
              speed={0.4}
              axis="x"
              style={{
                position: "absolute",
                inset: 0,
                display: "flex",
                alignItems: "center",
                pointerEvents: "none",
              }}
            >
              <div
                aria-hidden
                style={{
                  fontFamily: "var(--lp-font-serif)",
                  fontSize: 280,
                  fontWeight: 700,
                  color: "rgba(74,108,247,0.08)",
                  letterSpacing: "-0.04em",
                  whiteSpace: "nowrap",
                }}
              >
                START · START · START · START · START
              </div>
            </Parallax>
            <div
              aria-hidden
              style={{
                position: "absolute",
                top: -80,
                left: -60,
                width: 400,
                height: 400,
                borderRadius: "50%",
                background:
                  "radial-gradient(circle, rgba(74,108,247,0.35), transparent 65%)",
                filter: "blur(50px)",
              }}
            />
            <div
              aria-hidden
              style={{
                position: "absolute",
                bottom: -80,
                right: 200,
                width: 360,
                height: 360,
                borderRadius: "50%",
                background:
                  "radial-gradient(circle, rgba(255,179,152,0.4), transparent 65%)",
                filter: "blur(50px)",
              }}
            />

            <div style={{ position: "relative", maxWidth: 720 }}>
              <div
                className="lp-chip"
                style={{
                  marginBottom: 32,
                  height: 32,
                  padding: "0 14px",
                  fontSize: 12,
                  color: "var(--lp-blue)",
                  background: "var(--lp-glass-strong)",
                }}
              >
                <span
                  className="lp-dot"
                  style={{ background: "var(--lp-blue)" }}
                />
                今すぐはじめる
              </div>
              <h2
                style={{
                  fontFamily: "var(--lp-font-serif)",
                  fontSize: "clamp(48px, 7vw, 96px)",
                  lineHeight: 1.0,
                  margin: 0,
                  letterSpacing: "-0.04em",
                  fontWeight: 700,
                }}
              >
                今の応募、
                <br />
                <span className="lp-grad-text">
                  <CountUp to={30} />
                  社、1枚に
                </span>
                。
              </h2>
              <p
                style={{
                  fontSize: 17,
                  marginTop: 28,
                  color: "var(--lp-ink-2)",
                  lineHeight: 1.8,
                  maxWidth: 520,
                }}
              >
                Googleログインだけ。入力は最小限。まずは3社登録して、
                <br />
                直近の締切を整えるところから。
              </p>
              <div
                style={{
                  display: "flex",
                  gap: 14,
                  marginTop: 40,
                  flexWrap: "wrap",
                }}
              >
                <button
                  type="button"
                  className="lp-btn lp-btn-blue"
                  style={{ height: 58, padding: "0 28px", fontSize: 15 }}
                >
                  Googleで無料ではじめる →
                </button>
                <button
                  type="button"
                  className="lp-btn lp-btn-glass"
                  style={{ height: 58, padding: "0 28px", fontSize: 15 }}
                >
                  デモを見る
                </button>
              </div>
            </div>

            <div
              className="lp-float-y"
              style={{ position: "absolute", bottom: 40, right: 60 }}
            >
              <Mascot size={200} mood="cheer" accent="var(--lp-sun)" />
            </div>
          </div>
        </Reveal>
      </div>
    </section>
  );
}

// ── FOOTER ───────────────────────────────────────────────────────────────
function LPFooter() {
  const cols = [
    {
      t: "プロダクト",
      items: ["機能一覧", "Chrome拡張", "価格（無料）", "ロードマップ"],
    },
    {
      t: "サポート",
      items: ["FAQ", "お問い合わせ", "データエクスポート", "退会方法"],
    },
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
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: 10,
              marginBottom: 14,
            }}
          >
            <MiniMascot size={36} />
            <div
              style={{
                fontWeight: 800,
                fontSize: 22,
                fontFamily: "var(--lp-font-serif)",
                letterSpacing: "-0.02em",
              }}
            >
              Entr<span style={{ color: "var(--lp-blue)" }}>é</span>
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
            <div style={{ fontSize: 12, fontWeight: 800, marginBottom: 12 }}>
              {c.t}
            </div>
            <div
              style={{
                display: "grid",
                gap: 8,
                fontSize: 13,
                color: "rgba(255,255,255,0.75)",
              }}
            >
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
        <div>© 2026 Entré · シカくんと一緒に</div>
        <div className="lp-mono">v0.2.0-beta</div>
      </div>

      <style>{`
        @media (max-width: 880px) { .lp-footer-grid { grid-template-columns: 1fr 1fr !important; } }
      `}</style>
    </footer>
  );
}
