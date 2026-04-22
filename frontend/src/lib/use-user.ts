"use client";

import { useEffect, useState } from "react";
import { fetchCurrentUser, type AuthUser } from "./auth";

export type UserState =
  | { status: "loading"; user: null }
  | { status: "authenticated"; user: AuthUser }
  | { status: "guest"; user: null };

/**
 * 現在のユーザーをバックエンド /auth/me から取得する hook。
 * 初回マウント時に1回だけ問い合わせる。
 */
export function useUser(): UserState {
  const [state, setState] = useState<UserState>({ status: "loading", user: null });

  useEffect(() => {
    let cancelled = false;
    fetchCurrentUser()
      .then((user) => {
        if (cancelled) return;
        setState(user ? { status: "authenticated", user } : { status: "guest", user: null });
      })
      .catch(() => {
        if (cancelled) return;
        setState({ status: "guest", user: null });
      });
    return () => {
      cancelled = true;
    };
  }, []);

  return state;
}
