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

// タスクに、所属 entry の会社名・選考情報を添えたビュー用の型。
// backend には全タスク一覧 API が無いため、ここで entry ごとの tasks を集約して組み立てる。
export interface TaskWithContext extends TaskResponse {
  companyName?: string;
}

// ユーザーの全タスクを横断取得する。
// backend に「全タスク一覧」エンドポイントが無いため、
//   1. entries 一覧 (会社名 join 済み) を1回引く
//   2. entry ごとに tasks を取得 (Promise.all で並列化)
// の2段構えで集約する。entry 数ぶんのリクエストになるので、
// 件数が増えたら backend に GET /api/v1/tasks 集約エンドポイントを足すのが望ましい (follow-up)。
// 個別 entry の tasks 取得失敗は握りつぶし、取れたぶんだけ返す。
export async function listAllTasksServer(): Promise<TaskWithContext[]> {
  const entries = await listEntriesWithCompanyNamesServer();

  const perEntry = await Promise.all(
    entries.map(async (entry) => {
      const tasks = await listTasksByEntryServer(entry.id).catch(
        () => [] as TaskResponse[],
      );
      return tasks.map((task) => ({
        ...task,
        companyName: entry.companyName,
      }));
    }),
  );

  return perEntry.flat();
}
