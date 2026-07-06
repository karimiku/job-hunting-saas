"use client";

// initialTasks を SSR で受け取り、表示とトグル操作のみ担う Client Component。

import { useActionState, useMemo, useState, useTransition } from "react";
import { useFormStatus } from "react-dom";
import Link from "next/link";
import {
  ArrowRight,
  CalendarPlus,
  CheckCircle2,
  ClipboardList,
  Clock,
  Plus,
  Trash2,
} from "lucide-react";
import {
  createTaskFromTaskPageAction,
  deleteTaskAction,
  rescheduleTaskAction,
  setTaskStatusAction,
  type CreateTaskFormState,
} from "@/app/task/actions";
import type { TaskWithEntry } from "@/lib/api/server-resources";
import {
  companyDisplayName,
  type EntryResponse,
} from "@/lib/api/entries";
import { Confetti } from "./Confetti";
import { TaskTemplateChips, type TaskTemplate } from "./TaskTemplateChips";
import { addDays, DueDateQuickPicker } from "./DueDateQuickPicker";

interface Props {
  initialTasks: TaskWithEntry[];
  entries: EntryResponse[];
}

// 期日バッジは type (締切/予定) ではなく緊急度で色分けする
// (type の区別は sub テキストの「締切タスク」「予定」が担う)。
// 期日までの残り日数をカレンダー日で数える。単純な時刻差の floor だと、夕方に
// 「明日(翌0時UTC)」の締切を登録したとき差が24時間未満になり本日扱いになってしまう。
// 双方をローカル0時に丸めた日付差で数えることで、暦の上での日数と一致させる。
export function daysUntilDue(d: Date, now: Date): number {
  const due = new Date(d.getFullYear(), d.getMonth(), d.getDate());
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  return Math.round((due.getTime() - today.getTime()) / 86_400_000);
}

// 超過: 1-2日=pink、3日以上=より濃いpink-deep。本日=最も強いink、明日=pink-deep。
// それ以降=従来どおり (3日以内=amber、超=sky)。期日なし=sage。
export function dueLabel(dueDate: string | null, now: Date = new Date()): string {
  if (!dueDate) return "期日なし";
  const d = new Date(dueDate);
  // ISO 文字列 (YYYY-MM-DD...) を M/D に短縮表示。パースできなければ原文を返す。
  if (Number.isNaN(d.getTime())) return dueDate;
  const base = `${d.getMonth() + 1}/${d.getDate()}`;
  const days = daysUntilDue(d, now);
  if (days < 0) return `${base} ・${Math.abs(days)}日超過`;
  if (days === 0) return "本日締切";
  if (days === 1) return "明日締切";
  return base;
}

export function dueColor(dueDate: string | null, now: Date = new Date()): string {
  if (!dueDate) return "bg-sage";
  const d = new Date(dueDate);
  if (Number.isNaN(d.getTime())) return "bg-sage";
  const days = daysUntilDue(d, now);
  if (days < 0) return Math.abs(days) >= 3 ? "bg-pink-deep" : "bg-pink";
  if (days === 0) return "bg-ink";
  if (days === 1) return "bg-pink-deep";
  if (days <= 3) return "bg-amber";
  return "bg-sky";
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

export type TaskStatusFilter = "all" | "todo" | "week";

// 今日〜7日以内 (両端含む) を「今週締切」とする。期日なし・8日後以降・超過は含めない。
export function isDueThisWeek(dueDate: string | null, now: Date = new Date()): boolean {
  if (!dueDate) return false;
  const d = new Date(dueDate);
  if (Number.isNaN(d.getTime())) return false;
  const days = Math.floor((d.getTime() - now.getTime()) / 86_400_000);
  return days >= 0 && days <= 7;
}

export function filterTasksByStatusFilter(
  tasks: TaskWithEntry[],
  filter: TaskStatusFilter,
  now: Date = new Date(),
): TaskWithEntry[] {
  if (filter === "todo") return tasks.filter((task) => task.status !== "done");
  if (filter === "week") return tasks.filter((task) => isDueThisWeek(task.dueDate, now));
  return tasks;
}

export type TaskTypeFilter = "all" | "deadline" | "schedule";

export function filterTasksByTypeFilter(
  tasks: TaskWithEntry[],
  filter: TaskTypeFilter,
): TaskWithEntry[] {
  if (filter === "all") return tasks;
  return tasks.filter((task) => task.type === filter);
}

// 延期ワンタップの相対日付選択肢。今日は「延期」の選択肢としては不要なので含めない。
const RESCHEDULE_OPTIONS: { label: string; days: number }[] = [
  { label: "明日", days: 1 },
  { label: "+3日", days: 3 },
  { label: "+1週間", days: 7 },
];

export function TaskListView({ initialTasks, entries }: Props) {
  const [confetti, setConfetti] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [selectedEntryId, setSelectedEntryId] = useState<string>("all");
  const [statusFilter, setStatusFilter] = useState<TaskStatusFilter>("all");
  const [typeFilter, setTypeFilter] = useState<TaskTypeFilter>("all");
  const [isPending, startTransition] = useTransition();
  const [deletingIds, setDeletingIds] = useState<Record<string, boolean>>({});
  // 楽観更新用に「いまトグル中の taskId → 目標 status」を保持する。
  const [optimistic, setOptimistic] = useState<Record<string, "todo" | "done">>(
    {},
  );
  // 延期の楽観更新用に「いま延期中の taskId → 目標 dueDate (YYYY-MM-DD)」を保持する。
  const [optimisticDueDate, setOptimisticDueDate] = useState<
    Record<string, string>
  >({});

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

  const reschedule = (task: TaskWithEntry, dueDate: string) => {
    setError(null);
    setOptimisticDueDate((prev) => ({ ...prev, [task.id]: dueDate }));

    startTransition(async () => {
      const result = await rescheduleTaskAction(task.id, dueDate, task.entryId);
      if (!result.ok) {
        setOptimisticDueDate((prev) => {
          const next = { ...prev };
          delete next[task.id];
          return next;
        });
        setError(result.error ?? "タスクの延期に失敗しました");
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
  const entryFilteredTasks =
    selectedEntryId === "all"
      ? allVisibleTasks
      : allVisibleTasks.filter((task) => task.entryId === selectedEntryId);
  const typeFilteredTasks = filterTasksByTypeFilter(entryFilteredTasks, typeFilter);
  const tasks = filterTasksByStatusFilter(typeFilteredTasks, statusFilter);
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
          className="mb-2 rounded-lg bg-pink/40 px-3 py-2 text-[12px] font-semibold text-ink"
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
              <p className="text-[12px] font-extrabold">未完了を上から片づける</p>
              <p className="mt-0.5 text-[12px] text-ink-3">
                左の丸を押すと完了。期日が近い順に並びます。
              </p>
            </div>
            <span className="rounded-md bg-sage-soft px-2 py-1 text-[12px] font-bold text-sage">
              {remainingCount}件残り
            </span>
          </div>

          <TaskStatusFilterChips
            tasks={allVisibleTasks}
            selected={statusFilter}
            onSelect={setStatusFilter}
          />

          <TaskTypeFilterChips
            tasks={allVisibleTasks}
            selected={typeFilter}
            onSelect={setTypeFilter}
          />

          <EntryTaskFilter
            entries={entries}
            tasks={allVisibleTasks}
            selectedEntryId={selectedEntryId}
            onSelect={setSelectedEntryId}
          />

          {tasks.length === 0 ? (
            <div className="rounded-xl border border-dashed border-line bg-surface px-4 py-6 text-center">
              <p className="text-[12px] font-bold text-ink-2">
                この応募先にはタスクがありません
              </p>
              <p className="mt-1 text-[12px] text-ink-3">
                上のフォームで応募先を選ぶと、予定を追加できます。
              </p>
            </div>
          ) : (
            <div className="flex flex-col gap-3 lg:grid lg:grid-cols-2 lg:gap-3">
              {groupedTasks.map((group) => (
                <section
                  key={group.entryId}
                  className="rounded-xl border border-line bg-surface p-2.5"
                >
                  <div className="mb-2 flex items-center justify-between gap-2 px-0.5">
                    <Link
                      href={group.entryId === "unknown" ? "/entry" : `/entry/${group.entryId}`}
                      prefetch={false}
                      className="min-w-0 text-[12px] font-extrabold text-ink hover:text-sage"
                    >
                      <span className="truncate">{group.companyName}</span>
                    </Link>
                    <span className="shrink-0 rounded-md bg-cream px-2 py-0.5 text-[12px] font-black text-ink-3">
                      未完了 {group.openCount}
                    </span>
                  </div>
                  <ul className="flex flex-col gap-2">
                    {group.tasks.map((task) => {
                      const status = optimistic[task.id] ?? task.status;
                      const done = status === "done";
                      const dueDate = optimisticDueDate[task.id] ?? task.dueDate;
                      const displayTask =
                        dueDate === task.dueDate ? task : { ...task, dueDate };
                      return (
                        <TaskRow
                          key={task.id}
                          task={displayTask}
                          done={done}
                          isPending={isPending}
                          onToggle={toggle}
                          onDelete={deleteTask}
                          onReschedule={reschedule}
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
  onReschedule,
}: {
  task: TaskWithEntry;
  done: boolean;
  isPending: boolean;
  onToggle: (task: TaskWithEntry) => void;
  onDelete: (task: TaskWithEntry) => void;
  onReschedule: (task: TaskWithEntry, dueDate: string) => void;
}) {
  const [rescheduleOpen, setRescheduleOpen] = useState(false);

  return (
    <li
      className={`flex flex-col gap-1.5 rounded-lg border border-line bg-cream px-3 py-2.5 transition-opacity ${
        done ? "opacity-50" : ""
      }`}
    >
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={() => onToggle(task)}
          disabled={isPending}
          aria-pressed={done}
          aria-label={done ? "タスク未完了に戻す" : "タスク完了にする"}
          className={`grid h-6 w-6 shrink-0 place-items-center rounded-full text-[12px] text-white transition-colors focus:outline-none focus:ring-2 focus:ring-sage/30 disabled:opacity-60 ${
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
            <div className="flex min-w-0 items-start gap-1.5">
              <div className={`line-clamp-2 break-words text-[12px] font-semibold ${done ? "line-through" : ""}`}>
                {task.title}
              </div>
              <ArrowRight
                size={12}
                className="mt-0.5 shrink-0 text-ink-3 opacity-0 transition-opacity group-hover:opacity-100 group-focus:opacity-100"
                aria-hidden
              />
            </div>
            <div className="mt-0.5 truncate text-[12px] text-ink-3">
              {task.memo ? task.memo : task.type === "deadline" ? "締切タスク" : "予定"}
            </div>
          </div>
          <span
            className={`shrink-0 rounded-md px-1.5 py-0.5 font-mono text-[12px] font-bold text-white ${dueColor(task.dueDate)}`}
          >
            {dueLabel(task.dueDate)}
          </span>
        </Link>
        {!done && (
          <button
            type="button"
            onClick={() => setRescheduleOpen((prev) => !prev)}
            disabled={isPending}
            aria-expanded={rescheduleOpen}
            aria-label={`タスク「${task.title}」の期日を延期`}
            className="flex shrink-0 items-center gap-1 rounded-md border border-line px-1.5 py-1 text-[12px] font-bold text-ink-3 transition-colors hover:border-sage hover:text-sage focus:outline-none focus:ring-2 focus:ring-sage/30 disabled:opacity-60 sm:px-2"
          >
            <Clock size={13} aria-hidden />
            <span className="hidden sm:inline">延期</span>
          </button>
        )}
        <button
          type="button"
          onClick={() => onDelete(task)}
          disabled={isPending}
          aria-label={`タスク「${task.title}」を削除`}
          className="grid h-7 w-7 shrink-0 place-items-center rounded-md border border-line text-ink-3 transition-colors hover:border-pink-deep hover:text-pink-deep focus:outline-none focus:ring-2 focus:ring-pink-deep/30 disabled:opacity-60"
        >
          <Trash2 size={13} aria-hidden />
        </button>
      </div>
      {rescheduleOpen && (
        <div className="flex gap-1.5 pl-9">
          {RESCHEDULE_OPTIONS.map((option) => (
            <button
              key={option.label}
              type="button"
              onClick={() => {
                onReschedule(task, addDays(new Date(), option.days));
                setRescheduleOpen(false);
              }}
              className="rounded-md border border-line bg-surface px-2 py-1 text-[12px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
            >
              {option.label}
            </button>
          ))}
        </div>
      )}
    </li>
  );
}

const STATUS_FILTER_OPTIONS: { value: TaskStatusFilter; label: string }[] = [
  { value: "all", label: "すべて" },
  { value: "todo", label: "未完了のみ" },
  { value: "week", label: "今週締切" },
];

function TaskStatusFilterChips({
  tasks,
  selected,
  onSelect,
}: {
  tasks: TaskWithEntry[];
  selected: TaskStatusFilter;
  onSelect: (filter: TaskStatusFilter) => void;
}) {
  return (
    <div className="-mx-5 overflow-x-auto px-5 md:mx-0 md:px-0">
      <div className="flex min-w-max gap-1.5 pb-1">
        {STATUS_FILTER_OPTIONS.map((option) => {
          const selectedOption = selected === option.value;
          const count = filterTasksByStatusFilter(tasks, option.value).length;
          return (
            <button
              key={option.value}
              type="button"
              onClick={() => onSelect(option.value)}
              aria-pressed={selectedOption}
              className={`h-8 rounded-full border px-3 text-[12px] font-black transition-colors ${
                selectedOption
                  ? "border-sage bg-sage text-white"
                  : "border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
              }`}
            >
              {option.label} {count}
            </button>
          );
        })}
      </div>
    </div>
  );
}

const TYPE_FILTER_OPTIONS: { value: TaskTypeFilter; label: string }[] = [
  { value: "all", label: "すべて" },
  { value: "deadline", label: "締切" },
  { value: "schedule", label: "予定" },
];

function TaskTypeFilterChips({
  tasks,
  selected,
  onSelect,
}: {
  tasks: TaskWithEntry[];
  selected: TaskTypeFilter;
  onSelect: (filter: TaskTypeFilter) => void;
}) {
  return (
    <div className="-mx-5 overflow-x-auto px-5 md:mx-0 md:px-0">
      <div className="flex min-w-max gap-1.5 pb-1">
        {TYPE_FILTER_OPTIONS.map((option) => {
          const selectedOption = selected === option.value;
          const count = filterTasksByTypeFilter(tasks, option.value).length;
          return (
            <button
              key={option.value}
              type="button"
              onClick={() => onSelect(option.value)}
              aria-pressed={selectedOption}
              className={`h-8 rounded-full border px-3 text-[12px] font-black transition-colors ${
                selectedOption
                  ? "border-sage bg-sage text-white"
                  : "border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
              }`}
            >
              {option.label} {count}
            </button>
          );
        })}
      </div>
    </div>
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
          className={`h-8 rounded-full border px-3 text-[12px] font-black transition-colors ${
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
              className={`h-8 max-w-[180px] rounded-full border px-3 text-[12px] font-black transition-colors ${
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

  // action 完了で state 参照が変わったときだけ、controlled な title/type を追従させる
  // (useEffect は使わず、レンダー中に前回値と比較して同期する)。
  const [title, setTitle] = useState(values.title);
  const [type, setType] = useState(values.type);
  const [dueDate, setDueDate] = useState(values.dueDate);
  const [syncedState, setSyncedState] = useState(state);
  if (syncedState !== state) {
    setSyncedState(state);
    setTitle(values.title);
    setType(values.type);
    setDueDate(values.dueDate);
  }

  const applyTemplate = (template: TaskTemplate) => {
    setTitle(template.title);
    setType(template.type);
  };

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
          <p className="mt-0.5 text-[12px] leading-relaxed text-ink-3">
            応募先、内容、期日だけ入れます。
          </p>
        </div>
      </div>

      {entries.length === 0 ? (
        <div className="rounded-lg border border-dashed border-line bg-cream px-3 py-3 text-center">
          <p className="text-[12px] font-bold text-ink-2">先に応募先を登録してください</p>
          <p className="mt-1 text-[12px] leading-relaxed text-ink-3">
            タスクはどの企業の予定かを紐づけて管理します。
          </p>
          <Link
            href="/entry/new"
            prefetch={false}
            className="mt-3 inline-flex items-center gap-1.5 rounded-lg bg-sage px-3 py-1.5 text-[12px] font-bold text-white transition-transform hover:-translate-y-0.5"
          >
            <Plus size={13} aria-hidden />
            応募先を追加
          </Link>
        </div>
      ) : (
        <>
          <div className="mb-2">
            <span className="mb-1 block text-[12px] font-bold text-ink-2">よく使うタスク</span>
            <TaskTemplateChips onSelect={applyTemplate} />
          </div>

          <div className="grid gap-2 md:grid-cols-[1.2fr_1.3fr]">
            <label className="block">
              <span className="mb-1 block text-[12px] font-bold text-ink-2">応募先</span>
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
              <span className="mb-1 block text-[12px] font-bold text-ink-2">内容</span>
              <input
                name="title"
                aria-label="タスク名"
                value={title}
                onChange={(event) => setTitle(event.target.value)}
                required
                placeholder="例: ES提出"
                className="h-9 w-full rounded-md border border-line bg-cream px-2.5 text-[12px] font-semibold outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
              />
            </label>
          </div>

          <div className="mt-2 grid gap-2 md:grid-cols-[1fr_1fr]">
            <fieldset>
              <legend className="mb-1 block text-[12px] font-bold text-ink-2">種類</legend>
              <div className="flex gap-1.5">
                {[
                  ["deadline", "締切"],
                  ["schedule", "予定"],
                ].map(([value, label]) => (
                  <label
                    key={value}
                    className="flex-1 cursor-pointer rounded-md border border-line bg-cream px-2 py-1.5 text-center text-[12px] font-bold text-ink-2 transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage has-[:checked]:text-white"
                  >
                    <input
                      type="radio"
                      name="type"
                      value={value}
                      checked={type === value}
                      onChange={() => setType(value as "deadline" | "schedule")}
                      className="sr-only"
                    />
                    {label}
                  </label>
                ))}
              </div>
            </fieldset>
            <label className="block">
              <span className="mb-1 block text-[12px] font-bold text-ink-2">期日</span>
              <input
                name="dueDate"
                type="date"
                aria-label="期日"
                value={dueDate}
                onChange={(event) => setDueDate(event.target.value)}
                className="h-9 w-full rounded-md border border-line bg-cream px-2.5 font-mono text-[12px] font-semibold outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
              />
              <div className="mt-1.5">
                <DueDateQuickPicker onSelect={setDueDate} />
              </div>
            </label>
          </div>

          <label className="mt-2 block">
            <span className="mb-1 block text-[12px] font-bold text-ink-2">メモ</span>
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
              className="mt-2 rounded-md bg-pink/40 px-2.5 py-1.5 text-[12px] font-semibold text-ink"
            >
              {state.error}
            </p>
          )}
          {state.ok && (
            <p className="mt-2 rounded-md bg-sage-wash px-2.5 py-1.5 text-[12px] font-bold text-sage">
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
      className="inline-flex items-center gap-1.5 rounded-lg bg-sage px-3.5 py-2 text-[12px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40 disabled:opacity-60"
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
      <p className="mx-auto mt-1 max-w-[420px] text-[12px] leading-relaxed text-ink-2">
        {hasEntries
          ? "上のフォームから締切や面接予定を追加できます。"
          : "応募先を登録すると、その企業に紐づく締切や面接予定を追加できます。"}
      </p>
      {!hasEntries && (
        <div className="mt-4 flex flex-wrap justify-center gap-2">
          <Link
            href="/entry/new"
            prefetch={false}
            className="inline-flex items-center gap-1.5 rounded-lg border border-line bg-surface px-3 py-1.5 text-[12px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
          >
            <ClipboardList size={13} aria-hidden />
            応募先を追加
          </Link>
        </div>
      )}
    </div>
  );
}
