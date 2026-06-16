import { apiFetch } from "./client";
import type {
  SelectionFlowResponse,
  SelectionFlowSource,
  SelectionStageInput,
} from "../selection-flow";

export interface UpsertSelectionFlowInput {
  source: SelectionFlowSource | string;
  currentStagePosition?: number;
  confidence?: number | null;
  inboxClipId?: string | null;
  stages: SelectionStageInput[];
}

export async function getSelectionFlow(
  entryId: string,
): Promise<SelectionFlowResponse> {
  return apiFetch<SelectionFlowResponse>(
    `/api/v1/entries/${entryId}/selection-flow`,
  );
}

export async function upsertSelectionFlow(
  entryId: string,
  input: UpsertSelectionFlowInput,
): Promise<SelectionFlowResponse> {
  return apiFetch<SelectionFlowResponse>(
    `/api/v1/entries/${entryId}/selection-flow`,
    {
      method: "PUT",
      body: JSON.stringify(input),
    },
  );
}

export async function updateSelectionFlowCurrentStage(
  entryId: string,
  position: number,
): Promise<SelectionFlowResponse> {
  return apiFetch<SelectionFlowResponse>(
    `/api/v1/entries/${entryId}/selection-flow/current-stage`,
    {
      method: "PATCH",
      body: JSON.stringify({ position }),
    },
  );
}
