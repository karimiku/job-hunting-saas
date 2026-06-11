// Server Component。auth + entries + tasks + inbox を SSR で並列取得し、
// 集計済みデータを子コンポーネントに props で渡す。useEffect は使わない。

import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import {
  attachCompanyNamesToTasks,
  buildNavCounts,
  listEntriesWithCompanyNamesServer,
  listInboxClipsServer,
  listTasksServer,
} from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { DashboardEntries } from "@/components/entre/DashboardEntries";
import { DashboardQuests } from "@/components/entre/DashboardQuests";
import { DashboardNextAction } from "@/components/entre/DashboardNextAction";
import { DashboardStats } from "@/components/entre/DashboardStats";
import { SignOutButton } from "@/components/entre/SignOutButton";

export default async function DashboardPage() {
  // user / entries / tasks / inbox は独立なので並列取得 (cookies() は内部で memoize される)
  const [user, entries, rawTasks, clips] = await Promise.all([
    getCurrentUserServer(),
    listEntriesWithCompanyNamesServer().catch(() => []),
    listTasksServer().catch(() => []),
    listInboxClipsServer().catch(() => []),
  ]);
  if (!user) redirect("/login");

  const tasks = attachCompanyNamesToTasks(rawTasks, entries);
  const navCounts = buildNavCounts(entries, rawTasks, clips);

  const firstName = user.name.split(/[\s　]/)[0] || user.name;
  const openTasks = tasks.filter((t) => t.status === "todo").length;

  return (
    <AppShell userName={user.name} userSubtitle={user.email} navCounts={navCounts}>
      <div className="mx-auto max-w-[980px] px-4 py-5 md:px-8 md:py-8">
        <header className="mb-4 flex flex-col gap-3 md:mb-5 md:flex-row md:items-start md:justify-between">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight md:text-[28px]">
              {firstName}さんのホーム
            </h1>
            <p className="mt-1 text-[12px] leading-relaxed text-ink-2">
              迷ったら、ここに出ている一つの作業から進めます。
            </p>
          </div>

          <div className="flex shrink-0 items-center gap-2">
            <span className="rounded-full bg-sage-wash px-3.5 py-1.5 text-[11px] font-bold text-sage">
              未完了 {openTasks}
            </span>
            <SignOutButton />
          </div>
        </header>

        <DashboardNextAction
          inboxCount={navCounts.inbox}
          entryCount={entries.length}
          openTaskCount={openTasks}
        />

        <div className="mb-4 md:mb-5">
          <DashboardStats entries={entries} />
        </div>

        <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
          <DashboardEntries entries={entries} tasks={tasks} />
          <DashboardQuests tasks={tasks} />
        </div>
      </div>
    </AppShell>
  );
}
