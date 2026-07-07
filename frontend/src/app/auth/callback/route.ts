import { NextResponse } from "next/server";
import { createSupabaseServerClient } from "@/lib/supabase/server";

export async function GET(request: Request) {
  const requestUrl = new URL(request.url);
  const code = requestUrl.searchParams.get("code");
  const next = safeRedirectPath(requestUrl.searchParams.get("next"));
  const origin = trustedRedirectOrigin(requestUrl.origin);

  if (code) {
    try {
      const supabase = await createSupabaseServerClient();
      const { error } = await supabase.auth.exchangeCodeForSession(code);
      if (!error) {
        return NextResponse.redirect(`${origin}${next}`);
      }
    } catch {
      // Fall through to the login error redirect.
    }
  }

  return NextResponse.redirect(`${origin}/login?error=auth_callback`);
}

function safeRedirectPath(raw: string | null): string {
  if (!raw || !raw.startsWith("/") || raw.startsWith("//")) {
    return "/dashboard";
  }
  const target = new URL(raw, "http://localhost");
  if (target.origin !== "http://localhost") {
    return "/dashboard";
  }
  return `${target.pathname}${target.search}${target.hash}`;
}

// リダイレクト先の origin は攻撃者が制御できる x-forwarded-host / host ヘッダから
// 組み立てない。NEXT_PUBLIC_SITE_URL（確定した公開 URL）が設定されていればそれを使い、
// 未設定なら request.url 由来の origin（fallback、= 相対パスと同義でヘッダ非依存）を使う。
function trustedRedirectOrigin(fallback: string): string {
  const siteURL = process.env.NEXT_PUBLIC_SITE_URL?.trim();
  if (!siteURL) return fallback;
  try {
    return new URL(siteURL).origin;
  } catch {
    return fallback;
  }
}
