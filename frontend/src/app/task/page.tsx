// Server Component。ユーザーの全タスクを SSR で集約取得し、interactive 部分は Client に委譲する。
// backend に全タスク一覧 API が無いため、entries → entry ごとの tasks を server-resources で集約する。

import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import {
  listAllTasksServer,
  listEntriesWithCompanyNamesServer,
  listInboxClipsServer,
} from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { TaskListView } from "@/components/entre/TaskListView";

export default async function TaskPage() {
  const user = await getCurrentUserServer();
  if (!user) redirect("/login");

  // backend に全タスク一覧 API が無いため、entries (会社名 join 済み) を引いてから
  // entry ごとの tasks を集約する。取得失敗時は空一覧で表示 (UI 側で空状態を出す)。
  const entries = await listEntriesWithCompanyNamesServer().catch(() => []);
  const [tasks, clips] = await Promise.all([
    listAllTasksServer(entries).catch(() => []),
    listInboxClipsServer().catch(() => []),
  ]);

  // 既に取得済みの entries / tasks を再利用してサイドバーのバッジ数を組み立てる
  // (getNavCountsServer を呼ぶと entries/tasks を二重取得するため、ここでは手元の値から作る)。
  const navCounts = {
    entry: entries.length,
    task: tasks.filter((t) => t.status === "todo").length,
    inbox: clips.length,
  };

  return (
    <AppShell
      userName={user.name}
      userSubtitle="○○大学 4年"
      navCounts={navCounts}
    >
      <div className="relative mx-auto max-w-[700px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4">
          <h1 className="font-serif text-2xl font-extrabold tracking-tight">Task</h1>
          <p className="mt-0.5 text-[11px] text-ink-3">タスクや締切を1箇所で管理</p>
        </header>

        <TaskListView initialTasks={tasks} />
      </div>
    </AppShell>
  );
}
