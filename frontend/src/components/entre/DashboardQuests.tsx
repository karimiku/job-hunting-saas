// Server Component。実タスクから今日のタスクリストを組み立てる。
// タスクは entry 単位の API しか無いため、page 側で集約した TaskWithEntry[] を受け取る。
// useEffect は使わない（純粋な集計 + 表示のみ）。

import Link from "next/link";
import type { TaskWithEntry } from "@/lib/api/server-resources";
import { companyDisplayName } from "@/lib/api/entries";

export interface QuestItem {
  id: string;
  /** 期限ラベル（例: "5/30" / "期限なし"）。 */
  due: string;
  /** メインのタスク見出し（会社名 + タイトル）。 */
  label: string;
  /** サブテキスト（メモ）。 */
  sub: string;
  /** 期限の緊急度で決まるバッジ色クラス。 */
  color: string;
  done: boolean;
}

const MAX_QUESTS = 5;

function dueLabel(dueDate: string | null): string {
  if (!dueDate) return "期限なし";
  const d = new Date(dueDate);
  if (Number.isNaN(d.getTime())) return "期限なし";
  return `${d.getMonth() + 1}/${d.getDate()}`;
}

// 期限の近さで色を決める。過ぎている/今日=pink、3日以内=amber、それ以降=sky、期限なし=sage。
function dueColor(dueDate: string | null, now: Date): string {
  if (!dueDate) return "bg-sage";
  const d = new Date(dueDate);
  if (Number.isNaN(d.getTime())) return "bg-sage";
  const days = Math.floor((d.getTime() - now.getTime()) / 86_400_000);
  if (days <= 0) return "bg-pink";
  if (days <= 3) return "bg-amber";
  return "bg-sky";
}

// done を後ろ、未完了は期限昇順（期限なしは末尾）に並べて先頭 MAX_QUESTS 件を返す。
export function buildQuests(
  tasks: TaskWithEntry[],
  now: Date = new Date(),
): QuestItem[] {
  const sorted = [...tasks].sort((a, b) => {
    const aDone = a.status === "done" ? 1 : 0;
    const bDone = b.status === "done" ? 1 : 0;
    if (aDone !== bDone) return aDone - bDone;
    const at = a.dueDate ? new Date(a.dueDate).getTime() : Number.POSITIVE_INFINITY;
    const bt = b.dueDate ? new Date(b.dueDate).getTime() : Number.POSITIVE_INFINITY;
    return at - bt;
  });

  return sorted.slice(0, MAX_QUESTS).map((t) => {
    const company = companyDisplayName({ companyName: t.companyName });
    return {
      id: t.id,
      due: dueLabel(t.dueDate),
      label: `${company} ${t.title}`.trim(),
      sub: t.memo?.trim() || (t.type === "deadline" ? "締切タスク" : "予定"),
      color: dueColor(t.dueDate, now),
      done: t.status === "done",
    };
  });
}

export function questProgress(tasks: TaskWithEntry[]): number {
  if (tasks.length === 0) return 0;
  const done = tasks.filter((t) => t.status === "done").length;
  return Math.round((done / tasks.length) * 100);
}

/** 実タスクから組み立てた「今日のタスク」カード。 */
export function DashboardQuests({ tasks }: { tasks: TaskWithEntry[] }) {
  const quests = buildQuests(tasks);

  return (
    <div className="rounded-xl border border-line bg-surface p-5">
      <div className="mb-3 flex items-baseline justify-between">
        <h2 className="text-[14px] font-extrabold">今日のタスク</h2>
        <Link href="/task" className="text-[10px] font-bold text-sage">
          一覧
        </Link>
      </div>

      {quests.length === 0 ? (
        <div data-testid="quest-empty" className="py-6 text-center">
          <p className="text-[12px] font-bold text-ink-2">タスクはまだありません</p>
          <p className="mt-1 text-[11px] text-ink-3">
            Entryごとに締切や予定を追加すると、近い順に表示されます。
          </p>
          <Link
            href="/entry"
            className="mt-3 inline-flex rounded-lg border border-line bg-surface px-3 py-1.5 text-[11px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
          >
            Entryを確認
          </Link>
        </div>
      ) : (
        <ul>
          {quests.map((r, i) => (
            <li
              key={r.id}
              className={`flex items-center gap-3 py-2.5 ${
                i ? "border-t border-dashed border-line" : ""
              } ${r.done ? "opacity-50" : ""}`}
            >
              <span
                className={`grid h-[18px] w-[18px] place-items-center rounded-full text-[10px] text-white ${
                  r.done ? "border-[1.5px] border-sage bg-sage" : "border-[1.5px] border-line bg-transparent"
                }`}
              >
                {r.done ? "✓" : ""}
              </span>
              <div className="min-w-0 flex-1">
                <div className={`truncate text-xs font-semibold ${r.done ? "line-through" : ""}`}>
                  {r.label}
                </div>
                <div className="mt-0.5 truncate text-[10px] text-ink-3">{r.sub}</div>
              </div>
              <span
                className={`shrink-0 rounded-md px-2 py-0.5 font-mono text-[10px] font-bold text-white ${r.color}`}
              >
                {r.due}
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
