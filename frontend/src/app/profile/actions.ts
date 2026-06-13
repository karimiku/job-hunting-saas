"use server";

import { revalidatePath } from "next/cache";
import { serverFetch } from "@/lib/api/server";
import { ApiError } from "@/lib/api/client-types";
import type {
  AiAccessTokenResponse,
  CreateAiAccessTokenResponse,
} from "@/lib/api/aiTokens";

export interface CreateAiAccessTokenState {
  error?: string;
  token?: string;
  accessToken?: AiAccessTokenResponse;
  values?: {
    name: string;
  };
}

export interface RevokeAiAccessTokenState {
  error?: string;
  revokedId?: string;
}

function readField(form: FormData, name: string): string {
  const v = form.get(name);
  return typeof v === "string" ? v : "";
}

export async function createAiAccessTokenAction(
  _prev: CreateAiAccessTokenState,
  formData: FormData,
): Promise<CreateAiAccessTokenState> {
  const name = readField(formData, "name").trim();

  try {
    const created = await serverFetch<CreateAiAccessTokenResponse>(
      "/api/v1/ai/tokens",
      {
        method: "POST",
        body: JSON.stringify({ name }),
      },
    );
    revalidatePath("/profile");
    return {
      token: created.token,
      accessToken: created.accessToken,
      values: { name: "" },
    };
  } catch (err) {
    return {
      error: err instanceof ApiError ? err.message : "トークンの作成に失敗しました",
      values: { name },
    };
  }
}

export async function revokeAiAccessTokenAction(
  _prev: RevokeAiAccessTokenState,
  formData: FormData,
): Promise<RevokeAiAccessTokenState> {
  const tokenId = readField(formData, "tokenId").trim();
  if (!tokenId) {
    return { error: "トークンの指定が不正です" };
  }

  try {
    await serverFetch<void>(`/api/v1/ai/tokens/${tokenId}`, {
      method: "DELETE",
    });
    revalidatePath("/profile");
    return { revokedId: tokenId };
  } catch (err) {
    return {
      error: err instanceof ApiError ? err.message : "トークンの失効に失敗しました",
    };
  }
}
