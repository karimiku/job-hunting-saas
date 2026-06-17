import { describe, expect, it } from "vitest";
import {
  kanbanStageIndexOf,
  stageIndexOf,
} from "./entry-stage";

describe("stageIndexOf", () => {
  it("通常の選考進捗では other を先頭相当として扱う", () => {
    expect(stageIndexOf("application")).toBe(0);
    expect(stageIndexOf("offer")).toBe(5);
    expect(stageIndexOf("other")).toBe(0);
    expect(stageIndexOf("coding_test")).toBe(0);
  });
});

describe("kanbanStageIndexOf", () => {
  it("Kanban列の順序では other を最後の列として扱う", () => {
    expect(kanbanStageIndexOf("application")).toBe(0);
    expect(kanbanStageIndexOf("offer")).toBe(5);
    expect(kanbanStageIndexOf("other")).toBe(6);
    expect(kanbanStageIndexOf("coding_test")).toBe(6);
  });
});
