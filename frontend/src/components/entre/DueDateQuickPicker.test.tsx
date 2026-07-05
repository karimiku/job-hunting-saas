import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { addDays, DueDateQuickPicker } from "./DueDateQuickPicker";

describe("addDays", () => {
  it("+0日は基準日そのまま", () => {
    expect(addDays(new Date("2026-07-05T00:00:00"), 0)).toBe("2026-07-05");
  });

  it("+3日を計算する", () => {
    expect(addDays(new Date("2026-07-05T00:00:00"), 3)).toBe("2026-07-08");
  });

  it("+7日で月をまたぐ場合も計算する", () => {
    expect(addDays(new Date("2026-07-28T00:00:00"), 7)).toBe("2026-08-04");
  });
});

describe("DueDateQuickPicker", () => {
  it("今日／+3日／+1週間 のボタンを表示する", () => {
    render(<DueDateQuickPicker onSelect={vi.fn()} now={new Date("2026-07-05T00:00:00")} />);
    expect(screen.getByRole("button", { name: "今日" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "+3日" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "+1週間" })).toBeInTheDocument();
  });

  it("+3日を押すと基準日+3日を onSelect に渡す", async () => {
    const onSelect = vi.fn();
    render(<DueDateQuickPicker onSelect={onSelect} now={new Date("2026-07-05T00:00:00")} />);

    await userEvent.click(screen.getByRole("button", { name: "+3日" }));

    expect(onSelect).toHaveBeenCalledWith("2026-07-08");
  });

  it("+1週間を押すと基準日+7日を onSelect に渡す", async () => {
    const onSelect = vi.fn();
    render(<DueDateQuickPicker onSelect={onSelect} now={new Date("2026-07-05T00:00:00")} />);

    await userEvent.click(screen.getByRole("button", { name: "+1週間" }));

    expect(onSelect).toHaveBeenCalledWith("2026-07-12");
  });
});
