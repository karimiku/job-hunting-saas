import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { server } from "@/test/msw-server";
import { KanbanBoard } from "./KanbanBoard";

const API = "http://localhost:8080";
const e = (kind: string, source: string) => ({
  id: `e-${source}`,
  companyId: `c-${source}`,
  route: "本選考",
  source,
  status: "in_progress",
  stageKind: kind,
  stageLabel: kind,
  memo: "",
  createdAt: "x",
  updatedAt: "x",
});

describe("KanbanBoard", () => {
  it("ステージごとの 5 列にエントリーを振り分けて表示する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({
          entries: [
            e("application", "リクナビ"),
            e("interview", "マイナビ"),
            e("interview", "ONE CAREER"),
            e("offer", "OfferBox"),
          ],
        }),
      ),
    );
    render(<KanbanBoard />);
    await waitFor(() => {
      // application 列に 1, interview 列に 2, offer 列に 1
      expect(screen.getByTestId("column-count-application")).toHaveTextContent("1");
      expect(screen.getByTestId("column-count-interview")).toHaveTextContent("2");
      expect(screen.getByTestId("column-count-offer")).toHaveTextContent("1");
      expect(screen.getByTestId("column-count-document")).toHaveTextContent("0");
    });
  });
});
