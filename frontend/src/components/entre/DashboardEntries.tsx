import Link from "next/link";
import { ArrowRight, ClipboardList } from "lucide-react";
import {
  companyDisplayName,
  type EntryResponse,
} from "@/lib/api/entries";
import type { TaskWithEntry } from "@/lib/api/server-resources";
import {
  ENTRY_STATUS_LABEL,
  STAGE_BG,
  STAGE_ORDER,
  stageIndexOf,
} from "@/lib/entry-stage";

export interface DashboardEntryItem {
  id: string;
  company: string;
  stageKind: string;
  stageLabel: string;
  status: string;
  openTaskCount: number;
  nearestDue: string;
}

const MAX_ITEMS = 5;

function dueTime(value: string | null): number {
  if (!value) return Number.POSITIVE_INFINITY;
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? Number.POSITIVE_INFINITY : date.getTime();
}

function dueLabel(value: string | null): string {
  if (!value) return "期日なし";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "期日なし";
  return `${date.getMonth() + 1}/${date.getDate()}`;
}

export function buildDashboardEntries(
  entries: EntryResponse[],
  tasks: TaskWithEntry[],
): DashboardEntryItem[] {
  const tasksByEntry = new Map<string, TaskWithEntry[]>();
  for (const task of tasks) {
    if (task.status === "done") continue;
    const current = tasksByEntry.get(task.entryId) ?? [];
    current.push(task);
    tasksByEntry.set(task.entryId, current);
  }

  return [...entries]
    .sort((a, b) => {
      const aClosed = a.status === "rejected" || a.status === "withdrawn" ? 1 : 0;
      const bClosed = b.status === "rejected" || b.status === "withdrawn" ? 1 : 0;
      if (aClosed !== bClosed) return aClosed - bClosed;

      const aTasks = tasksByEntry.get(a.id) ?? [];
      const bTasks = tasksByEntry.get(b.id) ?? [];
      const aDue = Math.min(...aTasks.map((task) => dueTime(task.dueDate)));
      const bDue = Math.min(...bTasks.map((task) => dueTime(task.dueDate)));
      if (aDue !== bDue) return aDue - bDue;

      return stageIndexOf(b.stageKind) - stageIndexOf(a.stageKind);
    })
    .slice(0, MAX_ITEMS)
    .map((entry) => {
      const openTasks = tasksByEntry.get(entry.id) ?? [];
      const nearest = openTasks.reduce<TaskWithEntry | null>((current, task) => {
        if (!current) return task;
        return dueTime(task.dueDate) < dueTime(current.dueDate) ? task : current;
      }, null);
      return {
        id: entry.id,
        company: companyDisplayName(entry),
        stageKind: entry.stageKind,
        stageLabel: entry.stageLabel,
        status: entry.status,
        openTaskCount: openTasks.length,
        nearestDue: dueLabel(nearest?.dueDate ?? null),
      };
    });
}

export function DashboardEntries({
  entries,
  tasks,
}: {
  entries: EntryResponse[];
  tasks: TaskWithEntry[];
}) {
  const items = buildDashboardEntries(entries, tasks);

  return (
    <section className="rounded-xl border border-line bg-surface p-4 md:p-5">
      <div className="mb-3 flex items-center justify-between gap-3">
        <div>
          <h2 className="text-[14px] font-extrabold">進行中のEntry</h2>
          <p className="mt-0.5 text-[10px] text-ink-3">
            次に確認する企業を、タスク期限とステージで並べます。
          </p>
        </div>
        <Link href="/entry" prefetch={false} className="shrink-0 text-[10px] font-bold text-sage">
          一覧
        </Link>
      </div>

      {items.length === 0 ? (
        <div className="rounded-lg border border-dashed border-line bg-cream px-3 py-5 text-center">
          <ClipboardList className="mx-auto mb-2 text-sage" size={18} aria-hidden />
          <p className="text-[12px] font-bold text-ink-2">Entryはまだありません</p>
          <Link
            href="/entry/new"
            prefetch={false}
            className="mt-3 inline-flex rounded-lg bg-sage px-3 py-1.5 text-[11px] font-bold text-white"
          >
            Entryを追加
          </Link>
        </div>
      ) : (
        <ul className="flex flex-col gap-2">
          {items.map((item) => (
            <li key={item.id}>
              <Link
                href={`/entry/${item.id}`}
                prefetch={false}
                className="block rounded-lg border border-line bg-cream px-3 py-2.5 transition-colors hover:border-sage"
              >
                <div className="flex items-start justify-between gap-2">
                  <div className="min-w-0">
                    <div className="truncate text-[12px] font-extrabold">{item.company}</div>
                    <div className="mt-1 flex flex-wrap items-center gap-1.5 text-[9px] font-bold text-ink-3">
                      <span
                        className="rounded-sm px-1.5 py-0.5 text-white"
                        style={{ background: STAGE_BG[item.stageKind] ?? "var(--color-ink-3)" }}
                      >
                        {item.stageLabel}
                      </span>
                      <span>{ENTRY_STATUS_LABEL[item.status] ?? item.status}</span>
                      <span aria-hidden>·</span>
                      <span>未完了 {item.openTaskCount}</span>
                      <span aria-hidden>·</span>
                      <span>{item.nearestDue}</span>
                    </div>
                  </div>
                  <ArrowRight size={13} className="mt-0.5 shrink-0 text-ink-3" aria-hidden />
                </div>
                <div
                  className="mt-2 grid gap-0.5"
                  style={{ gridTemplateColumns: `repeat(${STAGE_ORDER.length}, minmax(0, 1fr))` }}
                  aria-hidden
                >
                  {STAGE_ORDER.map((kind, index) => (
                    <span
                      key={kind}
                      className="h-1.5 rounded-full"
                      style={{
                        background:
                          index <= stageIndexOf(item.stageKind)
                            ? STAGE_BG[kind]
                            : "var(--color-line-2)",
                      }}
                    />
                  ))}
                </div>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
