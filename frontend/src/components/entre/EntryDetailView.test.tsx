import { beforeEach, describe, expect, it, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { EntryDetailView } from "./EntryDetailView";
import type { EntryResponse } from "@/lib/api/entries";
import type { TaskResponse } from "@/lib/api/tasks";

// 更新系は Server Action 経由になったため、actions モジュールをモックして
// 呼び出し引数と楽観更新の UI を検証する (HTTP レイヤは actions 側のテストで担保)。
const {
  updateEntryAction,
  updateSelectionFlowCurrentStageAction,
  createTaskForEntryAction,
  deleteEntryAction,
  setTaskStatusAction,
  deleteTaskAction,
} =
  vi.hoisted(() => ({
    updateEntryAction: vi.fn(),
    updateSelectionFlowCurrentStageAction: vi.fn(),
    createTaskForEntryAction: vi.fn(),
    deleteEntryAction: vi.fn(),
    setTaskStatusAction: vi.fn(),
    deleteTaskAction: vi.fn(),
  }));

vi.mock("@/app/entry/actions", () => ({
  updateEntryAction,
  updateSelectionFlowCurrentStageAction,
  createTaskForEntryAction,
  deleteEntryAction,
}));
vi.mock("@/app/task/actions", () => ({ setTaskStatusAction, deleteTaskAction }));

const sample = (overrides: Partial<EntryResponse> = {}): EntryResponse => ({
  id: "e1",
  companyId: "c1",
  route: "本選考",
  source: "リクナビ",
  status: "in_progress",
  stageKind: "interview",
  stageLabel: "一次面接",
  memo: "テストメモ",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

const task = (overrides: Partial<TaskResponse> = {}): TaskResponse => ({
  id: "t1",
  entryId: "e1",
  title: "ES提出",
  type: "deadline",
  status: "todo",
  dueDate: "2026-05-30T00:00:00Z",
  memo: "",
  createdAt: "x",
  updatedAt: "x",
  ...overrides,
});

describe("EntryDetailView", () => {
  beforeEach(() => {
    updateEntryAction.mockReset().mockResolvedValue({ ok: true });
    updateSelectionFlowCurrentStageAction.mockReset().mockResolvedValue({
      ok: true,
      selectionFlow: {
        id: "flow1",
        entryId: "e1",
        source: "manual",
        currentStagePosition: 2,
        stages: [
          { id: "s1", position: 1, stageKind: "document", stageLabel: "ES提出", evidenceText: "" },
          { id: "s2", position: 2, stageKind: "interview", stageLabel: "一次面接", evidenceText: "" },
        ],
        createdAt: "x",
        updatedAt: "x",
      },
    });
    createTaskForEntryAction.mockReset();
    deleteEntryAction.mockReset().mockResolvedValue({ ok: true });
    setTaskStatusAction.mockReset().mockResolvedValue({ ok: true, status: "done" });
    deleteTaskAction.mockReset().mockResolvedValue({ ok: true });
  });

  it("initialEntry を表示する", () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);
    expect(screen.getByTestId("current-stage")).toHaveTextContent("一次面接");
    expect(screen.getByText("テストメモ")).toBeInTheDocument();
  });

  it("タスク追加フォームは既定で折りたたまれており、ボタンで開閉できる", async () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[task()]} />);

    expect(screen.queryByLabelText("タスク名")).not.toBeInTheDocument();
    expect(screen.getByText("ES提出")).toBeInTheDocument();

    const toggle = screen.getByRole("button", { name: "タスクを追加" });
    await userEvent.click(toggle);
    expect(screen.getByLabelText("タスク名")).toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: "閉じる" }));
    expect(screen.queryByLabelText("タスク名")).not.toBeInTheDocument();
  });

  it("フェーズと結果の関係を説明する注記を表示する", () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);
    expect(
      screen.getByText(/選考途中は上の「選考フェーズ」だけでOKです/),
    ).toBeInTheDocument();
  });

  it("ステージボタンの選択で更新 action が走り stageKind を任意更新できる", async () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);
    await userEvent.click(screen.getByRole("button", { name: "GD" }));

    await waitFor(() =>
      expect(updateEntryAction).toHaveBeenCalledWith("e1", {
        stageKind: "group",
        stageLabel: "GD",
        status: "in_progress",
      }),
    );
    expect(screen.getByTestId("current-stage")).toHaveTextContent("GD");
  });

  it("結果ステータスで落選を選択できる", async () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);
    await userEvent.click(screen.getByRole("button", { name: "落選" }));

    await waitFor(() =>
      expect(updateEntryAction).toHaveBeenCalledWith(
        "e1",
        expect.objectContaining({ status: "rejected" }),
      ),
    );
  });

  it("更新 action が失敗したら楽観更新をロールバックしてエラーを表示する", async () => {
    updateEntryAction.mockResolvedValue({ ok: false, error: "選考ステータスの更新に失敗しました" });

    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);
    await userEvent.click(screen.getByRole("button", { name: "GD" }));

    expect(await screen.findByText("選考ステータスの更新に失敗しました")).toBeInTheDocument();
    expect(screen.getByTestId("current-stage")).toHaveTextContent("一次面接");
  });

  it("内定到達時はスタンプを表示する", () => {
    render(
      <EntryDetailView
        initialEntry={sample({ stageKind: "offer", stageLabel: "内定" })}
        initialTasks={[]}
      />,
    );
    expect(screen.getByText("内定！")).toBeInTheDocument();
  });

  it("initialEntry が null のとき alert を表示する", () => {
    render(<EntryDetailView initialEntry={null} initialTasks={[]} />);
    expect(screen.getByRole("alert")).toBeInTheDocument();
  });

  it("会社名をヘッダの見出しに表示する", () => {
    render(
      <EntryDetailView initialEntry={sample({ companyName: "テスト商事" })} initialTasks={[]} />,
    );
    expect(screen.getByRole("heading", { name: "テスト商事" })).toBeInTheDocument();
  });

  it("応募元URLを詳細ヘッダから開ける", () => {
    render(
      <EntryDetailView
        initialEntry={sample({ sourceUrl: "https://job.rikunabi.com/2027/company/r123/" })}
        initialTasks={[]}
      />,
    );
    expect(
      screen.getByRole("link", { name: "https://job.rikunabi.com/2027/company/r123/" }),
    ).toHaveAttribute("href", "https://job.rikunabi.com/2027/company/r123/");
  });

  it("会社名が取得できないときはフォールバック見出しを表示する", () => {
    render(
      <EntryDetailView initialEntry={sample({ companyName: undefined })} initialTasks={[]} />,
    );
    expect(screen.getByRole("heading", { name: "（会社名未設定）" })).toBeInTheDocument();
  });

  it("Entry詳細からタスクを追加できる", async () => {
    createTaskForEntryAction.mockResolvedValue({
      ok: true,
      task: task({ id: "t-new", title: "一次面接準備" }),
    });

    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);

    await userEvent.click(screen.getByRole("button", { name: "タスクを追加" }));
    await userEvent.type(screen.getByLabelText("タスク名"), "一次面接準備");
    await userEvent.click(screen.getByRole("button", { name: "追加" }));

    await waitFor(() =>
      expect(createTaskForEntryAction).toHaveBeenCalledWith("e1", {
        title: "一次面接準備",
        type: "deadline",
        dueDate: undefined,
        memo: undefined,
      }),
    );
    expect(await screen.findByText("一次面接準備")).toBeInTheDocument();
  });

  it("定型チップをタップするとタスク名と種類欄に反映される", async () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);

    await userEvent.click(screen.getByRole("button", { name: "タスクを追加" }));
    await userEvent.click(screen.getByRole("button", { name: "一次面接" }));

    expect(screen.getByLabelText("タスク名")).toHaveValue("一次面接");
    expect(screen.getByRole("radio", { name: "予定" })).toBeChecked();
  });

  it("Entry詳細でタスクの完了状態を切り替えられる", async () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[task()]} />);

    await userEvent.click(screen.getByRole("button", { name: "タスク完了にする" }));

    await waitFor(() => expect(setTaskStatusAction).toHaveBeenCalledWith("t1", "done", "e1"));
  });

  it("Entry詳細でタスクを削除できる", async () => {
    vi.spyOn(window, "confirm").mockReturnValue(true);

    render(<EntryDetailView initialEntry={sample()} initialTasks={[task()]} />);

    await userEvent.click(screen.getByRole("button", { name: /タスク「ES提出」を削除/ }));

    await waitFor(() => expect(deleteTaskAction).toHaveBeenCalledWith("t1", "e1"));
    await waitFor(() =>
      expect(
        screen.queryByRole("button", { name: /タスク「ES提出」を削除/ }),
      ).not.toBeInTheDocument(),
    );
  });

  it("Entry詳細からEntryを削除できる", async () => {
    vi.spyOn(window, "confirm").mockReturnValue(true);

    render(
      <EntryDetailView
        initialEntry={sample({ companyName: "テスト商事" })}
        initialTasks={[]}
      />,
    );

    await userEvent.click(screen.getByRole("button", { name: "テスト商事 の応募先を削除" }));

    await waitFor(() => expect(deleteEntryAction).toHaveBeenCalledWith("e1"));
  });

  it("選考中の間は「結果を取り消す」導線を表示しない", () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);
    expect(
      screen.queryByRole("button", { name: /結果を取り消す/ }),
    ).not.toBeInTheDocument();
  });

  it("結果確定後は「結果を取り消す」導線から選考中に戻せる", async () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);
    await userEvent.click(screen.getByRole("button", { name: "落選" }));

    const revertButton = await screen.findByRole("button", {
      name: /結果を取り消す/,
    });
    await userEvent.click(revertButton);

    await waitFor(() =>
      expect(updateEntryAction).toHaveBeenLastCalledWith(
        "e1",
        expect.objectContaining({ status: "in_progress" }),
      ),
    );
  });

  it("選考フェーズと結果が別セクションの見出しで分離されている", () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);

    expect(screen.getByText("選考フェーズ")).toBeInTheDocument();
    expect(screen.getByText("今どの段階かを選びます")).toBeInTheDocument();
    expect(screen.getByText("結果（確定したら選ぶ）")).toBeInTheDocument();
    expect(screen.getByText("内定・お見送りなどが決まったら選びます")).toBeInTheDocument();
    // フェーズの「内定」と結果の「内定獲得」が文言で区別できる
    expect(screen.getByRole("button", { name: "内定" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "内定獲得" })).toBeInTheDocument();
  });

  it("Entry削除に失敗したらエラーを表示する", async () => {
    vi.spyOn(window, "confirm").mockReturnValue(true);
    deleteEntryAction.mockResolvedValue({ ok: false, error: "応募先の削除に失敗しました" });

    render(<EntryDetailView initialEntry={sample({ companyName: "テスト商事" })} initialTasks={[]} />);

    await userEvent.click(screen.getByRole("button", { name: "テスト商事 の応募先を削除" }));

    expect(await screen.findByText("応募先の削除に失敗しました")).toBeInTheDocument();
  });
});
