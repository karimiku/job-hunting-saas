import { describe, it, expect, vi, beforeEach, beforeAll, afterAll } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { dueColor, dueLabel, sortTasksForDisplay, TaskListView } from "./TaskListView";
import type { TaskWithEntry } from "@/lib/api/server-resources";
import type { EntryResponse } from "@/lib/api/entries";

// Server Action は next/headers を読むのでここでは mock し、呼び出し引数だけ検証する。
interface ActionResult {
  ok: boolean;
  status?: "todo" | "done";
  error?: string;
}
const setTaskStatusAction = vi.fn(
  async (
    taskId: string,
    status: "todo" | "done",
    entryId?: string,
  ): Promise<ActionResult> => {
    void taskId;
    void entryId;
    return { ok: true, status };
  },
);
const deleteTaskAction = vi.fn(async (taskId: string, entryId?: string) => {
  void taskId;
  void entryId;
  return { ok: true };
});
const createTaskFromTaskPageAction = vi.fn(
  async (prev: unknown, formData: FormData) => {
    void prev;
    void formData;
    return { ok: true };
  },
);
vi.mock("@/app/task/actions", () => ({
  createTaskFromTaskPageAction: (_prev: unknown, formData: FormData) =>
    createTaskFromTaskPageAction(_prev, formData),
  deleteTaskAction: (taskId: string, entryId?: string) =>
    deleteTaskAction(taskId, entryId),
  setTaskStatusAction: (
    taskId: string,
    status: "todo" | "done",
    entryId?: string,
  ) => setTaskStatusAction(taskId, status, entryId),
}));

// Confetti の trigger 値を記録し、発火有無を検証する。
const { confettiSpy } = vi.hoisted(() => ({ confettiSpy: vi.fn() }));
vi.mock("./Confetti", () => ({
  Confetti: ({ trigger }: { trigger: number }) => {
    confettiSpy(trigger);
    return null;
  },
}));

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

const entry = (overrides: Partial<EntryResponse> = {}): EntryResponse => ({
  id: "e1",
  companyId: "c1",
  companyName: "○○商事",
  route: "本選考",
  source: "マイナビ",
  status: "in_progress",
  stageKind: "application",
  stageLabel: "応募",
  memo: "",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

describe("TaskListView", () => {
  beforeAll(() => {
    // Date だけを固定し、userEvent が使う実タイマーには影響させない。
    vi.useFakeTimers({ toFake: ["Date"] });
    vi.setSystemTime(new Date("2026-05-27T00:00:00Z"));
  });

  afterAll(() => {
    vi.useRealTimers();
  });

  beforeEach(() => {
    setTaskStatusAction.mockClear();
    deleteTaskAction.mockClear();
    createTaskFromTaskPageAction.mockClear();
    confettiSpy.mockClear();
    setTaskStatusAction.mockImplementation(
      async (_taskId: string, status: "todo" | "done", _entryId?: string) => {
        void _entryId;
        return { ok: true, status };
      },
    );
  });

  it("タスクが無いとき空状態を表示する", () => {
    render(<TaskListView initialTasks={[]} entries={[entry()]} />);
    expect(screen.getByText(/タスクはまだありません/)).toBeInTheDocument();
  });

  it("タスクのタイトル・会社名・期日を表示する", () => {
    render(<TaskListView initialTasks={[task()]} entries={[entry()]} />);
    expect(screen.getByRole("link", { name: /ES提出/ })).toBeInTheDocument();
    expect(screen.getAllByText(/○○商事/).length).toBeGreaterThan(0);
    expect(screen.getByText("5/30")).toBeInTheDocument();
  });

  it("タスク名からタスク詳細へ遷移できる", () => {
    render(<TaskListView initialTasks={[task()]} entries={[entry()]} />);
    expect(screen.getByRole("link", { name: /ES提出/ })).toHaveAttribute(
      "href",
      "/task/t1",
    );
  });

  it("Entryフィルタで対象Entryのタスクだけ表示する", async () => {
    const user = userEvent.setup();
    render(
      <TaskListView
        initialTasks={[
          task({ id: "t1", entryId: "e1", title: "ES提出", companyName: "○○商事" }),
          task({ id: "t2", entryId: "e2", title: "一次面接", companyName: "△△銀行" }),
        ]}
        entries={[
          entry({ id: "e1", companyName: "○○商事" }),
          entry({ id: "e2", companyName: "△△銀行" }),
        ]}
      />,
    );

    await user.click(screen.getByRole("button", { name: /△△銀行 1/ }));

    expect(screen.queryByRole("link", { name: /ES提出/ })).not.toBeInTheDocument();
    expect(screen.getByRole("link", { name: /一次面接/ })).toBeInTheDocument();
    expect(screen.getByText("未完了 1")).toBeInTheDocument();
  });

  it("期日が無いタスクは「期日なし」と表示する", () => {
    render(<TaskListView initialTasks={[task({ dueDate: null })]} entries={[entry()]} />);
    expect(screen.getByText("期日なし")).toBeInTheDocument();
  });

  it("未完了タスクをトグルすると status=done で Server Action を呼び、成功後に祝福する", async () => {
    render(<TaskListView initialTasks={[task({ status: "todo" })]} entries={[entry()]} />);
    const toggle = screen.getByRole("button", { name: "タスク完了にする" });
    await userEvent.click(toggle);

    await waitFor(() =>
      expect(setTaskStatusAction).toHaveBeenCalledWith("t1", "done", "e1"),
    );
    // 楽観更新で aria-pressed が true になる
    await waitFor(() =>
      expect(
        screen.getByRole("button", { name: "タスク未完了に戻す" }),
      ).toHaveAttribute("aria-pressed", "true"),
    );
    await waitFor(() => expect(confettiFired()).toBe(true));
  });

  it("完了タスクをトグルすると status=todo で Server Action を呼ぶ", async () => {
    render(<TaskListView initialTasks={[task({ status: "done" })]} entries={[entry()]} />);
    const toggle = screen.getByRole("button", { name: "タスク未完了に戻す" });
    await userEvent.click(toggle);

    await waitFor(() =>
      expect(setTaskStatusAction).toHaveBeenCalledWith("t1", "todo", "e1"),
    );
  });

  it("Server Action が失敗したら楽観更新を巻き戻しエラーを表示し、祝福しない", async () => {
    setTaskStatusAction.mockImplementation(async () => ({
      ok: false,
      error: "タスクの更新に失敗しました",
    }));
    render(<TaskListView initialTasks={[task({ status: "todo" })]} entries={[entry()]} />);

    await userEvent.click(screen.getByRole("button", { name: "タスク完了にする" }));

    await waitFor(() => expect(screen.getByRole("alert")).toBeInTheDocument());
    // 巻き戻されて未完了 (aria-pressed=false) に戻る
    expect(
      screen.getByRole("button", { name: "タスク完了にする" }),
    ).toHaveAttribute("aria-pressed", "false");
    expect(confettiFired()).toBe(false);
  });

  it("Entry があるとタスク追加フォームを表示し、送信できる", async () => {
    const user = userEvent.setup();
    render(<TaskListView initialTasks={[]} entries={[entry()]} />);

    expect(screen.getByLabelText("Entry")).toHaveValue("e1");
    await user.type(screen.getByLabelText("タスク名"), "一次面接");
    await user.click(screen.getByRole("button", { name: /タスクを追加/ }));

    await waitFor(() => expect(createTaskFromTaskPageAction).toHaveBeenCalled());
    const fd = createTaskFromTaskPageAction.mock.calls[0][1] as FormData;
    expect(fd.get("entryId")).toBe("e1");
    expect(fd.get("title")).toBe("一次面接");
  });

  it("定型チップをタップするとタスク名と種類欄に反映される", async () => {
    const user = userEvent.setup();
    render(<TaskListView initialTasks={[]} entries={[entry()]} />);

    await user.click(screen.getByRole("button", { name: "一次面接" }));

    expect(screen.getByLabelText("タスク名")).toHaveValue("一次面接");
    expect(screen.getByRole("radio", { name: "予定" })).toBeChecked();
  });

  it("Entry が無いと追加フォームではなく応募先登録の導線を表示する", () => {
    render(<TaskListView initialTasks={[]} entries={[]} />);
    expect(screen.getByText("先に応募先を登録してください")).toBeInTheDocument();
    expect(screen.queryByLabelText("タスク名")).not.toBeInTheDocument();
  });

  it("削除ボタンで deleteTaskAction を呼び、成功時は一覧から消す", async () => {
    const user = userEvent.setup();
    vi.spyOn(window, "confirm").mockReturnValue(true);
    render(<TaskListView initialTasks={[task()]} entries={[entry()]} />);

    await user.click(screen.getByRole("button", { name: /タスク「ES提出」を削除/ }));

    await waitFor(() => expect(deleteTaskAction).toHaveBeenCalledWith("t1", "e1"));
    await waitFor(() =>
      expect(screen.queryByRole("link", { name: /ES提出/ })).not.toBeInTheDocument(),
    );
  });
});

describe("dueLabel / dueColor", () => {
  const now = new Date("2026-05-29T00:00:00Z");

  it("超過は「M/D ・n日超過」にし、1-2日はpink、3日以上はpink-deepにする", () => {
    expect(dueLabel("2026-05-28", now)).toBe("5/28 ・1日超過");
    expect(dueColor("2026-05-28", now)).toBe("bg-pink");
    expect(dueColor("2026-05-27", now)).toBe("bg-pink");
    expect(dueLabel("2026-05-26", now)).toBe("5/26 ・3日超過");
    expect(dueColor("2026-05-26", now)).toBe("bg-pink-deep");
  });

  it("本日締切は専用ラベルと最も強い色(ink)にする", () => {
    expect(dueLabel("2026-05-29", now)).toBe("本日締切");
    expect(dueColor("2026-05-29", now)).toBe("bg-ink");
  });

  it("明日締切は専用ラベルと次に強い色(pink-deep)にする", () => {
    expect(dueLabel("2026-05-30", now)).toBe("明日締切");
    expect(dueColor("2026-05-30", now)).toBe("bg-pink-deep");
  });

  it("3日以内はamber、それより先はskyにする", () => {
    expect(dueLabel("2026-05-31", now)).toBe("5/31");
    expect(dueColor("2026-05-31", now)).toBe("bg-amber");
    expect(dueColor("2026-06-02", now)).toBe("bg-sky");
  });

  it("期日なしはsage", () => {
    expect(dueLabel(null, now)).toBe("期日なし");
    expect(dueColor(null, now)).toBe("bg-sage");
  });
});

describe("sortTasksForDisplay", () => {
  it("未完了を先に、同じ状態では期日順に並べる", () => {
    const sorted = sortTasksForDisplay([
      task({ id: "done", status: "done", dueDate: "2026-05-01" }),
      task({ id: "late", status: "todo", dueDate: "2026-06-01" }),
      task({ id: "soon", status: "todo", dueDate: "2026-05-01" }),
      task({ id: "none", status: "todo", dueDate: null }),
    ]);

    expect(sorted.map((t) => t.id)).toEqual(["soon", "late", "none", "done"]);
  });
});
