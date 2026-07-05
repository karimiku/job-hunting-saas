"use client";

// Server から initialEntry / initialTasks を受け取る Client Component。
// 更新は Server Action 経由 (同一 origin なので CORS preflight が乗らない)。
// Action 内の revalidatePath がレスポンスに更新済み RSC ツリーを含めるため、
// router.refresh() による二重フルレンダーは不要。

import { useState, useTransition, type FormEvent } from "react";
import { CheckCircle2, ExternalLink, Plus, Trash2 } from "lucide-react";
import {
  companyDisplayName,
  entrySourceUrl,
  type EntryResponse,
} from "@/lib/api/entries";
import {
  ENTRY_STATUS_LABEL,
  STAGE_COLOR,
  STAGE_HINT,
  STAGE_LABEL,
  STAGE_ORDER,
  statusForStage,
  type StageKind,
} from "@/lib/entry-stage";
import { type TaskResponse } from "@/lib/api/tasks";
import type { SelectionFlowResponse, SelectionStageResponse } from "@/lib/selection-flow";
import {
  createTaskForEntryAction,
  deleteEntryAction,
  updateSelectionFlowCurrentStageAction,
  updateEntryAction,
} from "@/app/entry/actions";
import { deleteTaskAction, setTaskStatusAction } from "@/app/task/actions";
import { Confetti } from "./Confetti";
import { TaskTemplateChips, type TaskTemplate } from "./TaskTemplateChips";

const OUTCOME_STATUS = ["in_progress", "offered", "accepted", "rejected", "withdrawn"] as const;

interface Props {
  initialEntry: EntryResponse | null;
  initialTasks: TaskResponse[];
  initialSelectionFlow?: SelectionFlowResponse | null;
}

/** Entry 詳細 — ステージ進捗バー + 「進める →」 + 内定スタンプ + Tasks 表示。 */
export function EntryDetailView({
  initialEntry,
  initialTasks,
  initialSelectionFlow,
}: Props) {
  const [confetti, setConfetti] = useState(0);
  const [taskError, setTaskError] = useState<string | null>(null);
  const [entryError, setEntryError] = useState<string | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [isPending, startTransition] = useTransition();
  const [optimisticEntry, setOptimisticEntry] = useState<{
    stageKind: string;
    stageLabel: string;
    status: string;
  } | null>(null);
  const [selectionFlow, setSelectionFlow] = useState<SelectionFlowResponse | null>(
    initialSelectionFlow ?? null,
  );
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
  const [isTaskFormOpen, setIsTaskFormOpen] = useState(false);

  const applyTaskTemplate = (template: TaskTemplate) => {
    setTaskForm((prev) => ({ ...prev, title: template.title, type: template.type }));
  };

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
  // revalidate 後は作成済みタスクが initialTasks 側にも現れるため、id で重複排除する。
  const initialTaskIds = new Set(initialTasks.map((task) => task.id));
  const tasks = [
    ...initialTasks,
    ...createdTasks.filter((task) => !initialTaskIds.has(task.id)),
  ]
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
      const result = await updateEntryAction(e.id, next);
      if (!result.ok) {
        setOptimisticEntry(null);
        setEntryError(result.error ?? "選考ステータスの更新に失敗しました");
        return;
      }
      if (nextKind === "offer") setConfetti((n) => n + 1);
    });
  };

  const handleSelectFlowStage = (stage: SelectionStageResponse) => {
    const next = {
      stageKind: stage.stageKind,
      stageLabel: stage.stageLabel,
      status: statusForSelectionStage(stage.stageKind),
    };
    const previousFlow = selectionFlow;
    setEntryError(null);
    setOptimisticEntry(next);
    if (selectionFlow) {
      setSelectionFlow({
        ...selectionFlow,
        currentStagePosition: stage.position,
      });
    }
    startTransition(async () => {
      const result = await updateSelectionFlowCurrentStageAction(e.id, stage.position);
      if (!result.ok || !result.selectionFlow) {
        setSelectionFlow(previousFlow);
        setOptimisticEntry(null);
        setEntryError(result.error ?? "選考フローの更新に失敗しました");
        return;
      }
      setSelectionFlow(result.selectionFlow);
      if (stage.stageKind === "offer") setConfetti((n) => n + 1);
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
      const result = await updateEntryAction(e.id, next);
      if (!result.ok) {
        setOptimisticEntry(null);
        setEntryError(result.error ?? "結果ステータスの更新に失敗しました");
        return;
      }
      if (status === "offered" || status === "accepted") setConfetti((n) => n + 1);
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
      const result = await createTaskForEntryAction(e.id, {
        title,
        type: taskForm.type,
        dueDate: taskForm.dueDate ? `${taskForm.dueDate}T00:00:00.000Z` : undefined,
        memo: taskForm.memo.trim() || undefined,
      });
      if (!result.ok || !result.task) {
        setTaskError(result.error ?? "タスクの追加に失敗しました");
        return;
      }
      const created = result.task;
      setCreatedTasks((prev) => [created, ...prev]);
      setTaskForm({ title: "", type: taskForm.type, dueDate: "", memo: "" });
    });
  };

  const handleToggleTask = (task: TaskResponse) => {
    const current = optimisticTaskStatus[task.id] ?? task.status;
    const next = current === "done" ? "todo" : "done";
    setTaskError(null);
    setOptimisticTaskStatus((prev) => ({ ...prev, [task.id]: next }));
    startTransition(async () => {
      const result = await setTaskStatusAction(task.id, next, e.id);
      if (!result.ok) {
        setOptimisticTaskStatus((prev) => {
          const copy = { ...prev };
          delete copy[task.id];
          return copy;
        });
        setTaskError(result.error ?? "タスクの更新に失敗しました");
        return;
      }
      if (next === "done") setConfetti((n) => n + 1);
    });
  };

  const handleDeleteTask = (task: TaskResponse) => {
    if (!window.confirm(`「${task.title}」を削除しますか？`)) return;
    setTaskError(null);
    setDeletedTaskIds((prev) => ({ ...prev, [task.id]: true }));
    startTransition(async () => {
      const result = await deleteTaskAction(task.id, e.id);
      if (!result.ok) {
        setDeletedTaskIds((prev) => {
          const copy = { ...prev };
          delete copy[task.id];
          return copy;
        });
        setTaskError(result.error ?? "タスクの削除に失敗しました");
      }
    });
  };

  const handleDeleteEntry = () => {
    if (
      !window.confirm(
        `「${companyDisplayName(e)}」の応募先を削除しますか？関連するタスクも削除されます。`,
      )
    ) {
      return;
    }
    setDeleteError(null);
    startTransition(async () => {
      const result = await deleteEntryAction(e.id);
      if (!result.ok) {
        setDeleteError(result.error ?? "応募先の削除に失敗しました");
      }
    });
  };

  return (
    <div className="relative">
      {/* Header */}
      <div className="mb-4 flex items-start gap-3">
        <div className="min-w-0 flex-1">
          <h1 className="font-serif text-lg font-extrabold tracking-tight break-words">
            {companyDisplayName(e)}
          </h1>
          <p className="mt-0.5 text-[12px] text-ink-3">
            {e.source} · {e.route}
          </p>
          {sourceUrl && (
            <a
              href={sourceUrl}
              target="_blank"
              rel="noreferrer"
              className="mt-1 inline-flex max-w-full items-center gap-1 rounded-md border border-line bg-surface px-2 py-1 font-mono text-[12px] font-bold text-ink-3 transition-colors hover:border-sage hover:text-sage"
            >
              <span className="truncate">{sourceUrl}</span>
              <ExternalLink size={11} className="shrink-0" aria-hidden />
            </a>
          )}
        </div>
        <div className="flex shrink-0 items-center gap-2">
          {isOffer && (
            <div
              className="rounded-lg border-[2.5px] border-mint bg-mint/10 px-2.5 py-1.5 font-serif text-sm font-black text-mint"
              style={{ animation: "entre-stamp 0.6s cubic-bezier(0.2, 0.8, 0.4, 1) both" }}
            >
              内定！
            </div>
          )}
          <button
            type="button"
            onClick={handleDeleteEntry}
            disabled={isPending}
            aria-label={`${companyDisplayName(e)} の応募先を削除`}
            className="inline-flex h-8 items-center gap-1 rounded-md border border-line bg-surface px-2.5 text-[12px] font-bold text-ink-3 transition-colors hover:border-pink-deep hover:text-pink-deep focus:outline-none focus:ring-2 focus:ring-pink-deep/20 disabled:opacity-60"
          >
            <Trash2 size={14} aria-hidden />
            削除
          </button>
        </div>
      </div>
      {deleteError && (
        <p role="alert" className="mb-3 rounded-md bg-pink/40 px-2.5 py-1.5 text-[12px] font-semibold text-ink">
          {deleteError}
        </p>
      )}

      {/* Stage selector */}
      <section className="mb-3 rounded-xl border border-line bg-surface p-3">
        {/* 選考フェーズ（進捗の軸） */}
        <div>
          <div className="flex items-center justify-between gap-2">
            <p className="text-[12px] font-black text-ink">選考フェーズ</p>
            <p className="rounded-md border-[1.5px] border-sage bg-sage-wash px-2 py-0.5 text-[13px] font-black text-sage">
              現在: <span data-testid="current-stage">{e.stageLabel}</span>
            </p>
          </div>
          <p className="mt-0.5 text-[11px] text-ink-3">今どの段階かを選びます</p>

          {selectionFlow ? (
            <div className="mt-2 grid grid-cols-2 gap-2 md:grid-cols-3">
              {selectionFlow.stages.map((stage) => {
                const reached = stage.position <= selectionFlow.currentStagePosition;
                const selected = stage.position === selectionFlow.currentStagePosition;
                const color = stageColor(stage.stageKind);
                const hint = STAGE_HINT[stage.stageKind as StageKind];
                return (
                  <button
                    type="button"
                    key={stage.id}
                    onClick={() => handleSelectFlowStage(stage)}
                    disabled={isPending || selected}
                    aria-pressed={selected}
                    title={hint}
                    className={`relative grid min-h-11 place-items-center rounded-md border px-2 py-1.5 text-[12px] font-bold transition-transform enabled:hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/30 disabled:cursor-default ${
                      selected ? "scale-[1.04]" : ""
                    }`}
                    style={{
                      borderWidth: selected ? 3 : 1,
                      borderColor: selected ? color : "var(--color-line)",
                      background: reached ? color : "var(--color-line-2)",
                      color: reached ? "#fff" : "var(--color-ink-3)",
                      boxShadow: selected ? `0 0 0 4px ${color}55` : undefined,
                    }}
                  >
                    {selected && (
                      <span
                        aria-hidden
                        className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-ink px-2 py-0.5 text-[10px] font-black whitespace-nowrap text-white shadow-md"
                      >
                        現在ここ
                      </span>
                    )}
                    {stage.stageLabel}
                  </button>
                );
              })}
            </div>
          ) : (
            <div className="mt-2 grid grid-cols-3 gap-2 md:grid-cols-6">
              {STAGE_ORDER.map((kind, i) => {
                const reached = i <= currentIdx;
                const selected = e.stageKind === kind;
                const hint = STAGE_HINT[kind];
                return (
                  <button
                    type="button"
                    key={kind}
                    onClick={() => handleSelectStage(kind)}
                    disabled={isPending || selected}
                    aria-pressed={selected}
                    title={hint}
                    className={`relative grid min-h-11 place-items-center rounded-md border px-1.5 py-1.5 text-[12px] font-bold transition-transform enabled:hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/30 disabled:cursor-default ${
                      selected ? "scale-[1.04]" : ""
                    }`}
                    style={{
                      borderWidth: selected ? 3 : 1,
                      borderColor: selected ? STAGE_COLOR[kind] : "var(--color-line)",
                      background: reached ? STAGE_COLOR[kind] : "var(--color-line-2)",
                      color: reached ? "#fff" : "var(--color-ink-3)",
                      boxShadow: selected ? `0 0 0 4px ${STAGE_COLOR[kind]}55` : undefined,
                    }}
                  >
                    {selected && (
                      <span
                        aria-hidden
                        className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-ink px-2 py-0.5 text-[10px] font-black whitespace-nowrap text-white shadow-md"
                      >
                        現在ここ
                      </span>
                    )}
                    {STAGE_LABEL[kind]}
                  </button>
                );
              })}
            </div>
          )}
        </div>

        {/* 結果（生死の軸） */}
        <div className="mt-4 border-t border-dashed border-line pt-3">
          <p className="text-[12px] font-black text-ink">結果（確定したら選ぶ）</p>
          <p className="mt-0.5 text-[11px] text-ink-3">内定・お見送りなどが決まったら選びます</p>
          <p className="mt-0.5 text-[11px] text-ink-3">
            ※内定・お見送りが確定したときだけ選びます。選考途中は上の「選考フェーズ」だけでOKです
          </p>
          <div className="mt-2 grid grid-cols-2 gap-2 sm:grid-cols-5">
            {OUTCOME_STATUS.map((status) => {
              const selected = e.status === status;
              return (
                <button
                  key={status}
                  type="button"
                  onClick={() => handleSelectOutcome(status)}
                  disabled={isPending || selected}
                  aria-pressed={selected}
                  className={`min-h-11 rounded-md border px-1 text-[12px] font-black transition-colors focus:outline-none focus:ring-2 focus:ring-sage/30 ${
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
        </div>

        {entryError && (
          <p role="alert" className="mt-2 rounded-md bg-pink/40 px-2.5 py-1.5 text-[12px] font-semibold text-ink">
            {entryError}
          </p>
        )}
      </section>

      {/* Memo */}
      {e.memo && (
        <section className="mb-3 rounded-xl border border-line bg-cream-2 p-3">
          <p className="mb-1 text-[12px] font-bold text-sage">メモ</p>
          <p className="text-[12px] leading-relaxed text-ink-2">{e.memo}</p>
        </section>
      )}

      {/* Tasks */}
      <section className="mb-3 rounded-xl border border-line bg-surface p-3">
        <div className="mb-2 flex items-start justify-between gap-2">
          <div>
            <p className="text-[12px] font-bold">タスク</p>
            <p className="mt-0.5 text-[12px] text-ink-3">
              この応募先に必要な締切・予定を追加できます。
            </p>
          </div>
          <span className="rounded-md bg-sage-wash px-2 py-0.5 font-mono text-[12px] font-bold text-sage">
            {tasks.filter((task) => (optimisticTaskStatus[task.id] ?? task.status) === "todo").length}
          </span>
        </div>

        {tasks.length === 0 && (
          <p className="mb-3 rounded-lg border border-dashed border-line bg-cream px-3 py-4 text-center text-[12px] text-ink-3">
            まだタスクがありません。下の「タスクを追加」から締切や予定を追加できます。
          </p>
        )}
        {tasks.length > 0 && (
          <ul className="mb-3 flex flex-col gap-1.5">
            {tasks.map((task) => {
              const status = optimisticTaskStatus[task.id] ?? task.status;
              const done = status === "done";
              return (
                <li
                  key={task.id}
                  className={`flex items-center gap-2 rounded-md border border-line bg-cream px-2 py-1.5 text-[12px] ${
                    done ? "text-ink-3" : ""
                  }`}
                >
                  <button
                    type="button"
                    onClick={() => handleToggleTask(task)}
                    disabled={isPending}
                    aria-pressed={done}
                    aria-label={done ? "タスク未完了に戻す" : "タスク完了にする"}
                    className={`grid h-4 w-4 shrink-0 place-items-center rounded-full border-[1.5px] text-[12px] text-white ${
                      done ? "border-sage bg-sage" : "border-line bg-transparent"
                    }`}
                  >
                    {done ? <CheckCircle2 size={10} aria-hidden /> : null}
                  </button>
                  <span className={`min-w-0 flex-1 truncate ${done ? "line-through" : ""}`}>
                    {task.title}
                  </span>
                  {task.dueDate && (
                    <span className="shrink-0 font-mono text-[12px] text-ink-3">
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

        <button
          type="button"
          onClick={() => setIsTaskFormOpen((prev) => !prev)}
          aria-expanded={isTaskFormOpen}
          className="flex w-full items-center justify-center gap-1 rounded-lg border border-dashed border-line bg-cream px-2.5 py-2 text-[12px] font-bold text-sage transition-colors hover:border-sage"
        >
          <Plus
            size={12}
            aria-hidden
            className={`transition-transform ${isTaskFormOpen ? "rotate-45" : ""}`}
          />
          {isTaskFormOpen ? "閉じる" : "タスクを追加"}
        </button>

        {isTaskFormOpen && (
          <form
            onSubmit={handleCreateTask}
            className="mt-2 rounded-lg border border-line bg-cream p-2.5"
          >
            <div className="mb-2">
              <span className="mb-1 block text-[12px] font-bold text-ink-2">よく使うタスク</span>
              <TaskTemplateChips onSelect={applyTaskTemplate} />
            </div>

            <div className="grid gap-2 md:grid-cols-[1fr_auto]">
              <label className="block">
                <span className="mb-1 block text-[12px] font-bold text-ink-2">タスク名</span>
                <input
                  value={taskForm.title}
                  onChange={(event) =>
                    setTaskForm((prev) => ({ ...prev, title: event.target.value }))
                  }
                  placeholder="例: ES提出"
                  className="h-8 w-full rounded-md border border-line bg-surface px-2 text-[12px] font-semibold outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
                />
              </label>
              <fieldset>
                <legend className="mb-1 block text-[12px] font-bold text-ink-2">種類</legend>
                <div className="flex h-8 gap-1">
                  {[
                    ["deadline", "締切"],
                    ["schedule", "予定"],
                  ].map(([value, label]) => (
                    <label
                      key={value}
                      className="grid cursor-pointer place-items-center rounded-md border border-line bg-surface px-2 text-[12px] font-bold text-ink-2 transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage has-[:checked]:text-white"
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
                <p className="mt-1 text-[11px] text-ink-3">
                  締切＝提出物の期限 ／ 予定＝面接など日時
                </p>
              </fieldset>
            </div>
            <div className="mt-2 grid gap-2 md:grid-cols-[1fr_1.5fr_auto]">
              <label className="block">
                <span className="mb-1 block text-[12px] font-bold text-ink-2">期日</span>
                <input
                  type="date"
                  value={taskForm.dueDate}
                  onChange={(event) =>
                    setTaskForm((prev) => ({ ...prev, dueDate: event.target.value }))
                  }
                  className="h-8 w-full rounded-md border border-line bg-surface px-2 font-mono text-[12px] outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
                />
              </label>
              <label className="block">
                <span className="mb-1 block text-[12px] font-bold text-ink-2">メモ</span>
                <input
                  value={taskForm.memo}
                  onChange={(event) =>
                    setTaskForm((prev) => ({ ...prev, memo: event.target.value }))
                  }
                  placeholder="URL、持ち物、準備内容など"
                  className="h-8 w-full rounded-md border border-line bg-surface px-2 text-[12px] outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
                />
              </label>
              <button
                type="submit"
                disabled={isPending}
                className="self-end inline-flex h-8 items-center justify-center gap-1 rounded-md bg-sage px-2.5 text-[12px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 disabled:opacity-60"
              >
                <Plus size={12} aria-hidden />
                追加
              </button>
            </div>

            {taskError && (
              <p role="alert" className="mt-2 rounded-md bg-pink/40 px-2.5 py-1.5 text-[12px] font-semibold text-ink">
                {taskError}
              </p>
            )}
          </form>
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

function statusForSelectionStage(stageKind: string): string {
  return stageKind === "offer" ? "offered" : "in_progress";
}

function stageColor(stageKind: string): string {
  return STAGE_COLOR[stageKind as StageKind] ?? "var(--color-ink-3)";
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
