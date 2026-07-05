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
    expect(screen.getByText(/まだ応募先がありません/)).toBeInTheDocument();
  });

  it("会社名を主表示し、source も維持する", () => {
    render(<EntryListView entries={[E({ companyName: "テスト商事", source: "リクナビ" })]} />);
    expect(screen.getByText("テスト商事")).toBeInTheDocument();
    expect(screen.getByText("リクナビ")).toBeInTheDocument();
  });

  it("応募元URLがあるEntryは元ページへのリンクを表示する", () => {
    render(
      <EntryListView
        entries={[
          E({
            companyName: "テスト商事",
            sourceUrl: "https://job.rikunabi.com/2027/company/r123/",
          }),
        ]}
      />,
    );
    expect(
      screen.getByRole("link", { name: "テスト商事 の応募元ページを開く" }),
    ).toHaveAttribute("href", "https://job.rikunabi.com/2027/company/r123/");
  });

  it("会社名が取得できないときはフォールバック表示する", () => {
    render(<EntryListView entries={[E({ companyName: undefined })]} />);
    expect(screen.getByText("（会社名未設定）")).toBeInTheDocument();
  });
});
