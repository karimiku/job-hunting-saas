// Server Component から認証状態を取得するヘルパー。
// 401 は guest として返し、それ以外のエラーは投げる (page で error.tsx に拾わせる)。

import { cache } from "react";
import { serverFetch } from "./api/server";
import { ApiError, type AuthUser } from "./api/client-types";

/**
 * Server Component から `/auth/me` を叩いて現在のユーザーを返す。
 * 未ログイン (401) は null。
 */
export const getCurrentUserServer = cache(async (): Promise<AuthUser | null> => {
  try {
    return await serverFetch<AuthUser>("/auth/me");
  } catch (e) {
    if (e instanceof ApiError && e.unauthorized) return null;
    throw e;
  }
});
