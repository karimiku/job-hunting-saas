// 認証フロー:
//   1. Firebase redirect で Google サインイン → ID Token 取得
//   2. ID Token をバックエンドに POST → 検証 + Session Cookie 発行
//   3. 以降のリクエストは httpOnly Cookie で認証
//
// ID Token は「Session Cookie に交換する」以外の用途では使わない（XSS 面が小さくなる）。

import {
  getRedirectResult,
  signInWithRedirect,
  signOut as firebaseSignOut,
} from "firebase/auth";
import { getFirebaseAuth, googleProvider } from "./firebase";
import { apiFetch } from "./api";
import type { AuthUser } from "./api/client-types";

export type { AuthUser };

const authPerfPrefix = "entre.auth.google_redirect";

function now(): number {
  return typeof performance === "undefined" ? Date.now() : performance.now();
}

function mark(name: string): void {
  if (typeof performance === "undefined" || typeof performance.mark !== "function") return;
  performance.mark(`${authPerfPrefix}.${name}`);
}

function measure(name: string, start: string, end: string): void {
  if (typeof performance === "undefined" || typeof performance.measure !== "function") return;
  performance.measure(
    `${authPerfPrefix}.${name}`,
    `${authPerfPrefix}.${start}`,
    `${authPerfPrefix}.${end}`,
  );
}

function logAuthPerf(label: string, startedAt: number): void {
  const enabled =
    process.env.NODE_ENV !== "production" ||
    (typeof localStorage !== "undefined" &&
      localStorage.getItem("entre.authPerf") === "1");
  if (!enabled) return;

  console.info(`[auth] ${label}: ${Math.round(now() - startedAt)}ms`);
}

async function createBackendSession(idToken: string): Promise<AuthUser> {
  const startedAt = now();
  const res = await apiFetch("/auth/session", {
    method: "POST",
    body: JSON.stringify({ idToken }),
  });
  logAuthPerf("POST /auth/session", startedAt);

  if (!res.ok) {
    // バックエンド側に Session Cookie が作れなかった場合は Firebase 側も残さない。
    await firebaseSignOut(getFirebaseAuth()).catch(() => {});
    throw new Error(`session creation failed: ${res.status}`);
  }

  return (await res.json()) as AuthUser;
}

/**
 * Google サインインを redirect フローで開始する。
 *
 * popup フローはブラウザの Cross-Origin-Opener-Policy と相性が悪く、
 * Firebase SDK が popup.closed を確認するタイミングで console warning が出る。
 * redirect フローなら別windowを監視しないため、その警告を避けられる。
 */
export async function startGoogleRedirectSignIn(): Promise<void> {
  await signInWithRedirect(getFirebaseAuth(), googleProvider);
}

/**
 * Google redirect から戻った直後だけ Firebase の結果を回収し、
 * バックエンドの Session Cookie に交換する。
 */
export async function completeGoogleRedirectSignIn(): Promise<AuthUser | null> {
  mark("start");
  const startedAt = now();

  const redirectStartedAt = now();
  const result = await getRedirectResult(getFirebaseAuth());
  mark("redirect_result");
  measure("getRedirectResult", "start", "redirect_result");
  logAuthPerf("Firebase getRedirectResult", redirectStartedAt);
  if (!result) {
    logAuthPerf("Google redirect completion without result", startedAt);
    return null;
  }

  const tokenStartedAt = now();
  const idToken = await result.user.getIdToken();
  mark("id_token");
  measure("getIdToken", "redirect_result", "id_token");
  logAuthPerf("Firebase getIdToken", tokenStartedAt);

  const user = await createBackendSession(idToken);
  mark("backend_session");
  measure("createBackendSession", "id_token", "backend_session");
  measure("total", "start", "backend_session");
  logAuthPerf("Google redirect completion total", startedAt);
  return user;
}

/**
 * サインアウト。Firebase 側とバックエンド側の両方のセッションを失効させる。
 */
export async function signOut(): Promise<void> {
  await Promise.all([
    apiFetch("/auth/session", { method: "DELETE" }).catch(() => {}),
    firebaseSignOut(getFirebaseAuth()).catch(() => {}),
  ]);
}

/**
 * 現在ログイン中のユーザーを取得。未ログインなら null。
 */
export async function fetchCurrentUser(): Promise<AuthUser | null> {
  const res = await apiFetch("/auth/me");
  if (res.status === 401) return null;
  if (!res.ok) throw new Error(`/auth/me failed: ${res.status}`);
  return (await res.json()) as AuthUser;
}
