"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  ClipboardList,
  Columns3,
  Home,
  Inbox,
  PencilLine,
  UserCircle,
  type LucideIcon,
} from "lucide-react";

/** サイドバーのバッジに表示する実カウント。Server Component で集計して props で渡す。 */
export interface NavCounts {
  entry: number;
  task: number;
  inbox: number;
}

interface NavItem {
  k: string;
  l: string;
  icon: LucideIcon;
  href: string;
  /** どの NavCounts キーをバッジに出すか。未指定ならバッジ無し。 */
  countKey?: keyof NavCounts;
}

const NAV_ITEMS: NavItem[] = [
  { k: "home", l: "ホーム", icon: Home, href: "/dashboard" },
  { k: "entry", l: "Entry", icon: ClipboardList, href: "/entry", countKey: "entry" },
  { k: "kanban", l: "カンバン", icon: Columns3, href: "/kanban" },
  { k: "task", l: "タスク", icon: PencilLine, href: "/task", countKey: "task" },
  { k: "inbox", l: "保存箱", icon: Inbox, href: "/inbox", countKey: "inbox" },
];

export function Sidebar({
  userName = "ゲスト",
  userSubtitle = "",
  navCounts,
}: {
  userName?: string;
  userSubtitle?: string;
  /** Server Component から渡される実カウント。未指定の画面ではバッジを出さない。 */
  navCounts?: Partial<NavCounts>;
}) {
  const pathname = usePathname();

  return (
    <aside className="hidden md:flex w-[204px] shrink-0 flex-col gap-4 border-r border-line bg-cream px-3.5 py-4">
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
          const Icon = it.icon;
          const active =
            (it.href === "/dashboard" && pathname === "/dashboard") ||
            (it.href !== "/dashboard" && pathname?.startsWith(it.href));
          const count =
            it.countKey && navCounts ? navCounts[it.countKey] : undefined;
          return (
            <Link
              key={it.k}
              href={it.href}
              prefetch={false}
              className={`flex items-center gap-2.5 rounded-lg px-2.5 py-2.5 text-xs font-semibold transition-all ${
                active
                  ? "bg-sage text-white"
                  : "text-ink hover:bg-line-2"
              }`}
              aria-current={active ? "page" : undefined}
            >
              <Icon size={16} className="shrink-0" aria-hidden />
              <span className="min-w-0 flex-1 truncate">{it.l}</span>
              {count !== undefined && (
                <span
                  data-testid={`nav-count-${it.k}`}
                  className={`rounded-md px-1.5 py-px text-[9px] font-bold ${
                    active ? "bg-white/20 text-white" : "bg-sage-soft text-sage"
                  }`}
                >
                  {count}
                </span>
              )}
            </Link>
          );
        })}
      </nav>

      <div className="flex-1" />

      {/* User card */}
      <Link
        href="/profile"
        prefetch={false}
        className="flex items-center gap-2 rounded-lg px-1.5 py-2 transition-colors hover:bg-line-2"
      >
        <div className="grid h-[30px] w-[30px] place-items-center rounded-full bg-sage-soft text-sage">
          <UserCircle size={18} aria-hidden />
        </div>
        <div className="min-w-0 flex-1">
          <div className="text-[11px] font-bold truncate">{userName}</div>
          {userSubtitle && <div className="text-[9px] text-ink-3 truncate">{userSubtitle}</div>}
        </div>
      </Link>
    </aside>
  );
}

/** モバイル用ボトムタブバー。 */
export function MobileTabBar() {
  const pathname = usePathname();
  const tabs = [
    { k: "home", l: "ホーム", icon: Home, href: "/dashboard" },
    { k: "entry", l: "Entry", icon: ClipboardList, href: "/entry" },
    { k: "kanban", l: "ボード", icon: Columns3, href: "/kanban" },
    { k: "task", l: "タスク", icon: PencilLine, href: "/task" },
    { k: "inbox", l: "保存", icon: Inbox, href: "/inbox" },
  ];

  return (
    <nav className="fixed inset-x-0 bottom-0 z-30 flex h-[68px] border-t border-line bg-white/92 pt-1.5 backdrop-blur-xl md:hidden">
      {tabs.map((t) => {
        const Icon = t.icon;
        const active =
          (t.href === "/dashboard" && pathname === "/dashboard") ||
          (t.href !== "/dashboard" && pathname?.startsWith(t.href));
        return (
          <Link
            key={t.k}
            href={t.href}
            prefetch={false}
            className={`flex flex-1 flex-col items-center gap-0.5 transition-transform ${
              active ? "-translate-y-0.5 text-sage" : "text-ink-3"
            }`}
            aria-current={active ? "page" : undefined}
          >
            <Icon size={19} aria-hidden />
            <span className="text-[10px] font-bold">{t.l}</span>
            {active && <span className="mt-px h-1 w-1 rounded-full bg-sage" />}
          </Link>
        );
      })}
    </nav>
  );
}
