"use client";

import { createBrowserClient } from "@supabase/ssr";
import { getSupabaseBrowserEnv, hasSupabaseBrowserConfig } from "./env";

export { hasSupabaseBrowserConfig };

export function createSupabaseBrowserClient() {
  const { url, publishableKey } = getSupabaseBrowserEnv();
  return createBrowserClient(url, publishableKey);
}

export async function getSupabaseBrowserAccessToken(): Promise<string | null> {
  if (!hasSupabaseBrowserConfig()) return null;

  const supabase = createSupabaseBrowserClient();
  const { data, error } = await supabase.auth.getSession();
  if (error) return null;
  return data.session?.access_token ?? null;
}

