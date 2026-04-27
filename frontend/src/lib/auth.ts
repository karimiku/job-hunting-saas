// 認証フロー:
//   1. Firebase popup で Google サインイン → ID Token 取得
//   2. ID Token をバックエンドに POST → 検証 + Session Cookie 発行
//   3. 以降のリクエストは httpOnly Cookie で認証
//
// ID Token は「Session Cookie に交換する」以外の用途では使わない（XSS 面が小さくなる）。

import { signInWithPopup, signOut as firebaseSignOut } from "firebase/auth";
import { getFirebaseAuth, googleProvider } from "./firebase";
import { apiFetch } from "./api";
import type { AuthUser } from "./api/client-types";

export type { AuthUser };

/**
 * Google サインイン → バックエンドにセッション確立。
 * 成功時は AuthUser を返す。
 */
export async function signInWithGoogle(): Promise<AuthUser> {
  const result = await signInWithPopup(getFirebaseAuth(), googleProvider);
  const idToken = await result.user.getIdToken();

  const res = await apiFetch("/auth/session", {
    method: "POST",
    body: JSON.stringify({ idToken }),
  });

  if (!res.ok) {
    // バックエンド側に残った Firebase session cookie はないのでクライアント状態だけクリア
    await firebaseSignOut(getFirebaseAuth()).catch(() => {});
    throw new Error(`session creation failed: ${res.status}`);
  }

  return (await res.json()) as AuthUser;
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
