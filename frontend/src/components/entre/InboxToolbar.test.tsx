import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { InboxToolbar } from "./InboxToolbar";

describe("InboxToolbar", () => {
  it("検索欄に入力すると onQueryChange を呼ぶ", async () => {
    const user = userEvent.setup();
    const onQueryChange = vi.fn();
    render(
      <InboxToolbar
        query=""
        onQueryChange={onQueryChange}
        order="new"
        onOrderChange={vi.fn()}
        sources={[]}
        selectedSource={null}
        onSourceChange={vi.fn()}
      />,
    );

    await user.type(screen.getByPlaceholderText("タイトル・URL・会社名で検索"), "a");
    expect(onQueryChange).toHaveBeenCalledWith("a");
  });

  it("並び替えを変更すると onOrderChange を呼ぶ", async () => {
    const user = userEvent.setup();
    const onOrderChange = vi.fn();
    render(
      <InboxToolbar
        query=""
        onQueryChange={vi.fn()}
        order="new"
        onOrderChange={onOrderChange}
        sources={[]}
        selectedSource={null}
        onSourceChange={vi.fn()}
      />,
    );

    await user.selectOptions(screen.getByDisplayValue("新しい順"), "会社名順");
    expect(onOrderChange).toHaveBeenCalledWith("company");
  });

  it("source が1件以下ならチップを表示しない", () => {
    render(
      <InboxToolbar
        query=""
        onQueryChange={vi.fn()}
        order="new"
        onOrderChange={vi.fn()}
        sources={["マイナビ"]}
        selectedSource={null}
        onSourceChange={vi.fn()}
      />,
    );
    expect(screen.queryByRole("button", { name: "マイナビ" })).not.toBeInTheDocument();
  });

  it("sourceが2件以上あるとき「すべて」と各チップを表示し、押すと onSourceChange を呼ぶ", async () => {
    const user = userEvent.setup();
    const onSourceChange = vi.fn();
    render(
      <InboxToolbar
        query=""
        onQueryChange={vi.fn()}
        order="new"
        onOrderChange={vi.fn()}
        sources={["マイナビ", "リクナビ"]}
        selectedSource={null}
        onSourceChange={onSourceChange}
      />,
    );

    expect(screen.getByRole("button", { name: "すべて" })).toHaveAttribute("aria-pressed", "true");
    await user.click(screen.getByRole("button", { name: "リクナビ" }));
    expect(onSourceChange).toHaveBeenCalledWith("リクナビ");
  });
});
