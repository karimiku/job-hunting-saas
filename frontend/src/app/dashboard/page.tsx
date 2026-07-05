// Server Component。auth + entries + tasks + inbox を SSR 集約APIで取得し、
// 集計済みデータを子コンポーネントに props で渡す。useEffect は使わない。

import { cache } from "react";
import { redirect } from "next/navigation";
import { getAppPageDataServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { DashboardEntries } from "@/components/entre/DashboardEntries";
import { DashboardQuests } from "@/components/entre/DashboardQuests";
import { DashboardNextAction } from "@/components/entre/DashboardNextAction";
import { DashboardStats } from "@/components/entre/DashboardStats";
import { GettingStartedGuide } from "@/components/entre/GettingStartedGuide";
import { SignOutButton } from "@/components/entre/SignOutButton";

// React.cache で 1 リクエスト内 memoize。Date.now() 自体は impure だが、cache() で
// 包むことで「同一リクエストでは同じ値」を保証でき、components-and-hooks-must-be-pure 規則も満たす。
const getRenderedAt = cache(() => Date.now());

export default async function DashboardPage() {
  const pageData = await getAppPageDataServer();
  if (!pageData) redirect("/login");
  const { user, entries, tasks, navCounts } = pageData;

  const firstName = user.name.split(/[\s　]/)[0] || user.name;
  const openTasks = tasks.filter((t) => t.status === "todo").length;
  const renderedAt = getRenderedAt();
  // 超過判定はカレンダー日で行う（時刻差の floor だと、本日締切(翌0時UTC)を夕方に
  // 見たとき差が負になり本日ぶんまで超過に数えてしまうため、ローカル0時基準で比較）。
  const today = new Date(renderedAt);
  const startOfToday = new Date(
    today.getFullYear(),
    today.getMonth(),
    today.getDate(),
  ).getTime();
  const overdueTasks = tasks.filter((t) => {
    if (t.status !== "todo" || !t.dueDate) return false;
    const d = new Date(t.dueDate);
    const startOfDue = new Date(
      d.getFullYear(),
      d.getMonth(),
      d.getDate(),
    ).getTime();
    return startOfDue < startOfToday;
  }).length;
  const hasEntries = entries.length > 0;
  const hasTasks = tasks.length > 0;
  const hasDoneTasks = tasks.some((t) => t.status === "done");

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
            <span className="rounded-full bg-sage-wash px-3.5 py-1.5 text-[12px] font-bold text-sage">
              未完了 {openTasks}
            </span>
            <SignOutButton />
          </div>
        </header>

        {hasEntries ? (
          <>
            <DashboardNextAction
              inboxCount={navCounts.inbox}
              entryCount={entries.length}
              openTaskCount={openTasks}
              overdueTaskCount={overdueTasks}
            />

            <div className="mb-4 md:mb-5">
              <DashboardStats entries={entries} />
            </div>
          </>
        ) : (
          <GettingStartedGuide
            hasEntries={hasEntries}
            hasTasks={hasTasks}
            hasDoneTasks={hasDoneTasks}
          />
        )}

        {hasEntries && !hasTasks && (
          <GettingStartedGuide
            hasEntries={hasEntries}
            hasTasks={hasTasks}
            hasDoneTasks={hasDoneTasks}
            compact
          />
        )}

        <div className="grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
          <DashboardEntries entries={entries} tasks={tasks} />
          <DashboardQuests tasks={tasks} />
        </div>
      </div>
    </AppShell>
  );
}
