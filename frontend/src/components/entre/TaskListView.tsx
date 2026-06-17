"use client";

// initialTasks を SSR で受け取り、表示とトグル操作のみ担う Client Component。

import { useActionState, useMemo, useState, useTransition } from "react";
import { useFormStatus } from "react-dom";
import Link from "next/link";
import { ArrowRight, CalendarPlus, CheckCircle2, ClipboardList, Plus, Trash2 } from "lucide-react";
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
  const [confetti, setConfetti] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [selectedEntryId, setSelectedEntryId] = useState<string>("all");
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
      const result = await setTaskStatusAction(task.id, next, task.entryId);
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
      // action 内の revalidatePath("/task") がレスポンスに更新済み RSC ツリーを含めるため
      // router.refresh() は不要 (呼ぶと同じページをもう一度フルレンダーしてしまう)。
    });
  };

  const deleteTask = (task: TaskWithEntry) => {
    if (!window.confirm(`「${task.title}」を削除しますか？`)) return;
    setError(null);
    setDeletingIds((prev) => ({ ...prev, [task.id]: true }));
    startTransition(async () => {
      const result = await deleteTaskAction(task.id, task.entryId);
      if (!result.ok) {
        setDeletingIds((prev) => {
          const next = { ...prev };
          delete next[task.id];
          return next;
        });
        setError(result.error ?? "タスクの削除に失敗しました");
        return;
      }
    });
  };

  const entryById = useMemo(
    () => new Map(entries.map((entry) => [entry.id, entry])),
    [entries],
  );

  const allVisibleTasks = sortTasksForDisplay(initialTasks).filter(
    (task) => !deletingIds[task.id],
  );
  const tasks =
    selectedEntryId === "all"
      ? allVisibleTasks
      : allVisibleTasks.filter((task) => task.entryId === selectedEntryId);
  const groupedTasks = groupTasksByEntry(tasks, entryById);
  const remainingCount = tasks.filter(
    (task) => (optimistic[task.id] ?? task.status) === "todo",
  ).length;

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

      {allVisibleTasks.length === 0 ? (
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
              {remainingCount}件残り
            </span>
          </div>

          <EntryTaskFilter
            entries={entries}
            tasks={allVisibleTasks}
            selectedEntryId={selectedEntryId}
            onSelect={setSelectedEntryId}
          />

          {tasks.length === 0 ? (
            <div className="rounded-xl border border-dashed border-line bg-surface px-4 py-6 text-center">
              <p className="text-[12px] font-bold text-ink-2">
                このEntryにはタスクがありません
              </p>
              <p className="mt-1 text-[10px] text-ink-3">
                上のフォームで応募先を選ぶと、このEntryに予定を追加できます。
              </p>
            </div>
          ) : (
            <div className="flex flex-col gap-3">
              {groupedTasks.map((group) => (
                <section
                  key={group.entryId}
                  className="rounded-xl border border-line bg-surface p-2.5"
                >
                  <div className="mb-2 flex items-center justify-between gap-2 px-0.5">
                    <Link
                      href={group.entryId === "unknown" ? "/entry" : `/entry/${group.entryId}`}
                      prefetch={false}
                      className="min-w-0 text-[11px] font-extrabold text-ink hover:text-sage"
                    >
                      <span className="truncate">{group.companyName}</span>
                    </Link>
                    <span className="shrink-0 rounded-md bg-cream px-2 py-0.5 text-[9px] font-black text-ink-3">
                      未完了 {group.openCount}
                    </span>
                  </div>
                  <ul className="flex flex-col gap-2">
                    {group.tasks.map((task) => {
                      const status = optimistic[task.id] ?? task.status;
                      const done = status === "done";
                      return (
                        <TaskRow
                          key={task.id}
                          task={task}
                          done={done}
                          isPending={isPending}
                          onToggle={toggle}
                          onDelete={deleteTask}
                        />
                      );
                    })}
                  </ul>
                </section>
              ))}
            </div>
          )}
        </div>
      )}

      <Confetti trigger={confetti} count={22} />
    </div>
  );
}

function TaskRow({
  task,
  done,
  isPending,
  onToggle,
  onDelete,
}: {
  task: TaskWithEntry;
  done: boolean;
  isPending: boolean;
  onToggle: (task: TaskWithEntry) => void;
  onDelete: (task: TaskWithEntry) => void;
}) {
  return (
    <li
      className={`flex items-center gap-3 rounded-lg border border-line bg-cream px-3 py-2.5 transition-opacity ${
        done ? "opacity-50" : ""
      }`}
    >
      <button
        type="button"
        onClick={() => onToggle(task)}
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
      <Link
        href={`/task/${task.id}`}
        prefetch={false}
        className="group -my-1 flex min-w-0 flex-1 items-center gap-3 rounded-md py-1 pr-1 transition-colors hover:text-sage focus:outline-none focus:ring-2 focus:ring-sage/30"
      >
        <div className="min-w-0 flex-1">
          <div className="flex min-w-0 items-center gap-1.5">
            <div className={`truncate text-[12px] font-semibold ${done ? "line-through" : ""}`}>
              {task.title}
            </div>
            <ArrowRight
              size={12}
              className="shrink-0 text-ink-3 opacity-0 transition-opacity group-hover:opacity-100 group-focus:opacity-100"
              aria-hidden
            />
          </div>
          <div className="mt-0.5 truncate text-[10px] text-ink-3">
            {task.memo ? task.memo : task.type === "deadline" ? "締切タスク" : "予定"}
          </div>
        </div>
        <span
          className={`shrink-0 rounded-md px-2 py-0.5 font-mono text-[10px] font-bold text-white ${
            TYPE_BADGE[task.type] ?? "bg-sage"
          }`}
        >
          {formatDue(task.dueDate)}
        </span>
      </Link>
      <button
        type="button"
        onClick={() => onDelete(task)}
        disabled={isPending}
        aria-label={`タスク「${task.title}」を削除`}
        className="grid h-7 w-7 shrink-0 place-items-center rounded-md border border-line text-ink-3 transition-colors hover:border-pink-deep hover:text-pink-deep focus:outline-none focus:ring-2 focus:ring-pink-deep/30 disabled:opacity-60"
      >
        <Trash2 size={13} aria-hidden />
      </button>
    </li>
  );
}

function EntryTaskFilter({
  entries,
  tasks,
  selectedEntryId,
  onSelect,
}: {
  entries: EntryResponse[];
  tasks: TaskWithEntry[];
  selectedEntryId: string;
  onSelect: (entryId: string) => void;
}) {
  if (entries.length === 0) return null;

  const taskCountByEntry = new Map<string, number>();
  for (const task of tasks) {
    taskCountByEntry.set(task.entryId, (taskCountByEntry.get(task.entryId) ?? 0) + 1);
  }

  return (
    <div className="-mx-5 overflow-x-auto px-5 md:mx-0 md:px-0">
      <div className="flex min-w-max gap-1.5 pb-1">
        <button
          type="button"
          onClick={() => onSelect("all")}
          aria-pressed={selectedEntryId === "all"}
          className={`h-8 rounded-full border px-3 text-[10px] font-black transition-colors ${
            selectedEntryId === "all"
              ? "border-sage bg-sage text-white"
              : "border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
          }`}
        >
          すべて {tasks.length}
        </button>
        {entries.map((entry) => {
          const selected = selectedEntryId === entry.id;
          const count = taskCountByEntry.get(entry.id) ?? 0;
          return (
            <button
              key={entry.id}
              type="button"
              onClick={() => onSelect(entry.id)}
              aria-pressed={selected}
              className={`h-8 max-w-[180px] rounded-full border px-3 text-[10px] font-black transition-colors ${
                selected
                  ? "border-sage bg-sage text-white"
                  : "border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
              }`}
            >
              <span className="inline-block max-w-[120px] truncate align-bottom">
                {companyDisplayName(entry)}
              </span>{" "}
              {count}
            </button>
          );
        })}
      </div>
    </div>
  );
}

interface TaskGroup {
  entryId: string;
  companyName: string;
  openCount: number;
  tasks: TaskWithEntry[];
}

function groupTasksByEntry(
  tasks: TaskWithEntry[],
  entryById: Map<string, EntryResponse>,
): TaskGroup[] {
  const groups = new Map<string, TaskGroup>();
  for (const task of tasks) {
    const entry = entryById.get(task.entryId);
    const entryId = entry?.id ?? task.entryId ?? "unknown";
    const current =
      groups.get(entryId) ??
      {
        entryId,
        companyName: task.companyName ?? (entry ? companyDisplayName(entry) : "（会社名未設定）"),
        openCount: 0,
        tasks: [],
      };
    current.tasks.push(task);
    if (task.status === "todo") current.openCount += 1;
    groups.set(entryId, current);
  }
  return [...groups.values()];
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
          <p className="text-[12px] font-extrabold">タスクを追加</p>
          <p className="mt-0.5 text-[10px] leading-relaxed text-ink-3">
            応募先、内容、期日だけ入れます。
          </p>
        </div>
      </div>

      {entries.length === 0 ? (
        <div className="rounded-lg border border-dashed border-line bg-cream px-3 py-3 text-center">
          <p className="text-[11px] font-bold text-ink-2">先にEntryを追加してください</p>
          <p className="mt-1 text-[10px] leading-relaxed text-ink-3">
            タスクはどの企業の予定かを紐づけて管理します。
          </p>
          <Link
            href="/entry/new"
            prefetch={false}
            className="mt-3 inline-flex items-center gap-1.5 rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5"
          >
            <Plus size={13} aria-hidden />
            Entryを追加
          </Link>
        </div>
      ) : (
        <>
          <div className="grid gap-2 md:grid-cols-[1.2fr_1.3fr]">
            <label className="block">
              <span className="mb-1 block text-[10px] font-bold text-ink-2">応募先</span>
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
              <span className="mb-1 block text-[10px] font-bold text-ink-2">内容</span>
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
              <span className="mb-1 block text-[10px] font-bold text-ink-2">期日</span>
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
              タスクを追加しました。
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
      {pending ? "追加中…" : "タスクを追加"}
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
            prefetch={false}
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
