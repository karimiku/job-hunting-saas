"use client";

// Server から initialEntry / initialTasks を受け取る Client Component。
// 「進める →」ボタンで PATCH したあと router.refresh() で SSR を再評価する
// (SWR 的な戻し方をするわけではなく、Next.js の RSC 再フェッチに任せる)。

import { useState, useTransition, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { CheckCircle2, ExternalLink, Plus, Trash2 } from "lucide-react";
import {
  updateEntry,
  companyDisplayName,
  entrySourceUrl,
  type EntryResponse,
} from "@/lib/api/entries";
import {
  ENTRY_STATUS_LABEL,
  STAGE_COLOR,
  STAGE_LABEL,
  STAGE_ORDER,
  statusForStage,
  type StageKind,
} from "@/lib/entry-stage";
import {
  createTask,
  deleteTask,
  updateTask,
  type TaskResponse,
} from "@/lib/api/tasks";
import { Confetti } from "./Confetti";

const OUTCOME_STATUS = ["in_progress", "offered", "accepted", "rejected", "withdrawn"] as const;

interface Props {
  initialEntry: EntryResponse | null;
  initialTasks: TaskResponse[];
}

/** Entry 詳細 — ステージ進捗バー + 「進める →」 + 内定スタンプ + Tasks 表示。 */
export function EntryDetailView({ initialEntry, initialTasks }: Props) {
  const router = useRouter();
  const [confetti, setConfetti] = useState(0);
  const [taskError, setTaskError] = useState<string | null>(null);
  const [entryError, setEntryError] = useState<string | null>(null);
  const [isPending, startTransition] = useTransition();
  const [optimisticEntry, setOptimisticEntry] = useState<{
    stageKind: string;
    stageLabel: string;
    status: string;
  } | null>(null);
  const [createdTasks, setCreatedTasks] = useState<TaskResponse[]>([]);
  const [deletedTaskIds, setDeletedTaskIds] = useState<Record<string, boolean>>({});
  const [optimisticTaskStatus, setOptimisticTaskStatus] = useState<
    Record<string, TaskResponse["status"]>
  >({});
  const [taskForm, setTaskForm] = useState({
    title: "",
    type: "deadline" as TaskResponse["type"],
    dueDate: "",
    memo: "",
  });

  if (!initialEntry) {
    return (
      <p role="alert" className="rounded-lg bg-pink/40 p-3 text-[12px] font-semibold text-ink">
        詳細を読み込めませんでした
      </p>
    );
  }

  const e = {
    ...initialEntry,
    ...(optimisticEntry ?? {}),
  };
  const sourceUrl = entrySourceUrl(e);
  const currentIdx = STAGE_ORDER.indexOf(e.stageKind as (typeof STAGE_ORDER)[number]);
  const isOffer = e.stageKind === "offer" || e.status === "offered" || e.status === "accepted";
  const tasks = [...initialTasks, ...createdTasks]
    .filter((task) => !deletedTaskIds[task.id])
    .sort(compareTasks);

  const handleSelectStage = (nextKind: StageKind) => {
    const next = {
      stageKind: nextKind,
      stageLabel: STAGE_LABEL[nextKind],
      status: statusForStage(nextKind),
    };
    setEntryError(null);
    setOptimisticEntry(next);
    startTransition(async () => {
      try {
        await updateEntry(e.id, next);
        if (nextKind === "offer") setConfetti((n) => n + 1);
        router.refresh(); // Server Component を再評価して新しい entry を取得
      } catch {
        setOptimisticEntry(null);
        setEntryError("選考ステータスの更新に失敗しました");
      }
    });
  };

  const handleSelectOutcome = (status: string) => {
    const next =
      status === "offered" || status === "accepted"
        ? {
            stageKind: "offer",
            stageLabel: STAGE_LABEL.offer,
            status,
          }
        : {
            stageKind: e.stageKind,
            stageLabel: e.stageLabel,
            status,
          };
    setEntryError(null);
    setOptimisticEntry(next);
    startTransition(async () => {
      try {
        await updateEntry(e.id, next);
        if (status === "offered" || status === "accepted") setConfetti((n) => n + 1);
        router.refresh();
      } catch {
        setOptimisticEntry(null);
        setEntryError("結果ステータスの更新に失敗しました");
      }
    });
  };

  const handleCreateTask = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const title = taskForm.title.trim();
    if (!title) {
      setTaskError("タスク名は必須です");
      return;
    }
    setTaskError(null);
    startTransition(async () => {
      try {
        const created = await createTask(e.id, {
          title,
          type: taskForm.type,
          dueDate: taskForm.dueDate ? `${taskForm.dueDate}T00:00:00.000Z` : undefined,
          memo: taskForm.memo.trim() || undefined,
        });
        setCreatedTasks((prev) => [created, ...prev]);
        setTaskForm({ title: "", type: taskForm.type, dueDate: "", memo: "" });
        router.refresh();
      } catch {
        setTaskError("タスクの追加に失敗しました");
      }
    });
  };

  const handleToggleTask = (task: TaskResponse) => {
    const current = optimisticTaskStatus[task.id] ?? task.status;
    const next = current === "done" ? "todo" : "done";
    setTaskError(null);
    setOptimisticTaskStatus((prev) => ({ ...prev, [task.id]: next }));
    startTransition(async () => {
      try {
        await updateTask(task.id, { status: next });
        if (next === "done") setConfetti((n) => n + 1);
        router.refresh();
      } catch {
        setOptimisticTaskStatus((prev) => {
          const copy = { ...prev };
          delete copy[task.id];
          return copy;
        });
        setTaskError("タスクの更新に失敗しました");
      }
    });
  };

  const handleDeleteTask = (task: TaskResponse) => {
    if (!window.confirm(`「${task.title}」を削除しますか？`)) return;
    setTaskError(null);
    setDeletedTaskIds((prev) => ({ ...prev, [task.id]: true }));
    startTransition(async () => {
      try {
        await deleteTask(task.id);
        router.refresh();
      } catch {
        setDeletedTaskIds((prev) => {
          const copy = { ...prev };
          delete copy[task.id];
          return copy;
        });
        setTaskError("タスクの削除に失敗しました");
      }
    });
  };

  return (
    <div className="relative">
      {/* Header */}
      <div className="mb-4 flex items-center gap-3">
        <div className="min-w-0 flex-1">
          <h1 className="font-serif text-lg font-extrabold tracking-tight break-words">
            {companyDisplayName(e)}
          </h1>
          <p className="mt-0.5 text-[10px] text-ink-3">
            {e.source} · {e.route}
          </p>
          {sourceUrl && (
            <a
              href={sourceUrl}
              target="_blank"
              rel="noreferrer"
              className="mt-1 inline-flex max-w-full items-center gap-1 rounded-md border border-line bg-surface px-2 py-1 font-mono text-[10px] font-bold text-ink-3 transition-colors hover:border-sage hover:text-sage"
            >
              <span className="truncate">{sourceUrl}</span>
              <ExternalLink size={11} className="shrink-0" aria-hidden />
            </a>
          )}
        </div>
        {isOffer && (
          <div
            className="rounded-lg border-[2.5px] border-mint bg-mint/10 px-2.5 py-1.5 font-serif text-sm font-black text-mint"
            style={{ animation: "entre-stamp 0.6s cubic-bezier(0.2, 0.8, 0.4, 1) both" }}
          >
            内定！
          </div>
        )}
      </div>

      {/* Stage selector */}
      <section className="mb-3 rounded-xl border border-line bg-surface p-3">
        <div className="mb-2 flex items-center justify-between">
          <p className="text-[10px] font-bold text-ink-2">選考ステータス</p>
          <span className="rounded-full bg-cream px-2 py-0.5 text-[9px] font-black text-ink-3">
            {ENTRY_STATUS_LABEL[e.status] ?? e.status}
          </span>
        </div>
        <div className="grid grid-cols-3 gap-1.5 md:grid-cols-6">
          {STAGE_ORDER.map((kind, i) => {
            const reached = i <= currentIdx;
            const selected = e.stageKind === kind;
            return (
              <button
                type="button"
                key={kind}
                onClick={() => handleSelectStage(kind)}
                disabled={isPending || selected}
                aria-pressed={selected}
                className="grid min-h-9 place-items-center rounded-md border text-[10px] font-bold transition-transform enabled:hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/30 disabled:cursor-default"
                style={{
                  borderColor: selected ? STAGE_COLOR[kind] : "var(--color-line)",
                  background: reached ? STAGE_COLOR[kind] : "var(--color-line-2)",
                  color: reached ? "#fff" : "var(--color-ink-3)",
                }}
              >
                {STAGE_LABEL[kind]}
              </button>
            );
          })}
        </div>
        <div className="mt-2 grid grid-cols-5 gap-1">
          {OUTCOME_STATUS.map((status) => {
            const selected = e.status === status;
            return (
              <button
                key={status}
                type="button"
                onClick={() => handleSelectOutcome(status)}
                disabled={isPending || selected}
                aria-pressed={selected}
                className={`h-7 rounded-md border px-1 text-[9px] font-black transition-colors focus:outline-none focus:ring-2 focus:ring-sage/30 ${
                  selected
                    ? "border-sage bg-sage text-white"
                    : "border-line bg-cream text-ink-3 hover:border-sage hover:text-sage"
                }`}
              >
                {ENTRY_STATUS_LABEL[status]}
              </button>
            );
          })}
        </div>
        <p className="mt-2 text-[10px] font-bold text-sage">
          現在: <span data-testid="current-stage">{e.stageLabel}</span>
        </p>
        {entryError && (
          <p role="alert" className="mt-2 rounded-md bg-pink/40 px-2.5 py-1.5 text-[10px] font-semibold text-ink">
            {entryError}
          </p>
        )}
      </section>

      {/* Memo */}
      {e.memo && (
        <section className="mb-3 rounded-xl border border-line bg-cream-2 p-3">
          <p className="mb-1 text-[11px] font-bold text-sage">メモ</p>
          <p className="text-[11px] leading-relaxed text-ink-2">{e.memo}</p>
        </section>
      )}

      {/* Tasks */}
      <section className="mb-3 rounded-xl border border-line bg-surface p-3">
        <div className="mb-2 flex items-start justify-between gap-2">
          <div>
            <p className="text-[12px] font-bold">タスク</p>
            <p className="mt-0.5 text-[10px] text-ink-3">
              このEntryに必要な締切・予定を追加できます。
            </p>
          </div>
          <span className="rounded-md bg-sage-wash px-2 py-0.5 font-mono text-[10px] font-bold text-sage">
            {tasks.filter((task) => (optimisticTaskStatus[task.id] ?? task.status) === "todo").length}
          </span>
        </div>

        <form
          onSubmit={handleCreateTask}
          className="mb-3 rounded-lg border border-line bg-cream p-2.5"
        >
          <div className="grid gap-2 md:grid-cols-[1fr_auto]">
            <label className="block">
              <span className="mb-1 block text-[10px] font-bold text-ink-2">タスク名</span>
              <input
                value={taskForm.title}
                onChange={(event) =>
                  setTaskForm((prev) => ({ ...prev, title: event.target.value }))
                }
                placeholder="ES提出、面接準備、SPI受験"
                className="h-8 w-full rounded-md border border-line bg-surface px-2 text-[11px] font-semibold outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
              />
            </label>
            <fieldset>
              <legend className="mb-1 block text-[10px] font-bold text-ink-2">種類</legend>
              <div className="flex h-8 gap-1">
                {[
                  ["deadline", "締切"],
                  ["schedule", "予定"],
                ].map(([value, label]) => (
                  <label
                    key={value}
                    className="grid cursor-pointer place-items-center rounded-md border border-line bg-surface px-2 text-[10px] font-bold text-ink-2 transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage has-[:checked]:text-white"
                  >
                    <input
                      type="radio"
                      name="entry-detail-task-type"
                      value={value}
                      checked={taskForm.type === value}
                      onChange={() =>
                        setTaskForm((prev) => ({
                          ...prev,
                          type: value as TaskResponse["type"],
                        }))
                      }
                      className="sr-only"
                    />
                    {label}
                  </label>
                ))}
              </div>
            </fieldset>
          </div>
          <div className="mt-2 grid gap-2 md:grid-cols-[1fr_1.5fr_auto]">
            <label className="block">
              <span className="mb-1 block text-[10px] font-bold text-ink-2">期日</span>
              <input
                type="date"
                value={taskForm.dueDate}
                onChange={(event) =>
                  setTaskForm((prev) => ({ ...prev, dueDate: event.target.value }))
                }
                className="h-8 w-full rounded-md border border-line bg-surface px-2 font-mono text-[11px] outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
              />
            </label>
            <label className="block">
              <span className="mb-1 block text-[10px] font-bold text-ink-2">メモ</span>
              <input
                value={taskForm.memo}
                onChange={(event) =>
                  setTaskForm((prev) => ({ ...prev, memo: event.target.value }))
                }
                placeholder="URL、持ち物、準備内容など"
                className="h-8 w-full rounded-md border border-line bg-surface px-2 text-[11px] outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
              />
            </label>
            <button
              type="submit"
              disabled={isPending}
              className="self-end inline-flex h-8 items-center justify-center gap-1 rounded-md bg-sage px-2.5 text-[10px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 disabled:opacity-60"
            >
              <Plus size={12} aria-hidden />
              追加
            </button>
          </div>
        </form>

        {taskError && (
          <p role="alert" className="mb-2 rounded-md bg-pink/40 px-2.5 py-1.5 text-[10px] font-semibold text-ink">
            {taskError}
          </p>
        )}
        {tasks.length === 0 && (
          <p className="rounded-lg border border-dashed border-line bg-cream px-3 py-4 text-center text-[11px] text-ink-3">
            まだタスクがありません。上のフォームから締切や予定を追加できます。
          </p>
        )}
        {tasks.length > 0 && (
          <ul className="flex flex-col gap-1.5">
            {tasks.map((task) => {
              const status = optimisticTaskStatus[task.id] ?? task.status;
              const done = status === "done";
              return (
                <li
                  key={task.id}
                  className={`flex items-center gap-2 rounded-md border border-line bg-cream px-2 py-1.5 text-[11px] ${
                    done ? "text-ink-3" : ""
                  }`}
                >
                  <button
                    type="button"
                    onClick={() => handleToggleTask(task)}
                    disabled={isPending}
                    aria-pressed={done}
                    aria-label={done ? "タスク未完了に戻す" : "タスク完了にする"}
                    className={`grid h-4 w-4 shrink-0 place-items-center rounded-full border-[1.5px] text-[9px] text-white ${
                      done ? "border-sage bg-sage" : "border-line bg-transparent"
                    }`}
                  >
                    {done ? <CheckCircle2 size={10} aria-hidden /> : null}
                  </button>
                  <span className={`min-w-0 flex-1 truncate ${done ? "line-through" : ""}`}>
                    {task.title}
                  </span>
                  {task.dueDate && (
                    <span className="shrink-0 font-mono text-[9px] text-ink-3">
                      {formatTaskDue(task.dueDate)}
                    </span>
                  )}
                  <button
                    type="button"
                    onClick={() => handleDeleteTask(task)}
                    disabled={isPending}
                    aria-label={`タスク「${task.title}」を削除`}
                    className="grid h-6 w-6 shrink-0 place-items-center rounded-md text-ink-3 transition-colors hover:text-pink-deep disabled:opacity-60"
                  >
                    <Trash2 size={12} aria-hidden />
                  </button>
                </li>
              );
            })}
          </ul>
        )}
      </section>

      <Confetti trigger={confetti} count={28} />
    </div>
  );
}

function taskTime(task: TaskResponse): number {
  if (!task.dueDate) return Number.POSITIVE_INFINITY;
  const date = new Date(task.dueDate);
  return Number.isNaN(date.getTime()) ? Number.POSITIVE_INFINITY : date.getTime();
}

function compareTasks(a: TaskResponse, b: TaskResponse): number {
  const aDone = a.status === "done" ? 1 : 0;
  const bDone = b.status === "done" ? 1 : 0;
  if (aDone !== bDone) return aDone - bDone;
  return taskTime(a) - taskTime(b);
}

function formatTaskDue(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return `${date.getMonth() + 1}/${date.getDate()}`;
}
