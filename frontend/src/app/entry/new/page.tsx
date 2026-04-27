// Server Component。auth gate のみ担当。フォーム本体は Client (NewEntryForm) で
// useActionState を介して Server Action (./actions.ts) を呼ぶ。

import { redirect } from "next/navigation";
import Link from "next/link";
import { getCurrentUserServer } from "@/lib/auth-server";
import { AppShell } from "@/components/entre/AppShell";
import { NewEntryForm } from "./NewEntryForm";

export default async function NewEntryPage() {
  const user = await getCurrentUserServer();
  if (!user) redirect("/login");

  return (
    <AppShell userName={user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[640px] px-5 py-6 md:px-8 md:py-7">
        <Link
          href="/entry"
          className="mb-3 inline-flex items-center gap-1 text-[11px] font-semibold text-ink-3 hover:text-sage"
        >
          ‹ Entry 一覧
        </Link>

        <NewEntryForm />
      </div>
    </AppShell>
  );
}
