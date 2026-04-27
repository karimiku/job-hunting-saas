"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { Confetti } from "@/components/entre/Confetti";

interface Task {
  id: string;
  t: string;
  e: string;
  s: string;
  due: string;
  co: string;
  done: boolean;
}

const INITIAL_TASKS: Task[] = [
  { id: "t1", t: "14:00", e: "○○商事 一次面接(Web)", s: "明日締切 ES確認", due: "5/29", co: "bg-pink", done: false },
  { id: "t2", t: "23:59", e: "△△株式会社 ESを提出", s: "最終チェック", due: "5/28", co: "bg-amber", done: false },
  { id: "t3", t: "本日中", e: "◇◇テック SPI受験", s: "45分", due: "5/30", co: "bg-sky", done: false },
  { id: "t4", t: "17:00", e: "OBの田中さんへDM返信", s: "メッセージ準備済", due: "5/27", co: "bg-sage", done: true },
];

export default function TaskPage() {
  const router = useRouter();
  const state = useUser();
  const [tasks, setTasks] = useState(INITIAL_TASKS);
  const [confetti, setConfetti] = useState(0);

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  const toggle = (id: string) => {
    let triggeredDone = false;
    setTasks((prev) =>
      prev.map((t) => {
        if (t.id === id) {
          if (!t.done) triggeredDone = true;
          return { ...t, done: !t.done };
        }
        return t;
      }),
    );
    if (triggeredDone) setConfetti((n) => n + 1);
  };

  return (
    <AppShell userName={state.user.name} userSubtitle="○○大学 4年">
      <div className="relative mx-auto max-w-[700px] px-5 py-6 md:px-8 md:py-7">
        <header className="mb-4">
          <h1 className="font-serif text-2xl font-extrabold tracking-tight">Task</h1>
          <p className="mt-0.5 text-[11px] text-ink-3">タスクや締切を1箇所で管理</p>
        </header>

        <ul className="flex flex-col gap-2">
          {tasks.map((task) => (
            <li
              key={task.id}
              className={`flex items-center gap-3 rounded-xl border border-line bg-surface px-3 py-2.5 transition-opacity ${
                task.done ? "opacity-50" : ""
              }`}
            >
              <button
                type="button"
                onClick={() => toggle(task.id)}
                aria-pressed={task.done}
                aria-label={task.done ? "タスク未完了に戻す" : "タスク完了にする"}
                className={`grid h-5 w-5 shrink-0 place-items-center rounded-full text-[11px] text-white transition-colors ${
                  task.done ? "border-[1.5px] border-sage bg-sage" : "border-[1.5px] border-line bg-transparent"
                }`}
              >
                {task.done ? "✓" : ""}
              </button>
              <div className="min-w-0 flex-1">
                <div className={`text-[12px] font-semibold ${task.done ? "line-through" : ""}`}>
                  {task.e}
                </div>
                <div className="mt-0.5 text-[10px] text-ink-3">{task.s}</div>
              </div>
              <span
                className={`shrink-0 rounded-md px-2 py-0.5 font-mono text-[10px] font-bold text-white ${task.co}`}
              >
                {task.t}
              </span>
            </li>
          ))}
        </ul>

        <Confetti trigger={confetti} count={22} />
      </div>
    </AppShell>
  );
}
