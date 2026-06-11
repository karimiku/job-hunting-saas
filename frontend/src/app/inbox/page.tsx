// Server Component。
// 認証 + クリップ一覧を SSR で取得し、子の Client Component には完成した DTO を props で渡す。
// ここでは useEffect/useState を使わない。

import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import { ApiError } from "@/lib/api/client-types";
import {
  buildNavCounts,
  listInboxClipsServer,
  listCompaniesServer,
  listEntriesServer,
  listTasksServer,
} from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { InboxList } from "@/components/entre/InboxList";

export default async function InboxPage() {
  // user とデータは独立なので並列取得 (auth を待ってから始めると backend RTT が1段増える)。
  // clips はメインリソースなので 401 (未ログイン → 下で redirect) 以外は throw して error.tsx に拾わせる。
  const [user, clips, companies, entries, tasks] = await Promise.all([
    getCurrentUserServer(),
    listInboxClipsServer().catch((e) => {
      if (e instanceof ApiError && e.unauthorized) return [];
      throw e;
    }),
    listCompaniesServer().catch(() => []),
    listEntriesServer().catch(() => []),
    listTasksServer().catch(() => []),
  ]);
  if (!user) redirect("/login");
  const navCounts = buildNavCounts(entries, tasks, clips);

  return (
    <AppShell userName={user.name} userSubtitle={user.email} navCounts={navCounts}>
      <div className="mx-auto max-w-[800px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-start justify-between gap-3">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">保存箱</h1>
            <p className="mt-0.5 text-[11px] text-ink-3">
              保存した求人の一時置き場です。残す求人だけEntryに変換します。
            </p>
          </div>
          <div className="flex items-center gap-2 rounded-lg border border-line bg-surface px-2.5 py-2">
            <span className="font-mono text-[10px] font-bold text-sage">{clips.length}</span>
            <span className="text-[10px] font-semibold text-ink-3">保存中</span>
          </div>
        </header>

        <InboxList clips={clips} companies={companies} />
      </div>
    </AppShell>
  );
}
