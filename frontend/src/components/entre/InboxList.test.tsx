import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  filterClips,
  filterClipsBySource,
  findDuplicateGroups,
  InboxList,
  sortClips,
} from "./InboxList";
import type { InboxClipResponse } from "@/lib/api/inboxClips";
import type { CompanyResponse } from "@/lib/api/companies";

// Server Action はここでは検証しない（next/headers を読み込ませない）。
vi.mock("@/app/inbox/actions", () => ({
  convertInboxClipAction: vi.fn(async () => ({})),
  deleteInboxClipAction: vi.fn(async () => ({})),
}));

const companies: CompanyResponse[] = [];

function clip(overrides: Partial<InboxClipResponse> = {}): InboxClipResponse {
  return {
    id: "c1",
    url: "https://job.example.com/a",
    title: "○○商事 / 総合職",
    source: "マイナビ",
    guess: "○○商事",
    capturedAt: "2026-04-26T00:00:00Z",
    ...overrides,
  };
}

describe("filterClips", () => {
  const clips = [
    clip({ id: "1", title: "エンジニア採用", url: "https://a.example.com", guess: "アルファ商事" }),
    clip({ id: "2", title: "総合職募集", url: "https://beta.example.com", guess: "ベータ工業" }),
  ];

  it("空クエリなら全件を返す", () => {
    expect(filterClips(clips, "")).toEqual(clips);
    expect(filterClips(clips, "   ")).toEqual(clips);
  });

  it("タイトルの部分一致で絞り込む", () => {
    expect(filterClips(clips, "エンジニア").map((c) => c.id)).toEqual(["1"]);
  });

  it("URLの部分一致で絞り込む（大文字小文字を無視）", () => {
    expect(filterClips(clips, "BETA.EXAMPLE").map((c) => c.id)).toEqual(["2"]);
  });

  it("会社候補(guess)の部分一致で絞り込む", () => {
    expect(filterClips(clips, "アルファ").map((c) => c.id)).toEqual(["1"]);
  });

  it("一致しなければ空配列", () => {
    expect(filterClips(clips, "存在しない会社")).toEqual([]);
  });
});

describe("filterClipsBySource", () => {
  const clips = [
    clip({ id: "1", source: "マイナビ" }),
    clip({ id: "2", source: "リクナビ" }),
  ];

  it("source が null なら全件を返す", () => {
    expect(filterClipsBySource(clips, null)).toEqual(clips);
  });

  it("指定した source のみ返す", () => {
    expect(filterClipsBySource(clips, "リクナビ").map((c) => c.id)).toEqual(["2"]);
  });
});

describe("sortClips", () => {
  const clips = [
    clip({ id: "old", capturedAt: "2026-01-01T00:00:00Z", guess: "Zeta" }),
    clip({ id: "mid", capturedAt: "2026-02-01T00:00:00Z", guess: "Alpha" }),
    clip({ id: "new", capturedAt: "2026-03-01T00:00:00Z", guess: "Middle" }),
  ];

  it("新しい順（既定）は capturedAt 降順", () => {
    expect(sortClips(clips, "new").map((c) => c.id)).toEqual(["new", "mid", "old"]);
  });

  it("古い順は capturedAt 昇順", () => {
    expect(sortClips(clips, "old").map((c) => c.id)).toEqual(["old", "mid", "new"]);
  });

  it("会社名順は guess の辞書順", () => {
    expect(sortClips(clips, "company").map((c) => c.id)).toEqual(["mid", "new", "old"]);
  });

  it("元配列を破壊しない", () => {
    const before = [...clips];
    sortClips(clips, "old");
    expect(clips).toEqual(before);
  });
});

describe("findDuplicateGroups", () => {
  it("正規化した guess が一致するクリップが2件以上ある場合のみグループ化する", () => {
    const clips = [
      clip({ id: "1", guess: "株式会社サンプル" }),
      clip({ id: "2", guess: "サンプル" }),
      clip({ id: "3", guess: "他社" }),
    ];
    const groups = findDuplicateGroups(clips);
    expect(groups.size).toBe(1);
    const [ids] = [...groups.values()];
    expect(ids.map((c) => c.id).sort()).toEqual(["1", "2"]);
  });

  it("guess が空/未検出のクリップは無視する", () => {
    const clips = [clip({ id: "1", guess: "" }), clip({ id: "2", guess: "" })];
    expect(findDuplicateGroups(clips).size).toBe(0);
  });

  it("重複がなければ空の Map を返す", () => {
    const clips = [clip({ id: "1", guess: "アルファ" }), clip({ id: "2", guess: "ベータ" })];
    expect(findDuplicateGroups(clips).size).toBe(0);
  });
});

describe("InboxList", () => {
  it("クリップが0件なら空状態を表示する", () => {
    render(<InboxList clips={[]} companies={companies} />);
    expect(screen.getByText("クリップは空です")).toBeInTheDocument();
  });

  it("クリップが1件のときはツールバーを表示しない", () => {
    render(<InboxList clips={[clip()]} companies={companies} />);
    expect(screen.queryByPlaceholderText("タイトル・URL・会社名で検索")).not.toBeInTheDocument();
  });

  it("検索するとタイトル・URL・会社候補に一致する行だけ残る", async () => {
    const user = userEvent.setup();
    render(
      <InboxList
        clips={[
          clip({ id: "1", title: "アルファ採用ページ", guess: "アルファ商事" }),
          clip({ id: "2", title: "ベータ採用ページ", guess: "ベータ工業" }),
        ]}
        companies={companies}
      />,
    );

    expect(screen.getAllByRole("listitem")).toHaveLength(2);

    await user.type(screen.getByPlaceholderText("タイトル・URL・会社名で検索"), "ベータ");

    const items = screen.getAllByRole("listitem");
    expect(items).toHaveLength(1);
    expect(screen.getByText("ベータ採用ページ")).toBeInTheDocument();
    expect(screen.queryByText("アルファ採用ページ")).not.toBeInTheDocument();
  });

  it("同じ会社候補が複数あるクリップに件数バッジを表示する", () => {
    render(
      <InboxList
        clips={[
          clip({ id: "1", title: "求人A", guess: "株式会社サンプル" }),
          clip({ id: "2", title: "求人B", guess: "サンプル" }),
        ]}
        companies={companies}
      />,
    );

    expect(screen.getAllByText("同じ会社の候補が他に1件あります")).toHaveLength(2);
  });

  it("ソースチップで絞り込める", async () => {
    const user = userEvent.setup();
    render(
      <InboxList
        clips={[
          clip({ id: "1", title: "求人A", source: "マイナビ" }),
          clip({ id: "2", title: "求人B", source: "リクナビ" }),
        ]}
        companies={companies}
      />,
    );

    await user.click(screen.getByRole("button", { name: "リクナビ" }));
    expect(screen.getAllByRole("listitem")).toHaveLength(1);
    expect(screen.getByText("求人B")).toBeInTheDocument();
  });

  it("既存の変換・削除ボタンの文言は変わらない", () => {
    render(<InboxList clips={[clip()]} companies={companies} />);
    expect(screen.getByRole("button", { name: "応募先にする" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /削除/ })).toBeInTheDocument();
  });
});
