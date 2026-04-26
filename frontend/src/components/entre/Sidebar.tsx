"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Mascot } from "./Mascot";

interface NavItem {
  k: string;
  l: string;
  i: string;
  href: string;
  count?: number;
  dev?: boolean;
}

const NAV_ITEMS: NavItem[] = [
  { k: "home", l: "ホーム", i: "⌂", href: "/dashboard" },
  { k: "entry", l: "Entry", i: "▤", href: "/entry", count: 24 },
  { k: "kanban", l: "カンバン", i: "⊞", href: "/kanban" },
  { k: "roadmap", l: "就活ロードマップ", i: "⤴", href: "/roadmap" },
  { k: "task", l: "Task", i: "✓", href: "/task", count: 9 },
  { k: "inbox", l: "Inbox", i: "✉", href: "/inbox", count: 2 },
  { k: "profile", l: "プロフィール", i: "◉", href: "/profile" },
  { k: "es", l: "ESエディタ", i: "✎", href: "/es", dev: true },
];

export function Sidebar({ userName = "ゲスト", userSubtitle = "" }: { userName?: string; userSubtitle?: string }) {
  const pathname = usePathname();

  return (
    <aside className="hidden md:flex w-[220px] shrink-0 flex-col gap-3.5 border-r border-line bg-cream px-3.5 py-4">
      {/* Logo */}
      <div className="flex items-baseline gap-2 px-1.5">
        <span className="font-serif text-[22px] font-black italic tracking-tight">Entré</span>
        <span className="rounded-sm bg-sage-soft px-1.5 py-0.5 font-mono text-[8px] font-bold tracking-widest text-sage">
          BETA
        </span>
      </div>

      {/* Nav items */}
      <nav className="flex flex-col gap-0.5">
        {NAV_ITEMS.map((it) => {
          const active =
            (it.href === "/dashboard" && pathname === "/dashboard") ||
            (it.href !== "/dashboard" && pathname?.startsWith(it.href));
          return (
            <Link
              key={it.k}
              href={it.dev ? "#" : it.href}
              className={`flex items-center gap-2.5 rounded-lg px-2.5 py-2 text-xs font-semibold transition-all ${
                active
                  ? "bg-sage text-white"
                  : "text-ink hover:bg-line-2"
              } ${it.dev ? "opacity-55" : ""}`}
              aria-current={active ? "page" : undefined}
            >
              <span className="w-[18px] text-center text-sm">{it.i}</span>
              <span className="flex-1">{it.l}</span>
              {it.count !== undefined && (
                <span
                  className={`rounded-md px-1.5 py-px text-[9px] font-bold ${
                    active ? "bg-white/20 text-white" : "bg-sage-soft text-sage"
                  }`}
                >
                  {it.count}
                </span>
              )}
              {it.dev && (
                <span className="rounded-sm bg-cream-2 px-1.5 py-px text-[8px] font-bold text-amber-700">
                  開発中
                </span>
              )}
            </Link>
          );
        })}
      </nav>

      <div className="flex-1" />

      {/* Streak card */}
      <div className="rounded-[10px] border border-line bg-gradient-to-br from-cream-2 to-sage-wash px-3 py-2.5 text-[10px]">
        <div className="mb-1 flex items-center gap-1.5">
          <span className="text-sm" style={{ animation: "entre-wiggle 2s infinite" }}>🔥</span>
          <span className="text-[11px] font-extrabold">連続 7日</span>
        </div>
        <div className="text-[9px] text-ink-3">今週もがんばってますね！</div>
      </div>

      {/* User card */}
      <div className="flex items-center gap-2 px-1.5 py-2">
        <div className="grid h-[30px] w-[30px] place-items-center rounded-full bg-sage-soft">
          <Mascot size={26} />
        </div>
        <div className="min-w-0 flex-1">
          <div className="text-[11px] font-bold truncate">{userName}</div>
          {userSubtitle && <div className="text-[9px] text-ink-3 truncate">{userSubtitle}</div>}
        </div>
        <span className="text-[11px] text-ink-3">⚙</span>
      </div>
    </aside>
  );
}

/** モバイル用ボトムタブバー。 */
export function MobileTabBar() {
  const pathname = usePathname();
  const tabs = [
    { k: "home", l: "ホーム", i: "⌂", href: "/dashboard" },
    { k: "entry", l: "Entry", i: "▤", href: "/entry" },
    { k: "roadmap", l: "ロード", i: "⤴", href: "/roadmap" },
    { k: "inbox", l: "Inbox", i: "✉", href: "/inbox" },
    { k: "profile", l: "Me", i: "◉", href: "/profile" },
  ];

  return (
    <nav className="fixed inset-x-0 bottom-0 z-30 flex h-[68px] border-t border-line bg-white/92 pt-1.5 backdrop-blur-xl md:hidden">
      {tabs.map((t) => {
        const active =
          (t.href === "/dashboard" && pathname === "/dashboard") ||
          (t.href !== "/dashboard" && pathname?.startsWith(t.href));
        return (
          <Link
            key={t.k}
            href={t.href}
            className={`flex flex-1 flex-col items-center gap-0.5 transition-transform ${
              active ? "-translate-y-0.5 text-sage" : "text-ink-3"
            }`}
            aria-current={active ? "page" : undefined}
          >
            <span className="text-lg">{t.i}</span>
            <span className="text-[10px] font-bold">{t.l}</span>
            {active && <span className="mt-px h-1 w-1 rounded-full bg-sage" />}
          </Link>
        );
      })}
    </nav>
  );
}
