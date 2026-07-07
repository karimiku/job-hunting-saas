"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { DeerMascot } from "@/components/landing/mascot";
import { startGoogleRedirectSignIn } from "@/lib/auth";
import { hasSupabaseBrowserConfig } from "@/lib/supabase/client";

type LoginPhase = "idle" | "submitting" | "redirecting";

type LoginClientProps = {
  devLoginEnabled?: boolean;
};

export function LoginClient({ devLoginEnabled = false }: LoginClientProps) {
  const router = useRouter();
  const [phase, setPhase] = useState<LoginPhase>("idle");
  const [error, setError] = useState<string | null>(() => {
    if (typeof window === "undefined") {
      return null;
    }
    const params = new URLSearchParams(window.location.search);
    return params.get("error")
      ? "ログインに失敗しました。時間を置いてもう一度お試しください。"
      : null;
  });

  async function handleGoogleSignIn() {
    if (!hasSupabaseBrowserConfig()) {
      router.push("/");
      return;
    }

    setPhase("submitting");
    setError(null);

    try {
      await startGoogleRedirectSignIn("/dashboard");
      setPhase("redirecting");
    } catch (err) {
      setError(err instanceof Error ? err.message : "ログインに失敗しました");
      setPhase("idle");
    }
  }

  const loading = phase !== "idle";

  if (phase === "redirecting") {
    return (
      <LoginLoadingScreen title="Googleログインへ移動しています" />
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
      <header style={{ padding: "24px 32px" }}>
        <Link
          href="/"
          className="lp-serif"
          style={{
            fontWeight: 800,
            fontSize: 22,
            letterSpacing: "0",
            color: "var(--lp-ink)",
          }}
        >
          Entré
        </Link>
      </header>

      <main
        style={{
          flex: 1,
          display: "grid",
          placeItems: "center",
          padding: "24px",
        }}
      >
        <div style={{ width: "100%", maxWidth: 420, textAlign: "center" }}>
          <div style={{ display: "flex", justifyContent: "center", marginBottom: 28 }}>
            <DeerMascot size={120} mood="wave" tilt={-4} />
          </div>

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
              letterSpacing: "0",
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

          <div
            style={{
              background: "var(--lp-surface)",
              border: "1px solid var(--lp-line)",
              borderRadius: 20,
              padding: "28px 24px",
              boxShadow: "0 1px 0 rgba(43, 53, 48, 0.02)",
            }}
          >
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
              {phase === "submitting" ? "Googleへ移動中..." : "Google でログイン"}
            </button>
            {devLoginEnabled ? (
              <Link
                href="/dev/login"
                style={{
                  marginTop: 12,
                  width: "100%",
                  height: 40,
                  display: "inline-flex",
                  alignItems: "center",
                  justifyContent: "center",
                  borderRadius: 10,
                  border: "1px solid var(--lp-line)",
                  background: "transparent",
                  color: "var(--lp-ink-2)",
                  fontSize: 13,
                  fontWeight: 700,
                  textDecoration: "none",
                }}
              >
                開発用ログイン
              </Link>
            ) : null}

            <p
              style={{
                marginTop: 14,
                fontSize: 12,
                color: "var(--lp-ink-3)",
                lineHeight: 1.6,
              }}
            >
              取得するのはGoogleアカウントの氏名とメールアドレスだけ。メール本文や連絡先は読み取りません。無料で使えます。
            </p>

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

function LoginLoadingScreen({ title }: { title: string }) {
  return (
    <div
      className="lp-scope"
      style={{
        minHeight: "100vh",
        background: "var(--lp-cream)",
        color: "var(--lp-ink)",
        fontFamily: "var(--lp-font-jp)",
        display: "grid",
        placeItems: "center",
        padding: 24,
      }}
    >
      <main
        aria-live="polite"
        aria-busy="true"
        style={{
          width: "100%",
          maxWidth: 360,
          textAlign: "center",
          animation: "entre-fade-in 0.45s ease-out both",
        }}
      >
        <div
          style={{
            display: "inline-grid",
            placeItems: "center",
            position: "relative",
            marginBottom: 22,
          }}
        >
          <span
            aria-hidden
            style={{
              position: "absolute",
              width: 132,
              height: 132,
              borderRadius: "999px",
              background: "rgba(79, 110, 88, 0.14)",
              animation: "entre-pulse-ring 1.6s ease-out infinite",
            }}
          />
          <DeerMascot size={92} mood="sparkle" tilt={-3} />
        </div>

        <p
          className="lp-hand"
          style={{
            fontSize: 24,
            color: "var(--lp-sage-2)",
            lineHeight: 1,
            marginBottom: 8,
          }}
        >
          one moment
        </p>
        <h1
          className="lp-serif"
          style={{
            fontSize: 24,
            fontWeight: 800,
            letterSpacing: "0",
            lineHeight: 1.35,
            marginBottom: 8,
          }}
        >
          {title}
        </h1>
        <p
          style={{
            color: "var(--lp-ink-2)",
            fontSize: 13,
            lineHeight: 1.7,
          }}
        >
          セッションを安全に作成しています。
        </p>
      </main>
    </div>
  );
}
