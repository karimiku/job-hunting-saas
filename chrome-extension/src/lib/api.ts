// Entré バックエンド API クライアント。
// Chrome 拡張からは Session Cookie を共有できるよう host_permissions と credentials: include を使う。

const API_BASE =
  (import.meta.env.VITE_API_BASE_URL as string | undefined) ??
  "http://localhost:8080";

export interface CreateCompanyInput {
  name: string;
  memo?: string;
}

export interface CreateEntryInput {
  companyId: string;
  route: string;
  source: string;
  memo?: string;
  stageKind?: string;
  stageLabel?: string;
}

export interface CompanyResponse {
  id: string;
  name: string;
  memo: string;
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(init.headers ?? {}),
    },
  });
  if (!res.ok) {
    throw new Error(`API ${path} failed: ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export async function createCompany(input: CreateCompanyInput): Promise<CompanyResponse> {
  return request<CompanyResponse>("/api/v1/companies", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export async function createEntry(input: CreateEntryInput): Promise<{ id: string }> {
  return request<{ id: string }>("/api/v1/entries", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

interface UpdateEntryInput {
  stageKind?: string;
  stageLabel?: string;
}

async function updateEntry(id: string, input: UpdateEntryInput): Promise<void> {
  await request<unknown>(`/api/v1/entries/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
}

/** 会社+エントリーを 1 つの操作で保存する。
 *  popup でユーザーが選んだ stageKind/stageLabel は CreateEntryRequest が受けないため、
 *  default 以外を選んだときは create 後に PATCH で反映する (openapi 拡張は別 PR)。 */
export async function saveDetectedJob(input: {
  companyName: string;
  route: string;
  source: string;
  memo?: string;
  stageKind?: string;
  stageLabel?: string;
}): Promise<void> {
  const company = await createCompany({ name: input.companyName });
  const entry = await createEntry({
    companyId: company.id,
    route: input.route,
    source: input.source,
    memo: input.memo,
  });
  if (input.stageKind && input.stageKind !== "application") {
    await updateEntry(entry.id, {
      stageKind: input.stageKind,
      stageLabel: input.stageLabel,
    });
  }
}
