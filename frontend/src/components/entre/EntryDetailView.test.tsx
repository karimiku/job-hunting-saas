import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { server } from "@/test/msw-server";
import { EntryDetailView } from "./EntryDetailView";
import type { EntryResponse } from "@/lib/api/entries";

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
});
