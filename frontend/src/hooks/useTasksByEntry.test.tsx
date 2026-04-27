import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { server } from "@/test/msw-server";
import { useTasksByEntry } from "./useTasksByEntry";

const API = "http://localhost:8080";

describe("useTasksByEntry", () => {
  it("entryId 配下のタスクを取得する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries/e1/tasks`, () =>
        HttpResponse.json({
          tasks: [
            { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
          ],
        }),
      ),
    );
    const { result } = renderHook(() => useTasksByEntry("e1"));
    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.data).toHaveLength(1);
  });

  it("entryId が undefined のとき fetch しない", () => {
    const { result } = renderHook(() => useTasksByEntry(undefined));
    expect(result.current.loading).toBe(false);
  });
});
