// Server Component から使う型付きフェッチャー集。
// Client 用 (entries.ts / inboxClips.ts 等) と同じ型を返すが、Cookie を Next.js Server から
// backend に転送するため serverFetch を呼ぶ。
//
// SSR 化済み画面はここを使う。Client Component 配下のフックや mutate 系は引き続き client.ts。

import { serverFetch } from "./server";
import type { EntryResponse, ListEntriesParams } from "./entries";
import type { CompanyResponse } from "./companies";
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

export async function listCompaniesServer(): Promise<CompanyResponse[]> {
  const res = await serverFetch<{ companies: CompanyResponse[] }>(
    "/api/v1/companies",
  );
  return res.companies;
}

// entries に会社名を join した一覧。会社一覧を1回だけ引いて companyId→name で突き合わせる
// (N+1 を避けるため per-entry fetch はしない)。会社一覧取得に失敗しても entries は返す。
export async function listEntriesWithCompanyNamesServer(
  params: ListEntriesParams = {},
): Promise<EntryResponse[]> {
  const [entries, companies] = await Promise.all([
    listEntriesServer(params),
    listCompaniesServer().catch(() => [] as CompanyResponse[]),
  ]);
  const nameById = new Map(companies.map((c) => [c.id, c.name]));
  return entries.map((e) => ({ ...e, companyName: nameById.get(e.companyId) }));
}

// 単一 entry の会社名を引く。取得失敗時は undefined（UI 側でフォールバック表示）。
export async function getCompanyNameServer(
  companyId: string,
): Promise<string | undefined> {
  try {
    const company = await serverFetch<CompanyResponse>(
      `/api/v1/companies/${companyId}`,
    );
    return company.name;
  } catch {
    return undefined;
  }
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

export interface TaskWithEntry extends TaskResponse {
  /** タスクが属する entry の会社名 (join 済み、未設定なら undefined)。 */
  companyName?: string;
}

// 1人のユーザーの全タスクを1回の API で取得し、渡された entries から会社名を join する。
export async function listAllTasksServer(
  entries: Pick<EntryResponse, "id" | "companyName">[],
): Promise<TaskWithEntry[]> {
  if (entries.length === 0) return [];

  const res = await serverFetch<{ tasks: TaskResponse[] }>("/api/v1/tasks");
  const companyNameByEntryId = new Map(
    entries.map((entry) => [entry.id, entry.companyName]),
  );
  return res.tasks.map((task) => ({
    ...task,
    companyName: companyNameByEntryId.get(task.entryId),
  }));
}

export interface NavCounts {
  entry: number;
  task: number;
  inbox: number;
}

// サイドバーのバッジ用カウント。Entry / Inbox は一覧件数、Task は未完了タスク件数。
// どれか1つの取得に失敗しても 0 にフォールバックしてサイドバー描画は止めない。
export async function getNavCountsServer(): Promise<NavCounts> {
  const [entries, clips] = await Promise.all([
    listEntriesServer().catch(() => [] as EntryResponse[]),
    listInboxClipsServer().catch(() => [] as InboxClipResponse[]),
  ]);
  const tasks = await listAllTasksServer(entries).catch(
    () => [] as TaskWithEntry[],
  );
  return {
    entry: entries.length,
    task: tasks.filter((t) => t.status === "todo").length,
    inbox: clips.length,
  };
}
