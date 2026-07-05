import { createServerClient } from "@supabase/ssr";
import { cookies } from "next/headers";
import { getSupabaseBrowserEnv, hasSupabaseBrowserConfig } from "./env";

export function hasSupabaseServerConfig(): boolean {
  return hasSupabaseBrowserConfig();
}

export async function createSupabaseServerClient() {
  const { url, publishableKey } = getSupabaseBrowserEnv();
  const cookieStore = await cookies();

  return createServerClient(url, publishableKey, {
    cookies: {
      getAll() {
        return cookieStore.getAll();
      },
      setAll(cookiesToSet) {
        try {
          cookiesToSet.forEach(({ name, value, options }) => {
            cookieStore.set(name, value, options);
          });
        } catch {
          // Server Components cannot write cookies; proxy/route handlers refresh them.
        }
      },
    },
  });
}

export async function getSupabaseServerAccessToken(): Promise<string | null> {
  if (!hasSupabaseServerConfig()) return null;

  const supabase = await createSupabaseServerClient();
  const { data, error } = await supabase.auth.getSession();
  if (error) return null;
  return data.session?.access_token ?? null;
}
