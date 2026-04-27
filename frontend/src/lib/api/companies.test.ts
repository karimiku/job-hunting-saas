import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { server } from "@/test/msw-server";
import { listCompanies, createCompany, getCompany } from "./companies";

const API = "http://localhost:8080";

describe("companies API", () => {
  it("listCompanies は会社一覧を返す", async () => {
    server.use(
      http.get(`${API}/api/v1/companies`, () =>
        HttpResponse.json({
          companies: [
            { id: "c1", name: "○○商事", memo: "", createdAt: "x", updatedAt: "x" },
          ],
        }),
      ),
    );
    const result = await listCompanies();
    expect(result).toHaveLength(1);
    expect(result[0].name).toBe("○○商事");
  });

  it("createCompany は POST してレスポンスを返す", async () => {
    let body: Record<string, unknown> | null = null;
    server.use(
      http.post(`${API}/api/v1/companies`, async ({ request }) => {
        body = (await request.json()) as Record<string, unknown>;
        return HttpResponse.json(
          { id: "c-new", name: body.name, memo: "", createdAt: "x", updatedAt: "x" },
          { status: 201 },
        );
      }),
    );
    const result = await createCompany({ name: "新会社" });
    expect(body).toEqual({ name: "新会社" });
    expect(result.id).toBe("c-new");
  });

  it("getCompany は ID から取得", async () => {
    server.use(
      http.get(`${API}/api/v1/companies/c1`, () =>
        HttpResponse.json({ id: "c1", name: "○○商事", memo: "", createdAt: "x", updatedAt: "x" }),
      ),
    );
    const result = await getCompany("c1");
    expect(result.id).toBe("c1");
  });
});
