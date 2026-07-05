import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TASK_TEMPLATES, TaskTemplateChips } from "./TaskTemplateChips";

describe("TaskTemplateChips", () => {
  it("定型テンプレのチップをすべて表示する", () => {
    render(<TaskTemplateChips onSelect={vi.fn()} />);
    for (const template of TASK_TEMPLATES) {
      expect(screen.getByRole("button", { name: template.title })).toBeInTheDocument();
    }
  });

  it("チップをタップすると title と type を渡して onSelect を呼ぶ", async () => {
    const onSelect = vi.fn();
    render(<TaskTemplateChips onSelect={onSelect} />);

    await userEvent.click(screen.getByRole("button", { name: "一次面接" }));

    expect(onSelect).toHaveBeenCalledWith({ title: "一次面接", type: "schedule" });
  });
});
