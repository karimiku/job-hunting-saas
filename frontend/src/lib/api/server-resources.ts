// Server Component から使う型付きフェッチャー集。
// Client 用 (entries.ts / inboxClips.ts 等) と同じ型を返すが、Cookie を Next.js Server から
// backend に転送するため serverFetch を呼ぶ。
//
// SSR 化済み画面はここを使う。Client Component 配下のフックや mutate 系は引き続き client.ts。

import { serverFetch } from "./server";
import type { EntryResponse, ListEntriesParams } from "./entries";
import type { InboxClipResponse } from "./inboxClips";
import type { TaskResponse } from "./tasks";

export async function listEntriesServer(
  params: ListEntriesParams = {},
): Promise<EntryResponse[]> {
  const qs = new URLSearchParams();
  if (params.status) qs.set("status", params.status);
  if (params.stageKind) qs.set("stageKind", params.stageKind);
  if (params.source) qs.set("source", params.source);
  const path = `/api/v1/entries${qs.toString() ? `?${qs.toString()}` : ""}`;
  const res = await serverFetch<{ entries: EntryResponse[] }>(path);
  return res.entries;
}

export async function getEntryServer(id: string): Promise<EntryResponse> {
  return serverFetch<EntryResponse>(`/api/v1/entries/${id}`);
}

export async function listInboxClipsServer(): Promise<InboxClipResponse[]> {
  const res = await serverFetch<{ clips: InboxClipResponse[] }>(
    "/api/v1/inbox/clips",
  );
  return res.clips;
}

export async function listTasksByEntryServer(
  entryId: string,
): Promise<TaskResponse[]> {
  const res = await serverFetch<{ tasks: TaskResponse[] }>(
    `/api/v1/entries/${entryId}/tasks`,
  );
  return res.tasks;
}
