import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { EntryListView } from "./EntryListView";
import type { EntryResponse } from "@/lib/api/entries";

const E = (over: Partial<EntryResponse>): EntryResponse => ({
  id: "e1",
  companyId: "c1",
  route: "本選考",
  source: "リクナビ",
  status: "in_progress",
  stageKind: "interview",
  stageLabel: "一次面接",
  memo: "",
  createdAt: "x",
  updatedAt: "x",
  ...over,
});

describe("EntryListView", () => {
  it("渡されたエントリーをすべて表示する", () => {
    render(
      <EntryListView
        entries={[
          E({ id: "e1", stageLabel: "一次面接" }),
          E({ id: "e2", source: "マイナビ", stageKind: "document", stageLabel: "ES提出" }),
        ]}
      />,
    );
    expect(screen.getByText("一次面接")).toBeInTheDocument();
    expect(screen.getByText("ES提出")).toBeInTheDocument();
  });

  it("0件のとき空状態を表示する", () => {
    render(<EntryListView entries={[]} />);
    expect(screen.getByText(/まだエントリーがありません/)).toBeInTheDocument();
  });
});
