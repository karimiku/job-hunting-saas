import { apiFetch } from "./client";

export interface CompanyResponse {
  id: string;
  name: string;
  memo: string;
  createdAt: string;
  updatedAt: string;
}

export async function listCompanies(): Promise<CompanyResponse[]> {
  const res = await apiFetch<{ companies: CompanyResponse[] }>("/api/v1/companies");
  return res.companies;
}

export async function getCompany(id: string): Promise<CompanyResponse> {
  return apiFetch<CompanyResponse>(`/api/v1/companies/${id}`);
}

export interface CreateCompanyInput {
  name: string;
  memo?: string;
}

export async function createCompany(input: CreateCompanyInput): Promise<CompanyResponse> {
  return apiFetch<CompanyResponse>("/api/v1/companies", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export interface UpdateCompanyInput {
  name?: string;
  memo?: string;
}

export async function updateCompany(id: string, input: UpdateCompanyInput): Promise<CompanyResponse> {
  return apiFetch<CompanyResponse>(`/api/v1/companies/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
}

export async function deleteCompany(id: string): Promise<void> {
  await apiFetch<void>(`/api/v1/companies/${id}`, { method: "DELETE" });
}
