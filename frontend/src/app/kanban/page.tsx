// Server Component。entries を SSR で取得し、KanbanBoard (Client) に initial data を渡す。

import Link from "next/link";
import { redirect } from "next/navigation";
import { getAppPageDataServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { KanbanBoard } from "@/components/entre/KanbanBoard";
import { EntryViewSwitch } from "@/components/entre/EntryViewSwitch";
import { Plus } from "lucide-react";

export default async function KanbanPage() {
  const pageData = await getAppPageDataServer();
  if (!pageData) redirect("/login");
  const { user, entries, navCounts } = pageData;

  return (
    <AppShell userName={user.name} userSubtitle={user.email} navCounts={navCounts}>
      <div className="mx-auto max-w-[1400px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-baseline justify-between">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">ボード</h1>
            <p className="mt-0.5 text-[12px] text-ink-3">
              選考が進んだら、カードを次のフェーズへドラッグして動かします。
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
          <EntryViewSwitch active="board" />
        </div>

        <KanbanBoard initialEntries={entries} />
      </div>
    </AppShell>
  );
}
