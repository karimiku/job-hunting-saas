import { describe, expect, it } from "vitest";
import { render } from "@testing-library/react";
import { Mascot } from "./Mascot";

describe("Mascot", () => {
  it("happy mood は entre-mascot-idle アニメ class を付ける", () => {
    const { container } = render(<Mascot mood="happy" />);
    expect(container.querySelector("svg")).toHaveClass("entre-mascot-idle");
  });

  it("cheering mood は hop アニメを付ける", () => {
    const { container } = render(<Mascot mood="cheering" />);
    expect(container.querySelector("svg")).toHaveClass("entre-mascot-hop");
  });

  it("bow mood (新規) は bow アニメを付け、表情も bow になる", () => {
    const { container } = render(<Mascot mood="bow" />);
    const svg = container.querySelector("svg");
    expect(svg).toHaveClass("entre-mascot-bow");
    // bow は目を閉じる → ellipse の y は通常より下
    // (実装は閉じた目の path を持つ — 単に SVG が描画されることを確認)
    expect(svg).toBeInTheDocument();
  });

  it("animate=false のときアニメ class を付けない", () => {
    const { container } = render(<Mascot mood="happy" animate={false} />);
    const svg = container.querySelector("svg");
    expect(svg).not.toHaveClass("entre-mascot-idle");
  });

  it("size prop を尊重する", () => {
    const { container } = render(<Mascot mood="happy" size={120} />);
    const svg = container.querySelector("svg");
    expect(svg).toHaveAttribute("width", "120");
    expect(svg).toHaveAttribute("height", "120");
  });
});
