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
  Trash2,
} from "lucide-react";
import { deleteTaskAction, setTaskStatusAction } from "@/app/task/actions";
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
    if (!window.confirm(`「${task.title}」を削除しますか？`)) return;
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

  return (
    <div className="relative">
      <Link
        href="/task"
        prefetch={false}
        className="mb-3 inline-flex items-center gap-1 text-[11px] font-semibold text-ink-3 hover:text-sage"
      >
        <ArrowLeft size={13} aria-hidden />
        タスク一覧
      </Link>

      <section className="rounded-xl border border-line bg-surface p-4 shadow-card">
        <div className="mb-4 flex items-start justify-between gap-3">
          <div className="min-w-0 flex-1">
            <div className="mb-2 flex flex-wrap items-center gap-1.5">
              <span className="rounded-md bg-cream px-2 py-0.5 text-[10px] font-black text-ink-3">
                {TYPE_LABEL[task.type] ?? task.type}
              </span>
              <span
                className={`rounded-md px-2 py-0.5 text-[10px] font-black ${
                  done ? "bg-sage-wash text-sage" : "bg-pink/35 text-ink"
                }`}
              >
                {STATUS_LABEL[status] ?? status}
              </span>
            </div>
            <h1 className="font-serif text-xl font-extrabold tracking-tight break-words">
              {task.title}
            </h1>
          </div>
          <button
            type="button"
            onClick={toggleStatus}
            disabled={isPending}
            aria-pressed={done}
            className={`inline-flex shrink-0 items-center gap-1.5 rounded-lg px-3 py-2 text-[11px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40 disabled:opacity-60 ${
              done ? "bg-ink-3" : "bg-sage"
            }`}
          >
            <CheckCircle2 size={14} aria-hidden />
            {done ? "未完了に戻す" : "完了にする"}
          </button>
        </div>

        {error && (
          <p
            role="alert"
            className="mb-3 rounded-md bg-pink/40 px-3 py-2 text-[11px] font-semibold text-ink"
          >
            {error}
          </p>
        )}

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
            value={task.dueDate ? formatDate(task.dueDate) : "期日なし"}
          />
        </div>

        <div className="mt-3 rounded-lg border border-line bg-cream p-3">
          <div className="mb-1.5 flex items-center gap-1.5 text-[10px] font-bold text-ink-3">
            <ClipboardList size={13} aria-hidden />
            メモ
          </div>
          <p className="whitespace-pre-wrap break-words text-[12px] leading-relaxed text-ink-2">
            {task.memo || "メモはありません。"}
          </p>
        </div>

        <dl className="mt-3 grid gap-2 rounded-lg border border-dashed border-line bg-cream px-3 py-2 text-[10px] text-ink-3 md:grid-cols-2">
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
            disabled={isPending}
            className="inline-flex items-center gap-1.5 rounded-lg border border-line bg-surface px-3 py-2 text-[11px] font-bold text-ink-3 transition-colors hover:border-pink-deep hover:text-pink-deep focus:outline-none focus:ring-2 focus:ring-pink-deep/30 disabled:opacity-60"
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
      <div className="mb-1.5 flex items-center gap-1.5 text-[10px] font-bold text-ink-3">
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
