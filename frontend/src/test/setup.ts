import "@testing-library/jest-dom/vitest";
import { afterAll, afterEach, beforeAll, vi } from "vitest";
import { server } from "./msw-server";

// クライアント fetch のベースは本番では同一 origin の "/backend" (rewrite proxy) だが、
// jsdom の fetch (undici) は相対 URL を解決できないため、テストでは絶対 URL に上書きする。
// MSW ハンドラ (http://localhost:8080) もこの値に合わせている。
process.env.NEXT_PUBLIC_CLIENT_API_BASE = "http://localhost:8080";

// next/navigation を全テストで mock — useRouter / usePathname / useParams を提供する。
vi.mock("next/navigation", async () => {
  const actual = await vi.importActual<typeof import("next/navigation")>("next/navigation");
  return {
    ...actual,
    useRouter: () => ({
      push: vi.fn(),
      replace: vi.fn(),
      back: vi.fn(),
      forward: vi.fn(),
      refresh: vi.fn(),
      prefetch: vi.fn(),
    }),
    usePathname: () => "/",
    useParams: () => ({}),
    useSearchParams: () => new URLSearchParams(),
  };
});

// MSW: API モックを各テストで使い回す。
beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
