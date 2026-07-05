import { expect, test } from "@playwright/test";

test.describe("Onboarding flow (/onboarding)", () => {
  test("3 ステップ進めて 最初の応募先を登録する で /entry/new (未ログインなら /login) に遷移する", async ({
    page,
  }) => {
    await page.goto("/onboarding");

    await expect(page.getByText(/step 1 \/ 3/i)).toBeVisible();
    await expect(page.getByText("受けている企業を、1か所に集めます")).toBeVisible();
    await expect(page.getByRole("link", { name: "スキップしてホームへ" })).toBeVisible();

    // Step 1 → 2
    await page.getByRole("button", { name: /つぎへ/ }).click();
    await expect(page.getByText(/step 2 \/ 3/i)).toBeVisible();
    await expect(page.getByText("選考が進んだら、カードを動かすだけ")).toBeVisible();
    await expect(page.getByRole("link", { name: "スキップしてホームへ" })).toBeVisible();

    // Step 2 → 3
    await page.getByRole("button", { name: /つぎへ/ }).click();
    await expect(page.getByText(/step 3 \/ 3/i)).toBeVisible();
    await expect(page.getByText("締切はタスクに。ホームが毎朝の起点")).toBeVisible();
    await expect(
      page.getByRole("link", { name: "あとで登録する（ホームへ）" }),
    ).toBeVisible();

    // Step 3 → finish (最初の応募先を登録する)
    await page.getByRole("button", { name: "最初の応募先を登録する" }).click();
    // 認証ガードで /login にリダイレクトされ得るが、まずは /entry/new へ向かう挙動を確認
    await page.waitForURL(/\/(entry\/new|login)/);
  });

  test("スキップリンクで /dashboard に遷移する", async ({ page }) => {
    await page.goto("/onboarding");
    await page.getByRole("link", { name: "スキップしてホームへ" }).click();
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
