import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  EntryListView,
  filterEntries,
  filterEntriesByStage,
  filterEntriesByStatusGroup,
  nextTaskBadgeLabel,
  nextTaskForEntry,
  sortEntries,
} from "./EntryListView";
import type { EntryResponse } from "@/lib/api/entries";
import type { TaskWithEntry } from "@/lib/api/server-resources";

const T = (over: Partial<TaskWithEntry>): TaskWithEntry =>
  ({
    id: "t1",
    entryId: "e1",
    title: "ES提出",
    type: "deadline",
    status: "todo",
    dueDate: "2026-07-06T00:00:00.000Z",
    memo: "",
    companyName: "テスト社",
    ...over,
  }) as TaskWithEntry;

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

  it("進捗バーの近くに「Nステップ中M」を表示する（ステージ名はバッジと重複しないよう省く）", () => {
    render(
      <EntryListView
        entries={[E({ stageKind: "interview", stageLabel: "一次面接" })]}
      />,
    );
    expect(screen.getByText("6ステップ中4")).toBeInTheDocument();
    // ステージ名バッジは残るが、重複表示（"一次面接 ・ ..."）は出さない。
    expect(screen.queryByText(/一次面接 ・/)).not.toBeInTheDocument();
  });

  it("1件のときはツールバーを表示しない", () => {
    render(<EntryListView entries={[E({ companyName: "テスト商事" })]} />);
    expect(screen.queryByPlaceholderText("会社名・応募経路・媒体・メモで検索")).not.toBeInTheDocument();
  });

  it("2件以上のときはツールバーを表示し、検索すると絞り込む", async () => {
    const user = userEvent.setup();
    render(
      <EntryListView
        entries={[
          E({ id: "e1", companyName: "アルファ商事" }),
          E({ id: "e2", companyName: "ベータ工業" }),
        ]}
      />,
    );
    expect(screen.getByText("アルファ商事")).toBeInTheDocument();
    expect(screen.getByText("ベータ工業")).toBeInTheDocument();

    await user.type(
      screen.getByPlaceholderText("会社名・応募経路・媒体・メモで検索"),
      "アルファ",
    );
    expect(screen.getByText("アルファ商事")).toBeInTheDocument();
    expect(screen.queryByText("ベータ工業")).not.toBeInTheDocument();
  });

  it("ステージチップを押すとそのステージだけに絞り込む", async () => {
    const user = userEvent.setup();
    render(
      <EntryListView
        entries={[
          E({ id: "e1", companyName: "アルファ商事", stageKind: "document", stageLabel: "ES提出" }),
          E({ id: "e2", companyName: "ベータ工業", stageKind: "interview", stageLabel: "一次面接" }),
        ]}
      />,
    );

    await user.click(screen.getByRole("button", { name: "面接" }));
    expect(screen.getByText("ベータ工業")).toBeInTheDocument();
    expect(screen.queryByText("アルファ商事")).not.toBeInTheDocument();
  });

  it("状態チップを押すと結果の粒度で絞り込む", async () => {
    const user = userEvent.setup();
    render(
      <EntryListView
        entries={[
          E({ id: "e1", companyName: "アルファ商事", status: "in_progress" }),
          E({ id: "e2", companyName: "ベータ工業", status: "offered" }),
        ]}
      />,
    );

    await user.click(screen.getByRole("button", { name: "内定" }));
    expect(screen.getByText("ベータ工業")).toBeInTheDocument();
    expect(screen.queryByText("アルファ商事")).not.toBeInTheDocument();
  });

  it("絞り込みで0件になったら該当なしメッセージを表示する", async () => {
    const user = userEvent.setup();
    render(
      <EntryListView
        entries={[
          E({ id: "e1", companyName: "アルファ商事" }),
          E({ id: "e2", companyName: "ベータ工業" }),
        ]}
      />,
    );

    await user.type(
      screen.getByPlaceholderText("会社名・応募経路・媒体・メモで検索"),
      "存在しない会社名",
    );
    expect(screen.getByText("条件に一致する応募先がありません。")).toBeInTheDocument();
  });
});

describe("filterEntries", () => {
  it("空文字クエリはそのまま返す", () => {
    const entries = [E({ id: "e1" }), E({ id: "e2" })];
    expect(filterEntries(entries, "")).toEqual(entries);
    expect(filterEntries(entries, "   ")).toEqual(entries);
  });

  it("会社名・route・source・memoのいずれかに部分一致すれば残す", () => {
    const entries = [
      E({ id: "e1", companyName: "アルファ商事", route: "本選考", source: "リクナビ", memo: "" }),
      E({ id: "e2", companyName: "ベータ工業", route: "インターン", source: "マイナビ", memo: "紹介あり" }),
    ];
    expect(filterEntries(entries, "アルファ").map((e) => e.id)).toEqual(["e1"]);
    expect(filterEntries(entries, "インターン").map((e) => e.id)).toEqual(["e2"]);
    expect(filterEntries(entries, "マイナビ").map((e) => e.id)).toEqual(["e2"]);
    expect(filterEntries(entries, "紹介").map((e) => e.id)).toEqual(["e2"]);
  });

  it("大文字小文字を無視する", () => {
    const entries = [E({ id: "e1", source: "ABC Navi" })];
    expect(filterEntries(entries, "abc").map((e) => e.id)).toEqual(["e1"]);
  });

  it("一致しなければ空配列", () => {
    const entries = [E({ id: "e1", companyName: "アルファ商事" })];
    expect(filterEntries(entries, "ゼータ")).toEqual([]);
  });
});

describe("filterEntriesByStage", () => {
  it("stage が null ならそのまま返す", () => {
    const entries = [E({ id: "e1", stageKind: "document" }), E({ id: "e2", stageKind: "interview" })];
    expect(filterEntriesByStage(entries, null)).toEqual(entries);
  });

  it("指定ステージのみ残す", () => {
    const entries = [
      E({ id: "e1", stageKind: "document" }),
      E({ id: "e2", stageKind: "interview" }),
    ];
    expect(filterEntriesByStage(entries, "interview").map((e) => e.id)).toEqual(["e2"]);
  });

  it("未知の stageKind は other として扱う", () => {
    const entries = [E({ id: "e1", stageKind: "coding_test" })];
    expect(filterEntriesByStage(entries, "other").map((e) => e.id)).toEqual(["e1"]);
    expect(filterEntriesByStage(entries, "interview")).toEqual([]);
  });
});

describe("filterEntriesByStatusGroup", () => {
  const entries = [
    E({ id: "e1", status: "in_progress" }),
    E({ id: "e2", status: "offered" }),
    E({ id: "e3", status: "accepted" }),
    E({ id: "e4", status: "rejected" }),
    E({ id: "e5", status: "withdrawn" }),
  ];

  it("all はそのまま返す", () => {
    expect(filterEntriesByStatusGroup(entries, "all")).toEqual(entries);
  });

  it("in_progress は選考中のみ", () => {
    expect(filterEntriesByStatusGroup(entries, "in_progress").map((e) => e.id)).toEqual(["e1"]);
  });

  it("offer は offered/accepted をまとめる", () => {
    expect(filterEntriesByStatusGroup(entries, "offer").map((e) => e.id)).toEqual(["e2", "e3"]);
  });

  it("closed は rejected/withdrawn をまとめる", () => {
    expect(filterEntriesByStatusGroup(entries, "closed").map((e) => e.id)).toEqual(["e4", "e5"]);
  });
});

describe("sortEntries", () => {
  it("updated は updatedAt 降順", () => {
    const entries = [
      E({ id: "e1", updatedAt: "2026-01-01T00:00:00Z" }),
      E({ id: "e2", updatedAt: "2026-03-01T00:00:00Z" }),
      E({ id: "e3", updatedAt: "2026-02-01T00:00:00Z" }),
    ];
    expect(sortEntries(entries, "updated").map((e) => e.id)).toEqual(["e2", "e3", "e1"]);
  });

  it("company は会社名の辞書順", () => {
    const entries = [
      E({ id: "e1", companyName: "ベータ工業" }),
      E({ id: "e2", companyName: "アルファ商事" }),
    ];
    expect(sortEntries(entries, "company").map((e) => e.id)).toEqual(["e2", "e1"]);
  });

  it("stage は選考フェーズの進行順", () => {
    const entries = [
      E({ id: "e1", stageKind: "offer" }),
      E({ id: "e2", stageKind: "application" }),
      E({ id: "e3", stageKind: "interview" }),
    ];
    expect(sortEntries(entries, "stage").map((e) => e.id)).toEqual(["e2", "e3", "e1"]);
  });

  it("元配列を変更しない", () => {
    const entries = [E({ id: "e1", updatedAt: "2026-01-01T00:00:00Z" }), E({ id: "e2", updatedAt: "2026-02-01T00:00:00Z" })];
    const copy = [...entries];
    sortEntries(entries, "updated");
    expect(entries).toEqual(copy);
  });
});

describe("nextTaskForEntry", () => {
  it("その応募先の未完了タスクのうち最も近い期日を返す", () => {
    const tasks = [
      T({ id: "a", entryId: "e1", dueDate: "2026-07-10T00:00:00.000Z" }),
      T({ id: "b", entryId: "e1", dueDate: "2026-07-06T00:00:00.000Z" }),
      T({ id: "c", entryId: "e2", dueDate: "2026-07-01T00:00:00.000Z" }),
    ];
    expect(nextTaskForEntry("e1", tasks)?.id).toBe("b");
  });
  it("完了・期限なし・他応募先は除外し、無ければ null", () => {
    const tasks = [
      T({ id: "done", entryId: "e1", status: "done", dueDate: "2026-07-06T00:00:00.000Z" }),
      T({ id: "nodue", entryId: "e1", dueDate: null }),
      T({ id: "other", entryId: "e2" }),
    ];
    expect(nextTaskForEntry("e1", tasks)).toBeNull();
  });
});

describe("nextTaskBadgeLabel", () => {
  const now = new Date("2026-07-06T09:00:00Z");
  it("本日/明日/超過を暦日で判定する", () => {
    expect(nextTaskBadgeLabel(T({ dueDate: "2026-07-06T00:00:00.000Z", title: "ES" }), now)).toBe("本日締切（ES）");
    expect(nextTaskBadgeLabel(T({ dueDate: "2026-07-07T00:00:00.000Z", title: "面接" }), now)).toBe("明日締切（面接）");
    expect(nextTaskBadgeLabel(T({ dueDate: "2026-07-04T00:00:00.000Z", title: "SPI" }), now)).toBe("7/4 ・2日超過（SPI）");
    expect(nextTaskBadgeLabel(T({ dueDate: "2026-07-10T00:00:00.000Z", title: "説明会" }), now)).toBe("締切 7/10（説明会）");
  });
  it("タスクが無ければ null", () => {
    expect(nextTaskBadgeLabel(null, now)).toBeNull();
  });
});

describe("EntryListView 次締切バッジ", () => {
  it("カードに次タスクの締切を表示する", () => {
    render(
      <EntryListView
        entries={[E({ id: "e1", companyName: "サイバー社" })]}
        tasks={[T({ entryId: "e1", title: "ES提出", dueDate: "2026-07-10T00:00:00.000Z" })]}
      />,
    );
    expect(screen.getByText(/ES提出/)).toBeInTheDocument();
  });
});
