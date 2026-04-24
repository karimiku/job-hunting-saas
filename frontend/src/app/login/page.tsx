"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { signInWithGoogle } from "@/lib/auth";
import { useUser } from "@/lib/use-user";
import { DeerMascot } from "@/components/landing/mascot";

export default function LoginPage() {
  const router = useRouter();
  const state = useUser();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Session Cookie がまだ有効ならログインボタンを見せずに即遷移
  useEffect(() => {
    if (state.status === "authenticated") {
      router.replace("/dashboard");
    }
  }, [state.status, router]);

  async function handleGoogleSignIn() {
    // Firebase の公開 env がまだ設定されていない環境（本番投入前のプレビュー等）では
    // ログインを試みず LP に戻す。invalid-api-key のエラー画面を見せないため。
    if (!process.env.NEXT_PUBLIC_FIREBASE_API_KEY) {
      router.push("/");
      return;
    }
    setLoading(true);
    setError(null);
    try {
      await signInWithGoogle();
      router.push("/dashboard");
    } catch (err) {
      const code = (err as { code?: string })?.code;
      // ユーザーがポップアップを閉じた / ブラウザが多重ポップアップを抑制した = キャンセル扱い
      if (code === "auth/popup-closed-by-user" || code === "auth/cancelled-popup-request") {
        setLoading(false);
        return;
      }
      // API キー不正（env 差異など）は LP に逃がす
      if (code === "auth/invalid-api-key" || code === "auth/api-key-not-valid") {
        router.push("/");
        return;
      }
      if (code === "auth/popup-blocked") {
        setError("ブラウザがポップアップをブロックしました。ポップアップを許可して再度お試しください。");
      } else {
        console.error(err);
        setError(err instanceof Error ? err.message : "ログインに失敗しました");
      }
      setLoading(false);
    }
  }

  // 初期ロード・遷移中は背景だけ
  if (state.status !== "guest") {
    return (
      <div
        className="lp-scope"
        style={{ minHeight: "100vh", background: "var(--lp-cream)" }}
      />
    );
  }

  return (
    <div
      className="lp-scope"
      style={{
        minHeight: "100vh",
        background: "var(--lp-cream)",
        color: "var(--lp-ink)",
        fontFamily: "var(--lp-font-jp)",
        display: "flex",
        flexDirection: "column",
      }}
    >
      {/* 控えめなヘッダー */}
      <header style={{ padding: "24px 32px" }}>
        <Link
          href="/"
          className="lp-serif"
          style={{
            fontWeight: 800,
            fontSize: 22,
            letterSpacing: "-0.01em",
            color: "var(--lp-ink)",
          }}
        >
          Entré
        </Link>
      </header>

      {/* メイン */}
      <main
        style={{
          flex: 1,
          display: "grid",
          placeItems: "center",
          padding: "24px",
        }}
      >
        <div style={{ width: "100%", maxWidth: 420, textAlign: "center" }}>
          {/* マスコット */}
          <div style={{ display: "flex", justifyContent: "center", marginBottom: 28 }}>
            <DeerMascot size={120} mood="wave" tilt={-4} />
          </div>

          {/* タイトル */}
          <p
            className="lp-hand"
            style={{
              fontSize: 26,
              color: "var(--lp-sage-2)",
              marginBottom: 6,
              lineHeight: 1,
            }}
          >
            welcome back
          </p>
          <h1
            className="lp-serif"
            style={{
              fontSize: 34,
              fontWeight: 800,
              letterSpacing: "-0.02em",
              lineHeight: 1.25,
              marginBottom: 10,
            }}
          >
            Entré にログイン
          </h1>
          <p
            style={{
              fontSize: 14,
              color: "var(--lp-ink-2)",
              marginBottom: 36,
            }}
          >
            散らかった就活を、1枚に。
          </p>

          {/* カード */}
          <div
            style={{
              background: "var(--lp-surface)",
              border: "1px solid var(--lp-line)",
              borderRadius: 20,
              padding: "28px 24px",
              boxShadow: "0 1px 0 rgba(43, 53, 48, 0.02)",
            }}
          >
            {/* Google 公式 Light スタイル: 白背景 + 1px グレーボーダー + 4色 G + 規定文言 */}
            <button
              type="button"
              onClick={handleGoogleSignIn}
              disabled={loading}
              style={{
                width: "100%",
                height: 44,
                display: "inline-flex",
                alignItems: "center",
                justifyContent: "center",
                gap: 10,
                padding: "0 16px",
                background: "#FFFFFF",
                border: "1px solid #DADCE0",
                borderRadius: 10,
                color: "#1F1F1F",
                fontSize: 14,
                fontWeight: 500,
                fontFamily:
                  "Roboto, 'Noto Sans JP', -apple-system, BlinkMacSystemFont, sans-serif",
                cursor: loading ? "wait" : "pointer",
                transition: "background-color 150ms ease, box-shadow 150ms ease",
                opacity: loading ? 0.7 : 1,
              }}
              onMouseEnter={(e) => {
                (e.currentTarget as HTMLButtonElement).style.backgroundColor = "#F8F9FA";
                (e.currentTarget as HTMLButtonElement).style.boxShadow =
                  "0 1px 2px rgba(60, 64, 67, 0.15)";
              }}
              onMouseLeave={(e) => {
                (e.currentTarget as HTMLButtonElement).style.backgroundColor = "#FFFFFF";
                (e.currentTarget as HTMLButtonElement).style.boxShadow = "none";
              }}
            >
              <GoogleColoredG size={18} />
              {loading ? "サインイン中…" : "Google でログイン"}
            </button>

            {error ? (
              <p
                role="alert"
                style={{
                  marginTop: 16,
                  fontSize: 13,
                  color: "#B33A3A",
                }}
              >
                {error}
              </p>
            ) : null}
          </div>

          <p
            style={{
              marginTop: 20,
              fontSize: 12,
              color: "var(--lp-ink-3)",
              lineHeight: 1.6,
            }}
          >
            続行することで、
            <Link href="/terms" style={{ color: "var(--lp-ink-2)", textDecoration: "underline" }}>
              利用規約
            </Link>
            および
            <Link href="/privacy" style={{ color: "var(--lp-ink-2)", textDecoration: "underline" }}>
              プライバシーポリシー
            </Link>
            に同意したものとみなします。
          </p>
        </div>
      </main>
    </div>
  );
}

// Google 公式 G ロゴ（4色、改変禁止）
function GoogleColoredG({ size = 18 }: { size?: number }) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden>
      <path
        d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.76h3.56c2.08-1.92 3.28-4.74 3.28-8.09Z"
        fill="#4285F4"
      />
      <path
        d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.56-2.76c-.99.66-2.25 1.06-3.72 1.06-2.86 0-5.29-1.93-6.15-4.53H2.18v2.84A11 11 0 0 0 12 23Z"
        fill="#34A853"
      />
      <path
        d="M5.85 14.11A6.6 6.6 0 0 1 5.5 12c0-.73.13-1.44.35-2.11V7.05H2.18A11 11 0 0 0 1 12c0 1.78.43 3.46 1.18 4.95l3.67-2.84Z"
        fill="#FBBC05"
      />
      <path
        d="M12 5.38c1.62 0 3.06.56 4.2 1.65l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.05l3.67 2.84C6.71 7.3 9.14 5.38 12 5.38Z"
        fill="#EA4335"
      />
    </svg>
  );
}
