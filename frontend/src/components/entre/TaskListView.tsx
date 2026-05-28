"use client";

// Server から initialTasks を受け取る Client Component。
// データ取得は Server Component 側で完結し、ここは表示とトグル操作 (UI interactivity) のみ担う。
// チェックボックスのトグルは setTaskStatusAction (Server Action) を呼び、完了時に紙吹雪を出して
// router.refresh() で SSR を再評価する (楽観更新は startTransition 中の見た目だけに留める)。

import { useState, useTransition } from "react";
import { useRouter } from "next/navigation";
import { setTaskStatusAction } from "@/app/task/actions";
import type { TaskWithContext } from "@/lib/api/server-resources";
import { Confetti } from "./Confetti";

interface Props {
  initialTasks: TaskWithContext[];
}

// type ごとのバッジ色。deadline = 締切で目立つ色、schedule = 予定で落ち着いた色。
const TYPE_BADGE: Record<string, string> = {
  deadline: "bg-pink",
  schedule: "bg-sky",
};

function formatDue(dueDate: string | null): string {
  if (!dueDate) return "期日なし";
  // ISO 文字列 (YYYY-MM-DD...) を M/D に短縮表示。パースできなければ原文を返す。
  const d = new Date(dueDate);
  if (Number.isNaN(d.getTime())) return dueDate;
  return `${d.getMonth() + 1}/${d.getDate()}`;
}

export function TaskListView({ initialTasks }: Props) {
  const router = useRouter();
  const [confetti, setConfetti] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [isPending, startTransition] = useTransition();
  // 楽観更新用に「いまトグル中の taskId → 目標 status」を保持する。
  const [optimistic, setOptimistic] = useState<Record<string, "todo" | "done">>(
    {},
  );

  const toggle = (task: TaskWithContext) => {
    const next = task.status === "done" ? "todo" : "done";
    setError(null);
    setOptimistic((prev) => ({ ...prev, [task.id]: next }));
    if (next === "done") setConfetti((n) => n + 1);

    startTransition(async () => {
      const result = await setTaskStatusAction(task.id, next);
      if (!result.ok) {
        // 失敗時は楽観更新を巻き戻す。
        setOptimistic((prev) => {
          const next = { ...prev };
          delete next[task.id];
          return next;
        });
        setError(result.error ?? "タスクの更新に失敗しました");
        return;
      }
      router.refresh(); // SSR を再評価して最新の tasks を取得
    });
  };

  if (initialTasks.length === 0) {
    return (
      <div className="rounded-xl border border-line bg-surface px-4 py-8 text-center text-[12px] text-ink-3">
        まだタスクがありません。Entry の詳細からタスクを追加できます。
      </div>
    );
  }

  return (
    <div className="relative">
      {error && (
        <p
          role="alert"
          className="mb-2 rounded-lg bg-pink/40 px-3 py-2 text-[11px] font-semibold text-ink"
        >
          {error}
        </p>
      )}

      <ul className="flex flex-col gap-2">
        {initialTasks.map((task) => {
          const status = optimistic[task.id] ?? task.status;
          const done = status === "done";
          return (
            <li
              key={task.id}
              className={`flex items-center gap-3 rounded-xl border border-line bg-surface px-3 py-2.5 transition-opacity ${
                done ? "opacity-50" : ""
              }`}
            >
              <button
                type="button"
                onClick={() => toggle(task)}
                disabled={isPending}
                aria-pressed={done}
                aria-label={done ? "タスク未完了に戻す" : "タスク完了にする"}
                className={`grid h-5 w-5 shrink-0 place-items-center rounded-full text-[11px] text-white transition-colors disabled:opacity-60 ${
                  done
                    ? "border-[1.5px] border-sage bg-sage"
                    : "border-[1.5px] border-line bg-transparent"
                }`}
              >
                {done ? "✓" : ""}
              </button>
              <div className="min-w-0 flex-1">
                <div
                  className={`text-[12px] font-semibold ${done ? "line-through" : ""}`}
                >
                  {task.title}
                </div>
                <div className="mt-0.5 text-[10px] text-ink-3">
                  {task.companyName ?? "（会社名未設定）"}
                  {task.memo ? ` · ${task.memo}` : ""}
                </div>
              </div>
              <span
                className={`shrink-0 rounded-md px-2 py-0.5 font-mono text-[10px] font-bold text-white ${
                  TYPE_BADGE[task.type] ?? "bg-sage"
                }`}
              >
                {formatDue(task.dueDate)}
              </span>
            </li>
          );
        })}
      </ul>

      <Confetti trigger={confetti} count={22} />
    </div>
  );
}
