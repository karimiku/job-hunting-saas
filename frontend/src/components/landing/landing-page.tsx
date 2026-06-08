"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect } from "react";
import type { CSSProperties } from "react";

const productShots = {
  dashboard: { src: "/marketing/desktop-dashboard.png", width: 2880, height: 2048 },
  entry: { src: "/marketing/desktop-entry.png", width: 2880, height: 2048 },
  kanban: { src: "/marketing/desktop-kanban.png", width: 2880, height: 2048 },
  mobileEntry: { src: "/marketing/mobile-entry.png", width: 1170, height: 2532 },
  mobileKanban: { src: "/marketing/mobile-kanban.png", width: 1170, height: 2532 },
};

export function LandingPage() {
  useLpReveal();

  return (
    <div className="lp-scope lp-simple" style={{ minHeight: "100%" }}>
      <Nav />
      <main>
        <Hero />
        <CoreFlowSection />
        <CoreScreensSection />
        <FinalCTASection />
      </main>
      <Footer />
    </div>
  );
}

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

function Nav() {
  return (
    <nav className="lp-nav" aria-label="メインナビゲーション">
      <div className="lp-nav-inner">
        <a className="lp-logo-mark" href="#top" aria-label="Entré トップへ">
          <span className="name">
            <em>Entré</em>
          </span>
          <span className="beta">BETA</span>
        </a>
        <div className="lp-nav-links lp-hide-sm">
          <a href="#flow">流れ</a>
          <a href="#entry">Entry</a>
          <a href="#kanban">カンバン</a>
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

function Hero() {
  return (
    <section className="lp-simple-hero" id="top">
      <div className="lp-simple-hero-copy">
        <p className="lp-simple-kicker">Entry とカンバンだけに絞った就活管理</p>
        <h1>Entré</h1>
        <p className="lp-simple-lead">
          保存した求人を Entry にまとめて、選考状況はカンバンで動かす。
          締切、URL、メモ、次のタスクを一箇所に置けるシンプルな就活ボードです。
        </p>
        <div className="lp-simple-actions">
          <Link href="/login" className="lp-btn lp-btn-primary lp-btn-lg">
            <GoogleG />
            Googleで始める
          </Link>
          <a href="#flow" className="lp-btn lp-btn-secondary lp-btn-lg">
            画面を見る
          </a>
        </div>
      </div>

      <div className="lp-simple-hero-visual">
        <Image
          src={productShots.dashboard.src}
          alt="Entré のホーム画面。保存箱から Entry に変換する次の行動と今日のタスクを表示している。"
          width={productShots.dashboard.width}
          height={productShots.dashboard.height}
          priority
          unoptimized
          className="lp-simple-dashboard-shot"
        />
        <div className="lp-simple-mobile-float" aria-hidden>
          <Image
            src={productShots.mobileEntry.src}
            alt=""
            width={productShots.mobileEntry.width}
            height={productShots.mobileEntry.height}
            priority
            unoptimized
            className="lp-simple-mobile-shot"
          />
        </div>
      </div>
    </section>
  );
}

function CoreFlowSection() {
  const steps = [
    {
      label: "保存する",
      title: "気になる求人を保存箱へ",
      text: "あとで見る求人を一旦ためる。必要なものだけ Entry に変換します。",
    },
    {
      label: "まとめる",
      title: "Entry に情報を集約",
      text: "会社名、応募元URL、メモ、ステータス、次のタスクを同じ場所に置きます。",
    },
    {
      label: "動かす",
      title: "カンバンで選考を進める",
      text: "書類、テスト、面接、内定まで、今どこにあるかだけを見ます。",
    },
  ];

  return (
    <section className="lp-simple-section" id="flow">
      <div className="lp-simple-section-head lp-reveal">
        <p className="lp-simple-kicker">CORE FLOW</p>
        <h2>やることは3つだけ。</h2>
      </div>
      <div className="lp-simple-flow">
        {steps.map((step, index) => (
          <article
            className="lp-simple-flow-item lp-reveal"
            key={step.label}
            style={{ "--lp-reveal-i": index } as CSSProperties}
          >
            <span>{step.label}</span>
            <h3>{step.title}</h3>
            <p>{step.text}</p>
          </article>
        ))}
      </div>
    </section>
  );
}

function CoreScreensSection() {
  return (
    <section className="lp-simple-section lp-simple-screens">
      <ScreenFeature
        id="entry"
        eyebrow="ENTRY"
        title="応募先の情報は、Entry に全部置く。"
        body="会社ごとの応募フェーズ、応募経路、メモ、元URLを一覧で確認。探す時間を減らして、次の準備に集中できます。"
        image={productShots.entry}
        alt="Entré の Entry 一覧画面。応募先ごとの選考フェーズと応募経路を一覧表示している。"
      />
      <ScreenFeature
        id="kanban"
        eyebrow="KANBAN BOARD"
        title="選考状況は、カンバンで迷わない。"
        body="Entry をステージごとに並べて、進捗を俯瞰。今詰まっている列と次に動かす応募先がすぐ分かります。"
        image={productShots.kanban}
        alt="Entré のカンバン画面。Entry をエントリー、書類通過、面接、内定などの列で管理している。"
        mobileImage={productShots.mobileKanban}
        reverse
      />
    </section>
  );
}

function ScreenFeature({
  id,
  eyebrow,
  title,
  body,
  image,
  alt,
  mobileImage,
  reverse = false,
}: {
  id: string;
  eyebrow: string;
  title: string;
  body: string;
  image: { src: string; width: number; height: number };
  alt: string;
  mobileImage?: { src: string; width: number; height: number };
  reverse?: boolean;
}) {
  return (
    <div className={`lp-simple-screen ${reverse ? "is-reverse" : ""}`} id={id}>
      <div className="lp-simple-screen-copy lp-reveal">
        <p className="lp-simple-kicker">{eyebrow}</p>
        <h2>{title}</h2>
        <p>{body}</p>
      </div>
      <div className="lp-simple-screen-shot lp-reveal" style={{ "--lp-reveal-i": 1 } as CSSProperties}>
        <Image
          src={image.src}
          alt={alt}
          width={image.width}
          height={image.height}
          priority
          unoptimized
          className="lp-simple-product-shot"
        />
        {mobileImage ? (
          <Image
            src={mobileImage.src}
            alt=""
            width={mobileImage.width}
            height={mobileImage.height}
            priority
            unoptimized
            className="lp-simple-screen-mobile"
            aria-hidden
          />
        ) : null}
      </div>
    </div>
  );
}

function FinalCTASection() {
  return (
    <section className="lp-simple-final lp-reveal">
      <div>
        <p className="lp-simple-kicker">START</p>
        <h2>Entry とカンバンで、就活を軽くする。</h2>
      </div>
      <Link href="/login" className="lp-btn lp-btn-primary lp-btn-lg">
        <GoogleG />
        Googleで無料登録
      </Link>
    </section>
  );
}

function Footer() {
  return <footer className="lp-footer">© 2026 Entré</footer>;
}

function GoogleG() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" aria-hidden>
      <path
        fill="#fff"
        d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23zM5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18A10.99 10.99 0 001 12c0 1.77.42 3.44 1.18 4.93l3.66-2.84zM12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
      />
    </svg>
  );
}
