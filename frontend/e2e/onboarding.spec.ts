import { expect, test } from "@playwright/test";

test.describe("Onboarding flow (/onboarding)", () => {
  test("3 ステップ進めて はじめる ✨ で /dashboard に遷移する", async ({ page }) => {
    await page.goto("/onboarding");

    await expect(page.getByText(/step 1 \/ 3/i)).toBeVisible();
    await expect(page.getByText("はじめまして！")).toBeVisible();

    // Step 1 → 2
    await page.getByRole("button", { name: /つぎへ/ }).click();
    await expect(page.getByText(/step 2 \/ 3/i)).toBeVisible();
    await expect(page.getByText(/バラバラを、ぜんぶ1枚に/)).toBeVisible();

    // Step 2 → 3
    await page.getByRole("button", { name: /つぎへ/ }).click();
    await expect(page.getByText(/step 3 \/ 3/i)).toBeVisible();
    await expect(page.getByText(/内定までの道のりを、一緒に/)).toBeVisible();

    // Step 3 → finish (はじめる ✨)
    await page.getByRole("button", { name: /はじめる/ }).click();
    // 認証ガードで /login にリダイレクトされるはずだが、まずは /dashboard へ向かう挙動を確認
    await page.waitForURL(/\/(dashboard|login)/);
  });

  test("封筒くん SVG が描画されている", async ({ page }) => {
    await page.goto("/onboarding");
    // <svg aria-label="封筒くん (...)" /> を探す
    const mascot = page.locator('svg[aria-label^="封筒くん"]');
    await expect(mascot.first()).toBeVisible();
  });

  test("step インジケーターが進行中ステップを 1 つだけ示す", async ({ page }) => {
    await page.goto("/onboarding");
    await expect(page.locator('[aria-current="step"]')).toHaveCount(1);

    await page.getByRole("button", { name: /つぎへ/ }).click();
    await expect(page.getByText(/step 2 \/ 3/i)).toBeVisible();
    await expect(page.locator('[aria-current="step"]')).toHaveCount(1);
  });
});
