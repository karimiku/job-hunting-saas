import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { server } from "@/test/msw-server";
import { StatusBreakdown } from "./StatusBreakdown";

const API = "http://localhost:8080";
const e = (kind: string) => ({
  id: String(Math.random()),
  companyId: "c1",
  route: "本選考",
  source: "リクナビ",
  status: "in_progress",
  stageKind: kind,
  stageLabel: kind,
  memo: "",
  createdAt: "x",
  updatedAt: "x",
});

describe("StatusBreakdown", () => {
  it("ステージごとの件数を凡例として表示する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({
          entries: [e("interview"), e("interview"), e("document"), e("offer")],
        }),
      ),
    );
    render(<StatusBreakdown />);
    await waitFor(() => {
      expect(screen.getByTestId("count-interview")).toHaveTextContent("2");
      expect(screen.getByTestId("count-document")).toHaveTextContent("1");
      expect(screen.getByTestId("count-offer")).toHaveTextContent("1");
    });
  });
});
