"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  ClipboardList,
  Columns3,
  Home,
  Inbox,
  Map,
  PencilLine,
  Settings,
  UserCircle,
  type LucideIcon,
} from "lucide-react";
import { Mascot } from "./Mascot";

/** サイドバーのバッジに表示する実カウント。Server Component で集計して props で渡す。 */
export interface NavCounts {
  entry: number;
  task: number;
  inbox: number;
}

interface NavItem {
  k: string;
  l: string;
  helper: string;
  icon: LucideIcon;
  href: string;
  /** どの NavCounts キーをバッジに出すか。未指定ならバッジ無し。 */
  countKey?: keyof NavCounts;
  dev?: boolean;
}

const NAV_ITEMS: NavItem[] = [
  { k: "home", l: "ホーム", helper: "今日の状況", icon: Home, href: "/dashboard" },
  { k: "entry", l: "Entry", helper: "応募管理", icon: ClipboardList, href: "/entry", countKey: "entry" },
  { k: "kanban", l: "カンバン", helper: "選考フェーズ", icon: Columns3, href: "/kanban" },
  { k: "roadmap", l: "ロードマップ", helper: "進め方", icon: Map, href: "/roadmap" },
  { k: "task", l: "Task", helper: "締切と予定", icon: PencilLine, href: "/task", countKey: "task" },
  { k: "inbox", l: "Inbox", helper: "保存クリップ", icon: Inbox, href: "/inbox", countKey: "inbox" },
  { k: "profile", l: "プロフィール", helper: "アカウント", icon: UserCircle, href: "/profile" },
  { k: "es", l: "ESエディタ", helper: "準備中", icon: PencilLine, href: "/es", dev: true },
];

export function Sidebar({
  userName = "ゲスト",
  userSubtitle = "",
  navCounts,
}: {
  userName?: string;
  userSubtitle?: string;
  /** Server Component から渡される実カウント。未指定の画面ではバッジを出さない。 */
  navCounts?: NavCounts;
}) {
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
          const Icon = it.icon;
          const active =
            (it.href === "/dashboard" && pathname === "/dashboard") ||
            (it.href !== "/dashboard" && pathname?.startsWith(it.href));
          const count =
            it.countKey && navCounts ? navCounts[it.countKey] : undefined;
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
              <Icon size={16} className="shrink-0" aria-hidden />
              <span className="min-w-0 flex-1">
                <span className="block truncate">{it.l}</span>
                <span
                  className={`block truncate text-[9px] font-medium ${
                    active ? "text-white/70" : "text-ink-3"
                  }`}
                >
                  {it.helper}
                </span>
              </span>
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

      <div className="rounded-[10px] border border-line bg-surface px-3 py-2.5 text-[10px]">
        <div className="mb-1 text-[11px] font-extrabold">βチェック</div>
        <div className="leading-relaxed text-ink-3">
          {navCounts
            ? `Entry ${navCounts.entry}件 / 未完了Task ${navCounts.task}件 / Inbox ${navCounts.inbox}件`
            : "保存・Entry・Task の状態を各画面で確認できます。"}
        </div>
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
        <Settings size={14} className="text-ink-3" aria-hidden />
      </div>
    </aside>
  );
}

/** モバイル用ボトムタブバー。 */
export function MobileTabBar() {
  const pathname = usePathname();
  const tabs = [
    { k: "home", l: "ホーム", icon: Home, href: "/dashboard" },
    { k: "entry", l: "Entry", icon: ClipboardList, href: "/entry" },
    { k: "roadmap", l: "ロード", icon: Map, href: "/roadmap" },
    { k: "inbox", l: "Inbox", icon: Inbox, href: "/inbox" },
    { k: "profile", l: "Me", icon: UserCircle, href: "/profile" },
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
