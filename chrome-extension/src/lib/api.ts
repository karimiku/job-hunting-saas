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

/** 会社+エントリーを 1 つの操作で保存する。 */
export async function saveDetectedJob(input: {
  companyName: string;
  route: string;
  source: string;
  memo?: string;
}): Promise<void> {
  const company = await createCompany({ name: input.companyName });
  await createEntry({
    companyId: company.id,
    route: input.route,
    source: input.source,
    memo: input.memo,
  });
}
