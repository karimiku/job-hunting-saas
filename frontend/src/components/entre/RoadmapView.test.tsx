import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { server } from "@/test/msw-server";
import { RoadmapView } from "./RoadmapView";

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

describe("RoadmapView", () => {
  it("各マイルストーンの社数を API から計算する", async () => {
    server.use(
      http.get(`${API}/api/v1/entries`, () =>
        HttpResponse.json({
          entries: [
            e("application"),
            e("document"),
            e("document"),
            e("test"),
            e("interview"),
            e("interview"),
            e("interview"),
            e("offer"),
          ],
        }),
      ),
    );

    render(<RoadmapView />);
    await waitFor(() => {
      expect(screen.getByTestId("milestone-count-application")).toHaveTextContent("1");
      expect(screen.getByTestId("milestone-count-document")).toHaveTextContent("2");
      expect(screen.getByTestId("milestone-count-test")).toHaveTextContent("1");
      expect(screen.getByTestId("milestone-count-interview")).toHaveTextContent("3");
      expect(screen.getByTestId("milestone-count-offer")).toHaveTextContent("1");
    });
  });
});
