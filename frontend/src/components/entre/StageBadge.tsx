import { CSSProperties } from "react";

export type StageKind = "entry" | "doc" | "es" | "interview" | "offer";

interface StageBadgeProps {
  kind: StageKind;
  label?: string;
  size?: "sm" | "md";
}

const STAGE_LABEL: Record<StageKind, string> = {
  entry: "エントリー",
  doc: "書類選考",
  es: "ES提出",
  interview: "面接",
  offer: "内定",
};

const STAGE_BG: Record<StageKind, string> = {
  entry: "var(--color-stage-entry)",
  doc: "var(--color-stage-doc)",
  es: "var(--color-stage-es)",
  interview: "var(--color-stage-interview)",
  offer: "var(--color-stage-offer)",
};

/** 選考ステージ表示用のバッジ。 */
export function StageBadge({ kind, label, size = "md" }: StageBadgeProps) {
  const style: CSSProperties = {
    background: STAGE_BG[kind],
    color: "#fff",
    height: size === "sm" ? 18 : 22,
    padding: size === "sm" ? "0 6px" : "0 8px",
    fontSize: size === "sm" ? 9 : 10,
  };
  return (
    <span
      className="inline-flex items-center rounded font-bold leading-none"
      style={style}
    >
      {label ?? STAGE_LABEL[kind]}
    </span>
  );
}
