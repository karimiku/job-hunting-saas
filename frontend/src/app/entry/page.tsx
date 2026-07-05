// Server Component。エントリー一覧を SSR で取得して EntryListView に渡す。

import { redirect } from "next/navigation";
import Link from "next/link";
import { getAppPageDataServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { EntryListView } from "@/components/entre/EntryListView";
import { EntryViewSwitch } from "@/components/entre/EntryViewSwitch";
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
              応募先
              <span className="ml-2 align-middle text-[12px] font-bold text-ink-3">
                {entries.length}社
              </span>
            </h1>
            <p className="mt-0.5 text-[12px] text-ink-3">
              受けている企業の一覧です。選考がどこまで進んでいるかをここで確認します。
            </p>
          </div>
          <Link
            href="/entry/new"
            prefetch={false}
            className="inline-flex items-center gap-1.5 rounded-lg bg-sage px-3 py-1.5 text-[12px] font-bold text-white transition-transform hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40"
          >
            <Plus size={13} aria-hidden />
            応募先を追加
          </Link>
        </header>

        <div className="mb-4">
          <EntryViewSwitch active="list" />
        </div>

        <EntryListView entries={entries} />
      </div>
    </AppShell>
  );
}
