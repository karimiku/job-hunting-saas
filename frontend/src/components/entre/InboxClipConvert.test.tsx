import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { InboxClipConvert } from "./InboxClipConvert";
import type { InboxClipResponse } from "@/lib/api/inboxClips";

// Server Action はここでは検証しない（next/headers を読み込ませない）。
vi.mock("@/app/inbox/actions", () => ({
  convertInboxClipAction: vi.fn(async () => ({})),
}));

const clip: InboxClipResponse = {
  id: "clip1",
  url: "https://job.example.com/jobs/123",
  title: "○○商事 / 総合職 2026",
  source: "マイナビ",
  guess: "○○商事",
  capturedAt: "2026-04-26T00:00:00Z",
};

describe("InboxClipConvert", () => {
  it("初期状態では Entry化ボタンだけ表示しフォームは閉じている", () => {
    render(<InboxClipConvert clip={clip} />);
    expect(screen.getByRole("button", { name: /Entry化/ })).toBeInTheDocument();
    expect(screen.queryByLabelText(/会社名/)).not.toBeInTheDocument();
  });

  it("Entry化を押すとフォームが開き、clip の値が初期値になる", async () => {
    const user = userEvent.setup();
    render(<InboxClipConvert clip={clip} />);

    await user.click(screen.getByRole("button", { name: /Entry化/ }));

    expect(screen.getByLabelText(/会社名/)).toHaveValue("○○商事");
    expect(screen.getByLabelText(/ソース/)).toHaveValue("マイナビ");

    const memo = screen.getByLabelText(/メモ/) as HTMLTextAreaElement;
    expect(memo.value).toContain("○○商事 / 総合職 2026");
    expect(memo.value).toContain("https://job.example.com/jobs/123");

    const honsen = screen.getByRole("radio", { name: "本選考" }) as HTMLInputElement;
    expect(honsen.checked).toBe(true);

    // clipId が hidden で送られる
    const clipIdInput = document.querySelector<HTMLInputElement>('input[name="clipId"]');
    expect(clipIdInput?.value).toBe("clip1");
  });

  it("guess が空でも開ける（会社名は空初期値）", async () => {
    const user = userEvent.setup();
    render(<InboxClipConvert clip={{ ...clip, guess: "" }} />);
    await user.click(screen.getByRole("button", { name: /Entry化/ }));
    expect(screen.getByLabelText(/会社名/)).toHaveValue("");
  });
});
