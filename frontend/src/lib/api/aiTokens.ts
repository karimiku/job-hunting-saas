export interface AiAccessTokenResponse {
  id: string;
  name: string;
  tokenPrefix: string;
  createdAt: string;
  lastUsedAt: string | null;
  revokedAt: string | null;
}

export interface CreateAiAccessTokenResponse {
  token: string;
  accessToken: AiAccessTokenResponse;
}
