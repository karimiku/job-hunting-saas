// Server Component。エントリー一覧を SSR で取得して EntryListView に渡す。

import { redirect } from "next/navigation";
import Link from "next/link";
import { getCurrentUserServer } from "@/lib/auth-server";
import { listEntriesServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { EntryListView } from "@/components/entre/EntryListView";

export default async function EntryListPage() {
  const user = await getCurrentUserServer();
  if (!user) redirect("/login");

  const entries = await listEntriesServer();

  return (
    <AppShell userName={user.name} userSubtitle="○○大学 4年">
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

        <EntryListView entries={entries} />
      </div>
    </AppShell>
  );
}
