// Server Component。応募先0件、またはタスク0件のあいだホームに出す「はじめかた」チェックリスト。
// 完了状態は props から判定するだけの純粋表示（useEffect は使わない）。

import Link from "next/link";
import { ArrowRight } from "lucide-react";

export interface GettingStartedGuideProps {
  hasEntries: boolean;
  hasTasks: boolean;
  hasDoneTasks: boolean;
  compact?: boolean;
}

export interface GettingStartedStep {
  id: "entry" | "board" | "task";
  index: number;
  title: string;
  description: string;
  href: string;
  cta: string;
  done: boolean;
}

export function buildGettingStartedSteps({
  hasEntries,
  hasTasks,
}: Pick<GettingStartedGuideProps, "hasEntries" | "hasTasks">): GettingStartedStep[] {
  return [
    {
      id: "entry",
      index: 1,
      title: "応募先を登録する",
      description: "会社名だけでOK・30秒で終わります。",
      href: "/entry/new",
      cta: "登録する",
      done: hasEntries,
    },
    {
      id: "board",
      index: 2,
      title: "選考が進んだらボードでカードを動かす",
      description: "応募先の進み具合を、ドラッグでそのまま反映できます。",
      href: "/kanban",
      cta: "ボードを見る",
      done: hasEntries,
    },
    {
      id: "task",
      index: 3,
      title: "ES締切・面接日をタスクにする",
      description: "締切と面接日を登録すると、ホームに近い順で表示されます。",
      href: "/task",
      cta: "タスクを追加する",
      done: hasTasks,
    },
  ];
}

/** はじめかたガイド。3ステップのチェックリストカード。compact は応募先ありタスク無しの継続表示用。 */
export function GettingStartedGuide({
  hasEntries,
  hasTasks,
  hasDoneTasks,
  compact = false,
}: GettingStartedGuideProps) {
  const steps = buildGettingStartedSteps({ hasEntries, hasTasks });
  const doneCount = steps.filter((step) => step.done).length;

  return (
    <section
      data-testid="getting-started-guide"
      className={`mb-4 rounded-xl border border-sage/30 bg-sage-wash md:mb-5 ${
        compact ? "p-4" : "p-5 md:p-6"
      }`}
    >
      <div className="mb-3 flex items-baseline justify-between gap-3">
        <div>
          <h2 className={`font-extrabold text-ink ${compact ? "text-[14px]" : "text-[17px]"}`}>
            はじめかた
          </h2>
          {!compact && (
            <p className="mt-1 text-[12px] leading-relaxed text-ink-2">
              まずはこの3つだけ。3分で就活の全体が見えるようになります。
            </p>
          )}
        </div>
        <span className="shrink-0 text-[12px] font-bold text-sage">
          {doneCount}/{steps.length} 完了{hasDoneTasks ? "・タスク実行中" : ""}
        </span>
      </div>

      <ul className="flex flex-col gap-2">
        {steps.map((step) => (
          <li
            key={step.id}
            data-testid={`getting-started-step-${step.id}`}
            className={`flex items-center gap-3 rounded-lg border border-line bg-surface px-3 ${
              compact ? "py-2" : "py-2.5"
            } ${step.done ? "opacity-60" : ""}`}
          >
            <span
              aria-hidden
              className={`grid h-6 w-6 shrink-0 place-items-center rounded-full text-[12px] font-bold ${
                step.done
                  ? "border-[1.5px] border-sage bg-sage text-white"
                  : "border-[1.5px] border-line text-ink-3"
              }`}
            >
              {step.done ? "✓" : step.index}
            </span>
            <div className="min-w-0 flex-1">
              <div className={`text-[12px] font-bold text-ink ${step.done ? "line-through" : ""}`}>
                {step.title}
              </div>
              {!compact && (
                <div className="mt-0.5 text-[12px] text-ink-3">{step.description}</div>
              )}
            </div>
            {!step.done && (
              <Link
                href={step.href}
                prefetch={false}
                className="inline-flex shrink-0 items-center gap-1 text-[12px] font-bold text-sage"
              >
                {step.cta}
                <ArrowRight size={12} aria-hidden />
              </Link>
            )}
          </li>
        ))}
      </ul>
    </section>
  );
}
