type SupabaseBrowserEnv = {
  url: string;
  publishableKey: string;
};

export function hasSupabaseBrowserConfig(): boolean {
  return Boolean(
    process.env.NEXT_PUBLIC_SUPABASE_URL?.trim() &&
      process.env.NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY?.trim(),
  );
}

export function getSupabaseBrowserEnv(): SupabaseBrowserEnv {
  const url = process.env.NEXT_PUBLIC_SUPABASE_URL?.trim();
  const publishableKey = process.env.NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY?.trim();
  if (!url || !publishableKey) {
    throw new Error("Supabase public env is not configured");
  }
  return { url, publishableKey };
}

