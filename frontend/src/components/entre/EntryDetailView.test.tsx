import { http, HttpResponse } from "msw";
import { describe, expect, it, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { server } from "@/test/msw-server";
import { EntryDetailView } from "./EntryDetailView";
import type { EntryResponse } from "@/lib/api/entries";
import type { TaskResponse } from "@/lib/api/tasks";

const API = "http://localhost:8080";

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
  it("initialEntry を表示する", () => {
    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);
    expect(screen.getByText("一次面接")).toBeInTheDocument();
    expect(screen.getByText("テストメモ")).toBeInTheDocument();
  });

  it("「進める →」クリックで PATCH が走り stageKind が次に進む", async () => {
    let patchedBody: Record<string, unknown> | null = null;
    server.use(
      http.patch(`${API}/api/v1/entries/e1`, async ({ request }) => {
        patchedBody = (await request.json()) as Record<string, unknown>;
        return HttpResponse.json(sample({ stageKind: "group", stageLabel: "GD" }));
      }),
    );

    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);
    const advance = screen.getByRole("button", { name: /進める/ });
    await userEvent.click(advance);

    await waitFor(() => expect(patchedBody).not.toBeNull());
    expect(patchedBody).toMatchObject({ stageKind: "group" });
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
    let postedBody: Record<string, unknown> | null = null;
    server.use(
      http.post(`${API}/api/v1/entries/e1/tasks`, async ({ request }) => {
        postedBody = (await request.json()) as Record<string, unknown>;
        return HttpResponse.json(task({ id: "t-new", title: String(postedBody.title) }), {
          status: 201,
        });
      }),
    );

    render(<EntryDetailView initialEntry={sample()} initialTasks={[]} />);

    await userEvent.type(screen.getByLabelText("タスク名"), "一次面接準備");
    await userEvent.click(screen.getByRole("button", { name: "追加" }));

    await waitFor(() => expect(postedBody).not.toBeNull());
    expect(postedBody).toMatchObject({
      title: "一次面接準備",
      type: "deadline",
    });
    expect(await screen.findByText("一次面接準備")).toBeInTheDocument();
  });

  it("Entry詳細でタスクの完了状態を切り替えられる", async () => {
    let patchedBody: Record<string, unknown> | null = null;
    server.use(
      http.patch(`${API}/api/v1/tasks/t1`, async ({ request }) => {
        patchedBody = (await request.json()) as Record<string, unknown>;
        return HttpResponse.json(task({ status: "done" }));
      }),
    );

    render(<EntryDetailView initialEntry={sample()} initialTasks={[task()]} />);

    await userEvent.click(screen.getByRole("button", { name: "タスク完了にする" }));

    await waitFor(() => expect(patchedBody).toMatchObject({ status: "done" }));
  });

  it("Entry詳細でタスクを削除できる", async () => {
    vi.spyOn(window, "confirm").mockReturnValue(true);
    let deleted = false;
    server.use(
      http.delete(`${API}/api/v1/tasks/t1`, () => {
        deleted = true;
        return new HttpResponse(null, { status: 204 });
      }),
    );

    render(<EntryDetailView initialEntry={sample()} initialTasks={[task()]} />);

    await userEvent.click(screen.getByRole("button", { name: /タスク「ES提出」を削除/ }));

    await waitFor(() => expect(deleted).toBe(true));
    await waitFor(() => expect(screen.queryByText("ES提出")).not.toBeInTheDocument());
  });
});
