import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { server } from "@/test/msw-server";
import { listTasksByEntry, createTask, updateTask, deleteTask } from "./tasks";

const API = "http://localhost:8080";

describe("tasks API", () => {
  it("listTasksByEntry は entryId 配下のタスクを返す", async () => {
    server.use(
      http.get(`${API}/api/v1/entries/e1/tasks`, () =>
        HttpResponse.json({
          tasks: [
            { id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
          ],
        }),
      ),
    );
    const result = await listTasksByEntry("e1");
    expect(result).toHaveLength(1);
    expect(result[0].title).toBe("ES提出");
  });

  it("createTask は POST してレスポンスを返す", async () => {
    let body: Record<string, unknown> | null = null;
    server.use(
      http.post(`${API}/api/v1/entries/e1/tasks`, async ({ request }) => {
        body = (await request.json()) as Record<string, unknown>;
        return HttpResponse.json(
          { id: "t-new", entryId: "e1", title: body.title, type: body.type, status: "todo", dueDate: null, memo: "", createdAt: "x", updatedAt: "x" },
          { status: 201 },
        );
      }),
    );
    const result = await createTask("e1", { title: "ES提出", type: "deadline" });
    expect(body).toEqual({ title: "ES提出", type: "deadline" });
    expect(result.id).toBe("t-new");
  });

  it("updateTask は PATCH /tasks/:id で部分更新", async () => {
    let body: Record<string, unknown> | null = null;
    server.use(
      http.patch(`${API}/api/v1/tasks/t1`, async ({ request }) => {
        body = (await request.json()) as Record<string, unknown>;
        return HttpResponse.json({ id: "t1", entryId: "e1", title: "ES提出", type: "deadline", status: "done", dueDate: null, memo: "", createdAt: "x", updatedAt: "x2" });
      }),
    );
    const result = await updateTask("t1", { status: "done" });
    expect(body).toEqual({ status: "done" });
    expect(result.status).toBe("done");
  });

  it("deleteTask は DELETE /tasks/:id", async () => {
    let called = false;
    server.use(
      http.delete(`${API}/api/v1/tasks/t1`, () => {
        called = true;
        return new HttpResponse(null, { status: 204 });
      }),
    );
    await deleteTask("t1");
    expect(called).toBe(true);
  });
});
