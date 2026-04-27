import { apiFetch } from "./client";

export interface InboxClipResponse {
  id: string;
  url: string;
  title: string;
  source: string;
  guess: string;
  capturedAt: string;
}

export async function listInboxClips(): Promise<InboxClipResponse[]> {
  const res = await apiFetch<{ clips: InboxClipResponse[] }>("/api/v1/inbox/clips");
  return res.clips;
}

export interface CreateInboxClipInput {
  url: string;
  title: string;
  source: string;
  guess?: string;
}

export async function createInboxClip(input: CreateInboxClipInput): Promise<InboxClipResponse> {
  return apiFetch<InboxClipResponse>("/api/v1/inbox/clips", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function deleteInboxClip(id: string): Promise<void> {
  await apiFetch<void>(`/api/v1/inbox/clips/${id}`, { method: "DELETE" });
}
