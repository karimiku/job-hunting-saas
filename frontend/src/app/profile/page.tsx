// Server Component。/auth/me を SSR で取得し、表示。
// ログアウトボタンだけ Client (SignOutButton)。

import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";
import { SignOutButton } from "@/components/entre/SignOutButton";

export default async function ProfilePage() {
  const user = await getCurrentUserServer();
  if (!user) redirect("/login");

  return (
    <AppShell userName={user.name} userSubtitle={user.email}>
      <div className="mx-auto max-w-[800px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-center gap-4">
          <div className="grid h-16 w-16 place-items-center rounded-full bg-sage-soft">
            <Mascot size={48} mood="happy" />
          </div>
          <div className="flex-1">
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">{user.name}</h1>
            <p className="mt-0.5 text-[11px] text-ink-3">{user.email}</p>
            <p className="mt-0.5 text-[11px] text-ink-2">ログイン中のアカウント</p>
          </div>
        </header>

        <section className="mb-6 rounded-xl border border-line bg-surface p-5">
          <h2 className="mb-3 text-[13px] font-extrabold">アカウント情報</h2>
          <ul className="flex flex-col gap-1">
            <SettingsRow label="表示名" value={user.name} />
            <SettingsRow label="メール" value={user.email} />
            <SettingsRow label="ユーザーID" value={user.id} />
          </ul>
        </section>

        {/* Settings */}
        <section className="rounded-xl border border-line bg-surface p-5">
          <h2 className="mb-3 text-[13px] font-extrabold">設定</h2>
          <ul className="flex flex-col gap-1">
            <SettingsRow label="通知" value="ON" />
            <SettingsRow label="連携カレンダー" value="未設定" />
            <SettingsRow label="Chrome拡張" value="Inbox保存に対応" />
          </ul>
          <SignOutButton className="mt-4 w-full rounded-lg border border-line bg-surface py-2.5 text-[12px] font-bold text-ink-2 transition-colors hover:bg-line-2" />
        </section>
      </div>
    </AppShell>
  );
}

function SettingsRow({ label, value }: { label: string; value: string }) {
  return (
    <li className="flex items-center justify-between gap-3 border-b border-dashed border-line py-2 last:border-0">
      <span className="shrink-0 text-[12px] font-semibold">{label}</span>
      <span className="min-w-0 truncate text-right text-[11px] text-ink-3">{value}</span>
    </li>
  );
}
