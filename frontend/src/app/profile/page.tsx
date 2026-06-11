// Server Component。/auth/me を SSR で取得し、表示。
// ログアウトボタンだけ Client (SignOutButton)。

import { redirect } from "next/navigation";
import { UserCircle } from "lucide-react";
import { getCurrentUserServer } from "@/lib/auth-server";
import { serverFetch } from "@/lib/api/server";
import { ApiError } from "@/lib/api/client-types";
import type { AiAccessTokenResponse } from "@/lib/api/aiTokens";
import { AppShell } from "@/components/entre/AppShell";
import { SignOutButton } from "@/components/entre/SignOutButton";
import { AiAccessTokenPanel } from "@/components/entre/AiAccessTokenPanel";

export default async function ProfilePage() {
  const [user, tokenResult] = await Promise.all([
    getCurrentUserServer(),
    serverFetch<{ tokens: AiAccessTokenResponse[] }>("/api/v1/ai/tokens").catch(
      (e) => {
        if (e instanceof ApiError && e.unauthorized) {
          return { tokens: [] as AiAccessTokenResponse[] };
        }
        throw e;
      },
    ),
  ]);
  if (!user) redirect("/login");

  return (
    <AppShell userName={user.name} userSubtitle={user.email}>
      <div className="mx-auto max-w-[800px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-5 flex items-center gap-3">
          <div className="grid h-12 w-12 shrink-0 place-items-center rounded-full border border-line bg-surface text-sage">
            <UserCircle size={28} strokeWidth={1.8} aria-hidden="true" />
          </div>
          <div className="flex-1">
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">アカウント</h1>
            <p className="mt-1 text-[12px] text-ink-2">{user.email}</p>
          </div>
        </header>

        <section className="rounded-xl border border-line bg-surface p-5">
          <h2 className="mb-3 text-[13px] font-extrabold">アカウント情報</h2>
          <ul className="flex flex-col gap-1">
            <SettingsRow label="表示名" value={user.name} />
            <SettingsRow label="メール" value={user.email} />
            <SettingsRow label="ユーザーID" value={user.id} />
          </ul>
          <SignOutButton className="mt-5 w-full rounded-lg border border-line bg-surface py-2.5 text-[12px] font-bold text-ink-2 transition-colors hover:bg-line-2" />
        </section>

        <div className="mt-4">
          <AiAccessTokenPanel tokens={tokenResult.tokens} />
        </div>
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
