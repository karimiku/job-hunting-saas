import { apiFetch } from "./client";
export { ApiError } from "./client";

export interface EntryResponse {
  id: string;
  companyId: string;
  /** 会社名。backend レスポンスには無く、server-resources で companyId から join して埋める。 */
  companyName?: string;
  route: string;
  source: string;
  sourceUrl?: string;
  status: string;
  stageKind: string;
  stageLabel: string;
  memo: string;
  createdAt: string;
  updatedAt: string;
}

const NO_COMPANY_LABEL = "（会社名未設定）";

/** 会社名の主表示。join できなかった場合はフォールバック文言を返す。 */
export function companyDisplayName(
  entry: Pick<EntryResponse, "companyName">,
): string {
  return entry.companyName?.trim() || NO_COMPANY_LABEL;
}

export interface ListEntriesParams {
  status?: string;
  stageKind?: string;
  source?: string;
}

export async function listEntries(
  params: ListEntriesParams = {},
): Promise<EntryResponse[]> {
  const qs = new URLSearchParams();
  if (params.status) qs.set("status", params.status);
  if (params.stageKind) qs.set("stageKind", params.stageKind);
  if (params.source) qs.set("source", params.source);
  const path = `/api/v1/entries${qs.toString() ? `?${qs.toString()}` : ""}`;
  const res = await apiFetch<{ entries: EntryResponse[] }>(path);
  return res.entries;
}

export async function getEntry(id: string): Promise<EntryResponse> {
  return apiFetch<EntryResponse>(`/api/v1/entries/${id}`);
}

export interface CreateEntryInput {
  companyId: string;
  route: string;
  source: string;
  sourceUrl?: string;
  memo?: string;
}

export async function createEntry(
  input: CreateEntryInput,
): Promise<EntryResponse> {
  return apiFetch<EntryResponse>("/api/v1/entries", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export interface UpdateEntryInput {
  source?: string;
  sourceUrl?: string;
  status?: string;
  stageKind?: string;
  stageLabel?: string;
  memo?: string;
}

export async function updateEntry(
  id: string,
  input: UpdateEntryInput,
): Promise<EntryResponse> {
  return apiFetch<EntryResponse>(`/api/v1/entries/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
}

export async function deleteEntry(id: string): Promise<void> {
  await apiFetch<void>(`/api/v1/entries/${id}`, { method: "DELETE" });
}

export function entrySourceUrl(
  entry: Pick<EntryResponse, "sourceUrl" | "memo">,
): string | null {
  const direct = normalizeHttpsUrl(entry.sourceUrl);
  if (direct) return direct;

  for (const token of entry.memo.split(/\s+/)) {
    const candidate = normalizeHttpsUrl(token);
    if (candidate) return candidate;
  }
  return null;
}

function normalizeHttpsUrl(value: string | undefined): string | null {
  if (!value) return null;
  const trimmed = value.trim().replace(/[)\]、。,.]+$/, "");
  try {
    const url = new URL(trimmed);
    return url.protocol === "https:" ? url.href : null;
  } catch {
    return null;
  }
}
