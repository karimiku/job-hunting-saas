// Server Component。entry + tasks を SSR で並列取得し、interactive 部分は Client に委譲。

import Link from "next/link";
import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import {
  getEntryServer,
  getNavCountsServer,
  listCompaniesServer,
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

  // user / entry / tasks / companies は独立なので並列 fetch。
  // entry の取得失敗 (404 等) は EntryDetailView に渡す initial=null として扱い、UI 側でエラー表示。
  // 会社名は entry 取得後に単品 GET すると直列の RTT が1段増えるため、
  // 一覧を並列で引いて companyId で突き合わせる。
  const [user, entryRaw, tasks, companies, navCounts] = await Promise.all([
    getCurrentUserServer(),
    getEntryServer(entryId).catch((e) => {
      if (e instanceof ApiError) return null;
      throw e;
    }),
    listTasksByEntryServer(entryId).catch(() => []),
    listCompaniesServer().catch(() => []),
    getNavCountsServer(),
  ]);
  if (!user) redirect("/login");

  // entry が取れたら会社名を join する（見つからなければ UI 側でフォールバック表示）。
  const entry = entryRaw
    ? {
        ...entryRaw,
        companyName: companies.find((c) => c.id === entryRaw.companyId)?.name,
      }
    : null;

  return (
    <AppShell userName={user.name} userSubtitle={user.email} navCounts={navCounts}>
      <div className="mx-auto max-w-[700px] px-5 py-6 md:px-8 md:py-7">
        <Link
          href="/entry"
          prefetch={false}
          className="mb-3 inline-flex items-center gap-1 text-[11px] font-semibold text-ink-3 hover:text-sage"
        >
          ‹ Entry 一覧
        </Link>
        <EntryDetailView initialEntry={entry} initialTasks={tasks} />
      </div>
    </AppShell>
  );
}
