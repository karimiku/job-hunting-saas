"use client";

// initialTasks を SSR で受け取り、表示とトグル操作のみ担う Client Component。

import { useActionState, useState, useTransition } from "react";
import { useFormStatus } from "react-dom";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { CalendarPlus, CheckCircle2, ClipboardList, Plus, Trash2 } from "lucide-react";
import {
  createTaskFromTaskPageAction,
  deleteTaskAction,
  setTaskStatusAction,
  type CreateTaskFormState,
} from "@/app/task/actions";
import type { TaskWithEntry } from "@/lib/api/server-resources";
import {
  companyDisplayName,
  type EntryResponse,
} from "@/lib/api/entries";
import { Confetti } from "./Confetti";

interface Props {
  initialTasks: TaskWithEntry[];
  entries: EntryResponse[];
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

function taskSortValue(task: TaskWithEntry): number {
  if (!task.dueDate) return Number.POSITIVE_INFINITY;
  const d = new Date(task.dueDate);
  return Number.isNaN(d.getTime()) ? Number.POSITIVE_INFINITY : d.getTime();
}

export function sortTasksForDisplay(tasks: TaskWithEntry[]): TaskWithEntry[] {
  return [...tasks].sort((a, b) => {
    const aDone = a.status === "done" ? 1 : 0;
    const bDone = b.status === "done" ? 1 : 0;
    if (aDone !== bDone) return aDone - bDone;
    return taskSortValue(a) - taskSortValue(b);
  });
}

export function TaskListView({ initialTasks, entries }: Props) {
  const router = useRouter();
  const [confetti, setConfetti] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [isPending, startTransition] = useTransition();
  const [deletingIds, setDeletingIds] = useState<Record<string, boolean>>({});
  // 楽観更新用に「いまトグル中の taskId → 目標 status」を保持する。
  const [optimistic, setOptimistic] = useState<Record<string, "todo" | "done">>(
    {},
  );

  const toggle = (task: TaskWithEntry) => {
    const next = task.status === "done" ? "todo" : "done";
    setError(null);
    setOptimistic((prev) => ({ ...prev, [task.id]: next }));

    startTransition(async () => {
      const result = await setTaskStatusAction(task.id, next);
      if (!result.ok) {
        setOptimistic((prev) => {
          const next = { ...prev };
          delete next[task.id];
          return next;
        });
        setError(result.error ?? "タスクの更新に失敗しました");
        return;
      }
      // 成功後にだけ祝福する (失敗→ロールバック時に祝福が出るのを防ぐ)。
      if (next === "done") setConfetti((n) => n + 1);
      router.refresh();
    });
  };

  const deleteTask = (task: TaskWithEntry) => {
    if (!window.confirm(`「${task.title}」を削除しますか？`)) return;
    setError(null);
    setDeletingIds((prev) => ({ ...prev, [task.id]: true }));
    startTransition(async () => {
      const result = await deleteTaskAction(task.id);
      if (!result.ok) {
        setDeletingIds((prev) => {
          const next = { ...prev };
          delete next[task.id];
          return next;
        });
        setError(result.error ?? "タスクの削除に失敗しました");
        return;
      }
      router.refresh();
    });
  };

  const tasks = sortTasksForDisplay(initialTasks).filter(
    (task) => !deletingIds[task.id],
  );

  return (
    <div className="relative">
      <TaskCreatePanel entries={entries} />

      {error && (
        <p
          role="alert"
          className="mb-2 rounded-lg bg-pink/40 px-3 py-2 text-[11px] font-semibold text-ink"
        >
          {error}
        </p>
      )}

      {tasks.length === 0 ? (
        <TaskEmptyState hasEntries={entries.length > 0} />
      ) : (
        <div className="space-y-2">
          <div className="flex items-center justify-between rounded-xl border border-line bg-cream px-3 py-2">
            <div>
              <p className="text-[11px] font-extrabold">未完了を上から片づける</p>
              <p className="mt-0.5 text-[10px] text-ink-3">
                左の丸を押すと完了。期日が近い順に並びます。
              </p>
            </div>
            <span className="rounded-md bg-sage-soft px-2 py-1 text-[10px] font-bold text-sage">
              {tasks.filter((task) => (optimistic[task.id] ?? task.status) === "todo").length}件残り
            </span>
          </div>

          <ul className="flex flex-col gap-2">
            {tasks.map((task) => {
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
                    className={`grid h-6 w-6 shrink-0 place-items-center rounded-full text-[11px] text-white transition-colors focus:outline-none focus:ring-2 focus:ring-sage/30 disabled:opacity-60 ${
                      done
                        ? "border-[1.5px] border-sage bg-sage"
                        : "border-[1.5px] border-line bg-transparent"
                    }`}
                  >
                    {done ? <CheckCircle2 size={15} aria-hidden /> : ""}
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
                  <button
                    type="button"
                    onClick={() => deleteTask(task)}
                    disabled={isPending}
                    aria-label={`タスク「${task.title}」を削除`}
                    className="grid h-7 w-7 shrink-0 place-items-center rounded-md border border-line text-ink-3 transition-colors hover:border-pink-deep hover:text-pink-deep focus:outline-none focus:ring-2 focus:ring-pink-deep/30 disabled:opacity-60"
                  >
                    <Trash2 size={13} aria-hidden />
                  </button>
                </li>
              );
            })}
          </ul>
        </div>
      )}

      <Confetti trigger={confetti} count={22} />
    </div>
  );
}

function TaskCreatePanel({ entries }: { entries: EntryResponse[] }) {
  const initial: CreateTaskFormState = {
    values: {
      entryId: entries[0]?.id ?? "",
      title: "",
      type: "deadline",
      dueDate: "",
      memo: "",
    },
  };
  const [state, formAction] = useActionState(
    createTaskFromTaskPageAction,
    initial,
  );
  const values = state.values ?? initial.values!;

  return (
    <form
      action={formAction}
      className="mb-3 rounded-xl border border-line bg-surface p-3.5 shadow-card"
    >
      <div className="mb-3 flex items-start gap-2">
        <div className="grid h-8 w-8 shrink-0 place-items-center rounded-md bg-sage-wash text-sage">
          <CalendarPlus size={16} aria-hidden />
        </div>
        <div>
          <p className="text-[12px] font-extrabold">Taskを追加</p>
          <p className="mt-0.5 text-[10px] leading-relaxed text-ink-3">
            1つだけ決めればOKです。どのEntryの、何を、いつまでにやるかを登録します。
          </p>
        </div>
      </div>

      {entries.length === 0 ? (
        <div className="rounded-lg border border-dashed border-line bg-cream px-3 py-3 text-center">
          <p className="text-[11px] font-bold text-ink-2">先にEntryを追加してください</p>
          <p className="mt-1 text-[10px] leading-relaxed text-ink-3">
            Task はどの企業の予定かを紐づけて管理します。
          </p>
          <Link
            href="/entry/new"
            className="mt-3 inline-flex items-center gap-1.5 rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5"
          >
            <Plus size={13} aria-hidden />
            Entryを追加
          </Link>
        </div>
      ) : (
        <>
          <div className="mb-2 grid gap-1.5 md:grid-cols-3">
            {["1. Entryを選ぶ", "2. やることを書く", "3. 期日を入れる"].map((label) => (
              <div
                key={label}
                className="rounded-md border border-line bg-cream px-2 py-1.5 text-center text-[10px] font-bold text-ink-2"
              >
                {label}
              </div>
            ))}
          </div>

          <div className="grid gap-2 md:grid-cols-[1.2fr_1.3fr]">
            <label className="block">
              <span className="mb-1 block text-[10px] font-bold text-ink-2">どの応募先？</span>
              <select
                name="entryId"
                aria-label="Entry"
                defaultValue={values.entryId}
                className="h-9 w-full rounded-md border border-line bg-cream px-2 text-[12px] font-semibold outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
              >
                {entries.map((entry) => (
                  <option key={entry.id} value={entry.id}>
                    {companyDisplayName(entry)} / {entry.route}
                  </option>
                ))}
              </select>
            </label>
            <label className="block">
              <span className="mb-1 block text-[10px] font-bold text-ink-2">何をする？</span>
              <input
                name="title"
                aria-label="タスク名"
                defaultValue={values.title}
                required
                placeholder="ES提出、一次面接、SPI受験"
                className="h-9 w-full rounded-md border border-line bg-cream px-2.5 text-[12px] font-semibold outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
              />
            </label>
          </div>

          <div className="mt-2 grid gap-2 md:grid-cols-[1fr_1fr]">
            <fieldset>
              <legend className="mb-1 block text-[10px] font-bold text-ink-2">種類</legend>
              <div className="flex gap-1.5">
                {[
                  ["deadline", "締切"],
                  ["schedule", "予定"],
                ].map(([value, label]) => (
                  <label
                    key={value}
                    className="flex-1 cursor-pointer rounded-md border border-line bg-cream px-2 py-1.5 text-center text-[10px] font-bold text-ink-2 transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage has-[:checked]:text-white"
                  >
                    <input
                      type="radio"
                      name="type"
                      value={value}
                      defaultChecked={values.type === value}
                      className="sr-only"
                    />
                    {label}
                  </label>
                ))}
              </div>
            </fieldset>
            <label className="block">
              <span className="mb-1 block text-[10px] font-bold text-ink-2">いつまで？</span>
              <input
                name="dueDate"
                type="date"
                aria-label="期日"
                defaultValue={values.dueDate}
                className="h-9 w-full rounded-md border border-line bg-cream px-2.5 font-mono text-[12px] font-semibold outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
              />
            </label>
          </div>

          <label className="mt-2 block">
            <span className="mb-1 block text-[10px] font-bold text-ink-2">メモ</span>
            <input
              name="memo"
              defaultValue={values.memo}
              placeholder="持ち物、URL、準備内容など"
              className="h-9 w-full rounded-md border border-line bg-cream px-2.5 text-[12px] outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
            />
          </label>

          {state.error && (
            <p
              role="alert"
              className="mt-2 rounded-md bg-pink/40 px-2.5 py-1.5 text-[10px] font-semibold text-ink"
            >
              {state.error}
            </p>
          )}
          {state.ok && (
            <p className="mt-2 rounded-md bg-sage-wash px-2.5 py-1.5 text-[10px] font-bold text-sage">
              Taskを追加しました。
            </p>
          )}

          <div className="mt-3 flex justify-end">
            <TaskCreateSubmit />
          </div>
        </>
      )}
    </form>
  );
}

function TaskCreateSubmit() {
  const { pending } = useFormStatus();
  return (
    <button
      type="submit"
      disabled={pending}
      className="inline-flex items-center gap-1.5 rounded-lg bg-sage px-3.5 py-2 text-[11px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40 disabled:opacity-60"
    >
      <Plus size={13} aria-hidden />
      {pending ? "追加中…" : "Taskを追加"}
    </button>
  );
}

function TaskEmptyState({ hasEntries }: { hasEntries: boolean }) {
  return (
    <div className="rounded-xl border border-dashed border-line bg-surface px-4 py-8 text-center">
      <p className="font-serif text-base font-extrabold">タスクはまだありません</p>
      <p className="mx-auto mt-1 max-w-[420px] text-[11px] leading-relaxed text-ink-2">
        {hasEntries
          ? "上のフォームから締切や面接予定を追加できます。"
          : "Entry を作ると、その企業に紐づく締切や面接予定を追加できます。"}
      </p>
      {!hasEntries && (
        <div className="mt-4 flex flex-wrap justify-center gap-2">
          <Link
            href="/entry/new"
            className="inline-flex items-center gap-1.5 rounded-lg border border-line bg-surface px-3 py-1.5 text-[11px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
          >
            <ClipboardList size={13} aria-hidden />
            Entryを追加
          </Link>
        </div>
      )}
    </div>
  );
}
