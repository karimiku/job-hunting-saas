import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { server } from "@/test/msw-server";
import { useEntry } from "./useEntry";

const API = "http://localhost:8080";
const sample = { id: "e1", companyId: "c1", route: "本選考", source: "リクナビ", status: "in_progress", stageKind: "interview", stageLabel: "一次面接", memo: "", createdAt: "x", updatedAt: "x" };

describe("useEntry", () => {
  it("ID で 1件取得し data に入る", async () => {
    server.use(http.get(`${API}/api/v1/entries/e1`, () => HttpResponse.json(sample)));
    const { result } = renderHook(() => useEntry("e1"));
    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.data?.stageLabel).toBe("一次面接");
  });

  it("404 のとき error がセットされる", async () => {
    server.use(http.get(`${API}/api/v1/entries/missing`, () => HttpResponse.json({ message: "not found" }, { status: 404 })));
    const { result } = renderHook(() => useEntry("missing"));
    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.error).toBeDefined();
  });

  it("id が undefined のとき fetch しない", () => {
    const { result } = renderHook(() => useEntry(undefined));
    expect(result.current.loading).toBe(false);
    expect(result.current.data).toBeUndefined();
  });
});
