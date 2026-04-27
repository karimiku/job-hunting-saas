import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { RoadmapView } from "./RoadmapView";
import type { EntryResponse } from "@/lib/api/entries";

const e = (kind: string): EntryResponse => ({
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
  it("各マイルストーンの社数を集計する", () => {
    render(
      <RoadmapView
        entries={[
          e("application"),
          e("document"),
          e("document"),
          e("test"),
          e("interview"),
          e("interview"),
          e("interview"),
          e("offer"),
        ]}
      />,
    );
    expect(screen.getByTestId("milestone-count-application")).toHaveTextContent("1");
    expect(screen.getByTestId("milestone-count-document")).toHaveTextContent("2");
    expect(screen.getByTestId("milestone-count-test")).toHaveTextContent("1");
    expect(screen.getByTestId("milestone-count-interview")).toHaveTextContent("3");
    expect(screen.getByTestId("milestone-count-offer")).toHaveTextContent("1");
  });
});
