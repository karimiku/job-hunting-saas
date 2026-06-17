// Server Component から使う型付きフェッチャー集。
// Client 用 (entries.ts / inboxClips.ts 等) と同じ型を返すが、Cookie を Next.js Server から
// backend に転送するため serverFetch を呼ぶ。
//
// SSR 化済み画面はここを使う。Client Component 配下のフックや mutate 系は引き続き client.ts。

import { serverFetch } from "./server";
import { ApiError, type AuthUser } from "./client-types";
import type { EntryResponse, ListEntriesParams } from "./entries";
import type { CompanyResponse } from "./companies";
import type { InboxClipResponse } from "./inboxClips";
import type { SelectionFlowResponse } from "../selection-flow";
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

export async function getSelectionFlowServer(
  entryId: string,
): Promise<SelectionFlowResponse> {
  return serverFetch<SelectionFlowResponse>(
    `/api/v1/entries/${entryId}/selection-flow`,
  );
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
  return listEntriesServer(params);
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

export async function getTaskServer(taskId: string): Promise<TaskResponse> {
  return serverFetch<TaskResponse>(`/api/v1/tasks/${taskId}`);
}

export interface TaskWithEntry extends TaskResponse {
  /** タスクが属する entry の会社名 (join 済み、未設定なら undefined)。 */
  companyName?: string;
}

// 1人のユーザーの全タスクを1回の API で取得する。
export async function listTasksServer(): Promise<TaskResponse[]> {
  const res = await serverFetch<{ tasks: TaskResponse[] }>("/api/v1/tasks");
  return res.tasks;
}

export function attachCompanyNamesToTasks(
  tasks: TaskResponse[],
  entries: Pick<EntryResponse, "id" | "companyName">[],
): TaskWithEntry[] {
  const companyNameByEntryId = new Map(
    entries.map((entry) => [entry.id, entry.companyName]),
  );
  return tasks.map((task) => ({
    ...task,
    companyName: companyNameByEntryId.get(task.entryId),
  }));
}

interface TaskPageDataResponse {
  user: AuthUser;
  entries: EntryResponse[];
  tasks: TaskResponse[];
}

interface AppPageDataResponse {
  user: AuthUser;
  entries: EntryResponse[];
  tasks: TaskResponse[];
  clips: InboxClipResponse[];
  companies: CompanyResponse[];
}

export interface TaskPageData {
  user: AuthUser;
  entries: EntryResponse[];
  tasks: TaskWithEntry[];
}

export interface AppPageData {
  user: AuthUser;
  entries: EntryResponse[];
  tasks: TaskWithEntry[];
  clips: InboxClipResponse[];
  companies: CompanyResponse[];
  navCounts: NavCounts;
}

// /task の初期表示に必要な user / entries / tasks を1回の backend fetch で取得する。
export async function getTaskPageDataServer(): Promise<TaskPageData | null> {
  try {
    const data = await serverFetch<TaskPageDataResponse>("/api/v1/page-data/task");
    return {
      user: data.user,
      entries: data.entries,
      tasks: attachCompanyNamesToTasks(data.tasks, data.entries),
    };
  } catch (e) {
    if (e instanceof ApiError && e.unauthorized) return null;
    throw e;
  }
}

// dashboard / entry / kanban / inbox の初期表示に必要なデータを1回の backend fetch で取得する。
export async function getAppPageDataServer(): Promise<AppPageData | null> {
  try {
    const data = await serverFetch<AppPageDataResponse>("/api/v1/page-data/app");
    const tasks = attachCompanyNamesToTasks(data.tasks, data.entries);
    return {
      user: data.user,
      entries: data.entries,
      tasks,
      clips: data.clips,
      companies: data.companies,
      navCounts: buildNavCounts(data.entries, data.tasks, data.clips),
    };
  } catch (e) {
    if (e instanceof ApiError && e.unauthorized) return null;
    throw e;
  }
}

// 1人のユーザーの全タスクを1回の API で取得し、渡された entries から会社名を join する。
export async function listAllTasksServer(
  entries: Pick<EntryResponse, "id" | "companyName">[],
): Promise<TaskWithEntry[]> {
  if (entries.length === 0) return [];

  return attachCompanyNamesToTasks(await listTasksServer(), entries);
}

export interface NavCounts {
  entry: number;
  task: number;
  inbox: number;
}

export function buildNavCounts(
  entries: readonly EntryResponse[],
  tasks: readonly TaskResponse[],
  clips: readonly InboxClipResponse[],
): NavCounts {
  return {
    entry: entries.length,
    task: tasks.filter((t) => t.status === "todo").length,
    inbox: clips.length,
  };
}

// サイドバーのバッジ用カウント。Entry / Inbox は一覧件数、Task は未完了タスク件数。
// どれか1つの取得に失敗しても 0 にフォールバックしてサイドバー描画は止めない。
export async function getNavCountsServer(): Promise<NavCounts> {
  const pageData = await getAppPageDataServer().catch(() => null);
  return pageData?.navCounts ?? { entry: 0, task: 0, inbox: 0 };
}
