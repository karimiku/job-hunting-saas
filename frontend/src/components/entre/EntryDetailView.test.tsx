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
    expect(screen.getByText("一次面接")).toBeInTheDocument();
    expect(screen.getByText("テストメモ")).toBeInTheDocument();
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
    await waitFor(() => expect(screen.queryByText("ES提出")).not.toBeInTheDocument());
  });

  it("Entry詳細からEntryを削除できる", async () => {
    vi.spyOn(window, "confirm").mockReturnValue(true);

    render(
      <EntryDetailView
        initialEntry={sample({ companyName: "テスト商事" })}
        initialTasks={[]}
      />,
    );

    await userEvent.click(screen.getByRole("button", { name: "テスト商事 のEntryを削除" }));

    await waitFor(() => expect(deleteEntryAction).toHaveBeenCalledWith("e1"));
  });

  it("Entry削除に失敗したらエラーを表示する", async () => {
    vi.spyOn(window, "confirm").mockReturnValue(true);
    deleteEntryAction.mockResolvedValue({ ok: false, error: "Entryの削除に失敗しました" });

    render(<EntryDetailView initialEntry={sample({ companyName: "テスト商事" })} initialTasks={[]} />);

    await userEvent.click(screen.getByRole("button", { name: "テスト商事 のEntryを削除" }));

    expect(await screen.findByText("Entryの削除に失敗しました")).toBeInTheDocument();
  });
});
