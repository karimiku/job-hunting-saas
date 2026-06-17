// Server Component。タスク1件を取得し、詳細表示と操作を Client Component に委譲する。

import { notFound, redirect } from "next/navigation";
import { ApiError } from "@/lib/api/client-types";
import {
  getAppPageDataServer,
  getTaskServer,
  type TaskWithEntry,
} from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { TaskDetailView } from "@/components/entre/TaskDetailView";

interface Props {
  params: Promise<{ taskId: string }>;
}

export default async function TaskDetailPage({ params }: Props) {
  const { taskId } = await params;
  const [pageData, task] = await Promise.all([
    getAppPageDataServer(),
    getTaskServer(taskId).catch((e) => {
      if (e instanceof ApiError && (e.notFound || e.unauthorized)) return null;
      throw e;
    }),
  ]);

  if (!pageData) redirect("/login");
  if (!task) notFound();

  const entry = pageData.entries.find((item) => item.id === task.entryId) ?? null;
  const fallbackTask = pageData.tasks.find((item) => item.id === task.id);
  const taskWithEntry: TaskWithEntry = {
    ...task,
    companyName: fallbackTask?.companyName ?? entry?.companyName,
  };

  return (
    <AppShell
      userName={pageData.user.name}
      userSubtitle={pageData.user.email}
      navCounts={pageData.navCounts}
    >
      <div className="mx-auto max-w-[700px] px-5 py-6 md:px-8 md:py-7">
        <TaskDetailView task={taskWithEntry} entry={entry} />
      </div>
    </AppShell>
  );
}
