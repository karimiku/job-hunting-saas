// Server Component。entries を SSR で取得し、KanbanBoard (Client) に initial data を渡す。

import Link from "next/link";
import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import { listEntriesServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { KanbanBoard } from "@/components/entre/KanbanBoard";

export default async function KanbanPage() {
  const [user, entries] = await Promise.all([
    getCurrentUserServer(),
    listEntriesServer().catch(() => [] as never[]),
  ]);
  if (!user) redirect("/login");

  return (
    <AppShell userName={user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[1400px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-5 flex items-baseline justify-between">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">選考カンバン</h1>
            <p className="mt-0.5 text-[11px] text-ink-3">ステータスごとに整理</p>
          </div>
          <Link
            href="/entry/new"
            className="rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5"
          >
            ＋ エントリー
          </Link>
        </header>

        <KanbanBoard initialEntries={entries} />
      </div>
    </AppShell>
  );
}
