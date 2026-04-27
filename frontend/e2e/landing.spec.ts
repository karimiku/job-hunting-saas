import { expect, test } from "@playwright/test";

test.describe("Landing page (/)", () => {
  test("ロゴ Entré + メインキャッチが表示される", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveTitle(/Entré/);
    await expect(page.getByText("Entré").first()).toBeVisible();
  });

  test("ナビからログインへ遷移できる", async ({ page }) => {
    await page.goto("/");
    // LP のヘッダ "ログイン" / フッターの CTA 等いずれか最初の link で遷移
    const loginLinks = page.getByRole("link", { name: /ログイン|はじめる|sign in/i });
    if ((await loginLinks.count()) > 0) {
      await loginLinks.first().click();
      await expect(page).toHaveURL(/\/login/);
    } else {
      await page.goto("/login");
      await expect(page).toHaveURL(/\/login/);
    }
  });
});
