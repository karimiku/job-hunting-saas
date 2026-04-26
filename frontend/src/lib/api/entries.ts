import { apiFetch } from "./client";
export { ApiError } from "./client";

export interface EntryResponse {
  id: string;
  companyId: string;
  route: string;
  source: string;
  status: string;
  stageKind: string;
  stageLabel: string;
  memo: string;
  createdAt: string;
  updatedAt: string;
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
