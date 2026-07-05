import Link from "next/link";
import {
  ArrowRight,
  CheckSquare,
  ClipboardList,
  Inbox,
  type LucideIcon,
} from "lucide-react";

export interface DashboardNextActionInput {
  inboxCount: number;
  entryCount: number;
  openTaskCount: number;
}

export interface DashboardNextActionModel {
  title: string;
  body: string;
  href: string;
  cta: string;
  activeStep: "inbox" | "entry" | "task";
  Icon: LucideIcon;
}

export function getDashboardNextAction({
  inboxCount,
  entryCount,
  openTaskCount,
}: DashboardNextActionInput): DashboardNextActionModel {
  if (inboxCount > 0) {
    return {
      title: `保存箱の求人 ${inboxCount}件を応募先にする`,
      body: "残す求人だけ応募先に変換します。不要なものは削除して、管理対象を増やしすぎないようにします。",
      href: "/inbox",
      cta: "保存箱を開く",
      activeStep: "inbox",
      Icon: Inbox,
    };
  }

  if (entryCount === 0) {
    return {
      title: "最初の応募先を登録する",
      body: "管理したい企業だけを追加します。企業名、応募経路、いまの選考フェーズが分かれば十分です。",
      href: "/entry/new",
      cta: "応募先を追加",
      activeStep: "entry",
      Icon: ClipboardList,
    };
  }

  if (openTaskCount === 0) {
    return {
      title: "締切・予定を追加する",
      body: "ES締切、面接日、準備タスクを応募先に紐づけて、今日やることに出します。",
      href: "/task",
      cta: "タスクを追加",
      activeStep: "task",
      Icon: CheckSquare,
    };
  }

  return {
    title: `未完了タスク ${openTaskCount}件を確認する`,
    body: "近い締切から順に並んでいます。まず一番上のタスクだけ片づけます。",
    href: "/task",
    cta: "タスクを見る",
    activeStep: "task",
    Icon: CheckSquare,
  };
}

export function DashboardNextAction(props: DashboardNextActionInput) {
  const action = getDashboardNextAction(props);
  const ActionIcon = action.Icon;

  return (
    <section className="mb-4 rounded-xl border border-sage/30 bg-sage-wash p-4 md:mb-5 md:p-5">
      <div className="grid gap-3 md:grid-cols-[1fr_auto] md:items-center">
        <div className="flex items-start gap-3">
          <div className="grid h-11 w-11 shrink-0 place-items-center rounded-lg bg-surface text-sage shadow-sm">
            <ActionIcon size={20} aria-hidden />
          </div>
          <div className="min-w-0">
            <p className="text-[12px] font-black text-sage">
              次にやること
            </p>
            <h2 className="mt-0.5 text-[17px] font-extrabold text-ink">
              {action.title}
            </h2>
            <p className="mt-1 max-w-[560px] text-[12px] leading-relaxed text-ink-2">
              {action.body}
            </p>
          </div>
        </div>

        <Link
          href={action.href}
          prefetch={false}
          className="inline-flex h-10 items-center justify-center gap-1.5 rounded-lg bg-sage px-4 text-[12px] font-bold text-white transition-transform hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40"
        >
          {action.cta}
          <ArrowRight size={13} aria-hidden />
        </Link>
      </div>
    </section>
  );
}
