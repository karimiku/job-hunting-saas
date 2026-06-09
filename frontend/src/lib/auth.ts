// 認証フロー:
//   1. Firebase popup で Google サインイン → ID Token 取得
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

async function createBackendSession(idToken: string): Promise<AuthUser> {
  const res = await apiFetch("/auth/session", {
    method: "POST",
    body: JSON.stringify({ idToken }),
  });

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
  const result = await getRedirectResult(getFirebaseAuth());
  if (!result) return null;

  const idToken = await result.user.getIdToken();
  return createBackendSession(idToken);
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
