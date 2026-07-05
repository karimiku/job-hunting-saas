"use client";

import { ApiError, type AuthUser } from "./api/client-types";
import { apiFetch } from "./api/client";
import {
  createSupabaseBrowserClient,
  hasSupabaseBrowserConfig,
} from "./supabase/client";

export type { AuthUser };

function authCallbackUrl(next: string): string {
  if (typeof window === "undefined") {
    throw new Error("Google login must be started in the browser");
  }
  const url = new URL("/auth/callback", window.location.origin);
  url.searchParams.set("next", next);
  return url.toString();
}

/**
 * Google OAuth を Supabase Auth の PKCE flow で開始する。
 */
export async function startGoogleRedirectSignIn(next = "/dashboard"): Promise<void> {
  if (!hasSupabaseBrowserConfig()) {
    throw new Error("Supabase Auth is not configured");
  }

  const supabase = createSupabaseBrowserClient();
  const { error } = await supabase.auth.signInWithOAuth({
    provider: "google",
    options: {
      redirectTo: authCallbackUrl(next),
    },
  });
  if (error) throw error;
}

/**
 * サインアウト。Supabase session を破棄し、移行期間中の legacy backend session も best-effort で破棄する。
 */
export async function signOut(): Promise<void> {
  if (hasSupabaseBrowserConfig()) {
    const supabase = createSupabaseBrowserClient();
    await supabase.auth.signOut().catch(() => {});
  }
  await apiFetch("/auth/session", { method: "DELETE" }).catch(() => {});
}

/**
 * 現在ログイン中のユーザーを取得。未ログインなら null。
 */
export async function fetchCurrentUser(): Promise<AuthUser | null> {
  try {
    return await apiFetch<AuthUser>("/auth/me");
  } catch (error) {
    if (error instanceof ApiError && error.unauthorized) return null;
    throw error;
  }
}
