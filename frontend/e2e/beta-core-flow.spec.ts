import { expect, test } from "@playwright/test";

const API_BASE = process.env.PLAYWRIGHT_MOCK_API_BASE ?? "http://127.0.0.1:18080";

test.describe("Beta core flow — Inbox clip から Entry/Task 管理", () => {
  test("保存済み clip を Entry 化し、Entry/Kanban/Task で管理できる", async ({
    page,
    request,
  }, testInfo) => {
    const suffix = `${testInfo.project.name}-${Date.now()}`.replace(/\W+/g, "-");
    const company = `Codex Beta ${suffix}`;
    const title = `${company} 採用ページ`;
    const taskTitle = `ES提出 ${suffix}`;

    await page.context().addCookies([
      {
        name: "e2e-auth",
        value: "1",
        domain: "localhost",
        path: "/",
      },
    ]);

    const clip = await request.post(`${API_BASE}/api/v1/inbox/clips`, {
      headers: { cookie: "e2e-auth=1" },
      data: {
        url: `https://jobs.example.test/${suffix}`,
        title,
        source: "MockNavi",
        guess: company,
      },
    });
    expect(clip.ok()).toBeTruthy();

    await page.goto("/inbox");
    await expect(page.getByText(title)).toBeVisible();

    await page.getByRole("button", { name: "Entryとして管理" }).click();
    await expect(page.getByLabel("会社名")).toHaveValue(company);
    await expect(page.getByLabel("ソース")).toHaveValue("MockNavi");

    await page
      .getByRole("button", { name: /Entryを作成して開く/ })
      .click();
    await page.waitForURL(/\/entry\/[^/]+$/);
    await expect(page.getByRole("heading", { name: company })).toBeVisible();
    await expect(page.getByText("MockNavi · 本選考")).toBeVisible();

    await page.goto("/inbox");
    await expect(page.getByText(title)).toHaveCount(0);

    await page.goto("/entry");
    await expect(page.getByText(company).first()).toBeVisible();

    await page.goto("/kanban");
    await expect(page.getByText(company).first()).toBeVisible();

    await page.goto("/task");
    await page.getByLabel("タスク名").fill(taskTitle);
    await page.getByLabel("期日").fill("2026-06-15");
    await page.getByRole("button", { name: "Taskを追加" }).click();
    await expect(page.getByText("Taskを追加しました。")).toBeVisible();
    await expect(page.getByText(taskTitle)).toBeVisible();

    const completeButton = page.getByRole("button", {
      name: "タスク完了にする",
    });
    await completeButton.click();
    const taskRow = page.getByRole("listitem").filter({ hasText: taskTitle });
    await expect(
      taskRow.getByRole("button", { name: "タスク未完了に戻す" }),
    ).toHaveAttribute("aria-pressed", "true");
  });
});
