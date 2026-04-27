"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import Link from "next/link";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { EntryListView } from "@/components/entre/EntryListView";

export default function EntryListPage() {
  const router = useRouter();
  const state = useUser();

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  return (
    <AppShell userName={state.user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[900px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-baseline justify-between">
          <h1 className="font-serif text-2xl font-extrabold tracking-tight">
            Entry
          </h1>
          <Link
            href="/entry/new"
            className="rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5"
          >
            ＋ 追加
          </Link>
        </header>

        <EntryListView />
      </div>
    </AppShell>
  );
}
