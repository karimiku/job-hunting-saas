import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { StatusBreakdown } from "./StatusBreakdown";
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

describe("StatusBreakdown", () => {
  it("ステージごとの件数を凡例として表示する", () => {
    render(
      <StatusBreakdown
        entries={[e("interview"), e("interview"), e("document"), e("offer")]}
      />,
    );
    expect(screen.getByTestId("count-interview")).toHaveTextContent("2");
    expect(screen.getByTestId("count-document")).toHaveTextContent("1");
    expect(screen.getByTestId("count-offer")).toHaveTextContent("1");
  });
});
