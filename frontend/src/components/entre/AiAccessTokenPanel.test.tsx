import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { AiAccessTokenPanel } from "./AiAccessTokenPanel";
import type { AiAccessTokenResponse } from "@/lib/api/aiTokens";

// Server Action はここでは検証しない（next/headers を読み込ませない）。
vi.mock("@/app/profile/actions", () => ({
  createAiAccessTokenAction: vi.fn(async () => ({})),
  revokeAiAccessTokenAction: vi.fn(async () => ({})),
}));

describe("AiAccessTokenPanel", () => {
  it("既定では折りたたまれ、見出しのみ表示する", () => {
    render(<AiAccessTokenPanel tokens={[]} />);
    expect(screen.getByText("AI連携（上級者向け）")).toBeInTheDocument();
    expect(screen.getByText("使わない場合は設定不要です")).toBeInTheDocument();
    expect(screen.queryByText("AI連携トークン")).toBeNull();
  });

  it("見出しを押すと展開し、上級者向けの補足を表示する", async () => {
    const user = userEvent.setup();
    render(<AiAccessTokenPanel tokens={[]} />);
    await user.click(screen.getByRole("button", { name: /AI連携（上級者向け）/ }));

    expect(screen.getByText("AI連携トークン")).toBeInTheDocument();
    expect(
      screen.getByText(/AIアシスタント（Claude等）と連携したい人向けの機能です/),
    ).toBeInTheDocument();
    expect(screen.getByText("まだ発行していません")).toBeInTheDocument();
  });

  it("loadError があっても赤い警告調にせず、控えめな文言にする", async () => {
    const user = userEvent.setup();
    render(<AiAccessTokenPanel tokens={[]} loadError="いまは表示できませんでした。使う場合はあとで開いてください。" />);
    await user.click(screen.getByRole("button", { name: /AI連携（上級者向け）/ }));

    const message = screen.getByText(
      "いまは表示できませんでした。使う場合はあとで開いてください。",
    );
    expect(message).toBeInTheDocument();
    expect(message).not.toHaveAttribute("role", "alert");
    expect(message.className).not.toMatch(/pink/);
    expect(screen.getByText("いまは表示できません")).toBeInTheDocument();
  });

  it("トークンがある場合は一覧を表示する", async () => {
    const user = userEvent.setup();
    const tokens: AiAccessTokenResponse[] = [
      {
        id: "t1",
        name: "MCP用",
        tokenPrefix: "abcd1234",
        createdAt: "2026-04-01T00:00:00Z",
        lastUsedAt: null,
        revokedAt: null,
      },
    ];
    render(<AiAccessTokenPanel tokens={tokens} />);
    await user.click(screen.getByRole("button", { name: /AI連携（上級者向け）/ }));

    expect(screen.getByText("MCP用")).toBeInTheDocument();
    expect(screen.getByText("有効")).toBeInTheDocument();
  });
});
