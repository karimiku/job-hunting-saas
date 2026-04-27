import { expect, test } from "@playwright/test";

test.describe("Auth guards — 未ログイン時のリダイレクト", () => {
  // 認証必須ページは未ログイン状態で開くと /login に遷移する。
  const protectedRoutes = [
    "/dashboard",
    "/entry",
    "/kanban",
    "/roadmap",
    "/task",
    "/inbox",
    "/profile",
  ];

  for (const path of protectedRoutes) {
    test(`${path} は未ログインだと /login にリダイレクトされる`, async ({ page }) => {
      await page.goto(path);
      // useUser() はマウント後に状態判定するので、リダイレクトに少し時間がかかる
      await page.waitForURL(/\/login/, { timeout: 10_000 });
      await expect(page).toHaveURL(/\/login/);
    });
  }
});
