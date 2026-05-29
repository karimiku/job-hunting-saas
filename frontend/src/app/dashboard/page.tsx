// Server Component。auth + entries + tasks + inbox を SSR で並列取得し、
// 集計済みデータを子コンポーネントに props で渡す。useEffect は使わない。

import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import {
  listEntriesWithCompanyNamesServer,
  listAllTasksServer,
  getNavCountsServer,
} from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";
import { Reveal } from "@/components/entre/Reveal";
import { DashboardStats, summarizeEntries } from "@/components/entre/DashboardStats";
import { DashboardQuests } from "@/components/entre/DashboardQuests";
import { DashboardEncouragement } from "@/components/entre/DashboardEncouragement";
import { DashboardNextAction } from "@/components/entre/DashboardNextAction";
import { StatusBreakdown } from "@/components/entre/StatusBreakdown";
import { SignOutButton } from "@/components/entre/SignOutButton";

export default async function DashboardPage() {
  // user は独立、entries/navCounts も独立なので並列取得 (cookies() は内部で memoize される)
  const [user, entries, navCounts] = await Promise.all([
    getCurrentUserServer(),
    listEntriesWithCompanyNamesServer().catch(() => []),
    getNavCountsServer(),
  ]);
  if (!user) redirect("/login");

  // タスクは entry 単位 API しか無いので entries を引いてから集約する。
  const tasks = await listAllTasksServer(entries).catch(() => []);

  const firstName = user.name.split(/[\s　]/)[0] || user.name;
  const stats = summarizeEntries(entries);
  const openTasks = tasks.filter((t) => t.status === "todo").length;

  return (
    <AppShell userName={user.name} userSubtitle={user.email} navCounts={navCounts}>
      <div className="mx-auto max-w-[1100px] px-5 py-6 md:px-8 md:py-8">
        {/* Greeting + bowing mascot (看板) */}
        <header className="mb-5 flex flex-col gap-3 md:mb-6 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-3 animate-[entre-fade-in_0.6s_both]">
            <Mascot mood="bow" size={56} />
            <div>
              <p
                className="font-hand text-[22px] text-sage md:text-2xl"
                style={{ transform: "rotate(-1.5deg)", display: "inline-block" }}
              >
                Welcome back,
              </p>
              <h1 className="font-serif text-2xl font-extrabold tracking-tight md:text-[28px]">
                {firstName}
                <span className="ml-1 text-sm font-medium text-ink-2 md:text-base">
                  さん、今日の就活状況です
                </span>
              </h1>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <div className="rounded-full bg-sage-wash px-3.5 py-1.5 text-[11px] font-bold text-sage">
              未完了Task {openTasks}件
            </div>
            <SignOutButton />
          </div>
        </header>

        {/* Stats — SSR 集計 */}
        <section className="mb-5 md:mb-6">
          <DashboardStats entries={entries} />
        </section>

        <Reveal delay={100}>
          <DashboardNextAction
            inboxCount={navCounts.inbox}
            entryCount={entries.length}
            openTaskCount={openTasks}
          />
        </Reveal>

        {/* Quest + Status */}
        <div className="grid grid-cols-1 gap-4 md:grid-cols-[1.4fr_1fr]">
          {/* Today's quest — 実タスクから集計 */}
          <Reveal delay={150}>
            <DashboardQuests tasks={tasks} />
          </Reveal>

          {/* Status pie — SSR 集計 */}
          <Reveal delay={250}>
            <div className="rounded-xl border border-line bg-surface p-5">
              <h2 className="mb-3 text-[13px] font-extrabold">選考ステータス</h2>
              <StatusBreakdown entries={entries} />
            </div>
          </Reveal>
        </div>

        {/* Mascot encouragement — 実データから出し分け */}
        <Reveal delay={350}>
          <DashboardEncouragement
            interviewing={stats.interviewing}
            offered={stats.offered}
            openTasks={openTasks}
          />
        </Reveal>
      </div>
    </AppShell>
  );
}
