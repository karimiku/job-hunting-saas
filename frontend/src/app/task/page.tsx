// Server Component。ユーザーの全タスクを SSR で取得し、interactive 部分は Client に委譲する。

import { redirect } from "next/navigation";
import { getTaskPageDataServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { TaskListView } from "@/components/entre/TaskListView";

export default async function TaskPage() {
  const pageData = await getTaskPageDataServer();
  if (!pageData) redirect("/login");
  const { user, entries, tasks } = pageData;

  const navCounts = {
    entry: entries.length,
    task: tasks.filter((t) => t.status === "todo").length,
  };

  return (
    <AppShell
      userName={user.name}
      userSubtitle={user.email}
      navCounts={navCounts}
    >
      <div className="relative mx-auto max-w-[700px] px-5 py-6 md:px-8 md:py-7 lg:max-w-[1080px]">
        <header className="mb-4">
          <h1 className="font-serif text-2xl font-extrabold tracking-tight">
            タスク
            <span className="ml-2 align-middle text-[12px] font-bold text-ink-3">
              締切・予定
            </span>
          </h1>
          <p className="mt-0.5 text-[12px] text-ink-3">
            今日やることを上から片づけます。締切、面接予定、準備メモを応募先に紐づけます。
          </p>
        </header>

        <TaskListView initialTasks={tasks} entries={entries} />
      </div>
    </AppShell>
  );
}
