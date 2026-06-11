export const STAGE_ORDER = [
  "application",
  "document",
  "test",
  "interview",
  "group",
  "offer",
] as const;

export type StageKind = (typeof STAGE_ORDER)[number];

export const STAGE_LABEL: Record<StageKind, string> = {
  application: "エントリー",
  document: "書類",
  test: "テスト",
  interview: "面接",
  group: "GD",
  offer: "内定",
};

export const STAGE_COLOR: Record<StageKind, string> = {
  application: "var(--color-stage-entry)",
  document: "var(--color-stage-doc)",
  test: "var(--color-stage-es)",
  interview: "var(--color-stage-interview)",
  group: "var(--color-stage-interview-deep)",
  offer: "var(--color-stage-offer)",
};

export const STAGE_BG: Record<string, string> = {
  ...STAGE_COLOR,
  other: "var(--color-ink-3)",
};

export const ENTRY_STATUS_LABEL: Record<string, string> = {
  in_progress: "選考中",
  offered: "内定",
  accepted: "承諾",
  rejected: "落選",
  withdrawn: "辞退",
};

export function isStageKind(value: string): value is StageKind {
  return (STAGE_ORDER as readonly string[]).includes(value);
}

export function stageLabelOf(value: string): string {
  return isStageKind(value) ? STAGE_LABEL[value] : value;
}

export function stageIndexOf(value: string): number {
  const index = STAGE_ORDER.indexOf(value as StageKind);
  return index < 0 ? 0 : index;
}

export function statusForStage(stageKind: StageKind): string {
  return stageKind === "offer" ? "offered" : "in_progress";
}
