import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { server } from "@/test/msw-server";
import { useEntries } from "./useEntries";

const API = "http://localhost:8080";

describe("useEntries", () => {
  it("初期状態は loading=true / data=undefined", () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({ entries: [] }),
      ),
    );
    const { result } = renderHook(() => useEntries());
    expect(result.current.loading).toBe(true);
    expect(result.current.data).toBeUndefined();
    expect(result.current.error).toBeUndefined();
  });

  it("成功時に loading=false / data に配列がセットされる", async () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({
          entries: [
            { id: "e1", companyId: "c1", route: "本選考", source: "リクナビ", status: "in_progress", stageKind: "interview", stageLabel: "面接", memo: "", createdAt: "x", updatedAt: "x" },
          ],
        }),
      ),
    );
    const { result } = renderHook(() => useEntries());
    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.data).toHaveLength(1);
    expect(result.current.error).toBeUndefined();
  });

  it("失敗時に error がセットされ data は undefined のまま", async () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({ message: "boom" }, { status: 500 }),
      ),
    );
    const { result } = renderHook(() => useEntries());
    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.error).toBeDefined();
    expect(result.current.data).toBeUndefined();
  });
});
