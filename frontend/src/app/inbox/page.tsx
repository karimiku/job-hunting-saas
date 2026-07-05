// Server Component。
// 認証 + クリップ一覧を SSR で取得し、子の Client Component には完成した DTO を props で渡す。
// ここでは useEffect/useState を使わない。

import { redirect } from "next/navigation";
import { getAppPageDataServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { InboxList } from "@/components/entre/InboxList";

export default async function InboxPage() {
  const pageData = await getAppPageDataServer();
  if (!pageData) redirect("/login");
  const { user, clips, companies, navCounts } = pageData;

  return (
    <AppShell userName={user.name} userSubtitle={user.email} navCounts={navCounts}>
      <div className="mx-auto max-w-[800px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-start justify-between gap-3">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">保存箱</h1>
            <p className="mt-0.5 text-[12px] text-ink-3">
              保存した求人の一時置き場です。残したい求人だけ応募先にします。いらないものは削除。
            </p>
          </div>
          <div className="flex items-center gap-2 rounded-lg border border-line bg-surface px-2.5 py-2">
            <span className="font-mono text-[12px] font-bold text-sage">{clips.length}</span>
            <span className="text-[12px] font-semibold text-ink-3">保存中</span>
          </div>
        </header>

        <InboxList clips={clips} companies={companies} />
      </div>
    </AppShell>
  );
}
