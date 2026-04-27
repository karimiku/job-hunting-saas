// Server Component。
// 認証 + クリップ一覧を SSR で取得し、子の Client Component には完成した DTO を props で渡す。
// ここでは useEffect/useState を使わない。

import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import { listInboxClipsServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";
import { InboxList } from "@/components/entre/InboxList";

export default async function InboxPage() {
  const user = await getCurrentUserServer();
  if (!user) redirect("/login");

  const clips = await listInboxClipsServer();

  return (
    <AppShell userName={user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[800px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-baseline justify-between">
          <div>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">Inbox</h1>
            <p className="mt-0.5 text-[11px] text-ink-3">
              Chrome拡張から保存したクリップ
            </p>
          </div>
          <Mascot size={32} mood="wink" />
        </header>

        <InboxList clips={clips} />
      </div>
    </AppShell>
  );
}
