// Server Component。entry + tasks を SSR で並列取得し、interactive 部分は Client に委譲。

import Link from "next/link";
import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import {
  getEntryServer,
  listTasksByEntryServer,
} from "@/lib/api/server-resources";
import { ApiError } from "@/lib/api/client-types";
import { AppShell } from "@/components/entre/AppShell";
import { EntryDetailView } from "@/components/entre/EntryDetailView";

interface Props {
  params: Promise<{ entryId: string }>;
}

export default async function EntryDetailPage({ params }: Props) {
  const { entryId } = await params;
  const user = await getCurrentUserServer();
  if (!user) redirect("/login");

  // entry + tasks は独立なので並列 fetch。
  // 取得失敗 (404 等) は EntryDetailView に渡す initial=null として扱い、UI 側でエラー表示。
  const [entry, tasks] = await Promise.all([
    getEntryServer(entryId).catch((e) => {
      if (e instanceof ApiError) return null;
      throw e;
    }),
    listTasksByEntryServer(entryId).catch(() => []),
  ]);

  return (
    <AppShell userName={user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[700px] px-5 py-6 md:px-8 md:py-7">
        <Link
          href="/entry"
          className="mb-3 inline-flex items-center gap-1 text-[11px] font-semibold text-ink-3 hover:text-sage"
        >
          ‹ Entry 一覧
        </Link>
        <EntryDetailView initialEntry={entry} initialTasks={tasks} />
      </div>
    </AppShell>
  );
}
