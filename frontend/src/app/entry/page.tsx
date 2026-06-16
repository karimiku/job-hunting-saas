// Server Component。エントリー一覧を SSR で取得して EntryListView に渡す。

import { redirect } from "next/navigation";
import Link from "next/link";
import { getAppPageDataServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { EntryListView } from "@/components/entre/EntryListView";
import { Plus } from "lucide-react";

export default async function EntryListPage() {
  const pageData = await getAppPageDataServer();
  if (!pageData) redirect("/login");
  const { user, entries, navCounts } = pageData;

  return (
    <AppShell userName={user.name} userSubtitle={user.email} navCounts={navCounts}>
      <div className="mx-auto max-w-[900px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-start justify-between gap-3">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">
              Entry
              <span className="ml-2 align-middle text-[12px] font-bold text-ink-3">
                応募先 {entries.length}件
              </span>
            </h1>
            <p className="mt-0.5 text-[11px] text-ink-3">
              企業ごとに選考フェーズ、応募経路、次のタスクを整理します。
            </p>
          </div>
          <Link
            href="/entry/new"
            className="inline-flex items-center gap-1.5 rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40"
          >
            <Plus size={13} aria-hidden />
            Entryを追加
          </Link>
        </header>

        <EntryListView entries={entries} />
      </div>
    </AppShell>
  );
}
