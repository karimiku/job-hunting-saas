import Link from "next/link";
import { Columns3, List } from "lucide-react";

/** /entry（一覧）と /kanban（ボード）が同じ応募先データの別ビューであることを示す切替タブ。 */
export function EntryViewSwitch({ active }: { active: "list" | "board" }) {
  const tabs = [
    { key: "list" as const, label: "一覧", href: "/entry", Icon: List },
    { key: "board" as const, label: "ボード", href: "/kanban", Icon: Columns3 },
  ];
  return (
    <div
      role="tablist"
      aria-label="応募先の表示切替"
      className="inline-flex rounded-lg border border-line bg-surface p-0.5"
    >
      {tabs.map(({ key, label, href, Icon }) => {
        const selected = key === active;
        return (
          <Link
            key={key}
            href={href}
            prefetch={false}
            role="tab"
            aria-selected={selected}
            className={`inline-flex h-8 items-center gap-1.5 rounded-md px-3 text-[12px] font-bold transition-colors ${
              selected ? "bg-sage text-white" : "text-ink-3 hover:text-sage"
            }`}
          >
            <Icon size={14} aria-hidden />
            {label}
          </Link>
        );
      })}
    </div>
  );
}
