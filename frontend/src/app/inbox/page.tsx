// Server Component。
// 認証 + クリップ一覧を SSR で取得し、子の Client Component には完成した DTO を props で渡す。
// ここでは useEffect/useState を使わない。

import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import {
  listInboxClipsServer,
  getNavCountsServer,
  listCompaniesServer,
} from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";
import { InboxList } from "@/components/entre/InboxList";

export default async function InboxPage() {
  const user = await getCurrentUserServer();
  if (!user) redirect("/login");

  const [clips, navCounts, companies] = await Promise.all([
    listInboxClipsServer(),
    getNavCountsServer(),
    listCompaniesServer().catch(() => []),
  ]);

  return (
    <AppShell userName={user.name} userSubtitle={user.email} navCounts={navCounts}>
      <div className="mx-auto max-w-[800px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-start justify-between gap-3">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">Inbox</h1>
            <p className="mt-0.5 text-[11px] text-ink-3">
              Chrome拡張で保存した求人ページを確認し、管理対象の Entry に変換します。
            </p>
          </div>
          <div className="flex items-center gap-2 rounded-lg border border-line bg-surface px-2.5 py-2">
            <span className="font-mono text-[10px] font-bold text-sage">{clips.length}</span>
            <span className="text-[10px] font-semibold text-ink-3">保存中</span>
            <Mascot size={28} mood="wink" />
          </div>
        </header>

        <InboxList clips={clips} companies={companies} />
      </div>
    </AppShell>
  );
}
