import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { InboxClipConvert } from "./InboxClipConvert";
import type { InboxClipResponse } from "@/lib/api/inboxClips";
import type { CompanyResponse } from "@/lib/api/companies";

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

const companies: CompanyResponse[] = [
  {
    id: "co1",
    name: "株式会社○○商事",
    memo: "",
    createdAt: "2026-04-01T00:00:00Z",
    updatedAt: "2026-04-01T00:00:00Z",
  },
];

describe("InboxClipConvert", () => {
  it("初期状態では管理開始ボタンだけ表示しフォームは閉じている", () => {
    render(<InboxClipConvert clip={clip} companies={[]} />);
    expect(screen.getByRole("button", { name: /応募先にする/ })).toBeInTheDocument();
    expect(screen.queryByLabelText(/会社名/)).not.toBeInTheDocument();
  });

  it("管理開始ボタンを押すとフォームが開き、clip の値が初期値になる", async () => {
    const user = userEvent.setup();
    render(<InboxClipConvert clip={clip} companies={[]} />);

    await user.click(screen.getByRole("button", { name: /応募先にする/ }));

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
    render(<InboxClipConvert clip={{ ...clip, guess: "" }} companies={[]} />);
    await user.click(screen.getByRole("button", { name: /応募先にする/ }));
    expect(screen.getByLabelText(/会社名/)).toHaveValue("");
  });

  it("既存会社候補があると選択済みで表示し hidden に companyId を入れる", async () => {
    const user = userEvent.setup();
    render(<InboxClipConvert clip={clip} companies={companies} />);

    await user.click(screen.getByRole("button", { name: /応募先にする/ }));

    expect(screen.getByText("既存会社の候補")).toBeInTheDocument();
    expect(screen.getByRole("radio", { name: /株式会社○○商事/ })).toBeChecked();

    const existingCompanyInput = document.querySelector<HTMLInputElement>(
      'input[name="existingCompanyId"]',
    );
    expect(existingCompanyInput?.value).toBe("co1");
  });
});
