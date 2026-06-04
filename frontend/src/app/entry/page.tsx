// Server Component。エントリー一覧を SSR で取得して EntryListView に渡す。

import { redirect } from "next/navigation";
import Link from "next/link";
import { getCurrentUserServer } from "@/lib/auth-server";
import {
  listEntriesWithCompanyNamesServer,
  getNavCountsServer,
} from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { EntryListView } from "@/components/entre/EntryListView";
import { Plus } from "lucide-react";

export default async function EntryListPage() {
  const user = await getCurrentUserServer();
  if (!user) redirect("/login");

  const [entries, navCounts] = await Promise.all([
    listEntriesWithCompanyNamesServer(),
    getNavCountsServer(),
  ]);

  return (
    <AppShell userName={user.name} userSubtitle={user.email} navCounts={navCounts}>
      <div className="mx-auto max-w-[900px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-start justify-between gap-3">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">
              Entry
              <span className="ml-2 align-middle text-[12px] font-bold text-ink-3">
                応募先リスト
              </span>
            </h1>
            <p className="mt-0.5 text-[11px] text-ink-3">
              応募先ごとに選考フェーズ・応募経路・次のTaskを整理します。
            </p>
          </div>
          <Link
            href="/entry/new"
            className="inline-flex items-center gap-1.5 rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40"
          >
            <Plus size={13} aria-hidden />
            応募先を追加
          </Link>
        </header>

        <EntryListView entries={entries} />
      </div>
    </AppShell>
  );
}
