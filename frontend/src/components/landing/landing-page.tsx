"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect } from "react";
import type { CSSProperties, ReactNode } from "react";

// ═════════════════════════════════════════════════════════════════════════
// Entré LP v2 — React port
// v2 HTML の design を React/TSX に移植したもの。
// スタイルは landing.css（.lp-scope 配下）で定義している。
// ═════════════════════════════════════════════════════════════════════════

export function LandingPage() {
  useLpReveal();
  return (
    <div className="lp-scope" style={{ minHeight: "100%" }}>
      <Nav />
      <Hero />
      <ProblemSection />
      <MobileFirstSection />
      <ToolsSection />
      <StepsSection />
      <FAQSection />
      <FinalCTASection />
      <Footer />
    </div>
  );
}

// スクロール連動の fade-up。`.lp-reveal` が付いた要素を観測し、viewport に入ったら `is-visible` を付ける。
function useLpReveal() {
  useEffect(() => {
    const els = document.querySelectorAll<HTMLElement>(".lp-reveal");
    if (els.length === 0) return;
    const io = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            entry.target.classList.add("is-visible");
            io.unobserve(entry.target);
          }
        }
      },
      { threshold: 0.15, rootMargin: "0px 0px -60px 0px" }
    );
    els.forEach((el) => io.observe(el));
    return () => io.disconnect();
  }, []);
}

// ───────────────────────────────────────────────────────────────────────
// Nav
// ───────────────────────────────────────────────────────────────────────
function Nav() {
  return (
    <nav className="lp-nav">
      <div className="lp-nav-inner">
        <div className="lp-logo-mark">
          <span className="name">
            <em>Entré</em>
          </span>
          <span className="beta">BETA</span>
        </div>
        <div className="lp-nav-links lp-hide-sm">
          <a href="#features">機能</a>
          <a href="#how">使い方</a>
          <a href="#faq">よくある質問</a>
        </div>
        <div className="lp-nav-cta">
          <Link href="/login" className="lp-btn lp-btn-ghost lp-hide-sm">
            ログイン
          </Link>
          <Link href="/login" className="lp-btn lp-btn-primary">
            無料で始める
          </Link>
        </div>
      </div>
    </nav>
  );
}

// ───────────────────────────────────────────────────────────────────────
// Hero
// ───────────────────────────────────────────────────────────────────────
function Hero() {
  return (
    <section className="lp-hero">
      {/* decorative squiggles */}
      <svg className="lp-squiggle" style={{ top: 40, left: "18%" }} width="80" height="20" viewBox="0 0 80 20">
        <path d="M2 10 Q 12 2, 22 10 T 42 10 T 62 10 T 78 10" stroke="#C9CBB4" strokeWidth="1.4" fill="none" strokeLinecap="round" />
      </svg>
      <svg className="lp-squiggle" style={{ top: 110, right: "30%" }} width="30" height="30" viewBox="0 0 30 30">
        <path d="M15 4 L17 13 L26 15 L17 17 L15 26 L13 17 L4 15 L13 13 Z" fill="#D7B5A8" opacity="0.8" />
      </svg>
      <svg className="lp-squiggle" style={{ bottom: 20, left: "42%" }} width="60" height="12" viewBox="0 0 60 12">
        <path d="M2 6 Q 10 2, 18 6 T 34 6 T 58 6" stroke="#B9BFA5" strokeWidth="1.2" fill="none" strokeLinecap="round" />
      </svg>

      <div className="lp-hero-inner">
        <div>
          <div className="lp-hero-eyebrow">就活の、いちばんの味方に。</div>
          <h1 className="lp-hero-title">
            散らかった就活、
            <br />
            ぜんぶ
            <span className="lp-underline-hand">
              <span className="num">1枚</span>
              <svg viewBox="0 0 100 10" preserveAspectRatio="none">
                <path d="M2 6 Q 30 1, 60 5 T 98 4" stroke="#6B8A72" strokeWidth="2.2" fill="none" strokeLinecap="round" />
              </svg>
            </span>
            に。
          </h1>
          <p className="lp-hero-sub">
            エントリーの締切、選考ステータス、メモ、
            <br />
            URLまでひとまとめ。マイナビ、リクナビ、
            <br />
            ONE CAREER、企業サイト…あちこちの情報を
            <br />
            さっと保存できます。
          </p>
          <div className="lp-hero-actions">
            <Link href="/login" className="lp-btn lp-btn-primary lp-btn-lg">
              <GoogleG />
              Googleで無料ではじめる
            </Link>
          </div>
          <div className="lp-hero-chips">
            <HeroChip label="完全無料" />
            <HeroChip label="登録カンタン" />
            <HeroChip label="データは安全に保存" />
          </div>
        </div>

        <div className="lp-device-stack">
          {/* 装飾ステッカー — 浮遊・回転・ぷるぷる */}
          <svg
            className="lp-deco lp-deco-spin"
            style={{ left: -20, top: -20, width: 42, height: 42 }}
            viewBox="0 0 40 40"
            aria-hidden
          >
            <path d="M20 2 L23 16 L37 20 L23 24 L20 37 L17 24 L3 20 L17 16 Z" fill="#2B2A26" />
          </svg>
          <svg
            className="lp-deco lp-deco-star"
            style={{ left: 260, top: 0, width: 24, height: 24, animationDelay: "0.5s" }}
            viewBox="0 0 20 20"
            aria-hidden
          >
            <path d="M10 1 L12 8 L19 10 L12 12 L10 19 L8 12 L1 10 L8 8 Z" fill="#E9B9B0" />
          </svg>
          <span
            className="lp-deco lp-deco-star"
            style={{
              left: -80,
              top: 80,
              width: 10,
              height: 10,
              background: "#2B2A26",
              borderRadius: 9999,
              animationDelay: "1s",
            }}
            aria-hidden
          />
          <svg
            className="lp-deco lp-deco-spin"
            style={{ right: -10, top: 30, width: 34, height: 34, animationDuration: "25s" }}
            viewBox="0 0 40 40"
            aria-hidden
          >
            <path d="M20 4 L22 18 L36 20 L22 22 L20 36 L18 22 L4 20 L18 18 Z" fill="#D4BA82" />
          </svg>
          <svg
            className="lp-deco lp-deco-wiggle"
            style={{ left: -70, bottom: 140, width: 70, height: 20, animationDelay: "0.3s" }}
            viewBox="0 0 70 20"
            aria-hidden
          >
            <path
              d="M3 10 Q 10 2, 17 10 T 31 10 T 45 10 T 59 10 T 68 10"
              stroke="#6B8A72"
              strokeWidth="2"
              fill="none"
              strokeLinecap="round"
            />
          </svg>
          <svg
            className="lp-deco lp-deco-wiggle"
            style={{ right: -40, bottom: 100, width: 58, height: 44, animationDelay: "1.2s" }}
            viewBox="0 0 60 50"
            aria-hidden
          >
            <path d="M3 40 Q 20 20, 50 15" stroke="#2B2A26" strokeWidth="2" fill="none" strokeLinecap="round" />
            <path
              d="M42 10 L52 13 L46 22"
              stroke="#2B2A26"
              strokeWidth="2"
              fill="none"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
          <svg
            className="lp-deco lp-deco-bounce"
            style={{ left: 220, bottom: 50, width: 28, height: 28, animationDelay: "0.7s", zIndex: 6 }}
            viewBox="0 0 24 24"
            aria-hidden
          >
            <path
              d="M12 20 C 12 20, 3 13, 3 8 C 3 5, 5 3, 8 3 C 10 3, 11 4, 12 6 C 13 4, 14 3, 16 3 C 19 3, 21 5, 21 8 C 21 13, 12 20, 12 20 Z"
              fill="#E9B9B0"
              stroke="#2B2A26"
              strokeWidth="1.5"
              strokeLinejoin="round"
            />
          </svg>

          {/* bubble callout */}
          <div className="lp-bubble lp-bubble-1">
            PCも見やすく、
            <br />
            しっかり使える。
          </div>

          <PhoneMockup />
          <IMacMockup />

          <div style={{ position: "absolute", top: 110, right: -20, zIndex: 6 }} className="lp-float">
            <EnvelopeMascot size={78} mood="happy" />
          </div>
        </div>
      </div>
    </section>
  );
}

function HeroChip({ label }: { label: string }) {
  return (
    <span className="lp-hero-chip">
      <svg viewBox="0 0 12 12">
        <circle cx="6" cy="6" r="5" fill="none" stroke="#6B8A72" strokeWidth="1.5" />
        <path d="M3.5 6 L5 7.5 L8.5 4" stroke="#6B8A72" strokeWidth="1.5" fill="none" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
      {label}
    </span>
  );
}

// ───────────────────────────────────────────────────────────────────────
// Problem section
// ───────────────────────────────────────────────────────────────────────
function ProblemSection() {
  const items = [
    {
      icon: (
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
          <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
          <polyline points="14 2 14 8 20 8" />
          <path d="M8 13h8M8 17h5" />
        </svg>
      ),
      text: (
        <>
          情報がバラバラで、
          <br />
          管理が大変
        </>
      ),
    },
    {
      icon: (
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
          <rect x="3" y="4" width="18" height="18" rx="2" />
          <line x1="16" y1="2" x2="16" y2="6" />
          <line x1="8" y1="2" x2="8" y2="6" />
          <line x1="3" y1="10" x2="21" y2="10" />
        </svg>
      ),
      text: (
        <>
          締切や予定を
          <br />
          うっかり忘れてしまう
        </>
      ),
    },
    {
      icon: (
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
          <rect x="4" y="3" width="16" height="18" rx="2" />
          <line x1="8" y1="8" x2="16" y2="8" />
          <line x1="8" y1="12" x2="16" y2="12" />
          <line x1="8" y1="16" x2="12" y2="16" />
        </svg>
      ),
      text: (
        <>
          メモやURLを
          <br />
          探すのに時間がかかる
        </>
      ),
    },
    {
      icon: (
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
          <rect x="2" y="4" width="14" height="12" rx="2" />
          <rect x="10" y="10" width="12" height="10" rx="2" />
          <line x1="7" y1="18" x2="11" y2="18" />
        </svg>
      ),
      text: (
        <>
          PCとスマホで
          <br />
          情報が分散している
        </>
      ),
    },
  ];

  return (
    <section className="lp-section" style={{ paddingTop: 40 }}>
      <div className="lp-section-inner">
        <h2 className="lp-section-title">こんなお悩み、ありませんか？</h2>
        <div className="lp-problem-grid">
          {items.map((item, i) => (
            <div
              key={i}
              className="lp-problem-card lp-reveal"
              style={
                {
                  "--lp-reveal-i": i,
                  "--lp-tilt": `${[-1.2, 0.8, -0.6, 1][i]}deg`,
                } as CSSProperties
              }
            >
              <span className="lp-problem-tape" aria-hidden />
              <div className="lp-problem-icon">{item.icon}</div>
              <div className="lp-problem-text">{item.text}</div>
              <div className="lp-problem-solved">
                <span className="lp-problem-solved-check">✓</span> Entré で解決
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ───────────────────────────────────────────────────────────────────────
// Mobile First section
// ───────────────────────────────────────────────────────────────────────
function MobileFirstSection() {
  return (
    <section className="lp-device-detail" id="features">
      <div className="lp-section-inner">
        <div className="lp-device-row lp-reveal">
          <div className="lp-device-row-text">
            <div className="lp-device-eyebrow">FOR DESKTOP</div>
            <h2 className="lp-device-title">
              デスクで、
              <br />
              本腰を入れて。
            </h2>
            <p className="lp-device-sub">
              就活サイトを見ながら Chrome 拡張でワンクリック保存。
              選考ボードや一覧で全体を俯瞰しながら、腰を据えて応募を管理できます。
            </p>
            <ul className="lp-device-list">
              <li>
                <span className="lp-mf-check">✓</span> Chrome 拡張で URL・タイトルを1クリック保存
              </li>
              <li>
                <span className="lp-mf-check">✓</span> 選考ボード・一覧で全体を俯瞰
              </li>
              <li>
                <span className="lp-mf-check">✓</span> キーボードショートカットで高速操作
              </li>
            </ul>
          </div>
          <div className="lp-device-row-img is-desktop">
            <Image
              src="/PC.png"
              alt="Entré の PC ダッシュボード画面"
              width={3944}
              height={2564}
              className="lp-device-row-img-inner"
            />
          </div>
        </div>

        <div className="lp-device-sync-note">
          <span className="lp-device-sync-dash" />
          PC ↔ スマホ、データは自動同期
          <span className="lp-device-sync-dash" />
        </div>

        <div className="lp-device-row lp-device-row--reverse lp-reveal">
          <div className="lp-device-row-text">
            <div className="lp-device-eyebrow">FOR MOBILE</div>
            <h2 className="lp-device-title">
              通学中も、
              <br />
              寝る前も。
            </h2>
            <p className="lp-device-sub">
              通知で締切を逃さず、思いついたときにサッと URL を保存。
              タスクの完了や期限変更もスマホだけで完結します。
            </p>
            <ul className="lp-device-list">
              <li>
                <span className="lp-mf-check">✓</span> 締切・予定をプッシュ通知でお知らせ
              </li>
              <li>
                <span className="lp-mf-check">✓</span> タスク完了・期限変更もスマホだけで
              </li>
              <li>
                <span className="lp-mf-check">✓</span> URL をペーストしてサッと保存
              </li>
            </ul>
          </div>
          <div className="lp-device-row-img is-mobile">
            <Image
              src="/smartphone.png"
              alt="Entré のスマホダッシュボード画面"
              width={1857}
              height={3096}
              className="lp-device-row-img-inner"
            />
          </div>
        </div>
      </div>
    </section>
  );
}

// ───────────────────────────────────────────────────────────────────────
// 4 Tools section
// ───────────────────────────────────────────────────────────────────────
function ToolsSection() {
  const tools: Array<{
    icon: ReactNode;
    name: string;
    sub: string;
    desc: ReactNode;
    accent: string;
    preview: ReactNode;
  }> = [
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
          <rect x="4" y="3" width="16" height="18" rx="2" />
          <path d="M8 8h8M8 12h6M8 16h4" />
          <circle cx="17" cy="16" r="1" fill="#4F6E58" />
        </svg>
      ),
      name: "Entry",
      sub: "エントリー管理",
      desc: (
        <>
          企業や求人ごとに、
          <br />
          選考ステータスを一括で管理
        </>
      ),
      accent: "sage",
      preview: <EntryPreview />,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
          <path d="M9 11l3 3 8-8" />
          <path d="M20 12v7a2 2 0 01-2 2H6a2 2 0 01-2-2V6a2 2 0 012-2h11" />
        </svg>
      ),
      name: "Task",
      sub: "タスク管理",
      desc: (
        <>
          締切や直近予定を可視化、
          <br />
          リマインドで漏れない
        </>
      ),
      accent: "pink",
      preview: <TaskPreview />,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
          <path d="M4 4h16v16H4z" />
          <path d="M4 8l8 6 8-6" />
        </svg>
      ),
      name: "Inbox",
      sub: "通知・やること管理",
      desc: (
        <>
          企業からの連絡や
          <br />
          自分のやることを1画面で
        </>
      ),
      accent: "gold",
      preview: <InboxPreview />,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
          <rect x="3" y="3" width="8" height="8" rx="1" />
          <rect x="13" y="3" width="8" height="8" rx="1" />
          <rect x="3" y="13" width="8" height="8" rx="1" />
          <rect x="13" y="13" width="5" height="5" rx="1" />
          <path d="M18 16v5h3" />
        </svg>
      ),
      name: "Chrome拡張",
      sub: "ワンクリック登録",
      desc: (
        <>
          ブラウザからカンタンに
          <br />
          求人情報を保存
        </>
      ),
      accent: "sky",
      preview: <ChromePreview />,
    },
  ];

  return (
    <section className="lp-section">
      <div className="lp-section-inner">
        <h2 className="lp-section-title">就活を支える、4つのツール</h2>
        <div className="lp-tools-grid">
          {tools.map((t, i) => (
            <div
              key={t.name}
              className={`lp-tool-card lp-tool-card--${t.accent} lp-reveal`}
              style={{ "--lp-reveal-i": i } as CSSProperties}
            >
              <div className="lp-tool-icon">{t.icon}</div>
              <div className="lp-tool-name">{t.name}</div>
              <div className="lp-tool-sub">{t.sub}</div>
              <div className="lp-tool-desc">{t.desc}</div>
              <div className="lp-tool-preview">{t.preview}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function EntryPreview() {
  return (
    <div className="lp-preview-stack">
      <div className="lp-preview-entry-row">
        <span className="lp-preview-co">株式会社ミライ</span>
        <span className="lp-preview-pill lp-preview-pill--doc">書類選考</span>
      </div>
      <div className="lp-preview-entry-row">
        <span className="lp-preview-co">アクロス商事</span>
        <span className="lp-preview-pill lp-preview-pill--interview">面接</span>
      </div>
      <div className="lp-preview-entry-row">
        <span className="lp-preview-co">オリエンス</span>
        <span className="lp-preview-pill lp-preview-pill--offer">内定</span>
      </div>
    </div>
  );
}

function TaskPreview() {
  return (
    <div className="lp-preview-stack">
      <div className="lp-preview-task-row is-done">
        <span className="lp-preview-check is-done">✓</span>
        <span className="lp-preview-task-name">ES提出</span>
        <span className="lp-preview-date">完了</span>
      </div>
      <div className="lp-preview-task-row">
        <span className="lp-preview-check" />
        <span className="lp-preview-task-name">一次面接</span>
        <span className="lp-preview-date is-soon">今日</span>
      </div>
      <div className="lp-preview-task-row">
        <span className="lp-preview-check" />
        <span className="lp-preview-task-name">適性検査</span>
        <span className="lp-preview-date">4/28</span>
      </div>
    </div>
  );
}

function InboxPreview() {
  return (
    <div className="lp-preview-stack">
      <div className="lp-preview-notify">
        <span className="lp-preview-dot lp-preview-dot--danger" />
        <span className="lp-preview-task-name">ES提出</span>
        <span className="lp-preview-date is-soon">今日</span>
      </div>
      <div className="lp-preview-notify">
        <span className="lp-preview-dot lp-preview-dot--warning" />
        <span className="lp-preview-task-name">面接 14:00</span>
        <span className="lp-preview-date">明日</span>
      </div>
      <div className="lp-preview-inbox-badge">
        <span>📎 未整理</span>
        <span className="lp-preview-inbox-count">3</span>
      </div>
    </div>
  );
}

function ChromePreview() {
  return (
    <div className="lp-preview-browser">
      <div className="lp-preview-browser-bar">
        <span className="lp-preview-browser-dot" style={{ background: "#FF6159" }} />
        <span className="lp-preview-browser-dot" style={{ background: "#FFBD2E" }} />
        <span className="lp-preview-browser-dot" style={{ background: "#28C941" }} />
        <span className="lp-preview-browser-url">rikunabi.com/…</span>
      </div>
      <div className="lp-preview-browser-body">
        <div className="lp-preview-browser-save">
          <span>+</span> Entré に保存
        </div>
      </div>
    </div>
  );
}

// ───────────────────────────────────────────────────────────────────────
// Steps section
// ───────────────────────────────────────────────────────────────────────
function StepsSection() {
  return (
    <section className="lp-section" id="how" style={{ paddingTop: 40 }}>
      <div className="lp-section-inner">
        <h2 className="lp-section-title">使い方は、とってもシンプル。</h2>
        <div className="lp-steps-row">
          <StepCard
            num={1}
            title="アカウントを作成"
            desc={
              <>
                Googleで登録して
                <br />
                すぐに使えます。
              </>
            }
            visual={<GoogleFullColor />}
            hasArrow
          />
          <StepCard
            num={2}
            title="情報をまとめる"
            desc={
              <>
                エントリー情報・URLを
                <br />
                まとめて記録。
              </>
            }
            visual={
              <>
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.8">
                  <path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71" />
                  <path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71" />
                </svg>
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#4F6E58" strokeWidth="1.8">
                  <path d="M4 4h16v16H4z" />
                  <path d="M4 8l8 6 8-6" />
                </svg>
              </>
            }
            hasArrow
          />
          <StepCard
            num={3}
            title="ステータスを管理"
            desc={
              <>
                選考の進捗をひと目で
                <br />
                チェック。
              </>
            }
            visual={
              <div style={{ display: "flex", gap: 2, flexDirection: "column", width: "100%" }}>
                <div style={{ height: 5, background: "var(--sage-deep)", borderRadius: 2, width: "80%" }} />
                <div style={{ height: 5, background: "var(--accent-pink)", borderRadius: 2, width: "60%" }} />
                <div style={{ height: 5, background: "var(--accent-gold)", borderRadius: 2, width: "40%" }} />
              </div>
            }
            hasArrow
          />
          <StepCard
            num={4}
            title="通知を受け取る"
            desc={
              <>
                締切や予定をお知らせ、
                <br />
                うっかりを防止。
              </>
            }
            visual={
              <svg width="32" height="32" viewBox="0 0 32 32" fill="none" stroke="#4F6E58" strokeWidth="1.6">
                <path d="M10 28 V6 L18 10 L10 14" fill="#E9B9B0" />
                <line x1="10" y1="28" x2="10" y2="6" />
              </svg>
            }
          />
          <div className="lp-step-end-card lp-reveal" style={{ "--lp-reveal-i": 4 } as CSSProperties}>
            <EnvelopeMascot size={58} mood="wink" />
            <div className="lp-step-end-text">
              あなたの就活を、
              <br />
              やさしく自動化。
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

function StepCard({
  num,
  title,
  desc,
  visual,
  hasArrow,
}: {
  num: number;
  title: string;
  desc: ReactNode;
  visual: ReactNode;
  hasArrow?: boolean;
}) {
  return (
    <div className="lp-step-card lp-reveal" style={{ "--lp-reveal-i": num - 1 } as CSSProperties}>
      <div className="lp-step-num">{num}</div>
      <div className="lp-step-title">{title}</div>
      <div className="lp-step-visual">{visual}</div>
      <div className="lp-step-desc">{desc}</div>
      {hasArrow ? <div className="lp-step-arrow">›</div> : null}
    </div>
  );
}

// ───────────────────────────────────────────────────────────────────────
// FAQ section
// ───────────────────────────────────────────────────────────────────────
function FAQSection() {
  const faqs = [
    "本当に無料で使えますか？",
    "データはどこに保存されますか？",
    "対応しているデバイスは？",
  ];
  return (
    <section className="lp-section" id="faq" style={{ paddingTop: 40 }}>
      <div className="lp-section-inner">
        <h2 className="lp-section-title">よくある質問</h2>
        <div className="lp-faq-wrap">
          <div>
            <div className="lp-faq-list">
              {faqs.map((q, i) => (
                <div
                  key={q}
                  className="lp-faq-item lp-reveal"
                  style={{ "--lp-reveal-i": i } as CSSProperties}
                >
                  <span className="lp-faq-q">Q</span>
                  {q}
                </div>
              ))}
            </div>
            <div className="lp-faq-more">すべての質問を見る →</div>
          </div>
          <div className="lp-faq-side">
            <EnvelopeMascot size={68} mood="happy" />
            <div className="lp-faq-side-text">
              ご不明点は
              <br />
              サポートまで
              <br />
              お気軽にどうぞ！
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

// ───────────────────────────────────────────────────────────────────────
// Final CTA
// ───────────────────────────────────────────────────────────────────────
function FinalCTASection() {
  return (
    <section className="lp-section" style={{ paddingTop: 20, paddingBottom: 40 }}>
      <div className="lp-section-inner">
        <div className="lp-final-cta">
          <div>
            <h2>
              さあ、<em style={{ fontStyle: "italic" }}>Entré</em>で
              <br />
              就活を<span className="em">もっとシンプルに。</span>
            </h2>
          </div>
          <div className="lp-final-cta-actions">
            <Link href="/login" className="lp-btn lp-btn-primary lp-btn-lg">
              <GoogleG />
              Googleではじめる
            </Link>
            <EnvelopeMascot size={56} mood="happy" />
          </div>
          <div className="lp-final-chips">
            <span className="lp-hero-chip" style={{ background: "transparent" }}>
              <Check />完全無料
            </span>
            <span className="lp-hero-chip" style={{ background: "transparent" }}>
              <Check />登録カンタン
            </span>
            <span className="lp-hero-chip" style={{ background: "transparent" }}>
              <Check />データは安全に保存
            </span>
          </div>
        </div>
      </div>
    </section>
  );
}

function Check() {
  return (
    <svg viewBox="0 0 12 12" width="12" height="12">
      <circle cx="6" cy="6" r="5" fill="none" stroke="#4F6E58" strokeWidth="1.5" />
      <path d="M3.5 6 L5 7.5 L8.5 4" stroke="#4F6E58" strokeWidth="1.5" fill="none" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

function Footer() {
  return <footer className="lp-footer">© 2026 Entré — 就活管理をシンプルに</footer>;
}

// ═════════════════════════════════════════════════════════════════════════
// Mockup components (再利用)
// ═════════════════════════════════════════════════════════════════════════

function PhoneMockup() {
  return (
    <div className="lp-phone">
      <Image
        src="/smartphone.png"
        alt="Entré アプリのダッシュボード画面"
        width={1857}
        height={3096}
        priority
        className="lp-phone-img"
      />
    </div>
  );
}

function IMacMockup() {
  return (
    <div className="lp-imac-wrap">
      <Image
        src="/PC.png"
        alt="Entré アプリのダッシュボード画面（PC表示）"
        width={3944}
        height={2564}
        priority
        className="lp-imac-img"
      />
    </div>
  );
}

// ═════════════════════════════════════════════════════════════════════════
// Google logo / Envelope mascot
// ═════════════════════════════════════════════════════════════════════════

function GoogleG() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24">
      <path fill="#fff" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23zM5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18A10.99 10.99 0 001 12c0 1.77.42 3.44 1.18 4.93l3.66-2.84zM12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
    </svg>
  );
}

function GoogleFullColor() {
  return (
    <svg width="42" height="42" viewBox="0 0 42 42">
      <circle cx="21" cy="21" r="20" fill="#fff" stroke="#E6E1D3" />
      <path fill="#4285F4" d="M31 21.2c0-.7-.06-1.38-.18-2.03H21v3.84h5.34c-.23 1.24-.94 2.29-2 3v2.5h3.22C29.43 26.82 31 24.3 31 21.2z" />
      <path fill="#34A853" d="M21 31c2.7 0 4.97-.9 6.63-2.44l-3.22-2.5c-.9.6-2.04.95-3.41.95-2.62 0-4.84-1.77-5.63-4.15h-3.33v2.6C13.7 29 17.07 31 21 31z" />
      <path fill="#FBBC05" d="M15.37 22.86c-.2-.6-.31-1.24-.31-1.9s.11-1.3.31-1.9v-2.6h-3.33A9.99 9.99 0 0011 21c0 1.6.38 3.11 1.04 4.46l3.33-2.6z" />
      <path fill="#EA4335" d="M21 14.95c1.47 0 2.8.5 3.84 1.5l2.88-2.88C25.97 11.95 23.7 11 21 11c-3.93 0-7.3 2-8.96 4.94l3.33 2.6c.79-2.38 3-4.15 5.63-4.15z" />
    </svg>
  );
}

// ═════════════════════════════════════════════════════════════════════════
// EnvelopeMascot — v2 の封筒マスコット
// ═════════════════════════════════════════════════════════════════════════

type MascotMood = "happy" | "wink" | "sweat";

export function EnvelopeMascot({
  size = 78,
  mood = "happy",
  style,
}: {
  size?: number;
  mood?: MascotMood;
  style?: CSSProperties;
}) {
  return (
    <svg width={size} height={size} viewBox="0 0 100 100" style={style}>
      {/* 鹿の角（antlers） — 参考画像では葉ではなく小さな枝分かれした角 */}
      <path
        d="M42 28 L40 18 M40 18 L37 15 M40 18 L43 15"
        stroke="#8A6A4A"
        strokeWidth="1.6"
        fill="none"
        strokeLinecap="round"
      />
      <path
        d="M58 28 L60 18 M60 18 L63 15 M60 18 L57 15"
        stroke="#8A6A4A"
        strokeWidth="1.6"
        fill="none"
        strokeLinecap="round"
      />
      {/* envelope body */}
      <rect x="10" y="28" width="80" height="60" rx="8" fill="#FEF8E6" stroke="#2B2A26" strokeWidth="2" />
      <path d="M10 32 L50 60 L90 32" fill="none" stroke="#2B2A26" strokeWidth="2" strokeLinejoin="round" />
      {/* eyes */}
      {mood === "wink" ? (
        <>
          <path d="M38 66 L41 69 L44 66" stroke="#2B2A26" strokeWidth="1.6" fill="none" strokeLinecap="round" />
          <path d="M58 66 L61 69 L64 66" stroke="#2B2A26" strokeWidth="1.6" fill="none" strokeLinecap="round" />
        </>
      ) : (
        <>
          <circle cx="40" cy="68" r="2" fill="#2B2A26" />
          <circle cx="62" cy="68" r="2" fill="#2B2A26" />
        </>
      )}
      {/* mouth */}
      <path d="M44 76 Q51 80 58 76" stroke="#2B2A26" strokeWidth="1.5" fill="none" strokeLinecap="round" />
      {/* cheeks */}
      <circle cx="34" cy="74" r="2.5" fill="#E9B9B0" opacity="0.7" />
      <circle cx="68" cy="74" r="2.5" fill="#E9B9B0" opacity="0.7" />
      {/* optional sweat drop */}
      {mood === "sweat" ? (
        <path d="M82 38 Q80 42 82 45 Q84 42 82 38" fill="#7EC0D8" stroke="#2B2A26" strokeWidth="1" />
      ) : null}
    </svg>
  );
}
