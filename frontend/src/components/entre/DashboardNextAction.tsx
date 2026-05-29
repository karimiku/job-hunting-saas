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
      title: `保存クリップ ${inboxCount}件を応募先にする`,
      body: "会社名を確認して Entry にすると、選考ボードと締切管理に並びます。",
      href: "/inbox",
      cta: "Inboxで確認",
      activeStep: "inbox",
      Icon: Inbox,
    };
  }

  if (entryCount === 0) {
    return {
      title: "最初の応募先を追加する",
      body: "気になる企業を Entry に入れると、選考状況と予定を追えるようになります。",
      href: "/entry/new",
      cta: "Entryを追加",
      activeStep: "entry",
      Icon: ClipboardList,
    };
  }

  if (openTaskCount === 0) {
    return {
      title: "締切・予定を追加する",
      body: "ES締切、面接日、準備タスクを Entry に紐づけておくと見落としを防げます。",
      href: "/task",
      cta: "Taskを追加",
      activeStep: "task",
      Icon: CheckSquare,
    };
  }

  return {
    title: `未完了Task ${openTaskCount}件を確認する`,
    body: "近い締切から順に片づけると、今日の優先順位がはっきりします。",
    href: "/task",
    cta: "Taskを見る",
    activeStep: "task",
    Icon: CheckSquare,
  };
}

const STEPS: Array<{
  key: DashboardNextActionModel["activeStep"];
  label: string;
  caption: string;
  Icon: LucideIcon;
}> = [
  { key: "inbox", label: "保存", caption: "求人ページ", Icon: Inbox },
  { key: "entry", label: "応募先化", caption: "会社と選考", Icon: ClipboardList },
  { key: "task", label: "締切管理", caption: "予定と準備", Icon: CheckSquare },
];

export function DashboardNextAction(props: DashboardNextActionInput) {
  const action = getDashboardNextAction(props);
  const ActionIcon = action.Icon;

  return (
    <section className="mb-5 rounded-xl border border-sage/30 bg-sage-wash p-3.5 md:mb-6 md:p-4">
      <div className="grid gap-3 md:grid-cols-[1fr_auto] md:items-center">
        <div className="flex items-start gap-3">
          <div className="grid h-10 w-10 shrink-0 place-items-center rounded-lg bg-surface text-sage shadow-sm">
            <ActionIcon size={19} aria-hidden />
          </div>
          <div className="min-w-0">
            <p className="text-[10px] font-black text-sage">
              次にやること
            </p>
            <h2 className="mt-0.5 text-[15px] font-extrabold text-ink">
              {action.title}
            </h2>
            <p className="mt-1 max-w-[560px] text-[11px] leading-relaxed text-ink-2">
              {action.body}
            </p>
          </div>
        </div>

        <Link
          href={action.href}
          className="inline-flex h-9 items-center justify-center gap-1.5 rounded-lg bg-sage px-3.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40"
        >
          {action.cta}
          <ArrowRight size={13} aria-hidden />
        </Link>
      </div>

      <ol className="mt-3 grid gap-1.5 md:grid-cols-3" aria-label="応募管理の流れ">
        {STEPS.map((step, index) => {
          const StepIcon = step.Icon;
          const active = step.key === action.activeStep;
          return (
            <li
              key={step.key}
              className={`flex items-center gap-2 rounded-lg border px-2.5 py-2 ${
                active
                  ? "border-sage bg-surface text-ink"
                  : "border-line bg-cream/70 text-ink-3"
              }`}
            >
              <span
                className={`grid h-6 w-6 shrink-0 place-items-center rounded-md ${
                  active ? "bg-sage text-white" : "bg-surface text-ink-3"
                }`}
              >
                <StepIcon size={13} aria-hidden />
              </span>
              <span className="min-w-0">
                <span className="block text-[11px] font-extrabold">
                  {index + 1}. {step.label}
                </span>
                <span className="block truncate text-[9px] font-semibold">
                  {step.caption}
                </span>
              </span>
            </li>
          );
        })}
      </ol>
    </section>
  );
}
