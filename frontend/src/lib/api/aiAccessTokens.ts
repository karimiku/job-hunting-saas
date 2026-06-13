import { apiFetch } from "./client";

export interface AIAccessTokenItem {
  id: string;
  name: string;
  tokenPreview: string;
  lastUsedAt?: string | null;
  revokedAt?: string | null;
  createdAt: string;
}

export interface CreateAIAccessTokenResponse {
  token: string;
  item: AIAccessTokenItem;
}

export async function listAIAccessTokens(): Promise<AIAccessTokenItem[]> {
  const res = await apiFetch<{ tokens: AIAccessTokenItem[] }>("/api/v1/ai-access-tokens");
  return res.tokens;
}

export async function createAIAccessToken(name: string): Promise<CreateAIAccessTokenResponse> {
  return apiFetch<CreateAIAccessTokenResponse>("/api/v1/ai-access-tokens", {
    method: "POST",
    body: JSON.stringify({ name: name.trim() || undefined }),
  });
}

export async function revokeAIAccessToken(tokenId: string): Promise<void> {
  await apiFetch<void>(`/api/v1/ai-access-tokens/${tokenId}`, {
    method: "DELETE",
  });
}
