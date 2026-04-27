"use client";

import { useRouter, useParams } from "next/navigation";
import { useEffect } from "react";
import Link from "next/link";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { EntryDetailView } from "@/components/entre/EntryDetailView";

export default function EntryDetailPage() {
  const router = useRouter();
  const state = useUser();
  const { entryId } = useParams<{ entryId: string }>();

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  return (
    <AppShell userName={state.user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[700px] px-5 py-6 md:px-8 md:py-7">
        <Link
          href="/entry"
          className="mb-3 inline-flex items-center gap-1 text-[11px] font-semibold text-ink-3 hover:text-sage"
        >
          ‹ Entry 一覧
        </Link>
        <EntryDetailView entryId={entryId} />
      </div>
    </AppShell>
  );
}
