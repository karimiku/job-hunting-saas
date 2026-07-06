import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { EntryListToolbar } from "./EntryListToolbar";

describe("EntryListToolbar", () => {
  it("検索欄に入力すると onQueryChange を呼ぶ", async () => {
    const user = userEvent.setup();
    const onQueryChange = vi.fn();
    render(
      <EntryListToolbar
        query=""
        onQueryChange={onQueryChange}
        order="updated"
        onOrderChange={vi.fn()}
        stages={[]}
        selectedStage={null}
        onStageChange={vi.fn()}
        statusGroup="all"
        onStatusGroupChange={vi.fn()}
      />,
    );

    await user.type(screen.getByPlaceholderText("会社名・応募経路・媒体・メモで検索"), "a");
    expect(onQueryChange).toHaveBeenCalledWith("a");
  });

  it("並び替えを変更すると onOrderChange を呼ぶ", async () => {
    const user = userEvent.setup();
    const onOrderChange = vi.fn();
    render(
      <EntryListToolbar
        query=""
        onQueryChange={vi.fn()}
        order="updated"
        onOrderChange={onOrderChange}
        stages={[]}
        selectedStage={null}
        onStageChange={vi.fn()}
        statusGroup="all"
        onStatusGroupChange={vi.fn()}
      />,
    );

    await user.selectOptions(screen.getByDisplayValue("更新が新しい順"), "会社名順");
    expect(onOrderChange).toHaveBeenCalledWith("company");
  });

  it("状態チップは常に表示し、押すと onStatusGroupChange を呼ぶ", async () => {
    const user = userEvent.setup();
    const onStatusGroupChange = vi.fn();
    render(
      <EntryListToolbar
        query=""
        onQueryChange={vi.fn()}
        order="updated"
        onOrderChange={vi.fn()}
        stages={[]}
        selectedStage={null}
        onStageChange={vi.fn()}
        statusGroup="all"
        onStatusGroupChange={onStatusGroupChange}
      />,
    );

    expect(screen.getByRole("button", { name: "すべて" })).toHaveAttribute("aria-pressed", "true");
    await user.click(screen.getByRole("button", { name: "内定" }));
    expect(onStatusGroupChange).toHaveBeenCalledWith("offer");
  });

  it("ステージが1件以下ならチップを表示しない", () => {
    render(
      <EntryListToolbar
        query=""
        onQueryChange={vi.fn()}
        order="updated"
        onOrderChange={vi.fn()}
        stages={[{ value: "interview", label: "面接" }]}
        selectedStage={null}
        onStageChange={vi.fn()}
        statusGroup="all"
        onStatusGroupChange={vi.fn()}
      />,
    );
    expect(screen.queryByRole("button", { name: "面接" })).not.toBeInTheDocument();
  });

  it("ステージが2件以上あるとき「すべてのステージ」と各チップを表示し、押すと onStageChange を呼ぶ", async () => {
    const user = userEvent.setup();
    const onStageChange = vi.fn();
    render(
      <EntryListToolbar
        query=""
        onQueryChange={vi.fn()}
        order="updated"
        onOrderChange={vi.fn()}
        stages={[
          { value: "document", label: "書類" },
          { value: "interview", label: "面接" },
        ]}
        selectedStage={null}
        onStageChange={onStageChange}
        statusGroup="all"
        onStatusGroupChange={vi.fn()}
      />,
    );

    expect(screen.getByRole("button", { name: "すべてのステージ" })).toHaveAttribute(
      "aria-pressed",
      "true",
    );
    await user.click(screen.getByRole("button", { name: "面接" }));
    expect(onStageChange).toHaveBeenCalledWith("interview");
  });
});
