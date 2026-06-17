export const STAGE_ORDER = [
  "application",
  "document",
  "test",
  "interview",
  "group",
  "offer",
] as const;

export type StageKind = (typeof STAGE_ORDER)[number];
export type KanbanStageKind = StageKind | "other";

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

export const KANBAN_STAGE_ORDER = [...STAGE_ORDER, "other"] as const;

export const KANBAN_STAGE_LABEL: Record<KanbanStageKind, string> = {
  ...STAGE_LABEL,
  other: "その他",
};

export const KANBAN_STAGE_COLOR: Record<KanbanStageKind, string> = {
  ...STAGE_COLOR,
  other: "var(--color-ink-3)",
};

export const STAGE_BG: Record<string, string> = {
  ...KANBAN_STAGE_COLOR,
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

export function isKanbanStageKind(value: string): value is KanbanStageKind {
  return (KANBAN_STAGE_ORDER as readonly string[]).includes(value);
}

export function stageLabelOf(value: string): string {
  if (isKanbanStageKind(value)) return KANBAN_STAGE_LABEL[value];
  return value;
}

export function stageIndexOf(value: string): number {
  const index = STAGE_ORDER.indexOf(value as StageKind);
  return index < 0 ? 0 : index;
}

export function kanbanStageIndexOf(value: string): number {
  const index = KANBAN_STAGE_ORDER.indexOf(value as KanbanStageKind);
  return index < 0 ? KANBAN_STAGE_ORDER.indexOf("other") : index;
}

export function statusForStage(stageKind: StageKind): string {
  return stageKind === "offer" ? "offered" : "in_progress";
}

export function statusForKanbanStage(stageKind: KanbanStageKind): string {
  return stageKind === "offer" ? "offered" : "in_progress";
}
