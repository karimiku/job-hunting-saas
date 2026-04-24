"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { signOut } from "@/lib/auth";
import { useUser } from "@/lib/use-user";

export default function DashboardPage() {
  const router = useRouter();
  const state = useUser();
  const [signingOut, setSigningOut] = useState(false);

  // ゲストは /login にリダイレクト
  useEffect(() => {
    if (state.status === "guest") {
      router.replace("/login");
    }
  }, [state.status, router]);

  async function handleSignOut() {
    setSigningOut(true);
    try {
      await signOut();
      router.push("/login");
    } finally {
      setSigningOut(false);
    }
  }

  if (state.status !== "authenticated") {
    // loading / redirect 中はクリーム背景だけ出しておく
    return <div className="lp-scope" style={{ minHeight: "100vh", background: "var(--lp-cream)" }} />;
  }

  const user = state.user;

  return (
    <div
      className="lp-scope"
      style={{
        minHeight: "100vh",
        background: "var(--lp-cream)",
        color: "var(--lp-ink)",
        fontFamily: "var(--lp-font-jp)",
      }}
    >
      {/* ヘッダー */}
      <header
        style={{
          borderBottom: "1px solid var(--lp-line)",
          background: "var(--lp-cream-2)",
        }}
      >
        <div
          style={{
            maxWidth: 1200,
            margin: "0 auto",
            padding: "16px 24px",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          <Link
            href="/"
            style={{
              fontFamily: "var(--lp-font-serif)",
              fontWeight: 800,
              fontSize: 22,
              letterSpacing: "-0.01em",
              color: "var(--lp-ink)",
            }}
          >
            Entré
          </Link>

          <div style={{ display: "flex", alignItems: "center", gap: 14 }}>
            <span
              style={{
                fontSize: 14,
                color: "var(--lp-ink-2)",
              }}
            >
              {user.name}
            </span>
            <button
              type="button"
              onClick={handleSignOut}
              disabled={signingOut}
              className="lp-btn lp-btn-ghost"
              style={{ height: 36, padding: "0 14px", fontSize: 13 }}
            >
              {signingOut ? "…" : "ログアウト"}
            </button>
          </div>
        </div>
      </header>

      {/* メイン */}
      <main
        style={{
          maxWidth: 1100,
          margin: "0 auto",
          padding: "64px 24px 120px",
        }}
      >
        {/* ウェルカム */}
        <section style={{ marginBottom: 56 }}>
          <p
            className="lp-hand"
            style={{
              fontSize: 28,
              color: "var(--lp-sage-2)",
              marginBottom: 4,
            }}
          >
            おかえり、
          </p>
          <h1
            className="lp-serif"
            style={{
              fontSize: 44,
              fontWeight: 800,
              letterSpacing: "-0.02em",
              lineHeight: 1.2,
              marginBottom: 12,
            }}
          >
            {user.name} さん
          </h1>
          <p style={{ fontSize: 15, color: "var(--lp-ink-2)" }}>
            散らかった就活、ぜんぶこの1枚に集めていきましょう。
          </p>
        </section>

        {/* 選考リスト（空状態） */}
        <section
          className="lp-card"
          style={{
            background: "var(--lp-surface)",
            borderRadius: 20,
            border: "1px solid var(--lp-line)",
            padding: "48px 32px",
            textAlign: "center",
          }}
        >
          <div
            style={{
              width: 64,
              height: 64,
              margin: "0 auto 20px",
              borderRadius: 16,
              background: "var(--lp-sage-tint)",
              display: "grid",
              placeItems: "center",
              fontSize: 28,
            }}
            aria-hidden
          >
            🦌
          </div>
          <h2
            className="lp-serif"
            style={{ fontSize: 22, fontWeight: 700, marginBottom: 8 }}
          >
            まだ選考がありません
          </h2>
          <p
            style={{
              fontSize: 14,
              color: "var(--lp-ink-2)",
              marginBottom: 24,
              maxWidth: 420,
              marginLeft: "auto",
              marginRight: "auto",
            }}
          >
            応募した企業やエントリーを1件ずつ追加していきましょう。ここから先は順次機能を追加します。
          </p>
          <button
            type="button"
            className="lp-btn lp-btn-sage"
            disabled
            style={{ opacity: 0.6, cursor: "not-allowed" }}
          >
            + 選考を追加（開発中）
          </button>
        </section>

        {/* 今後追加予定 */}
        <section style={{ marginTop: 56 }}>
          <h3
            className="lp-serif"
            style={{
              fontSize: 18,
              fontWeight: 700,
              marginBottom: 20,
              color: "var(--lp-ink-2)",
            }}
          >
            これから追加する機能
          </h3>
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))",
              gap: 16,
            }}
          >
            <UpcomingCard
              tint="var(--lp-s-entry)"
              ink="var(--lp-s-entry-ink)"
              title="選考の一元管理"
              body="マイナビ・リクナビ・ワンキャリ、全部ここに。"
            />
            <UpcomingCard
              tint="var(--lp-s-interview)"
              ink="var(--lp-s-interview-ink)"
              title="AI メール取り込み"
              body="選考メールを貼るだけで日時を抽出。"
            />
            <UpcomingCard
              tint="var(--lp-s-doc)"
              ink="var(--lp-s-doc-ink)"
              title="タスクとリマインド"
              body="締切や面接を忘れないように通知。"
            />
            <UpcomingCard
              tint="var(--lp-s-offer)"
              ink="var(--lp-s-offer-ink)"
              title="Google カレンダー連携"
              body="登録した日程を自動でカレンダーへ。"
            />
          </div>
        </section>
      </main>
    </div>
  );
}

function UpcomingCard({
  tint,
  ink,
  title,
  body,
}: {
  tint: string;
  ink: string;
  title: string;
  body: string;
}) {
  return (
    <div
      style={{
        background: "var(--lp-surface)",
        border: "1px solid var(--lp-line)",
        borderRadius: 16,
        padding: 20,
      }}
    >
      <div
        style={{
          width: 36,
          height: 36,
          borderRadius: 10,
          background: tint,
          color: ink,
          display: "grid",
          placeItems: "center",
          fontSize: 14,
          fontWeight: 700,
          marginBottom: 12,
        }}
        aria-hidden
      >
        ✦
      </div>
      <h4 style={{ fontSize: 15, fontWeight: 700, marginBottom: 6 }}>{title}</h4>
      <p style={{ fontSize: 13, color: "var(--lp-ink-2)", lineHeight: 1.6 }}>{body}</p>
    </div>
  );
}
