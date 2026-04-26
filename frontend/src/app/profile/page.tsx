"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { signOut } from "@/lib/auth";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";

const BADGES = [
  { emoji: "🌱", label: "はじめの一歩", earned: true },
  { emoji: "🔥", label: "7日連続", earned: true },
  { emoji: "📝", label: "10社エントリー", earned: true },
  { emoji: "🤝", label: "面接デビュー", earned: true },
  { emoji: "🎉", label: "初内定", earned: false },
];

export default function ProfilePage() {
  const router = useRouter();
  const state = useUser();

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  const user = state.user;

  return (
    <AppShell userName={user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[800px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4 flex items-center gap-4">
          <div className="grid h-16 w-16 place-items-center rounded-full bg-sage-soft">
            <Mascot size={48} mood="happy" />
          </div>
          <div className="flex-1">
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">{user.name}</h1>
            <p className="mt-0.5 text-[11px] text-ink-3">{user.email}</p>
            <p className="mt-0.5 text-[11px] text-ink-2">○○大学 経済学部 4年</p>
          </div>
        </header>

        {/* Badges */}
        <section className="mb-6 rounded-xl border border-line bg-surface p-5">
          <h2 className="mb-3 text-[13px] font-extrabold">🏆 獲得バッジ</h2>
          <ul className="grid grid-cols-3 gap-2 md:grid-cols-5">
            {BADGES.map((b) => (
              <li
                key={b.label}
                className={`flex flex-col items-center gap-1 rounded-lg border p-3 text-center ${
                  b.earned
                    ? "border-line bg-cream"
                    : "border-line bg-line-2 opacity-50"
                }`}
              >
                <span className="text-2xl">{b.emoji}</span>
                <span className="text-[10px] font-bold">{b.label}</span>
              </li>
            ))}
          </ul>
        </section>

        {/* Settings */}
        <section className="rounded-xl border border-line bg-surface p-5">
          <h2 className="mb-3 text-[13px] font-extrabold">設定</h2>
          <ul className="flex flex-col gap-1">
            <SettingsRow label="通知" value="ON" />
            <SettingsRow label="連携カレンダー" value="未設定" />
            <SettingsRow label="Chrome拡張" value="インストール済" />
          </ul>
          <button
            type="button"
            onClick={async () => {
              await signOut();
              router.push("/login");
            }}
            className="mt-4 w-full rounded-lg border border-line bg-surface py-2.5 text-[12px] font-bold text-ink-2 transition-colors hover:bg-line-2"
          >
            ログアウト
          </button>
        </section>
      </div>
    </AppShell>
  );
}

function SettingsRow({ label, value }: { label: string; value: string }) {
  return (
    <li className="flex items-center justify-between border-b border-dashed border-line py-2 last:border-0">
      <span className="text-[12px] font-semibold">{label}</span>
      <span className="text-[11px] text-ink-3">{value}</span>
    </li>
  );
}
