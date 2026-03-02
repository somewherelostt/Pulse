"use client";

import { createBrowserClient } from "@supabase/ssr";

export function createClient() {
  // During SSR/prerender NEXT_PUBLIC_ vars may be absent — fall back to
  // placeholder strings so module initialisation never throws.  The real
  // values are injected by Vercel at runtime in the browser.
  const url = process.env.NEXT_PUBLIC_SUPABASE_URL ?? "https://placeholder.supabase.co";
  const key = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY ?? "placeholder-anon-key";
  return createBrowserClient(url, key);
}
