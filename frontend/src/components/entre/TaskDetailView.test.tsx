import { beforeEach, describe, expect, it, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TaskDetailView } from "./TaskDetailView";
import type { EntryResponse } from "@/lib/api/entries";
import type { TaskWithEntry } from "@/lib/api/server-resources";

const { push, setTaskStatusAction, deleteTaskAction, updateTaskAction } =
  vi.hoisted(() => ({
    push: vi.fn(),
    setTaskStatusAction: vi.fn(),
    deleteTaskAction: vi.fn(),
    updateTaskAction: vi.fn(),
  }));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push }),
}));

vi.mock("@/app/task/actions", () => ({
  setTaskStatusAction: (
    taskId: string,
    status: "todo" | "done",
    entryId?: string,
  ) => setTaskStatusAction(taskId, status, entryId),
  deleteTaskAction: (taskId: string, entryId?: string) =>
    deleteTaskAction(taskId, entryId),
  updateTaskAction: (
    taskId: string,
    input: Record<string, unknown>,
    entryId?: string,
  ) => updateTaskAction(taskId, input, entryId),
}));

const task = (overrides: Partial<TaskWithEntry> = {}): TaskWithEntry => ({
  id: "t1",
  entryId: "e1",
  title: "ES提出",
  type: "deadline",
  status: "todo",
  dueDate: "2026-06-01T00:00:00.000Z",
  memo: "提出前に誤字を見る",
  companyName: "テスト商事",
  createdAt: "2026-05-01T09:00:00.000Z",
  updatedAt: "2026-05-02T10:00:00.000Z",
  ...overrides,
});

const entry = (overrides: Partial<EntryResponse> = {}): EntryResponse => ({
  id: "e1",
  companyId: "c1",
  companyName: "テスト商事",
  route: "本選考",
  source: "マイナビ",
  status: "in_progress",
  stageKind: "document",
  stageLabel: "書類選考",
  memo: "",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

describe("TaskDetailView", () => {
  beforeEach(() => {
    push.mockReset();
    setTaskStatusAction.mockReset();
    deleteTaskAction.mockReset();
    updateTaskAction.mockReset();
    setTaskStatusAction.mockResolvedValue({ ok: true, status: "done" });
    deleteTaskAction.mockResolvedValue({ ok: true });
    updateTaskAction.mockResolvedValue({
      ok: true,
      task: { dueDate: "2026-06-10T00:00:00.000Z" },
    });
  });

  it("タスク詳細と紐づくEntry導線を表示する", () => {
    render(<TaskDetailView task={task()} entry={entry()} />);

    expect(screen.getByRole("heading", { name: "ES提出" })).toBeInTheDocument();
    expect(screen.getByText("締切")).toBeInTheDocument();
    expect(screen.getByText("未完了")).toBeInTheDocument();
    expect(screen.getByText("提出前に誤字を見る")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /応募先/ })).toHaveAttribute(
      "href",
      "/entry/e1",
    );
  });

  it("UTC 0時の期日をローカルタイムゾーンで前日にずらさず表示する", () => {
    const originalTZ = process.env.TZ;
    process.env.TZ = "America/Los_Angeles";

    try {
      render(
        <TaskDetailView
          task={task({ dueDate: "2026-06-01T00:00:00.000Z" })}
          entry={entry()}
        />,
      );

      expect(screen.getByText("2026/06/01")).toBeInTheDocument();
      expect(screen.queryByText("2026/05/31")).not.toBeInTheDocument();
    } finally {
      if (originalTZ === undefined) {
        delete process.env.TZ;
      } else {
        process.env.TZ = originalTZ;
      }
    }
  });

  it("完了切替で Server Action を呼び、表示を更新する", async () => {
    render(<TaskDetailView task={task()} entry={entry()} />);

    await userEvent.click(screen.getByRole("button", { name: "完了にする" }));

    await waitFor(() =>
      expect(setTaskStatusAction).toHaveBeenCalledWith("t1", "done", "e1"),
    );
    expect(screen.getByText("完了")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "未完了に戻す" })).toHaveAttribute(
      "aria-pressed",
      "true",
    );
  });

  it("完了切替が失敗したら状態を戻してエラーを表示する", async () => {
    setTaskStatusAction.mockResolvedValue({
      ok: false,
      error: "タスクの更新に失敗しました",
    });
    render(<TaskDetailView task={task()} entry={entry()} />);

    await userEvent.click(screen.getByRole("button", { name: "完了にする" }));

    await waitFor(() => expect(screen.getByRole("alert")).toBeInTheDocument());
    expect(screen.getByText("未完了")).toBeInTheDocument();
  });

  it("削除成功時は一覧へ戻る", async () => {
    vi.spyOn(window, "confirm").mockReturnValue(true);
    render(<TaskDetailView task={task()} entry={entry()} />);

    await userEvent.click(screen.getByRole("button", { name: "削除" }));

    await waitFor(() =>
      expect(deleteTaskAction).toHaveBeenCalledWith("t1", "e1"),
    );
    expect(push).toHaveBeenCalledWith("/task");
  });

  it("編集ボタンでフォームに切り替わり、保存で updateTaskAction を呼んで表示に戻る", async () => {
    render(<TaskDetailView task={task()} entry={entry()} />);

    await userEvent.click(screen.getByRole("button", { name: "編集" }));

    const titleInput = screen.getByLabelText("タスク名");
    await userEvent.clear(titleInput);
    await userEvent.type(titleInput, "ES提出（修正）");

    const memoInput = screen.getByLabelText("メモ");
    await userEvent.clear(memoInput);
    await userEvent.type(memoInput, "誤字を直した");

    await userEvent.click(screen.getByRole("button", { name: "予定" }));
    await userEvent.click(screen.getByRole("button", { name: "保存" }));

    await waitFor(() =>
      expect(updateTaskAction).toHaveBeenCalledWith(
        "t1",
        {
          title: "ES提出（修正）",
          type: "schedule",
          dueDate: "2026-06-01",
          memo: "誤字を直した",
        },
        "e1",
      ),
    );

    expect(
      screen.getByRole("heading", { name: "ES提出（修正）" }),
    ).toBeInTheDocument();
    expect(screen.getByText("誤字を直した")).toBeInTheDocument();
    expect(screen.getByText("予定")).toBeInTheDocument();
    expect(screen.queryByLabelText("タスク名")).not.toBeInTheDocument();
  });

  it("編集フォームでタスク名を空にすると保存できずエラーを表示する", async () => {
    render(<TaskDetailView task={task()} entry={entry()} />);

    await userEvent.click(screen.getByRole("button", { name: "編集" }));
    await userEvent.clear(screen.getByLabelText("タスク名"));
    await userEvent.click(screen.getByRole("button", { name: "保存" }));

    expect(screen.getByRole("alert")).toHaveTextContent("タスク名は必須です");
    expect(updateTaskAction).not.toHaveBeenCalled();
  });

  it("編集フォームのキャンセルで変更を破棄して表示モードに戻る", async () => {
    render(<TaskDetailView task={task()} entry={entry()} />);

    await userEvent.click(screen.getByRole("button", { name: "編集" }));
    await userEvent.clear(screen.getByLabelText("タスク名"));
    await userEvent.type(screen.getByLabelText("タスク名"), "破棄されるタイトル");
    await userEvent.click(screen.getByRole("button", { name: "キャンセル" }));

    expect(screen.getByRole("heading", { name: "ES提出" })).toBeInTheDocument();
    expect(updateTaskAction).not.toHaveBeenCalled();
  });
});
