import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TaskListView } from "./TaskListView";
import type { TaskWithEntry } from "@/lib/api/server-resources";

// Server Action は next/headers を読むのでここでは mock し、呼び出し引数だけ検証する。
interface ActionResult {
  ok: boolean;
  status?: "todo" | "done";
  error?: string;
}
const setTaskStatusAction = vi.fn(
  async (
    _taskId: string,
    status: "todo" | "done",
  ): Promise<ActionResult> => ({ ok: true, status }),
);
vi.mock("@/app/task/actions", () => ({
  setTaskStatusAction: (taskId: string, status: "todo" | "done") =>
    setTaskStatusAction(taskId, status),
}));

// Confetti をモックして、毎レンダーの trigger 値を記録する。
// 「Server Action 成功後にだけ祝福する」挙動 (trigger>0) を検証するため。
const { confettiSpy } = vi.hoisted(() => ({ confettiSpy: vi.fn() }));
vi.mock("./Confetti", () => ({
  Confetti: ({ trigger }: { trigger: number }) => {
    confettiSpy(trigger);
    return null;
  },
}));

// confetti が一度でも発火 (trigger>0) したか。
const confettiFired = () =>
  confettiSpy.mock.calls.some(([t]) => (t as number) > 0);

const task = (overrides: Partial<TaskWithEntry> = {}): TaskWithEntry => ({
  id: "t1",
  entryId: "e1",
  title: "ES提出",
  type: "deadline",
  status: "todo",
  dueDate: "2026-05-30",
  memo: "最終チェック",
  companyName: "○○商事",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

describe("TaskListView", () => {
  beforeEach(() => {
    setTaskStatusAction.mockClear();
    confettiSpy.mockClear();
    setTaskStatusAction.mockImplementation(
      async (_taskId: string, status: "todo" | "done") => ({ ok: true, status }),
    );
  });

  it("タスクが無いとき空状態を表示する", () => {
    render(<TaskListView initialTasks={[]} />);
    expect(screen.getByText(/まだタスクがありません/)).toBeInTheDocument();
  });

  it("タスクのタイトル・会社名・期日を表示する", () => {
    render(<TaskListView initialTasks={[task()]} />);
    expect(screen.getByText("ES提出")).toBeInTheDocument();
    expect(screen.getByText(/○○商事/)).toBeInTheDocument();
    expect(screen.getByText("5/30")).toBeInTheDocument();
  });

  it("期日が無いタスクは「期日なし」と表示する", () => {
    render(<TaskListView initialTasks={[task({ dueDate: null })]} />);
    expect(screen.getByText("期日なし")).toBeInTheDocument();
  });

  it("未完了タスクをトグルすると status=done で Server Action を呼び、成功後に祝福する", async () => {
    render(<TaskListView initialTasks={[task({ status: "todo" })]} />);
    const toggle = screen.getByRole("button", { name: "タスク完了にする" });
    await userEvent.click(toggle);

    await waitFor(() =>
      expect(setTaskStatusAction).toHaveBeenCalledWith("t1", "done"),
    );
    // 楽観更新で aria-pressed が true になる
    await waitFor(() =>
      expect(
        screen.getByRole("button", { name: "タスク未完了に戻す" }),
      ).toHaveAttribute("aria-pressed", "true"),
    );
    // Server Action 成功後に紙吹雪が発火する
    await waitFor(() => expect(confettiFired()).toBe(true));
  });

  it("完了タスクをトグルすると status=todo で Server Action を呼ぶ", async () => {
    render(<TaskListView initialTasks={[task({ status: "done" })]} />);
    const toggle = screen.getByRole("button", { name: "タスク未完了に戻す" });
    await userEvent.click(toggle);

    await waitFor(() =>
      expect(setTaskStatusAction).toHaveBeenCalledWith("t1", "todo"),
    );
  });

  it("Server Action が失敗したら楽観更新を巻き戻しエラーを表示し、祝福しない", async () => {
    setTaskStatusAction.mockImplementation(async () => ({
      ok: false,
      error: "タスクの更新に失敗しました",
    }));
    render(<TaskListView initialTasks={[task({ status: "todo" })]} />);

    await userEvent.click(screen.getByRole("button", { name: "タスク完了にする" }));

    await waitFor(() => expect(screen.getByRole("alert")).toBeInTheDocument());
    // 巻き戻されて未完了 (aria-pressed=false) に戻る
    expect(
      screen.getByRole("button", { name: "タスク完了にする" }),
    ).toHaveAttribute("aria-pressed", "false");
    // 失敗時は紙吹雪を出さない (祝福→エラーの壊れた UX を防ぐ)
    expect(confettiFired()).toBe(false);
  });
});
