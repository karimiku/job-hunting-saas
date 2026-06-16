import { STAGE_LABEL, type StageKind } from "./entry-stage";

export type SelectionFlowSource = "template" | "manual" | "ai_inbox" | "ai_paste";

export interface SelectionStageInput {
  stageKind: StageKind | "other";
  stageLabel: string;
  evidenceText?: string;
}

export interface SelectionStageResponse {
  id: string;
  position: number;
  stageKind: StageKind | "other" | string;
  stageLabel: string;
  evidenceText: string;
}

export interface SelectionFlowResponse {
  id: string;
  entryId: string;
  source: SelectionFlowSource | string;
  currentStagePosition: number;
  confidence?: number | null;
  inboxClipId?: string | null;
  stages: SelectionStageResponse[];
  createdAt: string;
  updatedAt: string;
}

export const DEFAULT_SELECTION_STAGES: SelectionStageInput[] = [
  { stageKind: "application", stageLabel: STAGE_LABEL.application },
  { stageKind: "document", stageLabel: STAGE_LABEL.document },
  { stageKind: "test", stageLabel: STAGE_LABEL.test },
  { stageKind: "interview", stageLabel: STAGE_LABEL.interview },
  { stageKind: "group", stageLabel: STAGE_LABEL.group },
  { stageKind: "offer", stageLabel: STAGE_LABEL.offer },
];

export function parseSelectionFlowText(raw: string): SelectionStageInput[] {
  return raw
    .split(/\n|→|>|、|,/)
    .map((label) => label.trim())
    .filter(Boolean)
    .slice(0, 20)
    .map((label) => ({
      stageKind: inferStageKind(label),
      stageLabel: label,
    }));
}

export function inferStageKind(label: string): SelectionStageInput["stageKind"] {
  if (/内定|オファー/i.test(label)) return "offer";
  if (/GD|グループ\s*ディスカッション|group/i.test(label)) return "group";
  if (/SPI|Web\s*テスト|適性|筆記|試験|テスト|コーディング|課題/i.test(label)) return "test";
  if (/ES|エントリーシート|書類/i.test(label)) return "document";
  if (/面接|面談|人事|技術|最終|一次|二次|三次/i.test(label)) return "interview";
  if (/応募|エントリー/i.test(label)) return "application";
  return "other";
}

export function flowPayloadFromForm(mode: string, customText: string) {
  const customStages = parseSelectionFlowText(customText);
  const useCustom = mode === "custom" && customStages.length > 0;
  return {
    source: useCustom ? "manual" : "template",
    currentStagePosition: 1,
    stages: useCustom ? customStages : DEFAULT_SELECTION_STAGES,
  };
}
