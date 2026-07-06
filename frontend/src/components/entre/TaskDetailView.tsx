"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, useTransition, type ReactNode } from "react";
import {
  ArrowLeft,
  Building2,
  CalendarClock,
  CheckCircle2,
  ClipboardList,
  Pencil,
  Trash2,
  X,
} from "lucide-react";
import {
  deleteTaskAction,
  setTaskStatusAction,
  updateTaskAction,
} from "@/app/task/actions";
import { DueDateQuickPicker } from "@/components/entre/DueDateQuickPicker";
import { companyDisplayName, type EntryResponse } from "@/lib/api/entries";
import type { TaskWithEntry } from "@/lib/api/server-resources";

interface Props {
  task: TaskWithEntry;
  entry: EntryResponse | null;
}

const TYPE_LABEL: Record<TaskWithEntry["type"], string> = {
  deadline: "締切",
  schedule: "予定",
};

const STATUS_LABEL: Record<TaskWithEntry["status"], string> = {
  todo: "未完了",
  done: "完了",
};

const ISO_DATE_ONLY = /^(\d{4})-(\d{2})-(\d{2})(?:T|$)/;

export function TaskDetailView({ task, entry }: Props) {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [status, setStatus] = useState<TaskWithEntry["status"]>(task.status);
  const [title, setTitle] = useState(task.title);
  const [taskType, setTaskType] = useState<TaskWithEntry["type"]>(task.type);
  const [dueDate, setDueDate] = useState<TaskWithEntry["dueDate"]>(task.dueDate);
  const [memo, setMemo] = useState(task.memo);
  const [isEditing, setIsEditing] = useState(false);
  const [draftTitle, setDraftTitle] = useState(title);
  const [draftType, setDraftType] = useState<TaskWithEntry["type"]>(taskType);
  const [draftDueDate, setDraftDueDate] = useState(toDateInputValue(dueDate));
  const [draftMemo, setDraftMemo] = useState(memo);
  const [isPending, startTransition] = useTransition();
  const done = status === "done";
  const companyName = companyDisplayName({
    companyName: task.companyName ?? entry?.companyName,
  });

  const toggleStatus = () => {
    const next = done ? "todo" : "done";
    const previous = status;
    setError(null);
    setStatus(next);

    startTransition(async () => {
      const result = await setTaskStatusAction(task.id, next, task.entryId);
      if (!result.ok) {
        setStatus(previous);
        setError(result.error ?? "タスクの更新に失敗しました");
      }
    });
  };

  const deleteTask = () => {
    if (!window.confirm(`「${title}」を削除しますか？`)) return;
    setError(null);
    startTransition(async () => {
      const result = await deleteTaskAction(task.id, task.entryId);
      if (!result.ok) {
        setError(result.error ?? "タスクの削除に失敗しました");
        return;
      }
      router.push("/task");
    });
  };

  const startEdit = () => {
    setError(null);
    setDraftTitle(title);
    setDraftType(taskType);
    setDraftDueDate(toDateInputValue(dueDate));
    setDraftMemo(memo);
    setIsEditing(true);
  };

  const cancelEdit = () => {
    setError(null);
    setIsEditing(false);
  };

  const saveEdit = () => {
    const trimmedTitle = draftTitle.trim();
    if (!trimmedTitle) {
      setError("タスク名は必須です");
      return;
    }
    setError(null);

    startTransition(async () => {
      const result = await updateTaskAction(
        task.id,
        {
          title: trimmedTitle,
          type: draftType,
          dueDate: draftDueDate || null,
          memo: draftMemo,
        },
        task.entryId,
      );
      if (!result.ok) {
        setError(result.error ?? "タスクの更新に失敗しました");
        return;
      }
      setTitle(trimmedTitle);
      setTaskType(draftType);
      setDueDate(result.task?.dueDate ?? (draftDueDate ? `${draftDueDate}T00:00:00.000Z` : null));
      setMemo(draftMemo);
      setIsEditing(false);
    });
  };

  return (
    <div className="relative">
      <Link
        href="/task"
        prefetch={false}
        className="mb-3 inline-flex items-center gap-1 text-[12px] font-semibold text-ink-3 hover:text-sage"
      >
        <ArrowLeft size={13} aria-hidden />
        タスク一覧
      </Link>

      <section className="rounded-xl border border-line bg-surface p-4 shadow-card">
        <div className="mb-4 flex items-start justify-between gap-3">
          <div className="min-w-0 flex-1">
            {!isEditing && (
              <>
                <div className="mb-2 flex flex-wrap items-center gap-1.5">
                  <span className="rounded-md bg-cream px-2 py-0.5 text-[12px] font-black text-ink-3">
                    {TYPE_LABEL[taskType] ?? taskType}
                  </span>
                  <span
                    className={`rounded-md px-2 py-0.5 text-[12px] font-black ${
                      done ? "bg-sage-wash text-sage" : "bg-pink/35 text-ink"
                    }`}
                  >
                    {STATUS_LABEL[status] ?? status}
                  </span>
                </div>
                <h1 className="font-serif text-xl font-extrabold tracking-tight break-words">
                  {title}
                </h1>
              </>
            )}
            {isEditing && (
              <p className="text-[12px] font-black text-ink-3">タスクを編集</p>
            )}
          </div>
          <div className="flex shrink-0 items-center gap-2">
            {!isEditing && (
              <button
                type="button"
                onClick={startEdit}
                disabled={isPending}
                className="inline-flex items-center gap-1.5 rounded-lg border border-line bg-surface px-3 py-2 text-[12px] font-bold text-ink-3 transition-colors hover:border-sage hover:text-sage focus:outline-none focus:ring-2 focus:ring-sage/30 disabled:opacity-60"
              >
                <Pencil size={13} aria-hidden />
                編集
              </button>
            )}
            <button
              type="button"
              onClick={toggleStatus}
              disabled={isPending || isEditing}
              aria-pressed={done}
              className={`inline-flex items-center gap-1.5 rounded-lg px-3 py-2 text-[12px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40 disabled:opacity-60 ${
                done ? "bg-ink-3" : "bg-sage"
              }`}
            >
              <CheckCircle2 size={14} aria-hidden />
              {done ? "未完了に戻す" : "完了にする"}
            </button>
          </div>
        </div>

        {error && (
          <p
            role="alert"
            className="mb-3 rounded-md bg-pink/40 px-3 py-2 text-[12px] font-semibold text-ink"
          >
            {error}
          </p>
        )}

        {isEditing ? (
          <div className="grid gap-3 rounded-lg border border-line bg-cream p-3">
            <label className="block">
              <span className="mb-1 block text-[12px] font-bold text-ink-3">
                タスク名
              </span>
              <input
                type="text"
                value={draftTitle}
                onChange={(e) => setDraftTitle(e.target.value)}
                disabled={isPending}
                className="w-full rounded-md border border-line bg-surface px-2.5 py-1.5 text-[13px] text-ink focus:outline-none focus:ring-2 focus:ring-sage/30"
              />
            </label>

            <div>
              <span className="mb-1 block text-[12px] font-bold text-ink-3">
                種類
              </span>
              <div className="flex gap-1.5">
                {(Object.keys(TYPE_LABEL) as TaskWithEntry["type"][]).map((t) => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => setDraftType(t)}
                    disabled={isPending}
                    aria-pressed={draftType === t}
                    className={`rounded-md px-2.5 py-1 text-[12px] font-bold transition-colors disabled:opacity-60 ${
                      draftType === t
                        ? "bg-sage text-white"
                        : "border border-line bg-surface text-ink-3 hover:border-sage hover:text-sage"
                    }`}
                  >
                    {TYPE_LABEL[t]}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <span className="mb-1 block text-[12px] font-bold text-ink-3">
                期日
              </span>
              <div className="flex flex-wrap items-center gap-2">
                <input
                  type="date"
                  value={draftDueDate}
                  onChange={(e) => setDraftDueDate(e.target.value)}
                  disabled={isPending}
                  className="rounded-md border border-line bg-surface px-2.5 py-1.5 text-[13px] text-ink focus:outline-none focus:ring-2 focus:ring-sage/30"
                />
                {draftDueDate && (
                  <button
                    type="button"
                    onClick={() => setDraftDueDate("")}
                    disabled={isPending}
                    className="text-[12px] font-bold text-ink-3 hover:text-pink-deep disabled:opacity-60"
                  >
                    期日をクリア
                  </button>
                )}
              </div>
              <div className="mt-1.5">
                <DueDateQuickPicker onSelect={setDraftDueDate} />
              </div>
            </div>

            <label className="block">
              <span className="mb-1 block text-[12px] font-bold text-ink-3">
                メモ
              </span>
              <textarea
                value={draftMemo}
                onChange={(e) => setDraftMemo(e.target.value)}
                disabled={isPending}
                rows={3}
                className="w-full resize-y rounded-md border border-line bg-surface px-2.5 py-1.5 text-[13px] text-ink focus:outline-none focus:ring-2 focus:ring-sage/30"
              />
            </label>

            <div className="flex justify-end gap-2">
              <button
                type="button"
                onClick={cancelEdit}
                disabled={isPending}
                className="inline-flex items-center gap-1 rounded-lg border border-line bg-surface px-3 py-2 text-[12px] font-bold text-ink-3 hover:text-ink disabled:opacity-60"
              >
                <X size={13} aria-hidden />
                キャンセル
              </button>
              <button
                type="button"
                onClick={saveEdit}
                disabled={isPending}
                className="inline-flex items-center gap-1 rounded-lg bg-sage px-3 py-2 text-[12px] font-bold text-white enabled:hover:-translate-y-0.5 disabled:opacity-60"
              >
                保存
              </button>
            </div>
          </div>
        ) : (
          <>
            <div className="grid gap-2 md:grid-cols-2">
              <DetailCard
                icon={<Building2 size={15} aria-hidden />}
                label="応募先"
                value={companyName}
                href={entry ? `/entry/${entry.id}` : "/entry"}
              />
              <DetailCard
                icon={<CalendarClock size={15} aria-hidden />}
                label="期日"
                value={dueDate ? formatDate(dueDate) : "期日なし"}
              />
            </div>

            <div className="mt-3 rounded-lg border border-line bg-cream p-3">
              <div className="mb-1.5 flex items-center gap-1.5 text-[12px] font-bold text-ink-3">
                <ClipboardList size={13} aria-hidden />
                メモ
              </div>
              <p className="whitespace-pre-wrap break-words text-[12px] leading-relaxed text-ink-2">
                {memo || "メモはありません。"}
              </p>
            </div>
          </>
        )}

        <dl className="mt-3 grid gap-2 rounded-lg border border-dashed border-line bg-cream px-3 py-2 text-[12px] text-ink-3 md:grid-cols-2">
          <div>
            <dt className="font-bold">作成</dt>
            <dd className="mt-0.5 font-mono">{formatDateTime(task.createdAt)}</dd>
          </div>
          <div>
            <dt className="font-bold">更新</dt>
            <dd className="mt-0.5 font-mono">{formatDateTime(task.updatedAt)}</dd>
          </div>
        </dl>

        <div className="mt-4 flex justify-end">
          <button
            type="button"
            onClick={deleteTask}
            disabled={isPending || isEditing}
            className="inline-flex items-center gap-1.5 rounded-lg border border-line bg-surface px-3 py-2 text-[12px] font-bold text-ink-3 transition-colors hover:border-pink-deep hover:text-pink-deep focus:outline-none focus:ring-2 focus:ring-pink-deep/30 disabled:opacity-60"
          >
            <Trash2 size={13} aria-hidden />
            削除
          </button>
        </div>
      </section>
    </div>
  );
}

function DetailCard({
  icon,
  label,
  value,
  href,
}: {
  icon: ReactNode;
  label: string;
  value: string;
  href?: string;
}) {
  const body = (
    <>
      <div className="mb-1.5 flex items-center gap-1.5 text-[12px] font-bold text-ink-3">
        {icon}
        {label}
      </div>
      <div className="truncate text-[12px] font-bold text-ink-2">{value}</div>
    </>
  );

  if (href) {
    return (
      <Link
        href={href}
        prefetch={false}
        className="block rounded-lg border border-line bg-cream p-3 transition-colors hover:border-sage hover:text-sage focus:outline-none focus:ring-2 focus:ring-sage/30"
      >
        {body}
      </Link>
    );
  }

  return <div className="rounded-lg border border-line bg-cream p-3">{body}</div>;
}

function formatDateTime(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("ja-JP", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

function formatDate(value: string): string {
  const dateOnly = formatIsoDateOnly(value);
  if (dateOnly) return dateOnly;

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("ja-JP", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    timeZone: "UTC",
  }).format(date);
}

// date input (YYYY-MM-DD) 用に ISO 文字列を変換する。formatIsoDateOnly と違い区切りは "-"。
function toDateInputValue(value: string | null): string {
  if (!value) return "";
  const match = ISO_DATE_ONLY.exec(value);
  if (!match) return "";
  const [, year, month, day] = match;
  return `${year}-${month}-${day}`;
}

function formatIsoDateOnly(value: string): string | null {
  const match = ISO_DATE_ONLY.exec(value);
  if (!match) return null;

  const [, year, month, day] = match;
  const parsed = new Date(Date.UTC(Number(year), Number(month) - 1, Number(day)));
  if (
    parsed.getUTCFullYear() !== Number(year) ||
    parsed.getUTCMonth() + 1 !== Number(month) ||
    parsed.getUTCDate() !== Number(day)
  ) {
    return null;
  }

  return `${year}/${month}/${day}`;
}
